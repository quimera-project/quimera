package checkbox

import (
	"fmt"
	"html/template"

	"github.com/ecoshub/jin"
	"github.com/quimera-project/quimera/internal/types/output"
	"github.com/quimera-project/quimera/internal/types/structure"
	"github.com/quimera-project/quimera/internal/utils/render"
)

type Checkbox struct {
	Query  output.Output
	Passed bool
}

// New returns a new Checkbox structure and any error encountered.
func New(json []byte) (structure.Structure, error) {
	p, err := jin.GetBool(json, "passed")
	if err != nil {
		return nil, fmt.Errorf("creating checkbox: getting \"passed\" attribute: %v", err)
	}
	jq, err := jin.Get(json, "query")
	if err != nil {
		return nil, fmt.Errorf("creating checkbox: getting \"query\" attribute: %v", err)
	}
	q, err := output.New(jq)
	if err != nil {
		return nil, fmt.Errorf("creating checkbox: creating \"query\": %v", err)
	}
	return &Checkbox{Query: q, Passed: p}, nil
}

// Render returns the rendered string of the Checkbox and any error encountered.
func (c *Checkbox) Render() (string, error) {
	return c.render(render.T.Render)
}

// Markdown returns the markdown representation of the Checkbox and any error encountered.
func (c *Checkbox) Markdown() (string, error) {
	return c.render(render.T.Markdown)
}

// Html returns the HTML representation of the Checkbox and any error encountered.
func (c *Checkbox) Html() (any, error) {
	var (
		r string
	)
	q, err := c.Query.Html()
	if err != nil {
		return nil, err
	}
	if c.Passed {
		r, err = render.Render(render.T.HTML.Checkbox.Ok, map[string]any{"Query": q})
		if err != nil {
			return nil, err
		}
	} else {
		r, err = render.Render(render.T.HTML.Checkbox.Failed, map[string]any{"Query": q})
		if err != nil {
			return nil, err
		}
	}
	return template.HTML(r), nil
}

// Render returns the rendered string of the Checkbox based on the selected model and any error encountered.
func (c *Checkbox) render(model *render.Model) (string, error) {
	if c.Passed {
		return render.Render(model.Checkbox.Ok, map[string]any{"Query": c.Query})
	}
	return render.Render(model.Checkbox.Failed, map[string]any{"Query": c.Query})
}
