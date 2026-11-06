package domain
import (
	"log"
	"os"
	"strings"
	"time"

	"ppg/internal/output"
	"ppg/internal/template"
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

func (p Page) writeToFile(pageData template.PageData) error {
	var filePath string
	if p.isMainPage() {
		filePath = "output/index.html"
	} else {
		output.CreateDir("output/" + p.getRelativePath())
		filePath = "output/" + p.getRelativePath() + "/index.html"
	}
	err := os.WriteFile(filePath, []byte(pageData.Parse()), 0644)
	if err != nil {
		log.Println("Error writing page ["+p.getTitle()+"] to the file:", err)
		return err
	}
	log.Println("Successfully wrote page ["+p.getTitle()+"] to the file:", filePath)

	return nil
}

func (p Page) getLanguage() string {
	return p.Metadata.Language
}
