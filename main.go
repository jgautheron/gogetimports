package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const usageDoc = `gogetimports: get a JSON-formatted list of imports

Usage:

  gogetimports ARGS <directory>

Flags:

  -only-third-parties  return only third party imports
`

var (
	flagThirdParties = flag.Bool("only-third-parties", false, "return only third party imports")

	// imports contains the list of import path.
	// filename:[]import path
	imports    = map[string][]string{}
	sourcePath = ""
)

func main() {
	flag.Parse()
	log.SetPrefix("gogetimports: ")

	args := flag.Args()
	if len(args) != 1 {
		usage()
	}
	sourcePath = args[0]

	if err := parseTree(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(imports)
}

func usage() {
	fmt.Fprintf(os.Stderr, usageDoc)
	os.Exit(1)
}

func parseTree() error {
	pathLen := len(sourcePath)
	// Parse recursively the given path if the recursive notation is found
	if pathLen >= 5 && sourcePath[pathLen-3:] == "..." {
		filepath.Walk(sourcePath[:pathLen-3], func(p string, f os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				// resume walking
				return nil
			}

			if f.IsDir() {
				parseDir(p)
			}
			return nil
		})
	} else {
		parseDir(sourcePath)
	}
	return nil
}

func parseDir(dir string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for fn, f := range pkg.Files {
			if _, ok := imports[fn]; !ok {
				imports[fn] = make([]string, 0)
			}
			for _, imprt := range f.Imports {
				// Cleanup the import path
				path := strings.Replace(imprt.Path.Value, `"`, "", 2)

				if *flagThirdParties && !isThirdParty(path) {
					continue
				}
				imports[fn] = append(imports[fn], path)
			}
		}
	}

	return nil
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
