package template

import (
	"strings"
	template "text/template"
	"time"
)

var path = "resource"

type PageData struct {
	Header       string
	Content      string
	Footer       string
	PageDetails  string
	HeadMetadata string
	MainPage     bool
}

type PageListData struct {
	Header       string
	Pages        []PageItemData
	Footer       string
	HeadMetadata string
}

type PageItemData struct {
	Title string
	Url   string
	Date  string
}
type MenuItem struct {
	Title string
	Url   string
}

type HeaderData struct {
	DomainUrl string
	Menu      []MenuItem
}

type PageDetails struct {
	PublishDateString string
}

type HeadMetadata struct {
	Title       string
	DomainUrl   string
	Author      string
	Description string
}
type FooterData struct {
	Data string
}

func ParseTemplate(pathToTemplate string, data any) string {
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

func GetHeader(headerData HeaderData) string {
	return ParseTemplate("header.html", headerData)
}
func GetFooter(data string) string {
	return ParseTemplate("footer.html", FooterData{Data: data})
}

func GetPageDetails(date int64) string {
	return ParseTemplate("page-metadata.html", PageDetails{PublishDateString: time.Unix(date, 0).Format("January 2, 2006")})
}

func GetHeadMetada(headMetadata HeadMetadata) string {
	return ParseTemplate("head.html", headMetadata)
}
