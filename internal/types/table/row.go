package table

import (
	pretty "github.com/jedib0t/go-pretty/v6/table"
	"github.com/quimera-project/quimera/internal/types/output"
)

type Row []output.Output

// render returns a pretty.Row with the corresponding rendered fields.
func (r Row) render() (row pretty.Row) {
	for _, o := range r {
		row = append(row, o.Render())
	}
	return
}

// markdown returns a pretty.Row with the corresponding markdown representation of the Row fields.
func (r Row) markdown() (row pretty.Row) {
	for _, o := range r {
		row = append(row, o.Markdown())
	}
	return
}
