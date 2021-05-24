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
	// HeadingRegex is the expression which will match non-title heading lines
	HeadingRegex = regexp.MustCompile("^([#]{2,})[ ]+(.+)")

	// ExistingTocError is thrown if the provided document already contains a mdtoc-generated
	// table of contents
	ExistingTocError = errors.New("document has existing table of contents")

	tocBegin = "<!---mdtoc begin--->"
	tocEnd   = "<!---mdtoc end--->"
)

// Bullet represents a single line in the table of contents
type Bullet struct {
	Indent int
	Text   string
	Link   string
}

// Toc stores table of contents metadata
type Toc struct {
	Bullets []Bullet
}

// Bytes returns a markdown formatted slice of bytes
func (t *Toc) Bytes() []byte {
	var buf []byte
	w := bytes.NewBuffer(buf)

	w.WriteString(fmt.Sprintf("%s\n", tocBegin))
	for _, b := range t.Bullets {
		w.WriteString(fmt.Sprintf("%s* [%s](#%s)\n", strings.Repeat(" ", b.Indent*2), b.Text, b.Link))
	}
	w.WriteString(fmt.Sprintf("%s\n\n", tocEnd))

	return w.Bytes()
}

// String returns a markdown formatted string
func (t *Toc) String() string {
	return string(t.Bytes())
}

// Add returns a copy of an existing document with a table of contents inserted
func Add(b []byte, toc *Toc, force bool) ([]byte, error) {
	var new []byte
	buf := bytes.NewBuffer(new)

	inOld := false
	newAdded := false

	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {

		// handle any previously existing toc's
		// begin comment, set flag and skip
		if strings.EqualFold(tocBegin, scanner.Text()) {
			if !force {
				return nil, ExistingTocError
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
		if !newAdded && HeadingRegex.FindSubmatch(scanner.Bytes()) != nil {
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

// Parse extacts table of contents attributes from an existing document
func Parse(b []byte) (*Toc, error) {
	toc := Toc{}

	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
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
		return &toc, err
	}

	return &toc, nil
}

// textToLink returns the header link formatted version of a string
//
// ex. `Header One Two` = `header-one-two`
func textToLink(s string) string {
	// TODO: find a more comprehensive/formally documented list of these
	rep := strings.NewReplacer(" ", "-", "/", "", ",", "", ".", "", "+", "", ":", "", ";", "")
	return strings.ToLower(rep.Replace(s))
}
