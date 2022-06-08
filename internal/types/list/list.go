package list

import (
	"fmt"

	"github.com/ecoshub/jin"
	"github.com/quimera-project/quimera/internal/types/structure"
	qjson "github.com/quimera-project/quimera/internal/utils/json"
	"github.com/quimera-project/quimera/internal/utils/render"
)

type List []structure.Structure

// New returns a new list structure and any error encountered.
func New(json []byte) (structure.Structure, error) {
	var l List
	err := jin.IterateArray(json, func(value []byte) (bool, error) {
		s, err := qjson.New(value)
		if err != nil {
			return false, fmt.Errorf("creating list: creating element: %v", err)
		}
		l = append(l, s.(structure.Structure))
		return true, nil
	})
	return l, err
}

// Render returns the rendered string of the list and any error encountered.
func (list List) Render() (string, error) {
	return render.TextTemplate(list, `{{range $index, $element := .}}{{if $index}}{{"\n"}}{{end}}{{$element.Render}}{{end}}`)
}

// Markdown returns the markdown representation of the list and any error encountered.
func (list List) Markdown() (string, error) {
	return render.TextTemplate(list, `{{range $index, $element := .}}{{if $index}}{{"\n"}}{{end}}- {{$element.Markdown}}{{end}}`)
}

// Html returns the HTML representation of the list and any error encountered.
func (list List) Html() (any, error) {
	return render.HtmlTemplate(list, `<ul>{{range .}}<li>{{.Html}}</li>{{end}}</ul>`)
}
