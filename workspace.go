package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

// key: original import name, value: new import name
var importMap = map[string]Import{}

func getCWD() string {
	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("jsonnetr could not determine cwd: %s", err.Error()))
		os.Exit(1)
	}
	return cwd
}

func copyFileToWorkspace(cwd string, rootFilename string, workspacePath string, localPath string) string {
	// TODO rz - make un-bad.
	srcPath := path.Join(cwd, localPath, rootFilename)
	destPath := path.Join(workspacePath, localPath, rootFilename)
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
		// os.RemoveAll(ws)
	}()
	logrus.WithField("path", ws).Debug("Created workspace")

	action(ws)
}

func importMapKeys() []string {
	keys := make([]string, 0, len(importMap))
	for k := range importMap {
		keys = append(keys, k)
	}
	return keys
}
