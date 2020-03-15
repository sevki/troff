package troff

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

// Macro is a troff macro
type Macro string

const (
	singleColumn            = "1C"
	doubleColumn            = "2C"
	beginAbstract           = "AB"
	endAbstract             = "AE"
	institution             = "AI"
	author                  = "AU"
	bold                    = "B"
	dateOnPage              = "DA"
	displayEnd              = "DE"
	displayStart            = "DS"
	tableEnd                = "TE"
	tableStart              = "TS"
	beginEquation           = "EQ"
	endEquation             = "EN"
	beginFootnote           = "FS"
	endFootnote             = "FE"
	italic                  = "I"
	beginAndIndentParagraph = "IP"
	endKeep                 = "KE"
	beginKeep               = "KF"
	startKeep               = "KS"
	increaseTypeSize        = "LG"
	leftAlignedParagraph    = "LP"
	changeDate              = "ND"
	numberedHeading         = "NH"
	normalType              = "NL"
	beginParagraph          = "PP"
	roman                   = "R"
	releasePaper            = "RP"
	endIndent               = "RE"
	beginIndent             = "RS"
	signature               = "SG"
	sectionHeading          = "SH"
	decreaseTypeSize        = "SM"
	title                   = "TL"
	tableHeading            = "TH"
	underline               = "UL"
	beginCodeBlock          = "P1"
	endCodeBlock            = "P2"
	html                    = "HTML"
	linebreak               = "\n"
)

func msPrint(ms Macro, args ...interface{}) []byte {
	frmt := ".%s"
	elems := []interface{}{ms}
	for _, arg := range args {
		switch arg.(type) {
		case string:
			frmt += " %s"
		case int:
			frmt += " %d"
		}
		elems = append(elems, arg)
	}
	frmt += "\n"
	return []byte(fmt.Sprintf(frmt, elems...))
}

// Display prints
//	.DS
//	<p>
//	.DE
func Display(w io.Writer, p []byte) {
	w.Write(msPrint(displayStart))
	w.Write(p)
	w.Write(msPrint(displayEnd))
}

// Table prints
//	.TS
//	.DS
//	<p>
//	.DE
//	.TE
func Table(w io.Writer, p []byte) {
	buf := bytes.NewBuffer(nil)
	tw := tabwriter.NewWriter(buf, 0, 0, 4, ' ', tabwriter.TabIndent)
	tw.Write(p)
	tw.Flush()
	tbl := buf.Bytes()
	i := bytes.IndexRune(tbl, '\n')
	header, body := tbl[:i], tbl[i+1:]
	buf = bytes.NewBuffer(nil)
	Esc(buf, msPrint(tableStart, "H"))
	Esc(buf, header)
	LineBreak(buf)
	TableHeading(buf)
	Esc(buf, body)
	Esc(buf, msPrint(tableEnd))
	Display(w, buf.Bytes())
}

// Esc should print escaped bytes, but it just writes it.
// TODO(sevki): fix!
func Esc(w io.Writer, p []byte) { w.Write(p) }

// TableHeading prints
//	.TH\n
func TableHeading(w io.Writer) { w.Write(msPrint(tableHeading)) }

// Indent prints
//	.RS\n
func Indent(w io.Writer) { w.Write(msPrint(beginIndent)) }

// Outdent prints
//	.RE\n
func Outdent(w io.Writer) { w.Write(msPrint(endIndent)) }

// IndentParagraph takes a <title> and prints
//	.IP <title>\n
func IndentParagraph(w io.Writer, title string) { w.Write(msPrint(beginAndIndentParagraph, title)) }

// LeftAlignedParagraph prints
//	.LP
func LeftAlignedParagraph(w io.Writer, title string, level int) {
	if level < 0 {
		w.Write(msPrint(leftAlignedParagraph, linebreak, title))
	} else {
		io.WriteString(w, linebreak)
		io.WriteString(w, title)
	}
}

// Bold prints
//	.B
func Bold(w io.Writer) { w.Write(msPrint(bold)) }

// Roman prints
//	.R
func Roman(w io.Writer) { w.Write(msPrint(roman)) }

// Italic prints
//	.I
func Italic(w io.Writer) { w.Write(msPrint(italic)) }

// Underline prints
//	.UL
func Underline(w io.Writer) { w.Write(msPrint(underline)) }

// CodeBlock prints codeblock
//	.P1
//	<codeblock>
//	.P2
func CodeBlock(w io.Writer, code []byte) {
	w.Write(msPrint(beginCodeBlock))
	io.WriteString(w, strings.TrimLeft(string(code), " "))
	w.Write(msPrint(endCodeBlock))
}

// HTML prints .HTML <attribute>\n
func HTML(w io.Writer, attribute string) { w.Write(msPrint(html, attribute)) }

// Title prints
//	.TL
//	<title>\n
func Title(w io.Writer, tl string) { w.Write(msPrint(title, linebreak, tl)) }

// Abstract prints
//	.AB
//	<abstract>
//	.AE
func Abstract(w io.Writer, abstract string) {
	w.Write(msPrint(beginAbstract, linebreak, strings.TrimLeft(abstract, " ")))
	w.Write(msPrint(endAbstract))
}

// AuthorBio prints
//	.AU
//	.I <a.Name>
//	.I <a.Email>
//	.AI <a.Affiliation>
func AuthorBio(w io.Writer, a Author) {
	w.Write(msPrint(author))
	w.Write(msPrint(italic, a.Name))
	w.Write(msPrint(italic, a.Email))
	if len(a.Affiliation) > 1 {
		w.Write(msPrint(institution, a.Affiliation))
	}
}

// ChangeDate changes the date to the given date by printing
// .ND  May  8,  1945
func ChangeDate(w io.Writer, d time.Time) {
	format := `Jan 2, 2006`
	w.Write(msPrint(changeDate, d.Format(format)))
}

func printTitleBlock(w io.Writer, bl *titleBlock) {
	HTML(w, bl.Title)
	Title(w, bl.Title)
	// Authors
	for _, author := range bl.Authors {
		AuthorBio(w, author)
	}
	ChangeDate(w, time.Time(bl.Date.Time))
	if len(bl.Abstract) > 0 {
		Abstract(w, bl.Abstract)
	}
	w.Write(msPrint(beginParagraph))
}

// LineBreak prints a new line
//	\n
func LineBreak(w io.Writer) { w.Write([]byte(linebreak)) }
