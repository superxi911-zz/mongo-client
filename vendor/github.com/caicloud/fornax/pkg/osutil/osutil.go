// Copyright 2015 caicloud authors. All rights reserved.

package osutil

import (
	"os"
	"os/user"

	"github.com/caicloud/fornax/pkg/log"
)

// GetHomeDir gets current user's home directory.
func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// GetStringEnvWithDefault get evironment value of 'name', and return provided
// default value if not found.
func GetStringEnvWithDefault(name, def string) string {
	if val := os.Getenv(name); val == "" {
		log.Infof("Env variant %s not found, using default value: %s", name, def)
		return def
	} else {
		log.Infof("Env variant %s found, using env value: %s", name, val)
		return val
	}
}

// OpenFile opens file from 'name', and create one if not exist.
func OpenFile(fileName string, flag int, perm os.FileMode) (*os.File, error) {
	var file *os.File
	var err error

	file, err = os.OpenFile(fileName, flag, perm)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(fileName)
		if err != nil {
			return nil, err
		}
	}

	return file, err
}

// IfFileExists returns true if the file exists.
func IfFileExists(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}
