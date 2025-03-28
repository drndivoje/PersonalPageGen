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
	"gopkg.in/yaml.v3"
)

type Page struct {
	Header      Header
	Title       string
	Path        string
	PublishDate int64
	Data        string
	Html        string
}

func (p Page) isMainPage() bool {
	return p.Path == "index"
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

func (p Page) GetAuthor() string {
	author, ok := p.Header.Attributes["author"]
	if !ok {
		return "Unknown"
	}
	return author.getValue()
}

func (p Page) GetDesctiption() string {
	description, ok := p.Header.Attributes["description"]
	if !ok {
		return "Unknown"
	}
	return description.getValue()
}

func (p Page) writeToOutput(config ConfigFile) error {

	filePath := "output/" + p.Path + "/index.html"
	pageMetaData := ""
	if p.PublishDate != 0 {
		pageMetaData = template.GetPageDetails(p.PublishDate, p.Header.getTags())
	}
	data := template.PageData{
		Header:       template.GetHeader(template.HeaderData{DomainUrl: config.getDomainUrl(), Menu: toMenuItem(config)}),
		Content:      p.Html,
		Footer:       template.GetFooter(config.Footer),
		PageDetails:  pageMetaData,
		HeadMetadata: getHeadMetaData(p, config.getDomainUrl()),
		MainPage:     p.isMainPage(),
		Tags:         p.Header.getTags(),
	}
	pageHtml := template.ParseTemplate("template.html", data)
	if p.isMainPage() {
		filePath = "output/index.html"
	} else {
		output.CreateDir("output/" + p.Path)
	}
	err := os.WriteFile(filePath, []byte(pageHtml), 0644)
	if err != nil {
		log.Println("Error writing page ["+p.Title+"] to the file:", err)
		return err
	}
	return nil
}

func (p Page) String() string {
	return fmt.Sprintf("Title: %s\nPath: %s\nDate: %d\n", p.Title, p.Path, p.PublishDate)
}

type WebSite struct {
	pages []Page

	blogPages []Page

	InputFolder string

	ConfigFile *ConfigFile
}

func (w *WebSite) WriteToOutputFolder() {
	for _, page := range w.pages {
		err := page.writeToOutput(*w.ConfigFile)
		if err != nil {
			return
		}
	}

	for _, page := range w.blogPages {
		err := page.writeToOutput(*w.ConfigFile)
		if err != nil {
			return
		}
	}

	output.CopyStaticFiles(w.InputFolder)

	var pageItems []template.PageItemData
	tagtoPagesMap := make(map[string][]template.PageItemData)

	for _, page := range w.blogPages {
		pageItem := template.PageItemData{
			Title: page.Title,
			Date:  time.Unix(page.PublishDate, 0).Format("January 2, 2006"),
			Url:   page.Path,
			Tags:  page.Header.getTags(),
		}
		pageItems = append(pageItems, pageItem)
		for _, tag := range page.Header.getTags() {
			tagtoPagesMap[tag] = append(tagtoPagesMap[tag], template.PageItemData{
				Title: page.Title,
				Date:  time.Unix(page.PublishDate, 0).Format("January 2, 2006"),
				Url:   page.Path,
				Tags:  page.Header.getTags(),
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

func writePageList(pages template.PageListData, path string) error {
	pageListHtml := template.ParseTemplate("list.html", pages)
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
		Footer: template.GetFooter(w.ConfigFile.Footer),
		HeadMetadata: template.GetHeadMetada(
			template.HeadMetadata{
				Title:       w.ConfigFile.Author + "'s Blog",
				Description: "Welcome to " + w.ConfigFile.Author + "'s Blog",
				DomainUrl:   w.getDomainUrl(),
				Author:      w.ConfigFile.Author,
			}),
	}
}

func getHeadMetaData(page Page, domainUrl string) string {
	return template.GetHeadMetada(
		template.HeadMetadata{
			Title:       page.Header.getTitle(),
			Description: page.GetDesctiption(),
			DomainUrl:   domainUrl,
			Author:      page.GetAuthor(),
		},
	)
}

func (w WebSite) getDomainUrl() string {
	return "http://" + w.ConfigFile.Domain
}

func (w *WebSite) String() string {
	return fmt.Sprintf("Pages: %v\nConfigFile: %v", len(w.pages), w.ConfigFile)
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

func parseHeader(content []byte) (Header, error) {

	const headerDelimiter = "+++"

	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 || lines[0] != headerDelimiter {
		return Header{}, fmt.Errorf("header does not start with +++")
	}

	header := Header{Attributes: make(map[string]HeaderAttribute)}
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

func parseHeaderAttribute(line string) (HeaderAttribute, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return HeaderAttribute{}, fmt.Errorf("empty line")
	}
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return HeaderAttribute{}, fmt.Errorf("cannot parse line: %s", line)
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if value == "" {
		return HeaderAttribute{}, fmt.Errorf("value cannot be an empty string")
	}
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		value = strings.Trim(value, "[]")
		values := strings.Split(value, ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}
		return HeaderAttribute{name: key, values: values}, nil
	} else if strings.HasPrefix(value, "[") || strings.HasSuffix(value, "]") {
		return HeaderAttribute{}, fmt.Errorf("invalid list format")
	}

	return HeaderAttribute{name: key, values: []string{value}}, nil
}

type Header struct {
	Attributes map[string]HeaderAttribute
}

type HeaderAttribute struct {
	name   string
	values []string
}

func (h HeaderAttribute) getValue() string {
	return h.values[0]
}

func (h HeaderAttribute) getValues() []string {
	return h.values
}

func (h Header) getTitle() string {
	title, ok := h.Attributes["title"]
	if !ok {
		return ""
	}
	return strings.ReplaceAll(title.getValue(), "\"", "")
}

func (h Header) getTags() []string {
	tags, ok := h.Attributes["tags"]
	if !ok {
		return []string{}
	}
	return tags.getValues()
}

func getDate(header Header) int64 {
	date, ok := header.Attributes["date"]
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

func NewWebSiteFromFolder(folderPath string) (*WebSite, error) {

	err := input.ValidateInputFolder(folderPath)
	if err != nil {
		return nil, err
	}
	config, err := loadConfig(folderPath)
	if err != nil {
		return nil, err
	}
	var pages []Page
	var blogPages []Page

	err = filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			relativePath, err := input.ParsePagePath(folderPath, path)
			if err != nil {
				return err
			}

			header, err := parseHeader(content)
			if err != nil {
				return err
			}

			pageContent, err := input.ParsePageContent(content, header.getTitle())
			if err != nil {
				return err
			}

			pageDate := getDate(header)

			htmlOutput := blackfriday.Run([]byte(pageContent))
			if htmlOutput != nil {
				page := Page{
					Header:      header,
					Title:       header.getTitle(),
					PublishDate: pageDate,
					Data:        string(content),
					Path:        relativePath,
					Html:        string(htmlOutput),
				}
				if strings.HasPrefix(relativePath, "blog") {
					blogPages = append(blogPages, page)
				} else {
					pages = append(pages, page)
				}

			}

		}
		return nil
	})
	sort.Slice(blogPages, func(i, j int) bool {
		return blogPages[i].PublishDate > blogPages[j].PublishDate
	})
	if err != nil {
		return nil, err
	}

	return &WebSite{
		pages:       pages,
		blogPages:   blogPages,
		InputFolder: folderPath,
		ConfigFile:  config,
	}, nil
}
