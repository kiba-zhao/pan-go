//go:build ignore
// +build ignore

package main

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

func buildLib(libPath string) error {
	println("build lib...")
	if err := os.MkdirAll(libPath, 0755); err != nil {
		if _, ok := err.(*os.PathError); !ok && !os.IsExist(err) {
			return err
		}
	}

	cmd := exec.Command(
		"go", "run", "golang.org/x/mobile/cmd/gomobile", "bind",
		"-target=android",
		"-androidapi=21",
		"-o", filepath.Join(libPath, "pan-go.gomobile.aar"),
		"./gomobile",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		println("build lib done.")
	}
	return err
}

func main() {
	println("bind android...")

	workspace, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	println("workspace path:" + workspace)

	libPath := filepath.Join(workspace, "mobile", "android", "servlet", "libs")
	err = buildLib(libPath)
	if err != nil {
		slog.Error(err.Error())
	}

	println("bind android done.")

}
