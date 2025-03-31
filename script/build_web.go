//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	println("build web assets...")

	workspace, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	println("workspace path:" + workspace)

	webDirPath := filepath.Join(workspace, "gui", "web")
	err = installDependencies(webDirPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = runBuild(webDirPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	println("build web assets done.")
}

func installDependencies(workspace string) error {
	println("install dependencies...")

	cmd := exec.Command("npm", "install")
	cmd.Dir = workspace
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		println("install dependencies done.")
	}
	return err
}

func runBuild(workspace string) error {
	println("run build...")

	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = workspace
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		println("run build done.")
	}
	return err
}
