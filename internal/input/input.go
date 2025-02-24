package input

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ValidateInputFolder(folderPath string) error {
	info, err := os.Stat(folderPath)
	if err != nil {
		return fmt.Errorf("error accessing folder: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("provided path is not a folder")
	}

	blogFolderPath := filepath.Join(folderPath, "blog")
	if _, err := os.Stat(blogFolderPath); os.IsNotExist(err) {
		return fmt.Errorf("blog folder does not exist in the provided path")
	}
	return nil
}

func ValidateConfigFile(folderPath string) ([]byte, error) {
	configPath := filepath.Join(folderPath, "config.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist")
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	return content, nil
}

func ParsePagePath(rootFolder, pageFolder string) (string, error) {
	relativePath, err := filepath.Rel(rootFolder, pageFolder)
	if err != nil {
		return "", err
	}
	relativePath = strings.TrimSuffix(relativePath, ".md")

	return relativePath, nil
}
func ParsePageContent(content []byte, title string) (string, error) {
	lines := strings.Split(string(content), "\n")
	headerEndIndex := -1
	headerCount := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "+++") {
			headerCount++
			if headerCount == 2 {
				headerEndIndex = i
				break
			}
		}
	}
	if headerEndIndex == -1 {
		return "", fmt.Errorf("header end not found")
	}
	return "# " + title + "\n" + strings.Join(lines[headerEndIndex+1:], "\n"), nil
}
