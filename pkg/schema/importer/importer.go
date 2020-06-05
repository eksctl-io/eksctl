package importer

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Importer retrieves an object representing a package
type Importer func(path string) (*ast.Object, error)

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

// copyGenDeclComments handles `type X struct {}` type declarations
// Comments on `GenDecl`s would otherwise be lost because after calling
// `NewPackage` we only have access to `TypeSpec`s _inside_ the `GenDecl`s
func copyGenDeclComments(scope *ast.Scope, fileMap map[string]*ast.File) {
	for _, f := range fileMap {
		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			// Check for `type X struct {}`
			if len(genDecl.Specs) != 1 && !genDecl.Lparen.IsValid() {
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
}

// NewImporter creates a memoizing function for importing packages
func NewImporter() Importer {
	importCache := make(map[string]*ast.Object)
	return func(path string) (pkg *ast.Object, err error) {
		if importCache[path] != nil {
			return importCache[path], nil
		}
		// Find out where our package is
		imported, err := build.Import(path, ".", build.FindOnly)
		if err != nil {
			return nil, err
		}
		dir := parseDir(imported.Dir)
		// Just take the first (and only) package from that directory
		for _, p := range dir {
			schemaPkg, _ := ast.NewPackage(token.NewFileSet(), p.Files, dummyImporter, nil)
			copyGenDeclComments(schemaPkg.Scope, p.Files)
			name := pkgName(path)

			importCache[path] = &ast.Object{
				Kind: ast.Pkg,
				Name: name,
				Decl: nil, // an ImportSpec could go here but we don't need one
				Data: schemaPkg.Scope,
				Type: nil,
			}
			return importCache[path], nil
		}
		panic(errors.New("Unreachable error, imported directory contained no packages"))
	}
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
	scope := importedPkg.Data.(*ast.Scope)
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
