//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/kris-nova/logger"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

// addDotImport works around the issue with ifacemaker (https://github.com/weaveworks/eksctl/issues/4925) by adding the
// missing import statement for the local package. It works by parsing the Go source and adding a dot import
// (. <packageName>) for the package. The generated source is further passed to `format.Source`.
func addDotImport(code []byte, packageName string) error {
	d := decorator.NewDecorator(token.NewFileSet())
	f, err := d.Parse(code)
	if err != nil {
		return err
	}

	for _, dl := range f.Decls {
		gd, ok := dl.(*dst.GenDecl)
		if !ok {
			logger.Warning("expected type %T; got %T", &dst.GenDecl{}, dl)
			continue
		}

		if gd.Tok == token.IMPORT {
			gd.Specs = append(gd.Specs, &dst.ImportSpec{
				Name: &dst.Ident{
					Name: ".",
				},
				Path: &dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(packageName)},
			})
		}
	}

	restorer := decorator.NewRestorer()
	var buf bytes.Buffer
	if err := restorer.Fprint(&buf, f); err != nil {
		return fmt.Errorf("error restoring source: %w", err)
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error formatting source: %w", src)
	}
	if _, err := os.Stdout.Write(src); err != nil {
		return fmt.Errorf("error writing source to STDOUT: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf("usage: add_import <package-to-import>"))
	}
	source, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(fmt.Errorf("unexpected error reading from STDIN: %w", err))
	}

	packageName := os.Args[1]
	if err = addDotImport(source, packageName); err != nil {
		panic(err)
	}
}
