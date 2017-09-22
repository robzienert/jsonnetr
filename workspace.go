package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
)

func getCWD() string {
	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr could not determine cwd: %s", err.Error()))
		os.Exit(1)
	}
	return cwd
}

// TODO rz - make unbad
func copyFileToWorkspace(cwd string, rootFilename string, workspacePath string, localPath string) string {
	var (
		srcPath  string
		destPath string
	)
	if strings.HasPrefix(rootFilename, localPath) {
		srcPath = path.Join(cwd, rootFilename)
		destPath = path.Join(workspacePath, rootFilename)
	} else {
		srcPath = path.Join(cwd, localPath, rootFilename)
		destPath = path.Join(workspacePath, localPath, rootFilename)
	}

	logrus.WithFields(logrus.Fields{
		"src":  srcPath,
		"dest": destPath,
	}).Debug("Copying file to workspace")

	rootFileSrc, err := os.Open(srcPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed opening root file: %s", err.Error()))
		os.Exit(1)
	}

	basedir, _ := path.Split(destPath)
	if err := os.MkdirAll(basedir, os.ModePerm); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed to create directories: %s", err.Error()))
		os.Exit(1)
	}

	rootFile, err := os.Create(destPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed creating root file copy: %s", err.Error()))
		os.Exit(1)
	}
	if _, err := io.Copy(rootFile, rootFileSrc); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed copying root file to workspace: %s", err.Error()))
		os.Exit(1)
	}
	if err := rootFile.Sync(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr failed to copy file contents to workspace: %s", err.Error()))
		os.Exit(1)
	}

	return destPath
}

func withWorkspace(action func(ws string)) {
	ws, err := ioutil.TempDir("", "jsonnetr")
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr could not create temporary workspace: %s", err.Error()))
		os.Exit(1)
	}
	defer func() {
		logrus.Debug("Cleaning up workspace")
		os.RemoveAll(ws)
	}()
	logrus.WithField("path", ws).Debug("Created workspace")

	action(ws)
}

func getRootWorkspaceFile(rootFilename string) string {
	for _, s := range sources {
		if s.OriginalPath == rootFilename {
			return s.WorkspacePath
		}
	}
	fmt.Fprintf(os.Stderr, "jsonnetr could not find root workspace file: %s\n", rootFilename)
	os.Exit(1)
	return "" // Dumb go.
}

func getWorkspaceImportByOriginalName(name string) string {
	for _, s := range sources {
		if s.OriginalPath == name {
			return s.WorkspacePath
		}
	}
	fmt.Fprintf(os.Stderr, "jsonnetr could not find workspace file for import: %s\n", name)
	os.Exit(1)
	return ""
}
