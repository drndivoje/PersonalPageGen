package template

import (
	"strings"
	template "text/template"
)

var path = "resource"

type Parsable interface {
	Parse() string
}

type PageData struct {
	MenuSection string
	Content     string
	Footer      string
	PageDetails string
	HeadSection string
	MainPage    bool
	Tags        []string
}

type PageListData struct {
	Header      string
	Pages       []PageItemData
	Footer      string
	HeadSection string
}

type PageItemData struct {
	Title string
	Url   string
	Date  string
	Tags  []string
}
type MenuItem struct {
	Title string
	Url   string
}

type HeaderData struct {
	DomainUrl string
	Menu      []MenuItem
}

type LocalPageItem struct {
	Path     string
	Language string
}

type PageDetails struct {
	Language          string
	PublishDateString string
	Tags              []string
	LocalePages       []LocalPageItem
}

type HeadSection struct {
	Title       string
	DomainUrl   string
	PagePath    string
	Author      string
	Description string
}
type FooterData struct {
	Data string
}


func ParseTemplate(pathToTemplate string, data Parsable) string {
	tmpl, err := template.ParseFiles(path + "/" + pathToTemplate)
	if err != nil {
		panic(err)
	}
	var output strings.Builder
	err = tmpl.Execute(&output, data)
	if err != nil {
		panic(err)
	}
	return output.String()
}

func (data PageData) Parse() string {
	return ParseTemplate("template.html", data)
}

func (data PageListData) Parse() string {
	return ParseTemplate("list.html", data)
}

func (data HeaderData) Parse() string {
	return ParseTemplate("header.html", data)
}

func (data FooterData) Parse() string {
	return ParseTemplate("footer.html", data)
}

func (data PageDetails) Parse() string {
	filteredLocalePages := make([]LocalPageItem, 0, len(data.LocalePages))
	for _, localePage := range data.LocalePages {
		if data.Language != localePage.Language {
			filteredLocalePages = append(filteredLocalePages, localePage)
		}
	}
	data.LocalePages = filteredLocalePages
	return ParseTemplate("page-metadata.html", data)
}

func (data HeadSection) Parse() string {
	return ParseTemplate("head.html", data)
}

func GetHeader(headerData HeaderData) string {
	return headerData.Parse()
}

func GetFooter(data string) string {
	return FooterData{Data: data}.Parse()
}

func GetPageDetails(date string, tags []string, localePages []LocalPageItem) string {
	return NewPageDetails("", date, tags, localePages).Parse()
}

func NewPageDetails(language string, publishDateString string, tags []string, localePages []LocalPageItem) PageDetails {
	if localePages == nil {
		localePages = []LocalPageItem{}
	}
	return PageDetails{
		Language:          language,
		PublishDateString: publishDateString,
		Tags:              tags,
		LocalePages:       localePages,
	}
}

func GetHeadSection(HeadSection HeadSection) string {
	return HeadSection.Parse()
}
