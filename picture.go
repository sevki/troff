package troff

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path"
	"strings"

	"sevki.org/x/debug"
)

// Position of the picture
type Position int

// PictureFlag sets the picture options like scale and outline
type PictureFlag int

const (
	// Left align the picture
	Left Position = iota
	// Center align the picture
	Center
	// Right align the picture
	Right

	// Outline the picture with a box.
	Outline PictureFlag = iota
	// Scale freely scales both picture dimensions.
	Scale
	// White out the area to be occupied by the picture.
	White
)

// Figure prints
//	.FG <link> <size>
func Figure(w io.Writer, link string, verticalSize string) {
	w.Write(msPrint(figure, strings.TrimLeft(link, " "), verticalSize))
}

// Picture is the mpicture macro
// https://9p.io/magic/man2html/6/mpictures
func Picture(w io.Writer, source string, height, width float32, position Position, offset float32, flags PictureFlag) {
	w.Write(msPrint(beginPicture, source, fmt.Sprintf("%.1fi, %.1fi", height, width)))
	w.Write(msPrint(endPicture))
}

func convertToPs(filename string, w int, h int) (string, error) {
	psname := strings.Replace(filename, path.Ext(filename), ".ps", -1)
	_, psname = path.Split(psname)

	args := []string{"convert"}

	if w > 0 && h > 0 {
		args = append(args, "-size", fmt.Sprintf("%dx%d", w, h))
	}
	args = append(args, filename, psname)
	cmd := exec.CommandContext(context.Background(), "magick", args...)
	stdOut := bytes.NewBuffer(nil)
	stdErr := bytes.NewBuffer(nil)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	fmt.Printf("running magick: %v\n", args)
	if err := cmd.Run(); err != nil {
		debug.Indent(stdErr, 1)
		fmt.Fprintf(stdErr, `magick: %v
args= %v
%s
`,
			err,
			args,
			stdErr.String(),
		)
		return "", errors.New(stdErr.String())
	}
	return psname, nil
}
func parseSize(s string) (w int, h int) {
	fmt.Sscanf(strings.TrimSpace(s), "%dx%d", &w, &h)
	return
}
