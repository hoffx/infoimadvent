package parser

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

// Parse is a method, that converts the markdown input into a sanitized html string
func Parse(input []byte) string {
	unsafe := blackfriday.Run(input)
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	return string(p.SanitizeBytes(unsafe))
}

// ParseAndProcess is a method, that converts the markdown input into a sanitized html string and
// provides a way to process this output with a slice of functions
func ParseAndProcess(input []byte, operations []func(*string) error) (string, error) {
	html := Parse(input)
	var err error
	for _, o := range operations {
		err = o(&html)
		if err != nil {
			return html, err
		}
	}
	return html, nil
}
