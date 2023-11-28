package parser

import (
	"fmt"
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

//go:generate go run ../cmd/kelpie generate --interfaces InterfaceFilter
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

func Parse(reader io.Reader, filter InterfaceFilter) ([]MockedInterface, error) {
	var interfaces []MockedInterface

	fileSet := token.NewFileSet()

	// For now ignore any errors on the grounds that we can still try to generate the mock,
	// even if the Go code won't compile.
	fileNode, _ := parser.ParseFile(fileSet, "", reader, parser.ParseComments)
	ast.Inspect(fileNode, func(n ast.Node) bool {
		if t, ok := n.(*ast.TypeSpec); ok {
			if t.Name.IsExported() {
				if typeSpecType, ok := t.Type.(*ast.InterfaceType); ok {
					if filter.Include(t.Name.Name) {
						mockedInterface := MockedInterface{
							Name:        t.Name.Name,
							PackageName: strings.ToLower(t.Name.Name),
						}

						for _, method := range typeSpecType.Methods.List {
							methodDefinition := MethodDefinition{
								// When are there multiple names?
								Name: method.Names[0].Name,
							}

							getTypeName := func(e ast.Expr) string {
								switch n := e.(type) {
								case *ast.Ident:
									return n.Name
								case *ast.ArrayType:
									return "[]" + n.Elt.(*ast.Ident).Name
								}

								panic(fmt.Sprintf("Unknown array element type %v. This is a bug in Kelpie!", e))
							}

							// TODO: check what situation would cause Type to not be ast.FuncType. Maybe ast.Bad?
							funcType := method.Type.(*ast.FuncType)
							for _, param := range funcType.Params.List {
								for _, paramName := range param.Names {
									methodDefinition.Parameters = append(methodDefinition.Parameters, ParameterDefinition{
										Name: paramName.Name,
										Type: getTypeName(param.Type),
									})
								}
							}

							if funcType.Results != nil {
								for _, result := range funcType.Results.List {
									if len(result.Names) > 0 {
										for _, resultName := range result.Names {
											methodDefinition.Results = append(methodDefinition.Results, ResultDefinition{
												Name: resultName.Name,
												Type: getTypeName(result.Type),
											})
										}
									} else {
										methodDefinition.Results = append(methodDefinition.Results, ResultDefinition{
											Type: getTypeName(result.Type),
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
