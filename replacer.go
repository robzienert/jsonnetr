package main

import (
	"fmt"
	"regexp"

	"github.com/Sirupsen/logrus"
)

type Replacer func(data []byte, oldImport, newImport string) []byte

func importReplacer(data []byte, oldImport, newImport string) []byte {
	logrus.WithFields(logrus.Fields{
		"old": oldImport,
		"new": newImport,
	}).Debug("Replacing import statement")
	r := regexp.MustCompile(fmt.Sprintf(replacementPatternFormat, oldImport))
	return []byte(r.ReplaceAllString(string(data), "import \""+newImport+"\""))
}
