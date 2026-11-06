package domain

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"ppg/internal/input"
	"ppg/internal/output"
	template "ppg/internal/template"

	blackfriday "github.com/russross/blackfriday/v2"
)

type WebSite struct {
	pages       []Page
	pagesMap    map[string][]Page
	InputFolder string
	ConfigFile  *ConfigFile
}

func (w *WebSite) parsePages() error {
	var err error
	err = filepath.WalkDir(w.InputFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			relativePath, err := input.ParsePagePath(w.InputFolder, path)
			if err != nil {
				return err
			}

			metadata, err := parseMetadata(content, relativePath)
			if err != nil {
				return err
			}

			pageContent, err := input.ParsePageContent(content, metadata.getTitle())
			if err != nil {
				return err
			}

			htmlOutput := blackfriday.Run([]byte(pageContent))
			if htmlOutput != nil {
				page := Page{
					Metadata: metadata,
					Data:     string(content),
					Html:     string(htmlOutput),
				}

				w.pagesMap[page.getKey()] = append(w.pagesMap[page.getKey()], page)
				if page.isBlogPage() && page.hasDefaultLanguage() {
					w.pages = insertSort(w.pages, page)
				}

			}

		}
		return nil
	})
	return err
}

func (w *WebSite) WriteToOutputFolder() {

	for path, pages := range w.pagesMap {
		writeAllPages(pages, *w.ConfigFile)
		fmt.Printf("%s: %d\n", path, len(pages))
	}

	output.CopyStaticFiles(w.InputFolder)

	var pageItems []template.PageItemData
	tagtoPagesMap := make(map[string][]template.PageItemData)

	for _, page := range w.pages {
		pageItem := template.PageItemData{
			Title: page.getTitle(),
			Date:  page.GetFormattedPublishDate(),
			Url:   page.getRelativePath(),
			Tags:  page.Metadata.getTags(),
		}
		pageItems = append(pageItems, pageItem)
		for _, tag := range page.Metadata.getTags() {
			tagtoPagesMap[tag] = append(tagtoPagesMap[tag], template.PageItemData{
				Title: page.getTitle(),
				Date:  page.GetFormattedPublishDate(),
				Url:   page.getRelativePath(),
				Tags:  page.Metadata.getTags(),
			})
		}
	}
	var pageListData = createPageListData(pageItems, *w)
	err := writePageList(pageListData, filepath.Join("output", "blog"))
	if err != nil {
		log.Println("Error writing blog pages", err)
		return
	}
	for tag, pages := range tagtoPagesMap {
		taggedPageLisData := createPageListData(pages, *w)
		err := writePageList(taggedPageLisData, filepath.Join("output", "tags", tag))
		if err != nil {
			log.Println("Error writing tag pages", err)
			return
		}
	}
	log.Println("Website generated successfully!")

}
func (w WebSite) getDomainUrl() string {
	return w.ConfigFile.getDomainUrl()
}

func (w *WebSite) String() string {
	return fmt.Sprintf("Pages: %v\nConfigFile: %v", len(w.pages), w.ConfigFile)
}
