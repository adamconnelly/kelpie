// Package parser contains the parser for reading Go files and building the interface definitions
// needed to generate mocks.
package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	"github.com/adamconnelly/kelpie/maps"
	"github.com/adamconnelly/kelpie/slices"
)

// ParsedPackage contains the information needed to generate mocks from parsing a package.
type ParsedPackage struct {
	// PackageDirectory is the directory containing the package source files.
	PackageDirectory string

	// Mocks are the mocks that were parsed from the package.
	Mocks []MockedInterface
}

// MockedInterface represents an interface that a mock should be generated for.
type MockedInterface struct {
	// Name contains the name of the interface.
	Name string

	// FullName contains the full name of the interface. For interfaces nested inside a struct,
	// this will contain the full dot-separated path to the interface, for example `UserService.ConfigRepository`.
	FullName string

	// PackageName contains the name of the package that the interface belongs to.
	PackageName string

	// Methods contains the list of methods in the interface.
	Methods []MethodDefinition

	// Imports contains the list of imports required by the mocked interface.
	Imports []string
}

// AnyMethodsHaveParameters returns true if at least one method in the interface has at least
// one parameter.
func (i MockedInterface) AnyMethodsHaveParameters() bool {
	for _, method := range i.Methods {
		if len(method.Parameters) > 0 {
			return true
		}
	}

	return false
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

	// IsVariadic indicates that this is the variable argument to a variadic function.
	IsVariadic bool

	// IsNonEmptyInterface indicates that the parameter type is an interface with at least one method.
	IsNonEmptyInterface bool
}

// ResultDefinition contains information about a method result.
type ResultDefinition struct {
	// Name is the name of the method result. This can be empty if the result is not named.
	Name string

	// Type is the type of the result.
	Type string
}

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
func Parse(packageName string, directory string, filter InterfaceFilter) (*ParsedPackage, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.NeedName | packages.NeedTypes | packages.NeedImports | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedFiles,
		Dir:   directory,
		Tests: true,
	}, "pattern="+packageName)
	if err != nil {
		return nil, errors.Wrap(err, "could not load type information")
	}

	var packageDirectory string
	interfaces := map[string]MockedInterface{}

	for _, p := range pkgs {
		if len(p.Syntax) > 0 && len(p.GoFiles) > 0 {
			sourceDirectory := filepath.Dir(p.GoFiles[0])
			if filepath.Base(sourceDirectory) == p.Name {
				packageDirectory = filepath.Dir(p.GoFiles[0])
			}
		}

		for _, fileNode := range p.Syntax {
			ast.Inspect(fileNode, func(n ast.Node) bool {
				if t, ok := n.(*ast.TypeSpec); ok {
					if t.Name.IsExported() {
						if interfaceType, ok := t.Type.(*ast.InterfaceType); ok {
							if filter.Include(t.Name.Name) {
								i := parseInterface(t.Name.Name, t.Name.Name, interfaceType, p, fileNode.Imports)
								interfaces[i.FullName] = i
							}
						} else if structType, ok := t.Type.(*ast.StructType); ok {
							for _, f := range structType.Fields.List {
								for _, i := range parseStructField(t, f, p, fileNode.Imports, filter) {
									interfaces[i.FullName] = i
								}
							}
						}
					}

					// As soon as we've found a type, we don't need to continue traversing down
					// this path of the tree since we'll already have gotten all the info we
					// need from the type node by now.
					return false
				}

				return true
			})
		}
	}

	return &ParsedPackage{PackageDirectory: packageDirectory, Mocks: maps.Values(interfaces)}, nil
}

func parseStructField(structNode *ast.TypeSpec, field *ast.Field, pkg *packages.Package, importSpecs []*ast.ImportSpec, filter InterfaceFilter) []MockedInterface {
	var interfaces []MockedInterface
	if structTypeInfo, ok := pkg.TypesInfo.Defs[structNode.Name]; ok {
		if interfaceType, ok := field.Type.(*ast.InterfaceType); ok {
			fullName := structTypeInfo.Name() + "." + field.Names[0].Name
			if filter.Include(fullName) {
				parsedInterface := parseInterface(field.Names[0].Name, fullName, interfaceType, pkg, importSpecs)
				interfaces = append(interfaces, parsedInterface)
			}
		} else if structType, ok := field.Type.(*ast.StructType); ok {
			for _, f := range structType.Fields.List {
				interfaces = append(interfaces, parseNestedStructField(structTypeInfo.Name()+"."+field.Names[0].Name+".", f, pkg, importSpecs, filter)...)
			}
		}
	}

	return interfaces
}

