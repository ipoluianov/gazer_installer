package main

import (
	"fmt"
	"os"

	"github.com/ipoluianov/gazer_installer/windows/tools"
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

	err = tools.Run("go", "temp/gazer_node", []string{"build", "-v", "-o", "../../bin/gazer_node.exe", "./main.go"})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	err = tools.Run("flutter", "temp/gazer_client", []string{"build", "windows"})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("copy flutter application to bin")

	err = tools.CopyDirectory("temp/gazer_client/build/windows/runner/Release/", "bin/")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("copy redist to bin")
	err = tools.Unzip("temp/gazer_client/redist/redist.zip", "bin/")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	bsPasswd, err := os.ReadFile("d:\\src\\codesign\\passwd.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	passwd := string(bsPasswd)

	fmt.Println("signing gazer_node")
	err = tools.Run("d:\\src\\codesign\\signtool.exe", "", []string{"sign", "/v", "/t", "http://timestamp.sectigo.com", "/f", "d:\\src\\codesign\\iip.pfx", "/p", passwd, "d:\\src\\github\\gazer_installer\\windows\\bin\\gazer_node.exe"})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("signing gazer_client")
	err = tools.Run("d:\\src\\codesign\\signtool.exe", "", []string{"sign", "/v", "/t", "http://timestamp.sectigo.com", "/f", "d:\\src\\codesign\\iip.pfx", "/p", passwd, "d:\\src\\github\\gazer_installer\\windows\\bin\\gazer_client.exe"})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("creating distributive ...")
	err = tools.Run("c:\\Program Files (x86)\\NSIS\\makensisw.exe", ".", []string{"install.nsi"})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("signing gazer_installer")
	err = tools.Run("d:\\src\\codesign\\signtool.exe", "", []string{"sign", "/v", "/t", "http://timestamp.sectigo.com", "/f", "d:\\src\\codesign\\iip.pfx", "/p", passwd, "d:\\src\\github\\gazer_installer\\windows\\GazerNode_2.4.9.exe"})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Complete")
}
