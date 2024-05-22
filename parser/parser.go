// Package parser contains the parser for reading Go files and building the interface definitions
// needed to generate mocks.
package parser

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	"github.com/adamconnelly/kelpie/slices"
)

// MockedInterface represents an interface that a mock should be generated for.
type MockedInterface struct {
	// Name contains the name of the interface.
	Name string

	// PackageName contains the name of the package that the interface belongs to.
	PackageName string

	// Methods contains the list of methods in the interface.
	Methods []MethodDefinition
}

// MethodDefinition defines a method in an interface.
type MethodDefinition struct {
	// Name is the name of the method.
	Name string

	// Parameters contains the parameters passed to the method.
	Parameters []ParameterDefinition

	// Results contains the method results.
	Results []ResultDefinition

	// Comment contains any comments added to the method.
	Comment string
}

// ParameterDefinition contains information about a method parameter.
type ParameterDefinition struct {
	// Name is the name of the parameter.
	Name string

	// Type is the parameter's type.
	Type string
}

// ResultDefinition contains information about a method result.
type ResultDefinition struct {
	// Name is the name of the method result. This can be empty if the result is not named.
	Name string

	// Type is the type of the result.
	Type string
}

//go:generate go run ../cmd/kelpie generate --interfaces InterfaceFilter

// InterfaceFilter is used to decide which interfaces mocks should be generated for.
type InterfaceFilter interface {
	// Include indicates that the specified interface should be included in the set of interfaces
	// to generate.
	Include(name string) bool
}

// IncludingInterfaceFilter is an InterfaceFilter that works based on an allow-list of interface
// names.
type IncludingInterfaceFilter struct {
	InterfacesToInclude []string
}

// Include returns true if the specified interface should be mocked, false otherwise.
func (f *IncludingInterfaceFilter) Include(name string) bool {
	return slices.Contains(f.InterfacesToInclude, func(n string) bool {
		return n == name
	})
}

// Parse parses the source contained in the reader.
func Parse(packageName string, directory string, filter InterfaceFilter) ([]MockedInterface, error) {
	var interfaces []MockedInterface

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedImports | packages.NeedSyntax,
		Dir:  directory,
	}, "pattern="+packageName)
	if err != nil {
		return nil, errors.Wrap(err, "could not load type information")
	}

	for _, p := range pkgs {
		for _, fileNode := range p.Syntax {
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
										Name:    method.Names[0].Name,
										Comment: strings.TrimSuffix(method.Doc.Text(), "\n"),
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
		}
	}

	return interfaces, nil
}

func getTypeName(e ast.Expr) string {
	switch n := e.(type) {
	case *ast.Ident:
		return n.Name
	case *ast.ArrayType:
		return "[]" + n.Elt.(*ast.Ident).Name
	case *ast.StarExpr:
		return "*" + getTypeName(n.X)
	}

	panic(fmt.Sprintf("Unknown array element type %v. This is a bug in Kelpie!", e))
}
