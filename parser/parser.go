// Package parser contains the parser for reading Go files and building the interface definitions
// needed to generate mocks.
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

// MockedInterface represents an interface that a mock should be generated for.
type MockedInterface struct {
	// Name contains the name of the interface.
	Name string

	// PackageName contains the name of the package that the interface belongs to.
	PackageName string

	// Methods contains the list of methods in the interface.
	Methods []MethodDefinition

	// Imports contains a list of imports required by the mocked interface.
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
func Parse(reader io.Reader, filter InterfaceFilter) ([]MockedInterface, error) {
	var interfaces []MockedInterface

	fileSet := token.NewFileSet()

	// For now ignore any errors on the grounds that we can still try to generate the mock,
	// even if the Go code won't compile.
	fileNode, _ := parser.ParseFile(fileSet, "", reader, parser.ParseComments)

	imports := make(map[string]string)

	ast.Inspect(fileNode, func(n ast.Node) bool {
		if i, ok := n.(*ast.ImportSpec); ok {
			importValue := i.Path.Value
			lastSlashIndex := strings.LastIndex(i.Path.Value, "/")
			importName := i.Path.Value[lastSlashIndex+1 : len(i.Path.Value)-1]
			if i.Name != nil && i.Name.Name != "" {
				importName = i.Name.Name
				importValue = importName + " " + i.Path.Value
			}

			imports[importName] = importValue
		}

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
									paramTypeImport, paramType := getImportAndType(param.Type, imports)
									if paramTypeImport != "" {
										mockedInterface.Imports = append(mockedInterface.Imports, paramTypeImport)
									}

									methodDefinition.Parameters = append(methodDefinition.Parameters, ParameterDefinition{
										Name: paramName.Name,
										Type: paramType,
									})
								}
							}

							if funcType.Results != nil {
								for _, result := range funcType.Results.List {
									resultTypeImport, resultType := getImportAndType(result.Type, imports)
									if resultTypeImport != "" {
										mockedInterface.Imports = append(mockedInterface.Imports, resultTypeImport)
									}

									if len(result.Names) > 0 {
										for _, resultName := range result.Names {
											methodDefinition.Results = append(methodDefinition.Results, ResultDefinition{
												Name: resultName.Name,
												Type: resultType,
											})
										}
									} else {
										methodDefinition.Results = append(methodDefinition.Results, ResultDefinition{
											Type: resultType,
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

func getImportAndType(e ast.Expr, imports map[string]string) (string, string) {
	switch n := e.(type) {
	case *ast.SelectorExpr:
		i := n.X.(*ast.Ident)

		return imports[i.Name], i.Name + "." + n.Sel.Name
	case *ast.StarExpr:
		i, t := getImportAndType(n.X, imports)

		return i, "*" + t
	case *ast.ArrayType:
		i, t := getImportAndType(n.Elt, imports)

		return i, "[]" + t
	}

	return "", getTypeName(e)
}

func getTypeName(e ast.Expr) string {
	if n, ok := e.(*ast.Ident); ok {
		return n.Name
	}

	panic(fmt.Sprintf("Unknown type %v. This is a bug in Kelpie!", e))
}
