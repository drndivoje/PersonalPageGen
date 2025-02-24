package main

import (
	"fmt"
	"os"

	"ppg/internal/domain"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <source_folder>")
		return
	}

	sourceFolder := os.Args[1]
	website, err := domain.NewWebSiteFromFolder(sourceFolder)
	if err != nil {
		fmt.Println("Error creating website:", err)
		return
	}

	err = os.Mkdir("output", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating output directory:", err)
		return
	}

	website.WriteToOutputFolder()

}
