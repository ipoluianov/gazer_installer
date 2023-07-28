package tools

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(zipPath, destDir string) error {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	for _, file := range zipFile.File {
		relPath := file.Name
		dstPath := filepath.Join(destDir, relPath)

		if file.FileInfo().IsDir() {
			os.MkdirAll(dstPath, file.Mode())
		} else {
			err = extractFile(file, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func extractFile(file *zip.File, destPath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
