// Package mdtoc generates a table of contents for an existing markdown document
package mdtoc

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	tocBegin  = "<!--mdtoc: begin-->"
	tocEnd    = "<!--mdtoc: end-->"
	tocIgnore = "<!--mdtoc: ignore-->"
)

var (
	// ErrExistingToc is thrown if the provided document already contains a mdtoc-generated
	// table of contents
	ErrExistingToc = errors.New("document has existing table of contents")

	// Version is the current library/CLI version
	//go:embed VERSION
	Version string

	// DefaultConfig defines the default configuration settings
	//
	// These field values will align with the default flag values from the CLI
	DefaultConfig = &Config{
		Force:          false,
		WithTocHeading: false,
		TocHeading:     DefaultTocHeading,
	}

	// DefaultTocHeading is the default heading applied when enabled
	DefaultTocHeading = "Table of Contents"

	// headingRegex is the expression which will match non-title heading lines
	headingRegex = regexp.MustCompile("^([#]{2,})[ ]+(.+)")
)

// Item represents a single line in the table of contents
type Item struct {
	Indent int
	Text   string
	Link   string
}

// Toc stores table of contents metadata
type Toc struct {
	Items  []Item
	Config *Config
}

// Config stores settings to be used when inserting a new Toc
type Config struct {
	Force          bool
	WithTocHeading bool
	TocHeading     string
}

// Bytes returns a markdown formatted slice of bytes
func (t *Toc) Bytes() []byte {
	var buf []byte
	w := bytes.NewBuffer(buf)

	w.WriteString(fmt.Sprintf("%s\n", tocBegin))
	if t.Config != nil && t.Config.WithTocHeading {
		w.WriteString(fmt.Sprintf("## %s %s\n\n", t.Config.TocHeading, tocIgnore))
	}
	for _, item := range t.Items {
		w.WriteString(fmt.Sprintf("%s* [%s](#%s)\n", strings.Repeat(" ", item.Indent*2), item.Text, item.Link))
	}
	w.WriteString(fmt.Sprintf("%s\n", tocEnd))

	return w.Bytes()
}

// String returns a markdown formatted string
func (t *Toc) String() string {
	return string(t.Bytes())
}

// Insert returns a copy of an existing document with a table of contents inserted
func Insert(b []byte, cfg *Config) ([]byte, error) {
	toc, err := New(b)
	if err != nil {
		return b, err
	}
	toc.Config = cfg

	var new []byte
	buf := bytes.NewBuffer(new)

	var inOld, newAdded bool

	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		// handle any previously existing toc's
		// begin comment, set flag and skip
		if strings.EqualFold(tocBegin, line) {
			if !cfg.Force {
				return nil, ErrExistingToc
			}
			inOld = true
			continue
		}
		// end comment, write toc in same location and reset flag
		if inOld && strings.EqualFold(line, tocEnd) {
			buf.Write(toc.Bytes())
			newAdded = true
			inOld = false
			continue
		}
		// old toc line, skip
		if inOld {
			continue
		}

		// when the first non-title heading is encoutered, insert new toc just before it
		if !newAdded && headingRegex.FindStringSubmatch(line) != nil {
			buf.Write(append(toc.Bytes(), []byte("\n")...))
			newAdded = true
		}

		buf.Write(append(scanner.Bytes(), []byte("\n")...))
	}

	if err := scanner.Err(); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}

// New extacts table of contents attributes from an existing document
func New(b []byte) (*Toc, error) {
	toc := Toc{Config: DefaultConfig}

	var inCodeBlock bool

	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		// handle code blocks to ensure `#` are not captured as headings
		// begin code block, set flag and skip
		if strings.HasPrefix(line, "```") && !inCodeBlock {
			inCodeBlock = true
			continue
		}
		// end code block, reset flag and skip
		if inCodeBlock && strings.HasPrefix(line, "```") {
			inCodeBlock = false
			continue
		}
		// code block line, skip
		if inCodeBlock {
			continue
		}

		m := headingRegex.FindStringSubmatch(line)
		if len(m) == 3 {
			// skip headings with ignore comments
			if strings.Contains(line, tocIgnore) {
				continue
			}

			// m[0]: Full regular expression match
			// m[1]: First match group (two or more `#` characters)
			// m[2]: Second match group (text of heading)
			toc.Items = append(toc.Items,
				Item{
					Indent: len(m[1]) - 2,
					Text:   m[2],
					Link:   textToLink(m[2]),
				})
		}
	}

	if err := scanner.Err(); err != nil {
		return &toc, err
	}

	toc.updateRepeatLinks()
	return &toc, nil
}

// textToLink returns the heading link formatted version of a string
//
// ex. `Heading One Two` = `heading-one-two`
func textToLink(s string) string {
	// TODO: find a more comprehensive/formally documented list of these
	rep := strings.NewReplacer(
		" ", "-",
		"/", "",
		",", "",
		".", "",
		"+", "",
		":", "",
		";", "",
		"`", "",
		`"`, "",
		`'`, "",
		"{", "",
		"}", "",
		"(", "",
		")", "",
	)
	return strings.ToLower(rep.Replace(s))
}

// updateRepeatLinks fixes the generated link text if the generated text is repeated
// in the same contents
func (t *Toc) updateRepeatLinks() {
	lookup := make(map[string]int, len(t.Items))

	for i, item := range t.Items {
		// if key already exists in the lookup, the  link text needs to append a `-n`,
		// where `n` is the number of previous occurrences. if the key does not already
		// exist, add a new key and set occurrences to 1.
		if val, ok := lookup[item.Link]; ok {
			key := item.Link // preserve the original lookup key
			t.Items[i].Link = fmt.Sprintf("%s-%d", item.Link, val)
			lookup[key]++
			continue
		}

		lookup[item.Link] = 1
	}
}
