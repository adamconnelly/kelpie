package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"

	"github.com/adamconnelly/kelpie/slices"
)

type MockedInterface struct {
	Name        string
	PackageName string
	Methods     []MethodDefinition
}

type MethodDefinition struct {
	Name       string
	Parameters []ParameterDefinition
	Results    []ResultDefinition
}

type ParameterDefinition struct {
	Name string
	Type string
}

type ResultDefinition struct {
	Name string
	Type string
}

//go:generate go run ../cmd/kelpie generate --source-file parser.go --package github.com/adamconnelly/kelpie/parser --interfaces InterfaceFilter
type InterfaceFilter interface {
	// Include indicates that the specified interface should be included in the set of interfaces
	// to generate.
	Include(name string) bool
}

type IncludingInterfaceFilter struct {
	InterfacesToInclude []string
}

func (f *IncludingInterfaceFilter) Include(name string) bool {
	return slices.Contains(f.InterfacesToInclude, func(n string) bool {
		return n == name
	})
}

func Parse(reader io.Reader, packageName string, filter InterfaceFilter) ([]MockedInterface, error) {
	var interfaces []MockedInterface

	fileSet := token.NewFileSet()
	// TODO: test error handling
	// TODO: handle void methods (seems to panic right now)
	fileNode, _ := parser.ParseFile(fileSet, "", reader, parser.ParseComments)
	ast.Inspect(fileNode, func(n ast.Node) bool {
		if t, ok := n.(*ast.TypeSpec); ok {
			if t.Name.IsExported() {
				if typeSpecType, ok := t.Type.(*ast.InterfaceType); ok {
					if filter.Include(packageName + "." + t.Name.Name) {
						mockedInterface := MockedInterface{
							Name:        t.Name.Name,
							PackageName: strings.ToLower(t.Name.Name),
						}

						for _, method := range typeSpecType.Methods.List {
							methodDefinition := MethodDefinition{
								// When are there multiple names?
								Name: method.Names[0].Name,
							}

							// TODO: check what situation would cause Type to not be ast.FuncType. Maybe ast.Bad?
							funcType := method.Type.(*ast.FuncType)
							for _, param := range funcType.Params.List {
								for _, paramName := range param.Names {
									methodDefinition.Parameters = append(methodDefinition.Parameters, ParameterDefinition{
										Name: paramName.Name,
										Type: param.Type.(*ast.Ident).Name,
									})
								}
							}

							if funcType.Results != nil {
								for _, result := range funcType.Results.List {
									if len(result.Names) > 0 {
										for _, resultName := range result.Names {
											methodDefinition.Results = append(methodDefinition.Results, ResultDefinition{
												Name: resultName.Name,
												Type: result.Type.(*ast.Ident).Name,
											})
										}
									} else {
										methodDefinition.Results = append(methodDefinition.Results, ResultDefinition{
											Type: result.Type.(*ast.Ident).Name,
										})
									}
								}
							}

							mockedInterface.Methods = append(mockedInterface.Methods, methodDefinition)
						}

						interfaces = append(interfaces, mockedInterface)
					}
				}
			}
		}

		return true
	})

	return interfaces, nil
}
