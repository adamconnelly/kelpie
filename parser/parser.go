// Package parser contains the parser for reading Go files and building the interface definitions
// needed to generate mocks.
package parser

import (
	"fmt"
	"go/ast"
	"strconv"
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

	// Imports contains the list of imports required by the mocked interface.
	Imports []string
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

//go:generate go run ../cmd/kelpie generate --package github.com/adamconnelly/kelpie/parser --interfaces InterfaceFilter

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
		Mode:  packages.NeedName | packages.NeedTypes | packages.NeedImports | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:   directory,
		Tests: true,
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
								importHelper := newImportHelper(p.TypesInfo, fileNode.Imports, p)
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
									for paramIndex, param := range funcType.Params.List {
										if len(param.Names) > 0 {
											for _, paramName := range param.Names {
												methodDefinition.Parameters = append(methodDefinition.Parameters, ParameterDefinition{
													Name: paramName.Name,
													Type: getTypeName(param.Type, p),
												})
											}
										} else {
											methodDefinition.Parameters = append(methodDefinition.Parameters, ParameterDefinition{
												Name: "_p" + strconv.Itoa(paramIndex),
												Type: getTypeName(param.Type, p),
											})
										}

										importHelper.AddImportsRequiredForType(param.Type)
									}

									if funcType.Results != nil {
										for _, result := range funcType.Results.List {
											if len(result.Names) > 0 {
												for _, resultName := range result.Names {
													methodDefinition.Results = append(methodDefinition.Results, ResultDefinition{
														Name: resultName.Name,
														Type: getTypeName(result.Type, p),
													})
												}
											} else {
												methodDefinition.Results = append(methodDefinition.Results, ResultDefinition{
													Type: getTypeName(result.Type, p),
												})
											}

											importHelper.AddImportsRequiredForType(result.Type)
										}
									}

									mockedInterface.Methods = append(mockedInterface.Methods, methodDefinition)
								}

								mockedInterface.Imports = importHelper.RequiredImports()

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

func getTypeName(e ast.Expr, p *packages.Package) string {
	switch n := e.(type) {
	case *ast.Ident:
		// Check if this is a type rather than, for example, a package name.
		if _, ok := p.TypesInfo.Types[e]; ok {
			if use, ok := p.TypesInfo.Uses[n]; ok {
				typePackage := use.Pkg()
				if typePackage != nil && typePackage.Path() == p.PkgPath {
					// If the type's package matches the package we're parsing, this is a reference
					// to a type in the same package. We'll need to adjust the type name to include
					// the package name so that it can be referenced correctly from the package
					// generated for the mock.
					return p.Name + "." + n.Name
				}
			}
		}

		return n.Name
	case *ast.ArrayType:
		elementType := getTypeName(n.Elt, p)
		return "[]" + elementType
	case *ast.StarExpr:
		return "*" + getTypeName(n.X, p)
	case *ast.SelectorExpr:
		packageName := getTypeName(n.X, p)

		return packageName + "." + n.Sel.Name
	case *ast.MapType:
		keyType := getTypeName(n.Key, p)
		valueType := getTypeName(n.Value, p)

		return "map[" + keyType + "]" + valueType
	case *ast.FuncType:
		var params []string
		for _, param := range n.Params.List {
			parameterNames := slices.Map(param.Names, func(i *ast.Ident) string { return i.Name })
			if len(parameterNames) > 0 {
				params = append(params, strings.Join(parameterNames, ", ")+" "+getTypeName(param.Type, p))
			} else {
				params = append(params, getTypeName(param.Type, p))
			}
		}

		var results []string
		if n.Results != nil {
			for _, result := range n.Results.List {
				resultNames := slices.Map(result.Names, func(i *ast.Ident) string { return i.Name })
				if len(resultNames) > 0 {
					results = append(results, strings.Join(resultNames, ", ")+" "+getTypeName(result.Type, p))
				} else {
					results = append(results, getTypeName(result.Type, p))
				}
			}
		}

		functionDefinition := "func(" + strings.Join(params, ", ") + ")"
		if len(results) > 0 {
			functionDefinition += " (" + strings.Join(results, ", ") + ")"
		}

		return functionDefinition
	case *ast.Ellipsis:
		return "..." + getTypeName(n.Elt, p)
	}

	panic(fmt.Sprintf("Unknown type %v. This is a bug in Kelpie!", e))
}
