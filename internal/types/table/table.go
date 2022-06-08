package table

import (
	"fmt"

	"github.com/ecoshub/jin"
	pretty "github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/types/output"
	"github.com/quimera-project/quimera/internal/types/structure"
	qjson "github.com/quimera-project/quimera/internal/utils/json"
	"github.com/quimera-project/quimera/internal/utils/render"
)

type Table struct {
	Header []string
	Body   []Row
}

var styles = map[string]pretty.Style{
	"StyleDefault":                    pretty.StyleDefault,
	"StyleBold":                       pretty.StyleBold,
	"StyleColoredBright":              pretty.StyleColoredBright,
	"StyleColoredDark":                pretty.StyleColoredDark,
	"StyleColoredBlackOnBlueWhite":    pretty.StyleColoredBlackOnBlueWhite,
	"StyleColoredBlackOnCyanWhite":    pretty.StyleColoredBlackOnCyanWhite,
	"StyleColoredBlackOnGreenWhite":   pretty.StyleColoredBlackOnGreenWhite,
	"StyleColoredBlackOnMagentaWhite": pretty.StyleColoredBlackOnMagentaWhite,
	"StyleColoredBlackOnYellowWhite":  pretty.StyleColoredBlackOnYellowWhite,
	"StyleColoredBlackOnRedWhite":     pretty.StyleColoredBlackOnRedWhite,
	"StyleColoredBlueWhiteOnBlack":    pretty.StyleColoredBlueWhiteOnBlack,
	"StyleColoredCyanWhiteOnBlack":    pretty.StyleColoredCyanWhiteOnBlack,
	"StyleColoredGreenWhiteOnBlack":   pretty.StyleColoredGreenWhiteOnBlack,
	"StyleColoredMagentaWhiteOnBlack": pretty.StyleColoredMagentaWhiteOnBlack,
	"StyleColoredRedWhiteOnBlack":     pretty.StyleColoredRedWhiteOnBlack,
	"StyleColoredYellowWhiteOnBlack":  pretty.StyleColoredYellowWhiteOnBlack,
	"StyleDouble":                     pretty.StyleDouble,
	"StyleLight":                      pretty.StyleLight,
	"StyleRounded":                    pretty.StyleRounded,
}

// New returns a new Table structure and any error encountered.
func New(json []byte) (structure.Structure, error) {
	h, err := header(json)
	if err != nil {
		return nil, fmt.Errorf("creating table: %v", err)
	}
	b, err := body(json)
	if err != nil {
		return nil, fmt.Errorf("creating table: %v", err)
	}
	return &Table{Header: h, Body: b}, nil
}

// header returns a Table header and any error encountered.
func header(json []byte) ([]string, error) {
	if err := qjson.DataTypeCheckPath(json, "header", "array"); err != nil {
		return nil, err
	}
	header, err := jin.GetStringArray(json, "header")
	if err != nil {
		return nil, fmt.Errorf("collecting header: %v", err)
	}
	return header, err
}

// body returns a Table body and any error encountered.
func body(json []byte) ([]Row, error) {
	var body []Row
	if err := qjson.DataTypeCheckPath(json, "body", "array"); err != nil {
		return nil, err
	}
	err := jin.IterateArray(json, func(value []byte) (bool, error) {
		if err := qjson.DataTypeCheck(value, "array"); err != nil {
			return false, err
		}
		var row Row
		err := jin.IterateArray(value, func(v []byte) (bool, error) {
			out, err := output.New(v)
			if err != nil {
				return false, fmt.Errorf("creating \"body\" row element: %v", err)
			}
			row = append(row, out)
			return true, nil
		})
		if err != nil {
			return false, err
		}
		body = append(body, row)
		return true, nil
	}, "body")
	return body, err
}

// Render returns the rendered string of the Table and any error encountered.
func (t *Table) Render() (string, error) {
	pt := pretty.NewWriter()
	pt.AppendHeader(t.formatHead())
	for _, row := range t.Body {
		pt.AppendRow(row.render())
	}
	if env.Config.Uncolorized {
		text.DisableColors()
	}
	style := styles[render.T.Render.Table]
	pt.SetStyle(style)
	return pt.Render(), nil
}

// Markdown returns the markdown representation of the Table and any error encountered.
func (t *Table) Markdown() (string, error) {
	pt := pretty.NewWriter()
	pt.AppendHeader(t.formatHead())
	for _, row := range t.Body {
		pt.AppendRow(row.markdown())
	}
	return pt.RenderMarkdown(), nil
}

// Html returns the HTML representation of the Table and any error encountered.
func (t *Table) Html() (any, error) {
	return render.HtmlTemplate(t, `<table><tr>{{range .Header}}<th>{{.}}</th>{{end}}</tr>{{range .Body}}<tr>{{range .}}<td>{{.Html}}</td>{{end}}</tr>{{end}}</table>`)
}

// formatHead returns a pretty.Row from the Table header.
func (t *Table) formatHead() (head pretty.Row) {
	for _, h := range t.Header {
		head = append(head, h)
	}
	return
}
