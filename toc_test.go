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
	codeblock, _ := os.ReadFile("testdata/codeblock.md")
	codeblockToc, _ := os.ReadFile("testdata/codeblock_toc.md")
	ignore, _ := os.ReadFile("testdata/ignore.md")
	ignoreToc, _ := os.ReadFile("testdata/ignore_toc.md")
	customHeading, _ := os.ReadFile("testdata/custom_heading.md")
	customHeadingToc, _ := os.ReadFile("testdata/custom_heading_toc.md")

	tt := []struct {
		name    string
		in      []byte
		cfg     *Config
		want    []byte
		wantErr error
	}{
		{"basic", basic, DefaultConfig, basicToc, nil},
		{"special", special, DefaultConfig, specialToc, nil},
		{"repeat", repeat, DefaultConfig, repeatToc, nil},
		{"code block", codeblock, DefaultConfig, codeblockToc, nil},
		{"ignore", ignore, DefaultConfig, ignoreToc, nil},
		{"existing without force", basicToc, DefaultConfig, nil, ErrExistingToc},
		{"existing with force", basicToc, &Config{Force: true}, basicToc, nil},
		{"custom heading", customHeading, &Config{WithTocHeading: true, TocHeading: "Contents"}, customHeadingToc, nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, gotErr := Insert(tc.in, tc.cfg)
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
				},
				Config: DefaultConfig,
			},
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
