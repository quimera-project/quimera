package pair

import (
	"fmt"
	"html/template"

	"github.com/ecoshub/jin"
	"github.com/quimera-project/quimera/internal/types/output"
	"github.com/quimera-project/quimera/internal/types/structure"
	"github.com/quimera-project/quimera/internal/utils/render"
)

type Pair struct {
	Key, Value output.Output
}

// New returns a new Pair structure and any error encountered.
func New(json []byte) (structure.Structure, error) {
	jk, err := jin.Get(json, "key")
	if err != nil {
		return nil, fmt.Errorf("creating pair: getting \"key\" attribute: %v", err)
	}
	k, err := output.New(jk)
	if err != nil {
		return nil, fmt.Errorf("creating pair: creating \"key\" attribute: %v", err)
	}
	jv, err := jin.Get(json, "value")
	if err != nil {
		return nil, fmt.Errorf("creating pair: getting \"value\" attribute: %v", err)
	}
	v, err := output.New(jv)
	if err != nil {
		return nil, fmt.Errorf("creating pair: creating \"value\" attribute: %v", err)
	}
	return &Pair{Key: k, Value: v}, nil
}

// Render returns the rendered string of the Pair and any error encountered.
func (p *Pair) Render() (string, error) {
	return p.render(render.T.Render)
}

// Markdown returns the markdown representation of the Pair and any error encountered.
func (p *Pair) Markdown() (string, error) {
	return p.render(render.T.Markdown)
}

// Html returns the HTML representation of the Pair and any error encountered.
func (p *Pair) Html() (any, error) {
	key, err := p.Key.Html()
	if err != nil {
		return nil, err
	}
	value, err := p.Value.Html()
	if err != nil {
		return nil, err
	}
	r, err := render.Render(render.T.HTML.Pair, map[string]any{"Key": key, "Value": value})
	if err != nil {
		return nil, err
	}
	return template.HTML(r), nil
}

// Render returns the rendered string of the Pair based on the selected model and any error encountered.
func (p *Pair) render(theme *render.Model) (string, error) {
	return render.Render(theme.Pair, map[string]any{"Pair": p})
}
