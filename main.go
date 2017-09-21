// TODO rz - retrieve imports async
// TODO rz - matcher plugins
// TODO rz - cleanup, omg, so bad
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
)

var version = "0.0.1-alpha"

type Replacer func(data []byte, oldImport, newImport string) ([]byte, error)

type Import struct {
	WorkingPath string
	Data        []byte
}

func main() {
	if _, ok := os.LookupEnv("JSONNETR_DEBUG"); ok {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("DEBUG LOGGING ENABLED")
	}
	printHelp()

	// TODO rz - Find all import statements (recursively), saving files to
	// tmp work directory, re-write imports w/ local references
	var rootFile string
	cwd := getCWD()
	withWorkspace(func(ws string) {
		resolvers := resolverRegistry{
			resolvers: []Resolver{
				localFileResolver{
					cwd:           cwd,
					workspacePath: ws,
					localPath:     filepath.Dir(os.Args[len(os.Args)-1]),
				},
			},
		}

		rootFile = copyFileToWorkspace(cwd, os.Args[len(os.Args)-1], ws, "")
		processImports(rootFile, resolvers, nil)
	})

	args := append([]string(nil), os.Args...)[1 : len(os.Args)-1]
	args = append(args, rootFile)
	logrus.WithField("args", strings.Join(args, " ")).Debug("Starting jsonnet call")

	var (
		out []byte
		err error
	)
	if out, err = exec.Command("jsonnet", args...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed to run jsonnet: %s", err.Error()))
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, string(out))
	os.Exit(0)
}

func processImports(filename string, resolvers resolverRegistry, replacer Replacer) {
	imports := findAllImports(filename)
	if len(imports) == 0 {
		return
	}
	resolvers.ResolveAll(imports, []string{"."})

	logrus.Debugf("Resolved %d imports", len(importMap))

	// for oldImport, imp := range importMap {
	// 	newData, err := replacer(imp.Data, oldImport, imp.WorkingPath)
	// 	if err != nil {
	// 		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed replacing import: %s", err.Error()))
	// 		os.Exit(1)
	// 	}
	// 	if err := ioutil.WriteFile(imp.WorkingPath, newData, 0644); err != nil {
	// 		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed writing parsed file: %s", err.Error()))
	// 		os.Exit(1)
	// 	}
	// }
}

func printHelp() {
	for _, arg := range os.Args {
		if arg == "--jsonntr-help" {
			println("jsonnetr " + version)
			println("")
			os.Exit(0)
		}
	}
}

func stringInSlice(s string, sl []string) bool {
	for _, it := range sl {
		if s == it {
			return true
		}
	}
	return false
}

func flattenImportMatches(matches [][]string) []string {
	var l []string
	for _, m := range matches {
		l = append(l, m[1])
	}
	return l
}
