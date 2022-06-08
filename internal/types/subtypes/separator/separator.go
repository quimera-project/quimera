package separator

import (
	"html/template"

	"github.com/quimera-project/quimera/internal/types/structure"
)

type Separator struct{}

// New returns a new Separator structure and any error encountered.
func New(json []byte) (structure.Structure, error) {
	return &Separator{}, nil
}

// Render returns the rendered string of the Separator and any error encountered.
func (s *Separator) Render() (string, error) {
	return "", nil
}

// Markdown returns the markdown representation of the Separator and any error encountered.
func (s *Separator) Markdown() (string, error) {
	return "---", nil
}

// Html returns the HTML representation of the Separator and any error encountered.
func (s *Separator) Html() (any, error) {
	return template.HTML("<br>"), nil
}
