package output

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func CreateDir(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating output directory:", err)
		return
	}
}

func CopyStaticFiles(inputFolder string) error {
	// Copy images
	srcImagesPath := filepath.Join(inputFolder, "images")
	dstImagesPath := "output/images"

	if err := CopyFiles(srcImagesPath, dstImagesPath); err != nil {
		return err
	}
	// Copy CSS
	if err := CopyFiles("resource/main.css", "output/css/main.css"); err != nil {
		return err
	}
	// Copy Icons
	if err := CopyFiles("resource/icons", "output/icons"); err != nil {
		return err
	}
	return nil
}

func CopyFiles(srcImagesPath, dstImagesPath string) error {
	return filepath.Walk(srcImagesPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(srcImagesPath, path)
			if err != nil {
				return err
			}
			dstPath := filepath.Join(dstImagesPath, relPath)
			CreateDir(filepath.Dir(dstPath))
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			err = os.WriteFile(dstPath, data, 0644)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
