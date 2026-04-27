package internal

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type Cfg struct {
	packageName string
	output      string
}

func DefaultConfig() Cfg {
	return Cfg{
		packageName: "gospecpaths",
	}
}

func (c Cfg) PackageName() string {
	if c.packageName == "" {
		return "gospecpaths"
	}

	return c.packageName
}

func (c Cfg) Output() string {
	if c.output == "" {
		return "gospecpaths.gen.go"
	}

	return c.output
}

func (c Cfg) SetPackageName(name string) Cfg {
	c.packageName = name

	return c
}

func (c Cfg) SetOutput(output string) Cfg {
	c.output = output

	return c
}

type Path struct {
	url      string
	pathItem *v3.PathItem
}

func NewPath(url string, item *v3.PathItem) Path {
	return Path{
		url:      url,
		pathItem: item,
	}
}

// IsValid returns whether the path is valid to be generated in the Path structs.
// The requirements are:
// - Get endpoint.
// - OperationId is filled.
func (p Path) IsValid() bool {
	if p.url == "" {
		return false
	}

	if p.pathItem == nil {
		return false
	}

	hasOperationID := p.pathItem.Get != nil && p.pathItem.Get.OperationId != ""

	return hasOperationID
}

// queryParam holds metadata about a query parameter for code generation.
type queryParam struct {
	name    string
	isArray bool
	explode bool
}

// AddEndpointStruct add the endpoint struct to the generated file.
func (p Path) AddEndpointStruct(f *jen.File) string {
	structName := p.getEndpointStructName()
	pathFields := make(map[string]jen.Code)

	var queryParams []queryParam

	for _, param := range p.pathItem.Get.Parameters {
		switch strings.ToLower(param.In) {
		case "path":
			paramName := param.Name
			pathFields[paramName] = jen.Id(paramName).String()
		case "query":
			isArray := false

			if param.Schema != nil {
				schema := param.Schema.Schema()
				if schema != nil && slices.Contains(schema.Type, "array") {
					isArray = true
				}
			}

			explode := true // OpenAPI default for form style
			if param.Explode != nil {
				explode = *param.Explode
			}

			queryParams = append(queryParams, queryParam{
				name:    param.Name,
				isArray: isArray,
				explode: explode,
			})
		}
	}

	f.Type().Id(structName).Struct()

	if len(queryParams) > 0 {
		f.Line()
		p.createQueryParamsStruct(f, structName, queryParams)
	}

	p.createPathFunction(f, structName, pathFields, len(queryParams) > 0)

	return structName
}

func (p Path) GetURL() string {
	return p.url
}

func (p Path) createPathFunction(f *jen.File, structName string, indexFields map[string]jen.Code, hasQueryParams bool) {
	fields := make([]jen.Code, 0, len(indexFields))

	body := make([]jen.Code, 0)

	body = append(body, jen.Id("message").Op(":=").Lit(p.url))
	for fieldName, field := range indexFields {
		c := jen.Id("message").
			Op("=").
			Qual("strings", "Replace").
			Call(jen.Id("message"), jen.Lit(fmt.Sprintf("{%s}", fieldName)), jen.Id(fieldName), jen.Lit(-1))
		body = append(body, c)
		fields = append(fields, field)
	}

	if hasQueryParams {
		queryParamsStructName := structName + "QueryParams"
		fields = append(fields, jen.Id("queryParams").Id(queryParamsStructName))
		body = append(body,
			jen.If(
				jen.
					Id("queryString").
					Op(":=").
					Id("queryParams").
					Dot("ToQueryString").
					Call().
					Op(";").
					Id("queryString").
					Op("!=").
					Lit(""),
			).
				Block(
					jen.
						Id("message").
						Op("=").
						Id("message").
						Op("+").
						Lit("?").
						Op("+").
						Id("queryString"),
				),
		)
	}

	body = append(body, jen.Return(jen.Id("message")))

	f.Func().Params(jen.Id("p").Id(structName)).Id("Path").Params(fields...).String().Block(
		body...,
	)
}

func (p Path) createQueryParamsStruct(f *jen.File, endpointStructName string, params []queryParam) {
	queryParamsStructName := endpointStructName + "QueryParams"

	fields := make([]jen.Code, 0, len(params))
	for _, param := range params {
		exportedName := exportName(param.name)
		if param.isArray {
			fields = append(fields, jen.Id(exportedName).Index().String())
		} else {
			fields = append(fields, jen.Id(exportedName).String())
		}
	}

	f.Type().Id(queryParamsStructName).Struct(fields...)
	f.Line()

	// Generate ToQueryString method
	body := make([]jen.Code, 0)
	body = append(body, jen.Id("values").Op(":=").Qual("net/url", "Values").Values())

	for _, param := range params {
		exportedName := exportName(param.name)
		if param.isArray {
			if param.explode {
				// Exploded array: ?tags=a&tags=b
				body = append(body,
					jen.For(jen.List(jen.Id("_"), jen.Id("v")).Op(":=").Range().Id("q").Dot(exportedName)).Block(
						jen.Id("values").Dot("Add").Call(jen.Lit(param.name), jen.Id("v")),
					),
				)
			} else {
				// Non-exploded array: ?fields=a,b,c
				body = append(body,
					jen.If(jen.Len(jen.Id("q").Dot(exportedName)).Op(">").Lit(0)).Block(
						jen.Id("values").Dot("Set").Call(
							jen.Lit(param.name),
							jen.Qual("strings", "Join").Call(jen.Id("q").Dot(exportedName), jen.Lit(",")),
						),
					),
				)
			}
		} else {
			body = append(body,
				jen.If(jen.Id("q").Dot(exportedName).Op("!=").Lit("")).Block(
					jen.Id("values").Dot("Set").Call(jen.Lit(param.name), jen.Id("q").Dot(exportedName)),
				),
			)
		}
	}

	body = append(body, jen.Return(jen.Id("values").Dot("Encode").Call()))

	f.Func().Params(jen.Id("q").Id(queryParamsStructName)).Id("ToQueryString").Params().String().Block(
		body...,
	)
	f.Line()
}

func exportName(name string) string {
	if len(name) == 0 {
		return name
	}

	runes := []rune(name)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func (p Path) getEndpointStructName() string {
	operationID := p.pathItem.Get.OperationId

	return fmt.Sprintf("%s%sEndpoint", strings.ToUpper(operationID[0:1]), operationID[1:])
}
