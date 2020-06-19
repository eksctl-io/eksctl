package importer

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Importer retrieves an object representing a package
type Importer func(path string) (PackageInfo, error)

func ignoreTestFiles(file os.FileInfo) bool {
	return !strings.HasSuffix(file.Name(), "_test.go")
}

// parseDir returns a map of packages
func parseDir(path string) map[string]*ast.Package {
	dir, err := parser.ParseDir(token.NewFileSet(), path, ignoreTestFiles, parser.ParseComments)
	if err != nil {
		panic(errors.Wrapf(err, "At least one error when parsing directory"))
	}
	return dir
}

func pkgName(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}

// dummyImporter returns the most basic Pkg object possible
// Borrowed from "go/doc/doc.go"
func dummyImporter(imports map[string]*ast.Object, path string) (*ast.Object, error) {
	pkg := imports[path]
	if pkg == nil {
		pkg = ast.NewObj(ast.Pkg, pkgName(path))
		pkg.Data = ast.NewScope(nil)
		imports[path] = pkg
	}
	return pkg, nil
}

func handleVariants(decl *ast.GenDecl) []*ast.ValueSpec {
	var variants = []*ast.ValueSpec{}
	for _, spec := range decl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			panic("Exected ValueSpec in Const GenDecl")
		}
		variants = append(variants, valueSpec)
	}
	return variants
}

// VariantMap groups constants together under a name
type VariantMap map[string][]*ast.ValueSpec

var regexpVariantDeclaration = regexp.MustCompile("[vV]alues for `(.*)`")

// handleGenDeclComments handles `type X struct {}` type declarations
// Comments on `GenDecl`s would otherwise be lost because after calling
// `NewPackage` we only have access to `TypeSpec`s _inside_ the `GenDecl`s
func handleGenDeclComments(scope *ast.Scope, fileMap map[string]*ast.File) VariantMap {
	var variants = make(VariantMap)
	for _, f := range fileMap {
		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if genDecl.Tok == token.CONST {
				if m := regexpVariantDeclaration.FindStringSubmatch(genDecl.Doc.Text()); m != nil {
					variants[m[1]] = handleVariants(genDecl)
				}
				continue
			}
			// Check for `type X struct {}`
			if genDecl.Tok == token.TYPE && len(genDecl.Specs) != 1 && !genDecl.Lparen.IsValid() {
				continue
			}
			typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
			if !ok {
				continue
			}
			name := typeSpec.Name.String()
			// It may not be necessary to us the scope to access things but
			// I don't want to assume these pointers point to the same obj
			scopedObj, ok := scope.Objects[name]
			if !ok {
				panic("Unreachable error, file declarations must be in the package scope")
			}
			scopedTypeSpec, ok := scopedObj.Decl.(*ast.TypeSpec)
			if !ok {
				panic("Unreachable error, scoped declaration must be the same type as the GenDecl")
			}
			if scopedTypeSpec.Doc.Text() == "" {
				scopedTypeSpec.Doc = genDecl.Doc
			}
		}

	}
	return variants
}

// PackageInfo holds all of the information we can understand about a package
type PackageInfo struct {
	Pkg      *ast.Object
	Variants VariantMap
}

// NewImporter creates a memoizing function for importing packages
func NewImporter(path string) (Importer, error) {
	importCache := make(map[string]PackageInfo)
	f := func(path string) (info PackageInfo, err error) {
		if cached, ok := importCache[path]; ok {
			return cached, nil
		}
		// Find out where our package is
		imported, err := build.Import(path, ".", build.FindOnly)
		if err != nil {
			return PackageInfo{}, err
		}
		dir := parseDir(imported.Dir)
		// Just take the first (and only) package from that directory
		for _, p := range dir {
			schemaPkg, _ := ast.NewPackage(token.NewFileSet(), p.Files, dummyImporter, nil)
			variants := handleGenDeclComments(schemaPkg.Scope, p.Files)
			name := pkgName(path)

			importCache[path] = PackageInfo{
				Pkg: &ast.Object{
					Kind: ast.Pkg,
					Name: name,
					Decl: nil, // an ImportSpec could go here but we don't need one
					Data: schemaPkg.Scope,
					Type: nil,
				},
				Variants: variants,
			}
			return importCache[path], nil
		}
		panic(errors.New("Unreachable error, imported directory contained no packages"))
	}
	pkg, err := f(path)
	if err != nil {
		return nil, err
	}
	importCache[""] = pkg
	return f, nil
}

func importPathFromSelector(it *ast.SelectorExpr) (path string, name string, err error) {
	// We assume we'll find an import on the lefthand side of the SelectorExpr
	importIdent := it.X.(*ast.Ident)
	if importIdent.Obj == nil {
		return "", "", errors.Errorf("Missing Obj for ident")
	}
	importSpec := importIdent.Obj.Decl.(*ast.ImportSpec)
	if importSpec.Path.Kind != token.STRING {
		return "", "", errors.Errorf("Cannot handle token of type %s as import path", importSpec.Path.Kind)
	}
	// Trim surrounding quotes
	importPath, err := strconv.Unquote(importSpec.Path.Value)
	if err != nil {
		panic(errors.Errorf("Impossible! string not quoted: %s", importSpec.Path.Value))
	}
	return importPath, it.Sel.Name, nil
}

// FindPkgObj takes a name like Struct and looks in the starting package
func (importer Importer) FindPkgObj(typeName string) (*ast.Object, bool) {
	importedPkg, err := importer("")
	if err != nil {
		panic(errors.Wrapf(err, "Error importing starting package!"))
	}

	scope := importedPkg.Pkg.Data.(*ast.Scope)
	obj, ok := scope.Objects[typeName]
	return obj, ok
}

// FindImportedTypeSpec takes a SelectorExpr `pkg.Struct` where `pkg` refers to
// an import and finds the corresponding TypeSpec
func (importer Importer) FindImportedTypeSpec(it *ast.SelectorExpr) (string, *ast.TypeSpec, error) {
	importPath, typeName, err := importPathFromSelector(it)
	if err != nil {
		return "", nil, errors.Wrapf(err, "couldn't get import path")
	}

	importedPkg, err := importer(importPath)
	if err != nil {
		return "", nil, errors.Wrapf(err, "couldn't handle struct field")
	}

	// Look for the righthand side of the SelectorExpr in our imported package
	scope := importedPkg.Pkg.Data.(*ast.Scope)
	obj, ok := scope.Objects[typeName]
	if !ok {
		return "", nil, errors.Errorf("Couldn't find object %s in imported package %s", it.Sel.Name, importPath)
	}
	inlineName := fmt.Sprintf("%s.%s", strings.ReplaceAll(importPath, "/", "|"), it.Sel.Name)
	typeSpec, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		return "", nil, errors.Errorf("Expected TypeSpec for %s", inlineName)
	}
	return inlineName, typeSpec, nil
}
