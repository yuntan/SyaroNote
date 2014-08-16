package main

import (
	. "github.com/OUCC/syaro/logger"
	"github.com/OUCC/syaro/setting"
	"github.com/OUCC/syaro/wikiio"

	"bufio"
	"bytes"
	"regexp"
)

func processWikiLink(b []byte, currentDir string) []byte {
	const RE_DOUBLE_BRACKET = "\\[\\[[^\\]]+\\]\\]"

	reader := bytes.NewReader(b)
	scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer

	re := regexp.MustCompile(RE_DOUBLE_BRACKET)

	for scanner.Scan() {
		line := scanner.Bytes()

		for {
			index := re.FindIndex(line)

			if len(index) != 0 { // tag found
				LoggerV.Println("processWikiLink: bracket tag found:",
					string(line[index[0]:index[1]]))

				name := line[index[0]+2 : index[1]-2]
				files := searchPage(string(name), currentDir)

				if len(files) != 0 { // page found
					LoggerV.Println("processWikiLink:", len(files), "pages found")
					LoggerV.Println("processWikiLink: select ", files[0].WikiPath())
					// TODO avoid ambiguous page
					line = embedLinkTag(line, index, name, files[0])

				} else { // page not found
					LoggerV.Println("processWikiLink: no page found")
					// TODO invalid link
					line = embedLinkTag(line, index, name, nil)
				}

			} else { // tag not found, so go next line
				break
			}
		}
		buffer.Write(line)
	}

	return buffer.Bytes()
}

func embedLinkTag(line []byte, tagIndex []int, linkname []byte, file *wikiio.WikiFile) []byte {
	if file == nil {
		// TODO file not found page
		return bytes.Join([][]byte{
			line[:tagIndex[0]],
			[]byte("<a class=\"notfound\" href=\""),
			[]byte(setting.UrlPrefix),
			[]byte("/404.html?name="),
			linkname,
			[]byte("\">"),
			linkname,
			[]byte("</a>"),
			line[tagIndex[1]:],
		}, nil)
	}
	return bytes.Join([][]byte{
		line[:tagIndex[0]],
		[]byte("<a href=\""),
		[]byte(file.URLPath()),
		[]byte("\">"),
		linkname,
		[]byte("</a>"),
		line[tagIndex[1]:],
	}, nil)
}

// TODO security check
func searchPage(name string, currentDir string) []*wikiio.WikiFile {
	if name == "" {
		return nil
	}

	// TODO
	// if filepath.IsAbs(name) {
	// 	// search name as absolute path
	// 	// example: /piyo /poyo/pyon.ext
	// 	return searchPageByAbsPath(name, currentDir)
	// } else if strings.Contains(name, "/") ||util.IsMarkdown(name) {
	// 	// search name as relative path
	// 	// example: ./hoge ../fuga.ext puyo.ext
	// 	return  searchPageByRelPath(name, currentDir)
	// } else {
	// 	// search name as base name
	// 	// example: abc
	return searchPageByBaseName(name)
	// }
}

func searchPageByBaseName(baseName string) []*wikiio.WikiFile {
	LoggerV.Printf("searchPageByBaseName(%s)", baseName)
	files, _ := wikiio.Search(baseName)
	return files
}
