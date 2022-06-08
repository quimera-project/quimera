package output

import (
	"fmt"
	"html/template"

	"github.com/ecoshub/jin"
	"github.com/quimera-project/quimera/internal/utils/live"
)

type Unit struct {
	Text  string
	Level string
}

// New returns a new Unit struct pointer and any error encountered.
//
// The Unit can be crafted from a json containing either a number, string, null, boolean or an object.
var getUnit = map[string]func(json []byte) (*Unit, error){
	"number":  func(json []byte) (*Unit, error) { return &Unit{Text: string(json), Level: "default"}, nil },
	"string":  func(json []byte) (*Unit, error) { return &Unit{Text: string(json), Level: "default"}, nil },
	"null":    func(json []byte) (*Unit, error) { return &Unit{Text: "", Level: "default"}, nil },
	"boolean": func(json []byte) (*Unit, error) { return &Unit{Text: string(json), Level: "default"}, nil },
	"object": func(json []byte) (*Unit, error) {
		tt, err := jin.GetType(json, "text")
		if err != nil {
			return nil, fmt.Errorf("retrieving \"text\" output type: %v", err)
		}
		if tt == "null" {
			return &Unit{Text: "", Level: ""}, nil
		}
		t, err := jin.GetString(json, "text")
		if err != nil {
			return nil, fmt.Errorf("creating \"text\": %v", err)
		}
		lt, err := jin.GetType(json, "level")
		if err != nil {
			return nil, fmt.Errorf("retrieving \"level\" output type: %v", err)
		}
		if lt == "null" {
			return &Unit{Text: t, Level: ""}, nil
		}
		l, err := jin.GetString(json, "level")
		if err != nil {
			return nil, fmt.Errorf("creating \"level\": %v", err)
		}
		return &Unit{Text: t, Level: l}, nil
	},
}

// Render returns the rendered string of the Unit.
func (unit *Unit) Render() string {
	switch unit.Level {
	case "high":
		return live.Printer.Color.Red(unit.Text).String()
	case "medium":
		return live.Printer.Color.Yellow(unit.Text).String()
	case "critical":
		return live.Printer.Color.BgYellow(unit.Text).Red().String()
	case "info":
		return live.Printer.Color.Cyan(unit.Text).String()
	case "missing":
		return live.Printer.Color.Gray(12-1, unit.Text).String()
	case "default":
		return unit.Text
	default:
		return unit.Text
	}
}

// Markdown returns the markdown representation of the Unit.
func (unit *Unit) Markdown() string {
	return unit.Text
}

// Html returns the HTML representation of the Unit and any error encountered.
func (unit *Unit) Html() (any, error) {
	switch unit.Level {
	case "high":
		return template.HTML(fmt.Sprintf("<span class=\"high\">%s</span>", unit.Text)), nil
	case "medium":
		return template.HTML(fmt.Sprintf("<span class=\"medium\">%s</span>", unit.Text)), nil
	case "critical":
		return template.HTML(fmt.Sprintf("<span class=\"critical\">%s</span>", unit.Text)), nil
	case "info":
		return template.HTML(fmt.Sprintf("<span class=\"info\">%s</span>", unit.Text)), nil
	case "missing":
		return template.HTML(fmt.Sprintf("<span class=\"missing\">%s</span>", unit.Text)), nil
	case "default":
		return template.HTML(fmt.Sprintf("<span>%s</span>", unit.Text)), nil
	default:
		return template.HTML(fmt.Sprintf("<span>%s</span>", unit.Text)), nil
	}
}
