package markdown

import (
	. "github.com/yuntan/blackfriday"
	"gopkg.in/yaml.v2"

	"bufio"
	"bytes"
	"regexp"
)

var (
	LinkWorker = func(b []byte) []byte {
		s := string(b)
		if len(s) < 5 {
			return nil
		}
		link := s[2 : len(s)-2]
		return []byte(`<a href="` + link + `">` + link + `</a>`)
	}

	reWikiLink = regexp.MustCompile(`\[\[[^\]]+\]\]`)
)

type MetaData struct {
	Title string   `yaml:"title"` // TODO
	Alias []string `yaml:"alias"`
	Tags  []string `yaml:"tags"`
}

func Convert(input []byte) []byte {
	if input == nil {
		return nil
	}

	// remove front matter
	sep := []byte("---\n")
	if bytes.HasPrefix(input, sep) {
		b := bytes.SplitN(input, sep, 3)
		if len(b) == 3 {
			input = b[2]
		}
	}

	// replace wikilink
	buf := new(bytes.Buffer)
	r := bytes.NewReader(input)
	scn := bufio.NewScanner(r)
	for scn.Scan() {
		// FIXME <pre>
		b := reWikiLink.ReplaceAllFunc(scn.Bytes(), LinkWorker)
		buf.Write(b)
		buf.WriteRune('\n')
	}
	input = buf.Bytes()

	htmlFlags := 0 |
		HTML_HREF_TARGET_BLANK |
		HTML_USE_XHTML |
		HTML_USE_SMARTYPANTS |
		// HTML_SMARTYPANTS_FRACTIONS |
		HTML_SMARTYPANTS_LATEX_DASHES |
		HTML_FOOTNOTE_RETURN_LINKS

	extensions := 0 |
		EXTENSION_NO_INTRA_EMPHASIS |
		EXTENSION_TABLES |
		EXTENSION_FENCED_CODE |
		EXTENSION_AUTOLINK |
		EXTENSION_STRIKETHROUGH |
		EXTENSION_SPACE_HEADERS |
		EXTENSION_FOOTNOTES |
		EXTENSION_HEADER_IDS |
		EXTENSION_AUTO_HEADER_IDS |
		EXTENSION_BACKSLASH_LINE_BREAK |
		EXTENSION_DEFINITION_LISTS |
		EXTENSION_LATEX_MATH

	renderer := HtmlRenderer(htmlFlags, "", "")
	return Markdown(input, renderer, extensions)
}

func Meta(input []byte) MetaData {
	ret := MetaData{}

	// get front matter
	sep := []byte("---\n")
	if !bytes.HasPrefix(input, sep) {
		return ret
	}
	b := bytes.SplitN(input, sep, 3)
	if len(b) != 3 {
		return ret
	}

	// parse front matter
	yaml.Unmarshal(b[1], &ret)
	return ret
}

func TOC(input []byte) []byte {
	// remove front matter
	sep := []byte("---\n")
	if bytes.HasPrefix(input, sep) {
		b := bytes.SplitN(input, sep, 3)
		if len(b) == 3 {
			input = b[2]
		}
	}

	htmlFlags := 0 |
		HTML_TOC |
		HTML_OMIT_CONTENTS |
		HTML_USE_XHTML

	extensions := 0 |
		EXTENSION_SPACE_HEADERS |
		EXTENSION_HEADER_IDS |
		EXTENSION_AUTO_HEADER_IDS

	return Markdown(input, HtmlRenderer(htmlFlags, "", ""), extensions)
}