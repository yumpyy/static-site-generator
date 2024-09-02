package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"regexp"
)

const (
    draftDir = "../draft/"
    publicDir = "../public/"
    templateFilePath = "../public/template.html"
)

func parseDraft(lines []string) (title, date, content string) {
	title = strings.TrimSpace(strings.Split(lines[0], ":")[1])
	date = strings.TrimSpace(strings.Split(lines[1], ":")[1])
    content = strings.Join(lines[3:], "\n")
	if date == "" {
		date = strings.ToLower(time.Now().Format("02-01-2006"))
	}

    return
}

func convertMarkdownToHtml(content string) string {
    mdHTMLTag := map[string]string {
        `\*(.*?)\*`   : "b",
        `_(.*?)_`     : "i",
        `~(.*?)~`     : "s",
        `^# (.*?)$`   : "h1",
        `^## (.*?)$`  : "h2",
        `^### (.*?)$` : "h3",
        `^- (.*?)$`   : "li",
    }

    var htmlContent string
        
    for pattern, tag := range mdHTMLTag {
        replacement := fmt.Sprintf("<%s>$1</%s>", tag, tag)
        htmlContent = regexp.MustCompile(pattern).ReplaceAllString(content, replacement)
    }

    return htmlContent
}

func generateHtmlFromTemplate(title, date, content string) string {
    templateReplacements := map[string]string {
        "title"   : title,
        "date"    : date,
        "content" : content,
    }

    templateFile, err := os.Open(templateFilePath)
    if err != nil {
        panic(err)
    }
    defer templateFile.Close()

    var templateLines []string
    scanner := bufio.NewScanner(templateFile)
    for scanner.Scan() {
        templateLines = append(templateLines, scanner.Text())
    }
     
    var htmlContent string
    for pattern, replacement := range templateReplacements {
        pattern = fmt.Sprintf(`{{(.*?)%s(.*?)}}`, pattern)
        htmlContent  = regexp.MustCompile(pattern).ReplaceAllString(content, replacement)
    }

    return htmlContent
}

func main() {
    draftFiles, err := os.ReadDir(draftDir)
    if err != nil {
        panic(err)
    }

    for _, draftFile := range draftFiles {
        if draftFile.IsDir() {
            continue
        }

        draftFileName := draftFile.Name()
        draftFilePath := draftDir + draftFileName
        draftFileHandle, err := os.Open(draftFilePath)
        if err != nil {
            panic(err)
        }
        defer draftFileHandle.Close()

        var draftLines []string
        scanner := bufio.NewScanner(draftFileHandle)
        for scanner.Scan() {
            draftLines = append(draftLines, scanner.Text())
        }

        title, date, content := parseDraft(draftLines)
        htmlContent := generateHtmlFromTemplate(title, date, convertMarkdownToHtml(content))

        publicFilePath := publicDir + strings.Split(draftFileName, ".")[0] + ".html"
        os.WriteFile(publicFilePath, []byte(htmlContent), 0644)
    }
}
