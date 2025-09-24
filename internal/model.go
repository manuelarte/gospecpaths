package internal

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type Cfg struct {
	packageName   string
	savedFilename string
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

func (c Cfg) SavedFilename() string {
	if c.savedFilename == "" {
		return "gospecpaths.gen.go"
	}

	return c.savedFilename
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

// AddEndpointStruct add the endpoint struct to the generated file.
func (p Path) AddEndpointStruct(f *jen.File) string {
	structName := p.getEndpointStructName()
	fields := make(map[string]jen.Code)

	for _, param := range p.pathItem.Get.Parameters {
		if strings.ToLower(param.In) == "path" {
			paramName := param.Name
			fields[paramName] = jen.Id(paramName).String()
		}
	}

	f.Type().Id(structName).Struct()
	p.createPathFunction(f, structName, fields)

	return structName
}

func (p Path) createPathFunction(f *jen.File, structName string, indexFields map[string]jen.Code) {
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

	body = append(body, jen.Return(jen.Id("message")))

	f.Func().Params(jen.Id("p").Id(structName)).Id("Path").Params(fields...).String().Block(
		body...,
	)
}

func (p Path) getEndpointStructName() string {
	operationID := p.pathItem.Get.OperationId

	return fmt.Sprintf("%s%sParam", strings.ToUpper(operationID[0:1]), operationID[1:])
}
