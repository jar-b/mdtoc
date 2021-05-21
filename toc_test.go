package mdtoc

import "testing"

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
		got := textToLink(tc.s)
		t.Run(tc.name, func(t *testing.T) {
			if got != tc.want {
				t.Fatalf("expected: %s want: %s", got, tc.want)
			}
		})
	}

}
