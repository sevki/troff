package troff

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path"
	"strings"

	blackfriday "github.com/russross/blackfriday/v2"
	yaml "gopkg.in/yaml.v2"
	"sevki.org/x/debug"
	"sevki.org/x/markdown"
)

// NewRenderer returns a blackfriday.Renderer that
// will print plan9 troff macros
// https://plan9.io/sys/doc/troff.pdf
func NewRenderer() markdown.Renderer { return &troffrenderer{} }

type troffrenderer struct {
	parsingTitle bool
	title        *markdown.Post
	list         []int
	buf          *bytes.Buffer
}

func (r *troffrenderer) Post() *markdown.Post { return r.title }

// RenderNode is the main rendering method. It will be called once for
// every leaf node and twice for every non-leaf node (first with
// entering=true, then with entering=false). The method should write its
// rendition of the node to the supplied writer w.
func (r *troffrenderer) RenderNode(w io.Writer, n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch n.Type {
	case blackfriday.Document:
		return blackfriday.GoToNext
	case blackfriday.Text:
		if r.parsingTitle {
			magic := []byte(`title:`)
			// BUG(sevki): we get things that are not title blocks here
			isTitle := strings.Index(string(n.Literal), string(magic)) == 0
			if err := yaml.Unmarshal(n.Literal, &r.title); isTitle && err != nil {
				panic(err.Error())
				return blackfriday.Terminate
			}
			r.parsingTitle = false
			return blackfriday.GoToNext
		}
		if len(n.Literal) > 0 {
			if r.buf != nil {
				fmt.Fprint(r.buf, strings.TrimLeft(string(n.Literal), " "))
			} else {
				fmt.Fprint(w, strings.TrimLeft(string(n.Literal), " "))
			}
		}
		return blackfriday.GoToNext
	case blackfriday.BlockQuote:
		if entering {
			LeftAlignedParagraph(w, string(n.Literal), 1)
		} else {
			LineBreak(w)
		}
		return blackfriday.GoToNext
	case blackfriday.HTMLBlock:
		if entering {
			HTML(w, string(n.Literal))
		}
		return blackfriday.GoToNext
	case blackfriday.Code:
		LineBreak(w)
		Bold(w, string(n.Literal))
		Roman(w)
		return blackfriday.GoToNext
	case blackfriday.HTMLSpan:
		return blackfriday.GoToNext
	case blackfriday.Paragraph:
		if entering {
			level := len(r.list)
			LeftAlignedParagraph(w, string(n.Literal), level)
		} else {
			LineBreak(w)
		}
		return blackfriday.GoToNext
	case blackfriday.TableHead, blackfriday.TableBody:
		return blackfriday.GoToNext
	case blackfriday.Table:
		if entering {
			r.buf = bytes.NewBuffer(nil)
		} else {
			Table(w, r.buf.Bytes())
			r.buf = nil
		}
		return blackfriday.GoToNext
	case blackfriday.TableRow:
		if !entering {
			fmt.Fprint(r.buf, linebreak)
		}
		return blackfriday.GoToNext
	case blackfriday.TableCell:
		if !entering {
			fmt.Fprint(r.buf, "\t")
		}
		return blackfriday.GoToNext
	case blackfriday.List:
		if entering {
			if len(r.list) > 0 {
				Indent(w)
			}
			r.list = append(r.list, 0)
		} else {
			if len(r.list) > 0 {
				r.list = r.list[:len(r.list)-1]
			}
			Outdent(w)
		}
		return blackfriday.GoToNext
	case blackfriday.Item:
		if entering {
			last := r.list[len(r.list)-1] + 1
			r.list[len(r.list)-1] = last
			bullet := ""
			if n.ListData.ListFlags&blackfriday.ListTypeOrdered != 0 {
				for _, n := range r.list {
					bullet += fmt.Sprintf("%d.", n)
				}
			} else {
				bullet += fmt.Sprintf("%c", rune(n.ListData.BulletChar))
			}
			IndentParagraph(w, bullet)
		}
		return blackfriday.GoToNext
	case blackfriday.Emph:
		if entering {
			Italic(w)
		} else {
			LineBreak(w)
			Roman(w)
		}
		return blackfriday.GoToNext
	case blackfriday.Strong:
		if entering {
			Bold(w)
		} else {
			LineBreak(w)
			Roman(w)
		}
		return blackfriday.GoToNext
	case blackfriday.Link:
		// TODO(sevki): find the link macro
		return blackfriday.GoToNext
	case blackfriday.Image:
		return blackfriday.SkipChildren

		if entering {
			s := strings.Split(string(n.LinkData.Destination), "=")
			filename := strings.TrimSpace(s[0])
			sizebit := "0x0"
			if len(s) > 1 {
				sizebit = s[1]
			}
			width, height := parseSize(sizebit)
			switch path.Ext(filename) {
			case ".ps":
			default:
				var err error
				if filename, err = convertToPs(filename, int(width), int(height)); err != nil {
					panic(err)
					CodeBlock(w, []byte(err.Error()))
					return blackfriday.Terminate
				}
			}

			if height < 0 {
				height = 1
			}
			if width < 0 {
				width = 1
			}
			Picture(w, filename, float32(width)/96, float32(height)/96, Left, -1.0, 0)
		}
		return blackfriday.SkipChildren
	case blackfriday.CodeBlock:
		CodeBlock(w, n.Literal)
		return blackfriday.GoToNext
	case blackfriday.Heading:
		r.parsingTitle = n.IsTitleblock
		if n.IsTitleblock {
			if entering {
				return blackfriday.GoToNext
			}
			printTitleBlock(w, r.title)
			return blackfriday.GoToNext
		}
		if entering && !n.IsTitleblock {
			fmt.Fprintln(w, ".SH")
		}
		if !entering {
			fmt.Fprintln(w, "")
		}
		return blackfriday.GoToNext
	default:
		panic(fmt.Sprintf("%s %s", n, n.Type))
	}

}

