package live

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/quimera-project/quimera/internal/env"
)

type color struct {
	aurora.Aurora
}

// newColor creates a new color and returns its pointer.
func newColor() *color {
	return &color{Aurora: aurora.NewAurora(!env.Config.Uncolorized)}
}

// Infof formats according to a format specifier and writes to stderr as an info message.
func (c *color) Infof(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", c.Cyan("⊡ Info:"), c.Cyan(fmt.Sprintf(format, a...)))
}

// Warningf formats according to a format specifier and writes to stderr as a warning message.
func (c *color) Warningf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", c.Yellow("☢ Warning:"), c.Yellow(fmt.Sprintf(format, a...)))
}

// Errorf formats according to a format specifier and writes to stderr as an error message.
func (c *color) Errorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", c.Red("✗ Error:"), c.Red(fmt.Sprintf(format, a...)))
}

// Fatalf formats according to a format specifier and writes to stderr as an error message.
//
// Quimera exits with a status code of 1.
func (c *color) Fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", c.Red("☠ Critical:"), c.Red(fmt.Sprintf(format, a...)))
	os.Exit(1)
}
