package utils

import (
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func CreateTar(src string) string {

	destFileSplitted := strings.Split(src, "/")
	destFile := strings.Join(destFileSplitted[:len(destFileSplitted)-1],"/") + "/saltshaker_states.tar.gz"

	lastFolder := destFileSplitted[len(destFileSplitted)-1]

	cmd := exec.Command("tar", "cvzf", "saltshaker_states.tar.gz", lastFolder + "/")
	cmd.Dir = src + "/../"
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	return destFile
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}