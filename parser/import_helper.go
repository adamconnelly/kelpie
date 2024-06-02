package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"slices"
	"strings"

	kslices "github.com/adamconnelly/kelpie/slices"
)

type importHelper struct {
	typesInfo             *types.Info
	packageNamesToImports map[string]string
	requiredImports       []string
}

func newImportHelper(typesInfo *types.Info, importSpecs []*ast.ImportSpec) *importHelper {
	if typesInfo == nil {
		panic("typesInfo cannot be nil")
	}

	packageNamesToImports := make(map[string]string, len(importSpecs))
	for _, i := range importSpecs {
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

	return &importHelper{
		typesInfo:             typesInfo,
		packageNamesToImports: packageNamesToImports,
	}
}

func (i *importHelper) AddImportsRequiredForType(e ast.Expr) {
	for _, identifier := range i.getPackageIdentifiers(e) {
		packageName := identifier.Name

		// First let's check if the package name is in the import map. This handles standard
		// expressions like `kelpie.Parser`.
		if imp, ok := i.packageNamesToImports[packageName]; ok {
			i.addImport(imp)
			continue
		}

		// If we couldn't find it, the package might be implicit because of a dot import. We
		// can use typesInfo.Uses to find the package for the identifier and look it up that way.
		if use, ok := i.typesInfo.Uses[identifier]; ok {
			pkg := use.Pkg()
			if pkg == nil {
				// If we found a usage but there's no package, it's a built-in type so no import
				// is required. Just ignore and continue.
				continue
			}

			if imp, ok := i.packageNamesToImports[pkg.Path()]; ok {
				i.addImport(imp)
				continue
			}
		}

		panic(fmt.Sprintf("Could not find import statement for identifier %v. This is a bug in Kelpie!", identifier))
	}
}

func (i *importHelper) RequiredImports() []string {
	// Make sure the imports are sorted so that the code generation is stable.
	slices.SortStableFunc(i.requiredImports, strings.Compare)

	return i.requiredImports
}

func (i *importHelper) getPackageIdentifiers(e ast.Expr) []*ast.Ident {
	identifiers := []*ast.Ident{}

	switch n := e.(type) {
	case *ast.Ident:
		identifiers = append(identifiers, n)
	case *ast.ArrayType:
		identifiers = append(identifiers, i.getPackageIdentifiers(n.Elt)...)
	case *ast.StarExpr:
		identifiers = append(identifiers, i.getPackageIdentifiers(n.X)...)
	case *ast.SelectorExpr:
		identifiers = append(identifiers, i.getPackageIdentifiers(n.X)...)
	case *ast.MapType:
		identifiers = append(identifiers, i.getPackageIdentifiers(n.Key)...)
		identifiers = append(identifiers, i.getPackageIdentifiers(n.Value)...)
	default:
		panic(fmt.Sprintf("Could not get package identifier from ast expression: %v. This is a bug in Kelpie!", e))
	}

	return identifiers
}

func (i *importHelper) addImport(imp string) {
	if !kslices.Contains(i.requiredImports, func(i string) bool { return imp == i }) {
		i.requiredImports = append(i.requiredImports, imp)
	}
}
