package mdtoc

import (
	"reflect"
	"strings"
	"testing"
)

func Test_textToLink(t *testing.T) {
	tt := []struct {
		name string
		s    string
		want string
	}{
		{"basic", "Heading One", "heading-one"},
		{"slash", "Heading/With/Slash", "headingwithslash"},
		{"underscore", "Heading_With_Underscore", "heading_with_underscore"},
		{"special characters", "Heading,0.1+2;3:4", "heading01234"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := textToLink(tc.s)
			if got != tc.want {
				t.Fatalf("expected: %s want: %s", got, tc.want)
			}
		})
	}

}

func TestParse(t *testing.T) {
	tt := []struct {
		name string
		in   string
		want *Toc
	}{
		{
			"basic",
			"# Title\n\n## Heading 1\n\n### Heading 2",
			&Toc{
				Heading: defaultTocHeading,
				Bullets: []Bullet{
					{Indent: 0, Text: "Heading 1", Link: "heading-1"},
					{Indent: 1, Text: "Heading 2", Link: "heading-2"},
				}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			f := strings.NewReader(tc.in)

			got, err := Parse(f)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected: %v want: %v", got, tc.want)
			}
		})
	}
}
