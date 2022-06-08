package terminal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/storage"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/render"
	"github.com/quimera-project/quimera/internal/utils/terminal"
)

// Output prints the Quimera output to stdout if env.Config.Silent is not true.
//
// If env.Config.Stdout is true, a file is created with the output.
func Output() {
	var (
		path = filepath.Join(env.Config.Output, "out")
		file *os.File
		err  error
	)
	if env.Config.Stdout {
		file, err = os.Create(path)
		if err != nil {
			live.Printer.Errorf("creating file \"%s\": %v", path, err)
		}
	}
	write(file, fmt.Sprintf("%s%s", title("Statistics"), storage.RenderStats()))
	for _, cat := range storage.SortedCategories() {
		write(file, title(fmt.Sprintf("%s CHECKS", cat)))
		for _, check := range storage.ChecksByCategory(cat) {
			c, err := check.Render()
			if err != nil {
				live.Printer.Errorf("%v\n", err)
			}
			write(file, fmt.Sprintf("%s\n", c))
		}
	}
}

// write writes the corresponding content to stdout if env.Config.Silent is not true and to the corresponding file if is supplied.
func write(file *os.File, content string) {
	if !env.Config.Silent {
		fmt.Print(content)
	}
	if file != nil {
		if _, err := file.WriteString(content); err != nil {
			live.Printer.Errorf("writing: %v", err)
		}
	}
}

// title generates the supplied msg as a title string.
func title(msg string) string {
	var text string
	if render.T.Render.Title.Up != "" {
		up, err := render.Render(render.T.Render.Title.Up, map[string]any{"Title": msg})
		if err != nil {
			live.Printer.Errorf("%v", err)
			return ""
		}
		text = fmt.Sprintf("%s\n", up)
	}
	if render.T.Render.Title.Left != "" {
		left, err := render.Render(render.T.Render.Title.Left, map[string]any{"Title": msg})
		if err != nil {
			live.Printer.Errorf("%v", err)
			return ""
		}
		text = fmt.Sprintf("%s%s", text, left)
	}
	text = fmt.Sprintf("%s %s", text, strings.ToUpper(msg))
	if render.T.Render.Title.Right != "" {
		right, err := render.Render(render.T.Render.Title.Right, map[string]any{"Title": msg})
		if err != nil {
			live.Printer.Errorf("%v", err)
			return ""
		}
		text = fmt.Sprintf("%s %s", text, right)
	}
	if render.T.Render.Title.Down != "" {
		down, err := render.Render(render.T.Render.Title.Down, map[string]any{"Title": msg})
		if err != nil {
			live.Printer.Errorf("%v", err)
			return ""
		}
		text = fmt.Sprintf("%s\n%s", text, down)
	}
	if render.T.Render.Title.Centered {
		text = terminal.CenterText(text)
	} else {
		text = fmt.Sprintf("%s\n", text)
	}
	return fmt.Sprintf("\n%s\n\n", text)
}
