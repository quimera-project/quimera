package terminal

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// CenterText center the supplied string on the terminal.
func CenterText(s string) string {
	var buf = &bytes.Buffer{}
	partes := strings.Split(s, "\n")
	width, _, _ := term.GetSize(0)
	n := (width - utf8.RuneCountInString(partes[0])) / 2
	for _, p := range partes {
		fmt.Fprintf(buf, "%s%s\n", strings.Repeat("\u0020", int(n)), p)
	}
	return buf.String()
}

// GetWidth returns the terminals witdh.
func GetWidth() int {
	width, _, _ := term.GetSize(0)
	return width
}
