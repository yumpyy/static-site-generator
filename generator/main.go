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

func generateHtmlFromTemplate(title, date string, content []string) string {
    modifiedContent := strings.Join(content, "")
    templateReplacements := map[string]string {
        "title"   : title,
        "date"    : date,
        "content" : modifiedContent,
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
     
    var htmlLines []string
    for _, line := range templateLines {
        line = strings.TrimSpace(line)
        for pattern, replacement := range templateReplacements {
            pattern = fmt.Sprintf(`{{(.*?)%s(.*?)}}`, pattern)
            line = regexp.MustCompile(pattern).ReplaceAllString(line, replacement)
        }
        htmlLines = append(htmlLines, fmt.Sprintf("%s\n", line))

    }
    htmlContent := strings.Join(htmlLines, "")

    return htmlContent
}

func parseDraft(lines []string) (title, date string, content []string) {
	title = strings.TrimSpace(strings.Split(lines[0], ":")[1])
	date = strings.TrimSpace(strings.Split(lines[1], ":")[1])
    content = lines[3:]
	if date == "" {
		date = strings.ToLower(time.Now().Format("02-01-2006"))
	}

    return
}

func convertMarkdownToHtml(content []string) []string {
    var htmlContent []string

    mdHTMLTag := map[string]string {
        `\*(.*?)\*`   : "b",
        `_(.*?)_`     : "i",
        `~(.*?)~`     : "s",
        `^# (.*?)$`   : "h1",
        `^## (.*?)$`  : "h2",
        `^### (.*?)$` : "h3",
        `^- (.*?)$`   : "li",
    }

    for _, line := range content {
        line = strings.TrimSpace(line)
        for pattern, tag := range mdHTMLTag {
            replacement := fmt.Sprintf("<%s>$1</%s>", tag, tag)
            line = regexp.MustCompile(pattern).ReplaceAllString(line, replacement)
        }
        htmlContent = append(htmlContent, fmt.Sprintf("%s\n", line))

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
