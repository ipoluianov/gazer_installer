package main

import (
	"fmt"
	"os"

	"github.com/ipoluianov/gazer_installer/linux_arm/tools"
)

func main() {
	version := "2.4.13"

	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "arm")

	var err error
	fmt.Println("Installer creation started ...")
	os.RemoveAll("temp")
	os.RemoveAll("bin")
	os.MkdirAll("temp", 0777)
	os.MkdirAll("bin", 0777)

	fmt.Println("git clone gazer_node")
	err = tools.Run("git", "temp", []string{"clone", "https://github.com/ipoluianov/gazer_node"})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("build gazer_node")

	err = tools.Run("go", "temp/gazer_node", []string{"build", "-v", "-o", "../../bin/gazer_node", "./main.go"})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("creating distributive ...")

	err = tools.GenerateInstallerFromDirectory("bin", "gazer_linux_arm_"+version+"_setup.sh")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Complete")
}
