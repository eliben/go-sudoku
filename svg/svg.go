package svg

import (
	"encoding/xml"
	"fmt"
	"io"
	"text/template"
)

type Canvas struct {
	writer io.Writer
}

var headerText = `<?xml version="1.0"?>
<svg width="{{.Width}}" height="{{.Height}}"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">
`

var headerTemplate = template.Must(template.New("header").Parse(headerText))

func New(writer io.Writer, width, height int) *Canvas {
	c := &Canvas{writer: writer}

	headerTemplate.Execute(writer, struct {
		Width, Height int
	}{Width: width, Height: height})

	return c
}

func (c *Canvas) End() {
	fmt.Fprintf(c.writer, "</svg>\n")
}

func (c *Canvas) Rect(x, y, width, height int, style string) {
	fmt.Fprintf(c.writer, `<rect x="%v" y="%v" width="%v" height="%v"`, x, y, width, height)
	if len(style) > 0 {
		fmt.Fprintf(c.writer, ` style="%s"`, style)
	}
	fmt.Fprintf(c.writer, "/>\n")
}

func (c *Canvas) Text(x, y int, text string, style string) {
	fmt.Fprintf(c.writer, `<text x="%v" y="%v"`, x, y)
	if len(style) > 0 {
		fmt.Fprintf(c.writer, ` style="%s"`, style)
	}
	fmt.Fprintf(c.writer, ">")
	xml.Escape(c.writer, []byte(text))
	fmt.Fprintf(c.writer, "</text>\n")
}
