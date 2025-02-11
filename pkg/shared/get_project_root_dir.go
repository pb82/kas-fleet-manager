//go:build !windows
// +build !windows

package shared

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
)

// GetProjectRootDir returns the root directory of the project.
// The root directory of the project is the directory that contains the go.mod file which contains
// the "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager" module name.
func GetProjectRootDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		glog.Fatal(err)
	}
	dirs := strings.Split(workingDir, "/")
	var goModPath string
	var rootPath string
	for _, d := range dirs {
		rootPath = rootPath + "/" + d
		goModPath = rootPath + "/go.mod"
		goModFile, err := ioutil.ReadFile(goModPath)
		if err != nil { // if the file doesn't exist, continue searching
			continue
		}
		// The project root directory is obtained based on the assumption that module name,
		// "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager", is contained in the 'go.mod' file.
		// Should the module name change in the code repo then it needs to be changed here too.
		if strings.Contains(string(goModFile), "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager") {
			break
		}
	}
	return rootPath
}
