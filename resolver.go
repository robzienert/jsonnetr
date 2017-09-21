package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
)

const (
	importPattern = "import\\s\"([a-zA-Z0-9_\\-\\?@\\.]+)\""
)

var importMatcher = regexp.MustCompile(importPattern)

func findAllImports(filename string) []string {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed reading file: %s", err.Error()))
		os.Exit(1)
	}

	logrus.WithField("path", path.Base(filename)).Debugf("Searching for imports: \n%s", string(dat))
	importMatches := importMatcher.FindAllStringSubmatch(string(dat), -1)
	if len(importMatches) == 0 {
		logrus.WithField("path", filename).Debug("No imports found")
		return []string{}
	}

	imports := flattenImportMatches(importMatches)
	logrus.WithField("imports", imports).Debugf("Found %d imports", len(imports))

	return imports
}

type Resolver interface {
	Name() string
	Supports(filepath string) bool
	Resolve(filepath string) (string, []byte, error)
}

type resolverRegistry struct {
	resolvers []Resolver
}

func (registry resolverRegistry) ResolveAll(imports []string, path []string) {
	// TODO detect cycles
	for _, imp := range imports {
		if _, ok := importMap[imp]; ok {
			continue
		}

		for _, r := range registry.resolvers {
			if r.Supports(imp) {
				workspaceFilepath, contents, err := r.Resolve(imp)
				if err != nil {
					fmt.Fprintf(os.Stderr, "jsonnetr could not resolve import '%s' with resolver '%s', at path '%s': %s", imp, r.Name(), strings.Join(path, "->"), err.Error())
					os.Exit(1)
				}
				importMap[imp] = Import{
					WorkingPath: workspaceFilepath,
					Data:        contents,
				}
				p := path
				p = append(path, imp)

				registry.ResolveAll(findAllImports(workspaceFilepath), p)
				goto nextImport
			}
		}
		fmt.Fprintf(os.Stderr, "jsonnetr could not find resolver for import '%s' at path '%s'\n", imp, strings.Join(path, "->"))
		os.Exit(1)
	nextImport:
	}
}

type localFileResolver struct {
	cwd           string
	workspacePath string
	localPath     string
}

func (r localFileResolver) Name() string {
	return "localFile"
}

func (r localFileResolver) Supports(filepath string) bool {
	u, err := url.Parse(filepath)
	if err != nil {
		return false
	}
	return u.Scheme == ""
}

func (r localFileResolver) Resolve(filepath string) (string, []byte, error) {
	newPath := copyFileToWorkspace(r.cwd, filepath, r.workspacePath, r.localPath)
	dat, err := ioutil.ReadFile(newPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed reading file: %s", err.Error()))
		return newPath, nil, err
	}

	return newPath, dat, nil
}
