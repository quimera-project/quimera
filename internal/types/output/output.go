package output

import (
	"fmt"
	"strings"

	"github.com/ecoshub/jin"
	"github.com/quimera-project/quimera/internal/utils/render"
)

type Output []*Unit

// New returns a new Output structure.
//
// It checks the type and crafts the corresponding Unit.
func New(json []byte) (Output, error) {
	if strings.TrimSpace(string(json)) == "" {
		return Output{&Unit{Text: "", Level: ""}}, nil
	}
	jt, err := jin.GetType(json)
	if err != nil {
		return nil, fmt.Errorf("retrieving output type: %v", err)
	}
	if jt == "array" {
		var out Output
		err := jin.IterateArray(json, func(value []byte) (bool, error) {
			vt, err := jin.GetType(value)
			if err != nil {
				return false, fmt.Errorf("retrieving output type: %v", err)
			}
			if vt == "array" {
				return false, fmt.Errorf("can not have an array of arrays in output")
			}
			u, err := getUnit[vt](value)
			if err != nil {
				return false, fmt.Errorf("getting unit: %v", err)
			}
			out = append(out, u)
			return true, nil
		})
		return out, err
	}
	u, err := getUnit[jt](json)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return Output{u}, nil
}

// Render returns the rendered string of the Output.
func (out Output) Render() string {
	r, _ := render.TextTemplate(out, `{{range $index, $element := .}}{{if $index}},{{end}}{{$element.Render}}{{end}}`)
	return r
}

// Markdown returns the markdown representation of the Output.
func (out Output) Markdown() string {
	r, _ := render.TextTemplate(out, `{{range $index, $element := .}}{{if $index}},{{end}}{{$element.Markdown}}{{end}}`)
	return r
}

// Html returns the HTML representation of the Output and any error encountered.
func (out Output) Html() (any, error) {
	return render.HtmlTemplate(out, `{{range $index, $element := .}} {{if $index}},{{end}} {{$element.Html}} {{end}}`)
}
