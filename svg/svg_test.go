package svg

import (
	"bytes"
	"strings"
	"testing"
)

// Sanity testing.
func TestSvgWrite(t *testing.T) {
	var buf bytes.Buffer

	width := 120
	height := 233
	x := 77
	y := 88

	canvas := New(&buf, width, height)
	canvas.Rect(x, y, 100, 200, "my style")
	canvas.Text(x+10, y+1, "hello", "")
	canvas.End()

	result := buf.String()

	want := `<?xml version="1.0"?>
<svg width="120" height="233"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">
<rect x="77" y="88" width="100" height="200" style="my style"/>
<text x="87" y="89">hello</text>
</svg>`

	if strings.TrimSpace(result) != strings.TrimSpace(want) {
		t.Errorf("got:\n %s\n\nwant:\n %v\n", result, want)
	}
}
