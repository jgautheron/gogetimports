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

const usageDoc = `gogetimports: get a JSON-formatted map of imports per file

Usage:

  gogetimports ARGS <directory>

Flags:

  -only-third-parties  return only third party imports
  -list                return a list instead of a map
  -ignore              ignore imports matching the given regular expression

Examples:

  gogetimports ./...
  gogetimports -only-third-parties $GOPATH/src/github.com/cockroachdb/cockroach
`

var (
	flagThirdParties = flag.Bool("only-third-parties", false, "return only third party imports")
	flagList         = flag.Bool("list", false, "return a list instead of a map")
	flagIgnore       = flag.String("ignore", "", "ignore imports matching the given regular expression")

	// imports contains the list of import path.
	// filename[]import path
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

	removeDuplicates := func(list []string) []string {
		encountered, result := map[string]bool{}, []string{}

		for el := range list {
			if _, ok := encountered[list[el]]; ok {
				continue
			}

			encountered[list[el]] = true
			result = append(result, list[el])
		}

		return result
	}

	var output interface{}
	if *flagList {
		lst := []string{}
		for _, mp := range imports {
			lst = append(lst, mp...)
		}
		output = removeDuplicates(lst)
	} else {
		output = imports
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(output)
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

				if len(*flagIgnore) != 0 {
					match, err := regexp.MatchString(*flagIgnore, path)
					if err != nil {
						return err
					}
					if match {
						continue
					}
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
