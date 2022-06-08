package raw

import (
	"fmt"
	"html/template"

	"github.com/ecoshub/jin"
	"github.com/quimera-project/quimera/internal/types/output"
	"github.com/quimera-project/quimera/internal/types/structure"
)

type Raw struct {
	output.Output
}

// New returns a new Raw structure and any error encountered.
func New(json []byte) (structure.Structure, error) {
	jc, err := jin.Get(json, "content")
	if err != nil {
		return nil, fmt.Errorf("creating raw: getting \"content\" attribute: %v", err)
	}
	c, err := output.New(jc)
	if err != nil {
		return nil, fmt.Errorf("creating raw: creating \"content\": %v", err)
	}
	return &Raw{Output: c}, nil
}

// Render returns the rendered string of the Raw and any error encountered.
func (r *Raw) Render() (string, error) {
	return r.Output.Render(), nil
}

// Markdown returns the markdown representation of the Raw and any error encountered.
func (r *Raw) Markdown() (string, error) {
	return r.Output.Markdown(), nil
}

// Html returns the HTML representation of the Raw and any error encountered.
func (r *Raw) Html() (any, error) {
	h, err := r.Output.Html()
	if err != nil {
		return nil, err
	}
	return template.HTML(fmt.Sprintf("<p>%s</p>", h)), nil
}
