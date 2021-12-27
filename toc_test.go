package mdtoc

import (
	"os"
	"reflect"
	"testing"
)

func TestInsert(t *testing.T) {
	// read testdata into memory
	basic, _ := os.ReadFile("testdata/basic.md")
	basicToc, _ := os.ReadFile("testdata/basic_toc.md")
	special, _ := os.ReadFile("testdata/special.md")
	specialToc, _ := os.ReadFile("testdata/special_toc.md")
	repeat, _ := os.ReadFile("testdata/repeat.md")
	repeatToc, _ := os.ReadFile("testdata/repeat_toc.md")

	tt := []struct {
		name    string
		in      []byte
		force   bool
		want    []byte
		wantErr error
	}{
		{"basic", basic, false, basicToc, nil},
		{"special", special, false, specialToc, nil},
		{"repeat", repeat, false, repeatToc, nil},
		{"existing without force", basicToc, false, nil, ErrExistingToc},
		{"existing with force", basicToc, true, basicToc, nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			toc, err := New(tc.in)
			if err != nil {
				t.Fatalf("parsing toc: %v", err)
			}

			got, gotErr := toc.Insert(tc.in, tc.force)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %s got: %s", string(got), string(tc.want))
			}
			if gotErr != tc.wantErr {
				t.Fatalf("expected: %v got: %v", tc.wantErr, gotErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tt := []struct {
		name string
		in   []byte
		want *Toc
	}{
		{
			"basic",
			[]byte("# Title\n\n## Heading 1\n\n### Heading 2"),
			&Toc{
				Items: []Item{
					{Indent: 0, Text: "Heading 1", Link: "heading-1"},
					{Indent: 1, Text: "Heading 2", Link: "heading-2"},
				}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := New(tc.in)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v got: %v", tc.want, got)
			}
		})
	}
}

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
			if tc.want != got {
				t.Fatalf("expected: %s got: %s", tc.want, got)
			}
		})
	}
}