func parseNestedStructField(prefix string, field *ast.Field, pkg *packages.Package, importSpecs []*ast.ImportSpec, filter InterfaceFilter) []MockedInterface {
	var interfaces []MockedInterface
	if interfaceType, ok := field.Type.(*ast.InterfaceType); ok {
		fullName := prefix + field.Names[0].Name
		if filter.Include(fullName) {
			parsedInterface := parseInterface(field.Names[0].Name, fullName, interfaceType, pkg, importSpecs)
			interfaces = append(interfaces, parsedInterface)
		}
	} else if structType, ok := field.Type.(*ast.StructType); ok {
		for _, f := range structType.Fields.List {
			interfaces = append(interfaces, parseNestedStructField(prefix+field.Names[0].Name+".", f, pkg, importSpecs, filter)...)
		}
	}

	return interfaces
}

func parseInterface(name, fullName string, i *ast.InterfaceType, p *packages.Package, imports []*ast.ImportSpec) MockedInterface {
	importHelper := newImportHelper(p.TypesInfo, imports, p)
	mockedInterface := MockedInterface{
		Name:        name,
		FullName:    fullName,
		PackageName: strings.ToLower(name),
	}

	for _, method := range i.Methods.List {
		methodDefinition := MethodDefinition{
			Name:    method.Names[0].Name,
			Comment: strings.TrimSuffix(method.Doc.Text(), "\n"),
		}

		funcType := method.Type.(*ast.FuncType)
		for paramIndex, param := range funcType.Params.List {
			if len(param.Names) > 0 {
				for _, paramName := range param.Names {
					typeInfo := getTypeInfo(param.Type, p)
					methodDefinition.Parameters = append(methodDefinition.Parameters, ParameterDefinition{
						Name:                paramName.Name,
						Type:                typeInfo.name,
						IsVariadic:          typeInfo.isVariadic,
						IsNonEmptyInterface: typeInfo.isNonEmptyInterface,
					})
				}
			} else {
				typeInfo := getTypeInfo(param.Type, p)
				methodDefinition.Parameters = append(methodDefinition.Parameters, ParameterDefinition{
					Name:                "_p" + strconv.Itoa(paramIndex),
					Type:                typeInfo.name,
					IsVariadic:          typeInfo.isVariadic,
					IsNonEmptyInterface: typeInfo.isNonEmptyInterface,
				})
			}

			importHelper.AddImportsRequiredForType(param.Type)
		}

		if funcType.Results != nil {
			for _, result := range funcType.Results.List {
				if len(result.Names) > 0 {
					typeInfo := getTypeInfo(result.Type, p)
					for _, resultName := range result.Names {
						methodDefinition.Results = append(methodDefinition.Results, ResultDefinition{
							Name: resultName.Name,
							Type: typeInfo.name,
						})
					}
				} else {
					typeInfo := getTypeInfo(result.Type, p)
					methodDefinition.Results = append(methodDefinition.Results, ResultDefinition{
						Type: typeInfo.name,
					})
				}

				importHelper.AddImportsRequiredForType(result.Type)
			}
		}

		mockedInterface.Methods = append(mockedInterface.Methods, methodDefinition)
	}

	mockedInterface.Imports = importHelper.RequiredImports()

	return mockedInterface
}

type typeInfo struct {
	name                string
	isVariadic          bool
	isNonEmptyInterface bool
}

func getTypeInfo(e ast.Expr, p *packages.Package) typeInfo {
	if ellipsis, ok := e.(*ast.Ellipsis); ok {
		return typeInfo{
			name:       getTypeName(ellipsis.Elt, p),
			isVariadic: true,
		}
	}

	return typeInfo{
		name:                getTypeName(e, p),
		isNonEmptyInterface: isNonEmptyInterface(e, p),
	}
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
	case *ast.InterfaceType:
		// This is maybe a bit of a simplification. We might need to actually take a look at the fields.
		return "interface{}"
	}

	panic(fmt.Sprintf("Unknown type %v. This is a bug in Kelpie!", e))
}

func isNonEmptyInterface(e ast.Expr, p *packages.Package) bool {
	if t, ok := p.TypesInfo.Types[e]; ok {
		if namedType, ok := t.Type.(*types.Named); ok {
			if i, ok := namedType.Underlying().(*types.Interface); ok {
				return !i.Empty()
			}
		}
	}

	return false
}
