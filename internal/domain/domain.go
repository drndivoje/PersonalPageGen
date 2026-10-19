package domain

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"ppg/internal/input"
	"ppg/internal/output"
	template "ppg/internal/template"

	blackfriday "github.com/russross/blackfriday/v2"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type Page struct {
	Metadata Metadata
	Data     string
	Html     string
}

func (p Page) getTitle() string {
	return p.Metadata.getTitle()
}

func (p Page) getPublishDate() int64 {
	return p.Metadata.getPublishDate()
}

func (p Page) GetFormattedPublishDate() string {
	return time.Unix(p.getPublishDate(), 0).Format("January 2, 2006")
}

func (p Page) getRelativePath() string {
	return p.Metadata.getResolvedPath()
}

func (p Page) getKey() string {
	return p.Metadata.RelativePath
}

func (p Page) isBlogPage() bool {
	return strings.HasPrefix(p.getRelativePath(), "blog/")
}

func (p Page) hasDefaultLanguage() bool {
	return p.Metadata.Language == "en"
}

func (p Page) isMainPage() bool {
	return p.getRelativePath() == "index"
}

func (p Page) GetDescription() string {
	description, ok := p.Metadata.Attributes["description"]
	if !ok {
		return p.getTitle()
	}
	return description.getValue()
}

func (p Page) getLocalPageItem() template.LocalPageItem {
	return template.LocalPageItem{
		Path:     p.getRelativePath(),
		Language: p.Metadata.Language,
	}
}

func (p Page) getLanguage() string {
	return p.Metadata.Language
}

func writeAllPages(pages []Page, config ConfigFile) error {
	var localPageItems []template.LocalPageItem
	for _, page := range pages {
		localPageItem := page.getLocalPageItem()
		localPageItems = append(localPageItems, localPageItem)
	}
	for _, page := range pages {
		filePath := "output/" + page.getRelativePath() + "/index.html"
		pageDetails := template.PageDetails{
			Language:          page.getLanguage(),
			PublishDateString: page.GetFormattedPublishDate(),
			Tags:              page.Metadata.getTags(),
			LocalePages:       localPageItems,
		}
		footer := template.GetFooter(config.GetFooter())
		HeadSection := getHeadSection(page, config)
		pageData := template.PageData{
			MenuSection: template.GetHeader(template.HeaderData{DomainUrl: config.getDomainUrl(), Menu: toMenuItem(config)}),
			Content:     page.Html,
			Footer:      footer,
			PageDetails: pageDetails.Parse(),
			HeadSection: HeadSection,
			MainPage:    page.isMainPage(),
			Tags:        page.Metadata.getTags(),
		}
		pageHtml := pageData.Parse()

		if page.isMainPage() {
			filePath = "output/index.html"
		} else {
			output.CreateDir("output/" + page.getRelativePath())
		}
		err := os.WriteFile(filePath, []byte(pageHtml), 0644)
		if err != nil {
			log.Println("Error writing page ["+page.getTitle()+"] to the file:", err)
			return err
		}
		log.Println("Successfully wrote page ["+page.getTitle()+"] to the file:", filePath)
	}
	return nil
}

type WebSite struct {
	pages       []Page
	pagesMap    map[string][]Page
	InputFolder string
	ConfigFile  *ConfigFile
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

func (w WebSite) getDomainUrl() string {
	return w.ConfigFile.getDomainUrl()
}

func (w *WebSite) String() string {
	return fmt.Sprintf("Pages: %v\nConfigFile: %v", len(w.pages), w.ConfigFile)
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

type ConfigFile struct {
	Main string `yaml:"site_name"`
	Menu []struct {
		Title string `yaml:"title"`
		Path  string `yaml:"path"`
	} `yaml:"menu"`
	Domain string `yaml:"domain"`
	Author string `yaml:"author"`
	Footer string `yaml:"footer"`
}

func (c ConfigFile) getDomainUrl() string {
	return "https://" + c.Domain
}

func (c ConfigFile) GetFooter() string {
	defaultFooter := "Built with <a href=\"https://github.com/drndivoje/PersonalPageGen\" target=\"_blank\" >Personal Page Gen</a>"
	if c.Footer == "" {
		return defaultFooter
	}
	return c.Footer + " - " + defaultFooter
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

func parseMetadata(content []byte, relativePath string) (Metadata, error) {

	const headerDelimiter = "+++"

	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 || lines[0] != headerDelimiter {
		return Metadata{}, fmt.Errorf("header does not start with +++")
	}

	relativePath, languageCode, err := parseRelativePath(relativePath)
	if err != nil {
		return Metadata{}, err
	}

	header := Metadata{
		Attributes:   make(map[string]MetadataAttribute),
		RelativePath: relativePath,
		Language:     languageCode,
	}
	for _, line := range lines[1:] {

		if line == headerDelimiter {
			break
		}
		headerAttribute, err := parseHeaderAttribute(line)
		if err != nil {
			log.Printf("Warning: cannot parse line: %s\n", err.Error())
			continue
		}
		header.Attributes[headerAttribute.name] = headerAttribute
	}
	return header, nil
}

type Metadata struct {
	Attributes   map[string]MetadataAttribute
	RelativePath string
	Language     string
}

type MetadataAttribute struct {
	name   string
	values []string
}

func (h MetadataAttribute) getValue() string {
	return h.values[0]
}

func (h MetadataAttribute) getValues() []string {
	return h.values
}

func (h Metadata) getTitle() string {
	title, ok := h.Attributes["title"]
	if !ok {
		return ""
	}
	return strings.ReplaceAll(title.getValue(), "\"", "")
}

func (h Metadata) getTags() []string {
	tags, ok := h.Attributes["tags"]
	if !ok {
		return []string{}
	}
	return tags.getValues()
}

func (h Metadata) getPublishDate() int64 {
	date, ok := h.Attributes["date"]
	if !ok {
		return 0
	}
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, date.getValue())
	if err != nil {
		log.Println("Error parsing date:", err)
		return 0
	}
	return int64(parsedDate.Unix())
}

func (h Metadata) getLanguage() string {
	if h.Language != "" {
		return h.Language
	}
	language, ok := h.Attributes["language"]
	if !ok {
		return "en"
	}
	return strings.ReplaceAll(language.getValue(), "\"", "")
}

// getResolvedPath returns the relative path of the metadata, including the language code if available.
func (h Metadata) getResolvedPath() string {
	if h.getLanguage() == "en" {
		return h.RelativePath
	}
	return h.RelativePath + "/" + h.getLanguage()
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
	log.Println("Writing page " + fmt.Sprintf("%d", len(pages.Pages)) + " list to", path)
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
	index := sort.Search(len(data), func(i int) bool { return data[i].getPublishDate() > el.getPublishDate() })
	data = append(data, Page{})
	copy(data[index+1:], data[index:])
	data[index] = el
	return data
}
