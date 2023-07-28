package tools

import (
	"io"
	"os"
	"path/filepath"
)

func CopyDirectory(srcDirPath, dstDirPath string) error {
	if err := os.MkdirAll(dstDirPath, 0755); err != nil {
		return err
	}

	fileList, err := getFileList(srcDirPath)
	if err != nil {
		return err
	}

	for _, srcPath := range fileList {
		relPath, err := filepath.Rel(srcDirPath, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDirPath, relPath)

		if isDir(srcPath) {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func getFileList(dirPath string) ([]string, error) {
	var fileList []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileList = append(fileList, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
