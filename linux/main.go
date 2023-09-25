package main

import (
	"fmt"
	"os"

	"github.com/ipoluianov/gazer_installer/linux/tools"
)

func main() {
	version := "2.4.19"

	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")

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

	err = tools.CopyFile("app_icon.ico", "bin/app_icon.ico")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("writing desktop file ...")
	desktopFileContent := `[Desktop Entry]
	Encoding=UTF-8
	Version=1.0
	Type=Application
	Terminal=false
	Exec=/usr/local/gazer/gazer_client
	Name=Gazer
	Icon=/usr/local/gazer/app_icon.ico`

	os.WriteFile("bin/gazer.desktop", []byte(desktopFileContent), 0666)

	fmt.Println("creating distributive ...")

	err = tools.GenerateInstallerFromDirectory("bin", "gazer_linux_"+version+"_setup.sh")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Complete")
}
