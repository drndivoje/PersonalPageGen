package domain

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"ppg/internal/input"
	"ppg/internal/output"
	template "ppg/internal/template"
	"sort"
	"strings"

	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func writeAllPages(pages []Page, config ConfigFile) error {
	var localPageItems []template.LocalPageItem
	for _, page := range pages {
		localPageItem := page.getLocalPageItem()
		localPageItems = append(localPageItems, localPageItem)
	}
	for _, page := range pages {
		pageDetails := template.PageDetails{
			Language:          page.getLanguage(),
			PublishDateString: page.GetFormattedPublishDate(),
			Tags:              page.Metadata.getTags(),
			LocalePages:       localPageItems,
		}
		footer := template.GetFooter(config.GetFooter())
		HeadSection := getHeadSection(page, config)
		var pageDetailsParsed string
		if page.isBlogPage() {
			pageDetailsParsed = pageDetails.Parse()
		} else {
			pageDetailsParsed = ""
		}
		pageData := template.PageData{
			MenuSection: template.GetHeader(template.HeaderData{DomainUrl: config.getDomainUrl(), Menu: toMenuItem(config)}),
			Content:     page.Html,
			Footer:      footer,
			PageDetails: pageDetailsParsed,
			HeadSection: HeadSection,
			MainPage:    page.isMainPage(),
			Tags:        page.Metadata.getTags(),
		}
		err := page.writeToFile(pageData)
		if err != nil {
			log.Println("Error writing page ["+page.getTitle()+"] to the file:", err)
			return err
		}
	}
	return nil
}

func createWebSite(folderPath string) (WebSite, error) {
	err := input.ValidateInputFolder(folderPath)
	if err != nil {
		return WebSite{}, err
	}
	config, err := loadConfig(folderPath)
	if err != nil {
		return WebSite{}, err
	}
	var pages []Page
	var pagesMap = make(map[string][]Page)

	return WebSite{
		pages:       pages,
		pagesMap:    pagesMap,
		InputFolder: folderPath,
		ConfigFile:  config,
	}, nil
}

func CreateWebSite(folderPath string) (WebSite, error) {
	var err error
	webSite, err := createWebSite(folderPath)
	if err != nil {
		return WebSite{}, err
	}
	err = webSite.parsePages()
	if err != nil {
		return WebSite{}, err
	}
	return webSite, nil
}

func getHeadSection(page Page, config ConfigFile) string {
	return template.GetHeadSection(
		template.HeadSection{
			Title:       page.Metadata.getTitle(),
			Description: page.GetDescription(),
			DomainUrl:   config.getDomainUrl(),
			Author:      config.Author,
			PagePath:    page.getRelativePath(),
		},
	)
}

func loadConfig(folderPath string) (*ConfigFile, error) {

	content, err := input.ValidateConfigFile(folderPath)
	if err != nil {
		return nil, err
	}

	var config ConfigFile
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}

func toMenuItem(config ConfigFile) []template.MenuItem {
	var menuItems []template.MenuItem
	for _, item := range config.Menu {
		menuItems = append(menuItems, template.MenuItem{
			Title: item.Title,
			Url:   config.getDomainUrl() + "/" + item.Path,
		})
	}
	return menuItems
}

func parseRelativePath(relativePath string) (string, string, error) {
	extension := filepath.Ext(relativePath)
	if extension == "" || len(extension) != 3 {
		//default language is en
		return relativePath, "en", nil
	}

	languageCode := strings.ToLower(strings.TrimPrefix(extension, "."))
	if !isValidLanguageCode(languageCode) {
		return "", "", fmt.Errorf("invalid language code in path %q: %s", relativePath, languageCode)
	}

	basePath := strings.TrimSuffix(relativePath, extension)
	return basePath, languageCode, nil
}

func isValidLanguageCode(code string) bool {
	region, err := language.ParseRegion(strings.ToUpper(code))
	if err != nil {
		return false
	}
	return region.IsCountry()
}

func parseHeaderAttribute(line string) (MetadataAttribute, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return MetadataAttribute{}, fmt.Errorf("empty line")
	}
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return MetadataAttribute{}, fmt.Errorf("cannot parse line: %s", line)
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if value == "" {
		return MetadataAttribute{}, fmt.Errorf("value cannot be an empty string")
	}
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		value = strings.Trim(value, "[]")
		values := strings.Split(value, ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}
		return MetadataAttribute{name: key, values: values}, nil
	} else if strings.HasPrefix(value, "[") || strings.HasSuffix(value, "]") {
		return MetadataAttribute{}, fmt.Errorf("invalid list format")
	}

	return MetadataAttribute{name: key, values: []string{value}}, nil
}

func writePageList(pages template.PageListData, path string) error {
	pageListHtml := pages.Parse()
	log.Println("Writing page "+fmt.Sprintf("%d", len(pages.Pages))+" list to", path)
	output.CreateDir(path)
	err := os.WriteFile(path+"/index.html", []byte(pageListHtml), 0644)
	if err != nil {
		return err
	}
	return nil
}

func createPageListData(pageItems []template.PageItemData, w WebSite) template.PageListData {
	return template.PageListData{
		Header: template.GetHeader(template.HeaderData{DomainUrl: w.getDomainUrl(), Menu: toMenuItem(*w.ConfigFile)}),
		Pages:  pageItems,
		Footer: template.GetFooter(w.ConfigFile.GetFooter()),
		HeadSection: template.GetHeadSection(
			template.HeadSection{
				Title:       w.ConfigFile.Author + "'s Blog",
				Description: "Welcome to " + w.ConfigFile.Author + "'s Blog",
				DomainUrl:   w.getDomainUrl(),
				Author:      w.ConfigFile.Author,
			}),
	}
}

func insertSort(data []Page, el Page) []Page {
	index := sort.Search(len(data), func(i int) bool { return data[i].getPublishDate() < el.getPublishDate() })
	data = append(data, Page{})
	copy(data[index+1:], data[index:])
	data[index] = el
	return data
}
