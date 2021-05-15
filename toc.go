// Package mdtoc generates a table of contents for an existing markdown document
package mdtoc

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	// HeadingRegex is the expression which will match heading lines
	HeadingRegex = regexp.MustCompile("^([#]{2,})[ ]+(.+)")

	defaultTocHeading = "## Table of Contents"
)

// Bullet represents a single line in the table of contents
type Bullet struct {
	Indent int
	Text   string
	Link   string
}

// Toc stores table of contents metadata
type Toc struct {
	Heading string
	Bullets []Bullet
}

// Bytes returns a markdown formatted slice of bytes
func (t *Toc) Bytes() []byte {
	var buf []byte
	w := bytes.NewBuffer(buf)

	w.WriteString(fmt.Sprintf("%s\n\n", t.Heading))
	for _, b := range t.Bullets {
		w.WriteString(fmt.Sprintf("%s* [%s](#%s)\n", strings.Repeat(" ", b.Indent*2), b.Text, b.Link))
	}

	return w.Bytes()
}

// String returns a markdown formatted string
func (t *Toc) String() string {
	return string(t.Bytes())
}

// AddToc inserts a table of contents into an existing markdown document
func AddToc(f *os.File, t Toc, overwrite bool) error {
	return nil
}

// ParseToc extacts table of contents attributes from an existing document
func ParseToc(f *os.File) (Toc, error) {
	toc := Toc{Heading: defaultTocHeading}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m := HeadingRegex.FindStringSubmatch(scanner.Text())
		if len(m) == 3 {
			// m[0]: Full regular expression match
			// m[1]: First match group (two or more `#` characters)
			// m[2]: Second match group (text of heading)
			toc.Bullets = append(toc.Bullets,
				Bullet{
					Indent: len(m[1]) - 2,
					Text:   m[2],
					Link:   textToLink(m[2]),
				})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return toc, nil
}

// textToLink returns the header link formatted version of a string
//
// ex. `Header One Two` = `header-one-two`
func textToLink(s string) string {
	rep := strings.NewReplacer(" ", "-", "/", "", ",", "", ".", "", "+", "", ":", "", ";", "")
	return strings.ToLower(rep.Replace(s))
}
