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
	return author
}

func (p Page) GetDesctiption() string {
	description, ok := p.Header.Attributes["description"]
	if !ok {
		return "Unknown"
	}
	return description
}

func (p Page) writeToOutput(config ConfigFile) error {

	filePath := "output/" + p.Path + "/index.html"
	pageMetaData := ""
	if p.PublishDate != 0 {
		pageMetaData = template.GetPageDetails(p.PublishDate)
	}
	data := template.PageData{
		Header:       template.GetHeader(template.HeaderData{DomainUrl: config.getDomainUrl(), Menu: toMenuItem(config)}),
		Content:      p.Html,
		Footer:       template.GetFooter(config.Footer),
		PageDetails:  pageMetaData,
		HeadMetadata: getHeadMetaData(p, config.getDomainUrl()),
		MainPage:     p.isMainPage(),
	}
	pageHtml := template.ParseTemplate("template.html", data)
	if p.isMainPage() {
		filePath = "output/index.html"
	} else {
		output.CreateDir("output/" + p.Path)
	}
	err := os.WriteFile(filePath, []byte(pageHtml), 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
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
	for _, page := range w.blogPages {
		pageItem := template.PageItemData{
			Title: page.Title,
			Date:  time.Unix(page.PublishDate, 0).Format("January 2, 2006"),
			Url:   page.Path,
		}
		pageItems = append(pageItems, pageItem)
	}
	var pageListData = template.PageListData{
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
	pageListHtml := template.ParseTemplate("list.html", pageListData)
	outputPath := filepath.Join("output", "blog", "index.html")
	output.CreateDir(filepath.Dir(outputPath))
	err := os.WriteFile(outputPath, []byte(pageListHtml), 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
		return
	}
	log.Println("Website generated successfully")

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

func (w *WebSite) getDomainUrl() string {
	return "http://" + w.ConfigFile.Domain
}

func (w WebSite) String() string {
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

	header := Header{Attributes: make(map[string]string)}
	for _, line := range lines[1:] {

		if line == headerDelimiter {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Warning: cannot parse line: %s\n", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		header.Attributes[key] = value
	}
	return header, nil
}

type Header struct {
	Attributes map[string]string
}

func (h Header) getTitle() string {
	title, ok := h.Attributes["title"]
	if !ok {
		return ""
	}
	title = strings.ReplaceAll(title, "\"", "")
	return title
}

func getDate(header Header) int64 {
	date, ok := header.Attributes["date"]
	if !ok {
		return 0
	}
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, date)
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
