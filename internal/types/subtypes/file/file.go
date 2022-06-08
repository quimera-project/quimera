package file

import (
	"fmt"
	"html/template"

	format "github.com/MichaelMure/go-term-text"
	"github.com/ecoshub/jin"
	"github.com/quimera-project/quimera/internal/types/structure"
	"github.com/quimera-project/quimera/internal/utils/render"
	"golang.org/x/term"
)

type File struct {
	Filename, Content string
}

// New returns a new File structure and any error encountered.
func New(json []byte) (structure.Structure, error) {
	n, err := jin.GetString(json, "filename")
	if err != nil {
		return nil, fmt.Errorf("creating file: getting \"filename\" attribute: %v", err)
	}
	c, err := jin.GetString(json, "content")
	if err != nil {
		return nil, fmt.Errorf("creating file: getting \"content\" attribute: %v", err)
	}
	return &File{Filename: n, Content: c}, nil
}

// Render returns the rendered string of the File and any error encountered.
func (f *File) Render() (string, error) {
	return f.render(render.T.Render)
}

// Markdown returns the markdown representation of the File and any error encountered.
func (f *File) Markdown() (string, error) {
	content, _ := format.Wrap(f.Content, 60,
		format.WrapIndent(""),
		format.WrapPad("> "),
	)
	return fmt.Sprintf("```%s```\n%s", f.Filename, content), nil
}

// Html returns the HTML representation of the File and any error encountered.
func (f *File) Html() (any, error) {
	return template.HTML(fmt.Sprintf("<div class=\"fieldset\"><h4>%s</h4><pre>%s</pre></div>", f.Filename, f.Content)), nil
}

// Render returns the rendered string of the File based on the selected model and any error encountered.
func (f *File) render(theme *render.Model) (string, error) {
	width, _, err := term.GetSize(0)
	if err != nil {
		return "", fmt.Errorf("rendering file: getting terminal size: %v", err)
	}
	indent, err := render.Render(render.T.Render.File.Indent, map[string]any{})
	if err != nil {
		return "", err
	}
	pad, err := render.Render(render.T.Render.File.Pad, map[string]any{})
	if err != nil {
		return "", err
	}
	content, _ := format.Wrap(f.Content, width,
		format.WrapIndent(fmt.Sprint(indent)),
		format.WrapPad(fmt.Sprint(pad)),
	)
	return render.Render(render.T.Render.File.Call, map[string]any{"Filename": f.Filename, "Content": content})
}
