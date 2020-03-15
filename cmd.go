package troff

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

// Command are MS commands
type Command string

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

func cmdPrint(cmd Command, args ...interface{}) []byte {
	frmt := ".%s"
	elems := []interface{}{cmd}
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
	w.Write(cmdPrint(displayStart))
	w.Write(p)
	w.Write(cmdPrint(displayEnd))
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
	Esc(buf, cmdPrint(tableStart, "H"))
	Esc(buf, header)
	LineBreak(buf)
	TableHeading(buf)
	Esc(buf, body)
	Esc(buf, cmdPrint(tableEnd))
	Display(w, buf.Bytes())
}

// Esc should print escaped bytes, but it just writes it.
// TODO(sevki): fix!
func Esc(w io.Writer, p []byte) { w.Write(p) }

// TableHeading prints
//	.TH\n
func TableHeading(w io.Writer) { w.Write(cmdPrint(tableHeading)) }

// Indent prints
//	.RS\n
func Indent(w io.Writer) { w.Write(cmdPrint(beginIndent)) }

// Outdent prints
//	.RE\n
func Outdent(w io.Writer) { w.Write(cmdPrint(endIndent)) }

// IndentParagraph takes a <title> and prints
//	.IP <title>\n
func IndentParagraph(w io.Writer, title string) { w.Write(cmdPrint(beginAndIndentParagraph, title)) }

// LeftAlignedParagraph prints
//	.LP
func LeftAlignedParagraph(w io.Writer, title string, level int) {
	if level < 0 {
		w.Write(cmdPrint(leftAlignedParagraph, linebreak, title))
	} else {
		io.WriteString(w, linebreak)
		io.WriteString(w, title)
	}
}

// Bold prints
//	.B
func Bold(w io.Writer) { w.Write(cmdPrint(bold)) }

// Roman prints
//	.R
func Roman(w io.Writer) { w.Write(cmdPrint(roman)) }

// Italic prints
//	.I
func Italic(w io.Writer) { w.Write(cmdPrint(italic)) }

// Underline prints
//	.UL
func Underline(w io.Writer) { w.Write(cmdPrint(underline)) }

// CodeBlock prints codeblock
//	.P1
//	<codeblock>
//	.P2
func CodeBlock(w io.Writer, code []byte) {
	w.Write(cmdPrint(beginCodeBlock))
	io.WriteString(w, strings.TrimLeft(string(code), " "))
	w.Write(cmdPrint(endCodeBlock))
}

// HTML prints .HTML <attribute>\n
func HTML(w io.Writer, attribute string) { w.Write(cmdPrint(html, attribute)) }

// Title prints
//	.TL
//	<title>\n
func Title(w io.Writer, tl string) { w.Write(cmdPrint(title, linebreak, tl)) }

// Abstract prints
//	.AB
//	<abstract>
//	.AE
func Abstract(w io.Writer, abstract string) {
	w.Write(cmdPrint(beginAbstract, linebreak, strings.TrimLeft(abstract, " ")))
	w.Write(cmdPrint(endAbstract))
}

// AuthorBio prints
//	.AU
//	.I <a.Name>
//	.I <a.Email>
//	.AI <a.Affiliation>
func AuthorBio(w io.Writer, a Author) {
	w.Write(cmdPrint(author))
	w.Write(cmdPrint(italic, a.Name))
	w.Write(cmdPrint(italic, a.Email))
	if len(a.Affiliation) > 1 {
		w.Write(cmdPrint(institution, a.Affiliation))
	}
}

// ChangeDate changes the date to the given date by printing
// .ND  May  8,  1945
func ChangeDate(w io.Writer, d time.Time) {
	format := `Jan 2, 2006`
	w.Write(cmdPrint(changeDate, d.Format(format)))
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
	w.Write(cmdPrint(beginParagraph))
}

// LineBreak prints a new line
//	\n
func LineBreak(w io.Writer) { w.Write([]byte(linebreak)) }
