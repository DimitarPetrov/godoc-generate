package main

import (
	"bufio"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const godocCommentFormat = "// %s missing godoc."

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current working directory: %v", err)
	}

	log.Print(fmt.Sprintf("Adding default go doc to each exported type/func recursively in %s", wd))

	if err := mapDirectory(wd, instrumentDir); err != nil {
		log.Fatalf("error while instrumenting current working directory: %v", err)
	}
}

func instrumentDir(path string) error {
	fset := token.NewFileSet()
	filter := func(info os.FileInfo) bool {
		return testsFilter(info) && generatedFilter(path, info)
	}
	pkgs, err := parser.ParseDir(fset, path, filter, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed parsing go files in directory %s: %v", path, err)
	}

	for _, pkg := range pkgs {
		if err := instrumentPkg(fset, pkg); err != nil {
			return err
		}
	}
	return nil
}

func instrumentPkg(fset *token.FileSet, pkg *ast.Package) error {
	for fileName, file := range pkg.Files {
		sourceFile, err := os.OpenFile(fileName, os.O_TRUNC|os.O_WRONLY, 0664)
		if err != nil {
			return fmt.Errorf("failed opening file %s: %v", fileName, err)
		}
		if err := instrumentFile(fset, file, sourceFile); err != nil {
			return fmt.Errorf("failed instrumenting file %s: %v", fileName, err)
		}
	}
	return nil
}

func instrumentFile(fset *token.FileSet, file *ast.File, out io.Writer) error {
	// Needed because ast does not support floating comments and deletes them.
	// In order to preserve all comments we just pre-parse it to dst which treats them as first class citizens.
	f, err := decorator.DecorateFile(fset, file)
	if err != nil {
		return fmt.Errorf("failed converting file from ast to dst: %v", err)
	}

	dst.Inspect(f, func(n dst.Node) bool {
		switch t := n.(type) {
		case *dst.FuncDecl:
			if t.Name.IsExported() && !containsGoDoc(t.Decs.Start.All(), t.Name.Name) {
				t.Decs.Start.Prepend(fmt.Sprintf(godocCommentFormat, t.Name.Name))
			}
		case *dst.GenDecl:
			if len(t.Specs) == 1 {
				switch s := t.Specs[0].(type) {
				case *dst.TypeSpec:
					if s.Name.IsExported() && !containsGoDoc(t.Decs.Start.All(), s.Name.Name) {
						t.Decs.Start.Prepend(fmt.Sprintf(godocCommentFormat, s.Name.Name))
					}
					return true
				case *dst.ValueSpec:
					if s.Names[0].IsExported() && !containsGoDoc(t.Decs.Start.All(), s.Names[0].Name) {
						t.Decs.Start.Prepend(fmt.Sprintf(godocCommentFormat, s.Names[0].Name))
					}
					return true
				default:
					return true
				}
			}
			for _, spec := range t.Specs {
				switch s := spec.(type) {
				case *dst.TypeSpec:
					if s.Name.IsExported() && !containsGoDoc(s.Decs.Start.All(), s.Name.Name) {
						s.Decs.Start.Prepend(fmt.Sprintf(godocCommentFormat, s.Name.Name))
					}
				case *dst.ValueSpec:
					if s.Names[0].IsExported() && !containsGoDoc(s.Decs.Start.All(), s.Names[0].Name) {
						s.Decs.Start.Prepend(fmt.Sprintf(godocCommentFormat, s.Names[0].Name))
					}
				}
			}
		}
		return true
	})
	return decorator.Fprint(out, f)
}

func containsGoDoc(decs []string, name string) bool {
	for _, dec := range decs {
		if strings.HasPrefix(dec, "// "+name) {
			return true
		}
	}
	return false
}

// Filter excluding go test files from directory
func testsFilter(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

// Filter excluding generated go files from directory.
// Generated file is considered a file which matches one of the following:
// 1. The name of the file contains "generated"
// 2. First line of the file contains "generated" or "GENERATED"
func generatedFilter(path string, info os.FileInfo) bool {
	if strings.Contains(info.Name(), "generated") {
		return false
	}

	f, err := os.Open(path + "/" + info.Name())
	if err != nil {
		panic(fmt.Sprintf("Failed opening file %s: %v", path+"/"+info.Name(), err))
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	line := scanner.Text()

	if strings.Contains(line, "generated") || strings.Contains(line, "GENERATED") {
		return false
	}
	return true
}

func mapDirectory(dir string, operation func(string) error) error {
	return filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == "vendor" {
				return filepath.SkipDir
			}

			if info.IsDir() {
				return operation(path)
			}
			return nil
		})
}
