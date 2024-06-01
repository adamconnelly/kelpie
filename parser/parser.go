// Package parser contains the parser for reading Go files and building the interface definitions
// needed to generate mocks.
package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	kslices "github.com/adamconnelly/kelpie/slices"
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

func (m *MockedInterface) addImports(imports []string) {
	for _, imp := range imports {
		if !kslices.Contains(m.Imports, func(i string) bool { return imp == i }) {
			m.Imports = append(m.Imports, imp)
		}
	}
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
	return kslices.Contains(f.InterfacesToInclude, func(n string) bool {
		return n == name
	})
}

// Parse parses the source contained in the reader.
func Parse(packageName string, directory string, filter InterfaceFilter) ([]*MockedInterface, error) {
	var interfaces []*MockedInterface

	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.NeedTypes | packages.NeedImports | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:   directory,
		Tests: true,
	}, "pattern="+packageName)
	if err != nil {
		return nil, errors.Wrap(err, "could not load type information")
	}

	for _, p := range pkgs {
		for _, fileNode := range p.Syntax {
			packageNamesToImports := make(map[string]string, len(fileNode.Imports))
			for _, i := range fileNode.Imports {
				if i.Name != nil && i.Path.Value != "" {
					packageName := i.Name.Name
					if packageName == "." {
						packageName = strings.Trim(i.Path.Value, `"`)
					}

					packageNamesToImports[packageName] = i.Name.Name + ` ` + i.Path.Value
				} else if i.Name != nil {
					packageNamesToImports[i.Name.Name] = `"` + i.Name.Name + `"`
				} else if i.Path.Value != "" {
					packageName := filepath.Base(strings.Trim(i.Path.Value, `"`))
					packageNamesToImports[packageName] = i.Path.Value
				}
			}

			ast.Inspect(fileNode, func(n ast.Node) bool {
				if t, ok := n.(*ast.TypeSpec); ok {
					if t.Name.IsExported() {
						if typeSpecType, ok := t.Type.(*ast.InterfaceType); ok {
							if filter.Include(t.Name.Name) {
								mockedInterface := &MockedInterface{
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

										mockedInterface.addImports(getRequiredImports(param.Type, packageNamesToImports, p.TypesInfo))
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

											mockedInterface.addImports(getRequiredImports(result.Type, packageNamesToImports, p.TypesInfo))
										}
									}

									mockedInterface.Methods = append(mockedInterface.Methods, methodDefinition)
								}

								slices.SortStableFunc(mockedInterface.Imports, strings.Compare)

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
	case *ast.SelectorExpr:
		packageName := getTypeName(n.X)

		return packageName + "." + n.Sel.Name
	}

	panic(fmt.Sprintf("Unknown type %v. This is a bug in Kelpie!", e))
}

func getPackageIdentifiers(e ast.Expr) []*ast.Ident {
	identifiers := []*ast.Ident{}

	switch n := e.(type) {
	case *ast.Ident:
		identifiers = append(identifiers, n)
	case *ast.ArrayType:
		identifiers = append(identifiers, getPackageIdentifiers(n.Elt)...)
	case *ast.StarExpr:
		identifiers = append(identifiers, getPackageIdentifiers(n.X)...)
	case *ast.SelectorExpr:
		identifiers = append(identifiers, getPackageIdentifiers(n.X)...)
	default:
		panic(fmt.Sprintf("Could not get package identifier from ast expression: %v. This is a bug in Kelpie!", e))
	}

	return identifiers
}

func getRequiredImports(e ast.Expr, packageNamesToImports map[string]string, typesInfo *types.Info) []string {
	requiredImports := []string{}
	for _, identifier := range getPackageIdentifiers(e) {
		packageName := identifier.Name

		// First let's check if the package name is in the import map. This handles standard
		// expressions like `kelpie.Parser`.
		if i, ok := packageNamesToImports[packageName]; ok {
			requiredImports = append(requiredImports, i)
			continue
		}

		// If we couldn't find it, the package might be implicit because of a dot import. We
		// can use typesInfo.Uses to find the package for the identifier and look it up that way.
		if use, ok := typesInfo.Uses[identifier]; ok {
			pkg := use.Pkg()
			if pkg == nil {
				// If we found a usage but there's no package, it's a built-in type so no import
				// is required. Just ignore and continue.
				continue
			}

			if i, ok := packageNamesToImports[pkg.Path()]; ok {
				requiredImports = append(requiredImports, i)
				continue
			}
		}

		panic(fmt.Sprintf("Could not find import statement for identifier %v. This is a bug in Kelpie!", identifier))
	}

	return requiredImports
}
