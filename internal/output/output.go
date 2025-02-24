package output

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func CreateDir(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {
		log.Printf("Failed to create directory %s: %v", path, err)
		return
	}
}

func CopyStaticFiles(inputFolder string) error {
	// Copy images
	srcImagesPath := filepath.Join(inputFolder, "images")
	dstImagesPath := "output/images"
	if _, err := os.Stat(srcImagesPath); os.IsNotExist(err) {
		log.Println("Source images folder does not exist:", srcImagesPath)
	} else {
		if err := copyFiles(srcImagesPath, dstImagesPath); err != nil {
			return err
		}
	}

	// Copy CSS
	if err := copyFiles("resource/main.css", "output/css/main.css"); err != nil {
		return err
	}
	// Copy Icons
	if err := copyFiles("resource/icons", "output/icons"); err != nil {
		return err
	}
	return nil
}

func copyFiles(srcImagesPath, dstImagesPath string) error {
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
