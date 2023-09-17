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
	// TODO: rename to Results
	Returns []ReturnDefinition
}

type ParameterDefinition struct {
	Name string
	Type string
}

type ReturnDefinition struct {
	Name string
	Type string
}

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
	fileNode, _ := parser.ParseFile(fileSet, "", reader, parser.ParseComments)
	ast.Inspect(fileNode, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			if t.Name.IsExported() {
				switch typeSpecType := t.Type.(type) {
				case *ast.InterfaceType:
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

							for _, result := range funcType.Results.List {
								if len(result.Names) > 0 {
									for _, resultName := range result.Names {
										methodDefinition.Returns = append(methodDefinition.Returns, ReturnDefinition{
											Name: resultName.Name,
											Type: result.Type.(*ast.Ident).Name,
										})
									}
								} else {
									methodDefinition.Returns = append(methodDefinition.Returns, ReturnDefinition{
										Type: result.Type.(*ast.Ident).Name,
									})
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
