// Package mdtoc generates a table of contents for an existing markdown document
package mdtoc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// ErrExistingToc is thrown if the provided document already contains a mdtoc-generated
	// table of contents
	ErrExistingToc = errors.New("document has existing table of contents")

	// headingRegex is the expression which will match non-title heading lines
	headingRegex = regexp.MustCompile("^([#]{2,})[ ]+(.+)")

	tocBegin = "<!---mdtoc begin--->"
	tocEnd   = "<!---mdtoc end--->"
)

// Item represents a single line in the table of contents
type Item struct {
	Indent int
	Text   string
	Link   string
}

// Toc stores table of contents metadata
type Toc struct {
	Items []Item
}

// Bytes returns a markdown formatted slice of bytes
func (t *Toc) Bytes() []byte {
	var buf []byte
	w := bytes.NewBuffer(buf)

	w.WriteString(fmt.Sprintf("%s\n", tocBegin))
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
func Insert(b []byte, force bool) ([]byte, error) {
	toc, err := New(b)
	if err != nil {
		return b, err
	}

	var new []byte
	buf := bytes.NewBuffer(new)

	var inOld, newAdded bool

	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {

		// handle any previously existing toc's
		// begin comment, set flag and skip
		if strings.EqualFold(tocBegin, scanner.Text()) {
			if !force {
				return nil, ErrExistingToc
			}
			inOld = true
			continue
		}
		// end comment, reset flag and skip
		if inOld && strings.EqualFold(scanner.Text(), tocEnd) {
			inOld = false
			continue
		}
		// old toc line, skip
		if inOld {
			continue
		}

		// when the first non-title heading is encoutered, insert new toc just before it
		if !newAdded && headingRegex.FindSubmatch(scanner.Bytes()) != nil {
			buf.Write(toc.Bytes())
			newAdded = true
		}

		buf.Write(scanner.Bytes())
		buf.Write([]byte("\n"))
	}

	if err := scanner.Err(); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}

// New extacts table of contents attributes from an existing document
func New(b []byte) (*Toc, error) {
	toc := Toc{}

	var inCodeBlock bool

	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {

		// handle code blocks to ensure `#` are not captured as headings
		// begin code block, set flag and skip
		if strings.HasPrefix(scanner.Text(), "```") && !inCodeBlock {
			inCodeBlock = true
			continue
		}
		// end code block, reset flag and skip
		if inCodeBlock && strings.HasPrefix(scanner.Text(), "```") {
			inCodeBlock = false
			continue
		}
		// code block line, skip
		if inCodeBlock {
			continue
		}

		m := headingRegex.FindStringSubmatch(scanner.Text())
		if len(m) == 3 {
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

// textToLink returns the header link formatted version of a string
//
// ex. `Header One Two` = `header-one-two`
func textToLink(s string) string {
	// TODO: find a more comprehensive/formally documented list of these
	rep := strings.NewReplacer(" ", "-", "/", "", ",", "", ".", "", "+", "", ":", "", ";", "", "`", "", `"`, "", `'`, "", "{", "", "}", "")
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
