// TODO rz - retrieve imports async
// TODO rz - matcher plugins
// TODO rz - cleanup, omg, so bad
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
)

const (
	importPattern            = "import\\s\"([a-zA-Z0-9_\\-\\?@\\./:]+)\""
	replacementPatternFormat = "import\\s\"(%s)\""
)

var (
	version       = "0.0.1-alpha"
	importMatcher = regexp.MustCompile(importPattern)
)

func main() {
	if _, ok := os.LookupEnv("JSONNETR_DEBUG"); ok {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("DEBUG LOGGING ENABLED")
	}
	printHelp()

	if len(os.Args) == 0 {
		printHelp()
		os.Exit(1)
	}

	cwd := getCWD()
	withWorkspace(func(ws string) {
		resolvers := resolverRegistry{
			resolvers: []Resolver{
				localFileResolver{
					cwd:           cwd,
					workspacePath: ws,
					localPath:     filepath.Dir(os.Args[len(os.Args)-1]),
				},
				httpResolver{
					workspacePath: ws,
				},
			},
		}

		processImports(os.Args[len(os.Args)-1], resolvers, importReplacer)
		runJsonnet()
	})

	os.Exit(0)
}

func processImports(filename string, resolvers resolverRegistry, replacer Replacer) {
	resolvers.ResolveAll(filename, []string{"."})

	logrus.Debugf("Resolved %d imports", len(sources)-1)
	for _, src := range sources {
		logrus.WithField("path", src.WorkspacePath).Debug("Updating template imports")
		newData := src.Data
		for _, imp := range src.Imports {
			newData = replacer(newData, imp, getWorkspaceImportByOriginalName(imp))
		}
		if err := ioutil.WriteFile(src.WorkspacePath, newData, 0644); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed writing parsed file: %s", err.Error()))
			os.Exit(1)
		}
	}
}

func runJsonnet() {
	args := append([]string(nil), os.Args...)[1 : len(os.Args)-1]
	args = append(args, getRootWorkspaceFile(os.Args[len(os.Args)-1]))
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
