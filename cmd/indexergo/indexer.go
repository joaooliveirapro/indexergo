package indexergo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Document struct {
	URL         string
	Ranking     float64
	QueryWeight map[string]float64
}

type Indexer struct {
	URLsFilePath        string
	URLsList            []string
	PageInfo            *PageInfo
	LookByQuerySelector []string
}

/* Indexer logic */
func (i *Indexer) IndexDocuments() error {
	var urls []string
	if len(i.URLsFilePath) != 0 {
		// Read URLs file
		var err error
		urls, err = ReadURLsFile(i.URLsFilePath)
		if err != nil {
			return fmt.Errorf("[error] error=%v", err)
		}
	} else {
		urls = i.URLsList
	}
	// Iterate each URL
	for _, url := range urls {
		pageinfo, err := NewPageInfo(url)
		if err != nil {
			return fmt.Errorf("[error] url=%s error=%v", url, err)
		}
		// Attach pageinfo to indexer
		i.PageInfo = pageinfo

		// Prevent word concatenation due to how goquery ingores <br> tags when
		// calling .Text() method.
		html := strings.ReplaceAll(pageinfo.HTTPResponse.HTML, "<br>", "\n")

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return fmt.Errorf("[error] url=%s goquery=%v", url, err)
		}
		// Join all the text content in the selection
		var textToIndex strings.Builder
		for _, query := range i.LookByQuerySelector {
			doc.Find(query).Each(func(i int, s *goquery.Selection) {
				textToIndex.WriteString(s.Text()) // Errors are ignored
			})
		}
		// Run HTML tag frequency on  i.PageInfo.HTTPResponse.HTML document
		tags := HTMLTagsFrequency(i.PageInfo.HTTPResponse.HTML)
		i.PageInfo.HTMLTags = tags

		// Run tokens frequency (on selection if provided)
		var tokens map[string]int
		if textToIndex.Len() > 0 {
			tokens = ContentTokensFrequency(textToIndex.String())
		} else {
			// Default to using HTML if not provided
			tokens = ContentTokensFrequency(i.PageInfo.HTTPResponse.HTML)
		}
		i.PageInfo.ContentTokens = tokens
		// Append JSON
		const IndexJSONFilePath = "index.json"
		err = AppendPageInfoToJson(IndexJSONFilePath, i.PageInfo)
		if err != nil {
			return fmt.Errorf("[error] url=%s appending to JSON: %v", url, err)
		}
	}
	return nil
}

func (i *Indexer) Search(query string) ([]Document, error) {
	docs, err := ReadPageInfoCache("index.json")
	if err != nil {
		return nil, fmt.Errorf("[error] failed to marshal JSON: %v", err)
	}
	searchTokens := strings.Split(query, " ")
	tfidfResults := CalculateTFIDF(&docs, searchTokens)
	return tfidfResults, nil
}

/*
Reads URLs file line by line and returns a list of valid URLs.
URLs are parsed and validated using net/url. Invalid URLs are skipped.
*/
func ReadURLsFile(filePath string) ([]string, error) {
	// Open file
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("[error] failed to open file: %v", err)
	}
	defer file.Close()
	// Read file
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("[error] reading file: %v", err)
	}
	// Append all lines to list
	/* assumes all lines are valid URLs as parsed by net/url. Invalid URLs are skipped */
	urls := []string{}
	for scanner.Scan() {
		url_ := scanner.Text()
		parsedURL, err := url.Parse(url_)
		URLisValid := err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
		if URLisValid {
			urls = append(urls, url_)
		}
	}
	return urls, nil
}

func ReadPageInfoCache(filepath string) ([]PageInfo, error) {
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("[error] reading file: %v", err)

	}
	allPagesInfo := []PageInfo{}
	if len(fileData) > 0 { // Handle empty file case
		if err := json.Unmarshal(fileData, &allPagesInfo); err != nil {
			return nil, fmt.Errorf("[error] parsing JSON: %v", err)
		}
	}
	return allPagesInfo, nil
}

/*
Reads index JSON file and update its contents with the new PageInfo pointer data
*/
func AppendPageInfoToJson(filepath string, pageInfo *PageInfo) error {
	allPagesInfo, err := ReadPageInfoCache(filepath)
	if err != nil {
		return fmt.Errorf("[error] failed to read JSON: %v", err)
	}
	allPagesInfo = append(allPagesInfo, *pageInfo)
	updatedData, err := json.MarshalIndent(allPagesInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("[error] failed to marshal JSON: %v", err)
	}
	err = os.WriteFile(filepath, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("[error] failed to write to file: %v", err)
	}
	return err
}

/*
Updates i.PageInfo.HTMLTags map with the count of each HTML tag found.
Uses regex to find tags: `<(?<tagname>\w*)[^>]*>`
*/
func HTMLTagsFrequency(html string) map[string]int {
	tags := map[string]int{}
	// Regex to capture <tagname>(content)<
	r := regexp.MustCompile(`<(?<tagname>\w*)[^>]*>`)
	for _, match := range r.FindAllStringSubmatch(html, -1) {
		tagName := match[1]
		if len(tagName) == 0 { // Skip empty strings
			continue
		}
		// Increment frequency
		tags[tagName]++
	}
	return tags
}

/*
Updates i.PageInfo.ContentTokens map with the count of each token (valid words) found in content.
Tokens that contain numbers are skipped. Punctuation is removed.
*/
func ContentTokensFrequency(content string) map[string]int {
	tokens := map[string]int{}
	// Lower case content
	content = strings.ToLower(content)
	// Remove all NON words (punctuation) - keeps "\n" and "\n".
	content = string(regexp.MustCompile(`[^\w]`).ReplaceAll([]byte(content), []byte(" ")))
	for _, w := range strings.Fields(content) { // strings.Fields trims and splits by spaces
		// Skip empty strings
		if len(w) == 0 {
			continue
		}
		// Skip tokens with numbers (usually these are not valid words)
		if regexp.MustCompile(`\d`).Match([]byte(w)) {
			continue
		}
		// Increment frequency
		tokens[w]++
	}
	return tokens
}
