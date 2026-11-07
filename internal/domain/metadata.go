package domain

import (
	"fmt"
	"log"
	"strings"
	"time"
)

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
