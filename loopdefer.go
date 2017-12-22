// loopdefer reports defers that are called from loops.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fset := token.NewFileSet()
	for _, name := range args {
		f, err := parser.ParseFile(fset, name, nil, 0)
		if err != nil {
			return err
		}
		v := &vis{fset: fset}
		ast.Walk(v, f)
	}
	return nil
}

type vis struct {
	fset       *token.FileSet
	start, end token.Pos
}

func (v *vis) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	// fmt.Printf("%s:\t%T\t%+[2]v\n", v.fset.Position(node.Pos()), node)
	switch x := node.(type) {
	case *ast.ForStmt:
		return &vis{fset: v.fset, start: x.Body.Lbrace, end: x.Body.Rbrace}
	case *ast.FuncLit, *ast.FuncDecl:
		return &vis{fset: v.fset}
	case *ast.DeferStmt:
		if v.start < x.Defer && x.Defer < v.end {
			fmt.Printf("%s:\tdefer use in a loop\n", v.fset.Position(x.Defer))
		}
	}
	return v
}

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: loopdefer file.go [otherfile.go] ...")
		flag.PrintDefaults()
	}
}