// RenderHeader is a method that allows the renderer to produce some
// content preceding the main body of the output document. The header is
// understood in the broad sense here. For example, the default HTML
// renderer will write not only the HTML document preamble, but also the
// table of contents if it was requested.
//
// The method will be passed an entire document tree, in case a particular
// implementation needs to inspect it to produce output.
//
// The output should be written to the supplied writer w. If your
// implementation has no header to write, supply an empty implementation.
func (r *troffrenderer) RenderHeader(w io.Writer, n *blackfriday.Node) {}

// RenderFooter is a symmetric counterpart of RenderHeader.
func (r *troffrenderer) RenderFooter(w io.Writer, n *blackfriday.Node) {
	fmt.Fprintln(w, ".SG")

}

func renderTroff(payload []byte) ([]byte, error) {
	args := []string{"-Tutf", "-ms", "-mpictures"}

	buf := bytes.NewBuffer(payload)
	cmd := exec.CommandContext(context.Background(), "/usr/lib/plan9/bin/troff", args...)
	cmd.Stdin = buf
	stdOut := bytes.NewBuffer(nil)
	stdErr := bytes.NewBuffer(nil)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	fmt.Printf("running troff: %v\n", args)
	if err := cmd.Run(); err != nil {
		debug.Indent(stdErr, 1)
		log.Printf(`troff: %v
args= %v
%s
`,
			err,
			args,
			stdErr.String(),
		)
		return nil, err
	}
	return stdOut.Bytes(), nil
}

func tr2ps(payload []byte) ([]byte, error) {
	args := []string{}

	buf := bytes.NewBuffer(payload)
	cmd := exec.CommandContext(context.Background(), "tr2post", args...)
	cmd.Stdin = buf
	stdOut := bytes.NewBuffer(nil)
	stdErr := bytes.NewBuffer(nil)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	fmt.Printf("running tr2post: %v\n", args)
	if err := cmd.Run(); err != nil {
		debug.Indent(stdErr, 1)
		log.Printf(`troff: %v
args= %v
%s
`,
			err,
			args,
			stdErr.String(),
		)
		return nil, err
	}
	return stdOut.Bytes(), nil
}

func renderps(payload []byte) ([]byte, error) {

	buf := bytes.NewBuffer(payload)
	args := []string{"-", "example.pdf"}

	cmd := exec.CommandContext(context.Background(), "ps2pdf", args...)
	cmd.Stdin = buf
	stdOut := bytes.NewBuffer(nil)
	stdErr := bytes.NewBuffer(nil)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	fmt.Printf("running ps2pdf: %v\n", args)
	if err := cmd.Run(); err != nil {
		debug.Indent(stdErr, 1)
		log.Printf(`ps2pdf: %v
args= %v
%s
`,
			err,
			args,
			stdErr.String(),
		)
		return nil, err
	}
	return stdOut.Bytes(), nil
}
