package question

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ecoshub/jin"
	"github.com/quimera-project/quimera/internal/types/output"
	"github.com/quimera-project/quimera/internal/types/structure"
	"github.com/quimera-project/quimera/internal/utils/render"
)

type Question struct {
	Query  string
	Answer output.Output
}

const (
	max_spaces = 40
)

// New returns a new Question structure and any error encountered.
func New(json []byte) (structure.Structure, error) {
	q, err := jin.GetString(json, "query")
	if err != nil {
		return nil, fmt.Errorf("creating question: getting \"query\" attribute: %v", err)
	}
	ja, err := jin.Get(json, "answer")
	if err != nil {
		return nil, fmt.Errorf("creating question: getting \"answer\" attribute: %v", err)
	}
	a, err := output.New(ja)
	if err != nil {
		return nil, fmt.Errorf("creating question: creating \"answer\" attribute: %v", err)
	}
	return &Question{Query: q, Answer: a}, nil
}

// Render returns the rendered string of the Question and any error encountered.
func (q *Question) Render() (string, error) {
	return q.render(render.T.Render)
}

// Markdown returns the markdown representation of the Question and any error encountered.
func (q *Question) Markdown() (string, error) {
	return q.render(render.T.Markdown)
}

// Html returns the HTML representation of the Question and any error encountered.
func (q *Question) Html() (any, error) {
	h, err := q.Answer.Html()
	if err != nil {
		return nil, err
	}
	return template.HTML(fmt.Sprintf("<p>%s%s%s</p>", q.Query, strings.Repeat(".", max_spaces-len(q.Query)), h)), nil
}

// Render returns the rendered string of the Question based on the selected model and any error encountered.
func (q *Question) render(theme *render.Model) (string, error) {
	return render.Render(theme.Question, map[string]any{"Question": q})
}

// Union returns the union string which separates the question query from the answer.
func (question *Question) Union(characther string) string {
	return strings.Repeat(characther, max_spaces-len(question.Query))
}
