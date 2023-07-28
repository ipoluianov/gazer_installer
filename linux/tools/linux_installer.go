package tools

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func GenerateInstallerFromDirectory(dirPath, outputPath string) error {
	tarGzPath := outputPath + ".tar.gz"
	if err := createTarGz(dirPath, tarGzPath); err != nil {
		return err
	}

	tarGzContent, err := os.ReadFile(tarGzPath)
	if err != nil {
		return err
	}
	encodedTarGzContent := base64.StdEncoding.EncodeToString(tarGzContent)

	scriptFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer scriptFile.Close()

	writer := bufio.NewWriter(scriptFile)

	fmt.Fprintln(writer, "#!/bin/bash")
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "mkdir /usr/local/gazer/")
	fmt.Fprintln(writer, "sudo /usr/local/gazer/gazer_node -stop")
	fmt.Fprintln(writer, "sudo /usr/local/gazer/gazer_node -uninstall")

	fmt.Fprintln(writer, "echo '"+encodedTarGzContent+"' | base64 -d | tar -xz -C /usr/local/gazer/")
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "chmod -R 755 /usr/local/gazer/")
	fmt.Fprintln(writer, "find /usr/local/gazer/ -type f -exec chmod 644 {} +")
	fmt.Fprintln(writer, "chmod 777 /usr/local/gazer/gazer_node")
	fmt.Fprintln(writer, "chmod 777 /usr/local/gazer/gazer_client")
	fmt.Fprintln(writer, "sudo /usr/local/gazer/gazer_node -install")
	fmt.Fprintln(writer, "sudo /usr/local/gazer/gazer_node -start")
	fmt.Fprintln(writer, "yes | cp /usr/local/gazer/gazer.desktop /usr/share/applications/")

	fmt.Fprintln(writer)

	writer.Flush()

	fmt.Printf("Installer was created: %s\n", outputPath)
	return nil
}

func createTarGz(dirPath, outputPath string) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
