package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
)

var sources = []Source{}

type Source struct {
	OriginalPath  string
	WorkspacePath string
	Data          []byte
	Imports       []string
}

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
	// TODO this signature is dumb
	Resolve(filepath string) (string, []byte, error)
}

type resolverRegistry struct {
	resolvers []Resolver
}

func (registry resolverRegistry) ResolveAll(sourceFilename string, path []string) {
	// TODO don't resolve if it already exists in sources list
	for _, r := range registry.resolvers {
		if r.Supports(sourceFilename) {
			workspaceFilepath, contents, err := r.Resolve(sourceFilename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "jsonnetr could not resolve source file '%s' with resolver '%s' at path '%s': %s", sourceFilename, r.Name(), strings.Join(path, "->"), err.Error())
				os.Exit(1)
			}
			imports := findAllImports(workspaceFilepath)
			sources = append(sources, Source{
				OriginalPath:  sourceFilename,
				WorkspacePath: workspaceFilepath,
				Data:          contents,
				Imports:       imports,
			})

			for _, imp := range imports {
				// TODO rz - only resolve if not already resolved
				registry.ResolveAll(imp, append(path, sourceFilename))
			}
			return
		}
	}

	fmt.Fprintf(os.Stderr, "jsonnsetr could not find resolver for source '%s' at path '%s'\n", sourceFilename, strings.Join(path, "->"))
	os.Exit(1)
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

type httpResolver struct {
	workspacePath string
}

func (r httpResolver) Name() string {
	return "http"
}

func (r httpResolver) Supports(filepath string) bool {
	u, err := url.Parse(filepath)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func (r httpResolver) Resolve(filepath string) (string, []byte, error) {
	resp, err := http.Get(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed getting remote http file: %s", err.Error()))
		return "", nil, err
	}
	defer resp.Body.Close()

	f, err := ioutil.TempFile(r.workspacePath, "http")
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed creating temp file for http import: %s", err.Error()))
		return "", nil, err
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed copying http import contents to file: %s", err.Error()))
		return f.Name(), nil, err
	}

	dat, err := ioutil.ReadFile(f.Name())
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed reading file: %s", err.Error()))
		return f.Name(), nil, err
	}

	return f.Name(), dat, nil
}
