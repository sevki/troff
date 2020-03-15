package troff

import (
	"bytes"
	"reflect"
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
			t.Log(string(pdf))
			t.Fail()
		})
	}
}

func Test_postTime_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		t       *postTime
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("postTime.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_troffrenderer_RenderNode(t *testing.T) {
	type args struct {
		n        *blackfriday.Node
		entering bool
	}
	tests := []struct {
		name  string
		r     *troffrenderer
		args  args
		want  blackfriday.WalkStatus
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if got := tt.r.RenderNode(w, tt.args.n, tt.args.entering); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("troffrenderer.RenderNode() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("troffrenderer.RenderNode() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_troffrenderer_RenderHeader(t *testing.T) {
	type args struct {
		n *blackfriday.Node
	}
	tests := []struct {
		name  string
		r     *troffrenderer
		args  args
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.r.RenderHeader(w, tt.args.n)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("troffrenderer.RenderHeader() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_troffrenderer_RenderFooter(t *testing.T) {
	type args struct {
		n *blackfriday.Node
	}
	tests := []struct {
		name  string
		r     *troffrenderer
		args  args
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.r.RenderFooter(w, tt.args.n)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("troffrenderer.RenderFooter() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_renderTroff(t *testing.T) {
	type args struct {
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderTroff(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderTroff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("renderTroff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tr2ps(t *testing.T) {
	type args struct {
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tr2ps(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("tr2ps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tr2ps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_renderps(t *testing.T) {
	type args struct {
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderps(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("renderps() = %v, want %v", got, tt.want)
			}
		})
	}
}
