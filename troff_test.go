package troff

import (
	"fmt"
	"testing"

	blackfriday "github.com/russross/blackfriday/v2"
)

func m(s string) []byte { return []byte(s) }

func TestRender(t *testing.T) {
	tests := map[string]struct {
		payload []byte
	}{
		"title": {payload: m(`% title: markdown troff renderer
% authors:
% - name: Sevki Hasirci
%   email: s@sevki.org
%   affiliation: funemployment
% date: Sat Feb 22 15:18:37 GMT 2020
% 
% tags: [acme, plan9]
% abstract: Creating documents using plan9 troff 

# Heading 1
## Heading 2
### Heading 3

![glenda](glenda.svg =32x32)

*italic*

__bold__

plain text

1. First ordered list item
2. Another item
	* Unordered sub-list. 
1. Actual numbers don't matter, just that it's a number
	1. Ordered sub-list
4. And another item.

	You can have properly indented paragraphs within list items. Notice the blank line above, and the leading spaces (at least one, but we'll use three here to also align the raw Markdown).

	To have a line break without a paragraph, you will need to use two trailing spaces.
	Note that this line is separate, but within the same paragraph.
	(This is contrary to the typical GFM line break behaviour, where trailing spaces are not required.)

* Unordered list can use asterisks
- Or minuses
+ Or pluses

some text followed by

		code blocks
		with multiple lines


Name|Age|Country
--------|-----|-----
Bob|27|tur
Alice|23|uk	

[github.com/sevki/mdms](https://github.com/sevki/mdms)

`)},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			out := blackfriday.Run(test.payload,
				blackfriday.WithRenderer(&troffrenderer{}),
				blackfriday.WithExtensions(blackfriday.Titleblock|blackfriday.CommonExtensions),
			)
			t.Log(string(out))
			bytez, err := renderTroff(out)
			if err != nil {
				t.Log(err)
				t.FailNow()
				return
			}
			psbytez, err := tr2ps(bytez)
			if err != nil {
				t.Log(err)
				t.FailNow()
				return
			}
			pdf, err := renderps(psbytez)
			if err != nil {
				t.Log(err)
				t.FailNow()
				return
			}
			_ = pdf
		})
	}
}

func ExampleRun() {
	t := []byte(`% title: markdown to troff
% authors:
% - name: Sevki Hasirci
%   email: s@sevki.org
%   affiliation: funemployment
% date: Sat Feb 22 15:18:37 GMT 2020
% 
% tags: [plan9, troff, markdown]
% abstract: Creating documents using plan9 troff 

# Hello, Troff!
`)
	out := blackfriday.Run(t,
		blackfriday.WithRenderer(&troffrenderer{}),
		blackfriday.WithExtensions(blackfriday.Titleblock|blackfriday.CommonExtensions),
	)
	fmt.Print(string(out))

	// Output:
	// .HTML markdown troff renderer
	// .TL
	//  markdown troff renderer
	// .AU
	// .I Sevki Hasirci
	// .I s@sevki.org
	// .AI funemployment
	// .ND Feb 22, 2020
	// .AB
	//  Creating documents using plan9 troff
	// .AE
	// .PP
	// .SH
	// Hello troff!!
	// .SG
}
