package main

import (
	"fmt"
	"os"

	"github.com/ipoluianov/gazer_installer/linux/tools"
)

func main() {
	var err error
	fmt.Println("Installer creation started ...")
	os.RemoveAll("temp")
	os.RemoveAll("bin")
	os.MkdirAll("temp", 0777)
	os.MkdirAll("bin", 0777)

	fmt.Println("git clone gazer_client")
	err = tools.Run("git", "temp", []string{"clone", "https://github.com/ipoluianov/gazer_client"})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

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

	err = tools.Run("flutter", "temp/gazer_client", []string{"build", "linux"})
	//err = tools.Run("flutter", "temp/gazer_client", []string{"build", "windows"})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("copy flutter application to bin")

	err = tools.CopyDirectory("temp/gazer_client/build/linux/x64/release/bundle/", "bin/")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("creating distributive ...")

	err = tools.GenerateInstallerFromDirectory("bin", "gazer_linux_setup.sh")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Complete")
}
