package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
)

const usage = `gogetimports: get a json-formatted list of imports

Usage:

  gogetimports -path <directory>

Flags:

  -path                path to be scanned for imports
  -only-third-parties  return only third party imports

`

var (
	flagPath         = flag.String("path", "./", "path to be scanned for imports")
	flagThirdParties = flag.Bool("only-third-parties", false, "return only third party imports")

	// imports contains the list of import path.
	// filename:[]import path
	imports = map[string][]string{}
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()
	log.SetPrefix("gogetimports: ")

	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := parseTree(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func parseTree() error {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, *flagPath, nil, 0)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for fn, f := range pkg.Files {
			ast.Walk(&ImportVisitor{
				packageName: pkg.Name,
				fileName:    fn,
			}, f)
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(imports)

	return nil
}

// ImportVisitor carries the package name and file name
// for passing it to the imports map.
type ImportVisitor struct {
	packageName, fileName string
}

// Visit matches the ImportSpec and fills the imports map.
func (v *ImportVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch t := node.(type) {
		case *ast.ImportSpec:
			// Cleanup the import path
			path := strings.Replace(t.Path.Value, `"`, "", 2)

			if *flagThirdParties && !isThirdParty(path) {
				return v
			}

			_, ok := imports[v.fileName]
			if !ok {
				imports[v.fileName] = make([]string, 0)
			}

			imports[v.fileName] = append(imports[v.fileName], path)
		}
	}
	return v
}

// isThirdParty determines if the given import path is a third party or not.
// It's safe to assume that if the first path of the import path looks like a domain name,
// then we're dealing with a third party.
func isThirdParty(path string) bool {
	r, err := regexp.Compile(`^(\w+)\.(\w+)/`)
	if err != nil {
		return false
	}
	return r.MatchString(path)
}
