package render

import (
	"fmt"
	"strings"

	"github.com/quimera-project/quimera/internal/env"
	qfs "github.com/quimera-project/quimera/internal/utils/fs"
	"github.com/quimera-project/quimera/internal/utils/live"
	"gopkg.in/yaml.v3"
)

var Variables = map[string]any{}

type Icon struct {
	Light, Dark string
}

type Theme struct {
	Statistics string                       `yaml:"statistics"`
	Emojis     map[string]string            `yaml:"emojis"`
	Css        string                       `yaml:"css"`
	Icons      map[string]map[string]string `yaml:"icons"`
	HTML       *Model                       `yaml:"html"`
	Markdown   *Model                       `yaml:"markdown"`
	Obsidian   string                       `yaml:"obsidian"`
	Render     *Model                       `yaml:"render"`
}

type Model struct {
	Title struct {
		Up       string `yaml:"up"`
		Left     string `yaml:"left"`
		Right    string `yaml:"right"`
		Down     string `yaml:"down"`
		Centered bool   `yaml:"centered"`
	} `yaml:"title"`
	Check struct {
		Info string `yaml:"info"`
		Name struct {
			Failed string `yaml:"failed"`
			Ok     string `yaml:"ok"`
		} `yaml:"title"`
	} `yaml:"check"`
	Checkbox struct {
		Failed string `yaml:"failed"`
		Ok     string `yaml:"ok"`
	} `yaml:"checkbox"`
	File struct {
		Call   string `yaml:"call"`
		Indent string `yaml:"indent"`
		Pad    string `yaml:"pad"`
	} `yaml:"file"`
	Pair     string `yaml:"pair"`
	Question string `yaml:"question"`
	Table    string `yaml:"table"`
}

var T *Theme

// Init initializes the T struct. It always starts by unmarshalling the "default" theme, then checks the user's supplied theme.
//
// It also initializes the Variables map with the corresponding colors and render functions.
func Init() {
	theme, err := qfs.Templates.ReadTemplateTheme("default")
	if err != nil {
		live.Printer.Fatalf("reading %s.yaml: %v", env.Config.Theme, err)
	}
	if err := yaml.Unmarshal(theme, &T); err != nil {
		live.Printer.Fatalf("%v\n", err)
	}
	if env.Config.Theme != "default" && env.Config.Theme != "" {
		theme, err := qfs.Templates.ReadTemplateTheme(env.Config.Theme)
		if err != nil {
			live.Printer.Fatalf("reading %s.yaml: %v", env.Config.Theme, err)
		}
		if err := yaml.Unmarshal(theme, &T); err != nil {
			live.Printer.Fatalf("%v\n", err)
		}
	}

	Variables = map[string]any{
		"Render":        fmt.Sprintf,
		"Repeat":        strings.Repeat,
		"Bold":          live.Printer.Color.Bold,
		"Faint":         live.Printer.Color.Faint,
		"Italic":        live.Printer.Color.Italic,
		"Underline":     live.Printer.Color.Underline,
		"Blink":         live.Printer.Color.Blink,
		"Reverse":       live.Printer.Color.Reverse,
		"Hidden":        live.Printer.Color.Hidden,
		"StrikeThrough": live.Printer.Color.StrikeThrough,
		"Overlined":     live.Printer.Color.Overlined,
		"Gray": func(v any) string {
			return live.Printer.Color.Gray(12-1, v).String()
		},
		"Black":         live.Printer.Color.Black,
		"Red":           live.Printer.Color.Red,
		"Green":         live.Printer.Color.Green,
		"Yellow":        live.Printer.Color.Yellow,
		"Blue":          live.Printer.Color.Blue,
		"Magenta":       live.Printer.Color.Magenta,
		"Cyan":          live.Printer.Color.Cyan,
		"White":         live.Printer.Color.White,
		"BrightBlack":   live.Printer.Color.BrightBlack,
		"BrightRed":     live.Printer.Color.BrightRed,
		"BrightGreen":   live.Printer.Color.BrightGreen,
		"BrightYellow":  live.Printer.Color.BrightYellow,
		"BrightBlue":    live.Printer.Color.BrightBlue,
		"BrightMagenta": live.Printer.Color.BrightMagenta,
		"BrightCyan":    live.Printer.Color.BrightCyan,
		"BrightWhite":   live.Printer.Color.BrightWhite,
		"BgGray": func(v any) string {
			return live.Printer.Color.Gray(20-1, v).BgGray(8 - 1).String()
		},
		"BgBlack":         live.Printer.Color.BgBlack,
		"BgRed":           live.Printer.Color.BgRed,
		"BgGreen":         live.Printer.Color.BgGreen,
		"BgYellow":        live.Printer.Color.BgYellow,
		"BgBlue":          live.Printer.Color.BgBlue,
		"BgMagenta":       live.Printer.Color.BgMagenta,
		"BgCyan":          live.Printer.Color.BgCyan,
		"BgWhite":         live.Printer.Color.BgWhite,
		"BgBrightBlack":   live.Printer.Color.BgBrightBlack,
		"BgBrightRed":     live.Printer.Color.BgBrightRed,
		"BgBrightGreen":   live.Printer.Color.BgBrightGreen,
		"BgBrightYellow":  live.Printer.Color.BgBrightYellow,
		"BgBrightBlue":    live.Printer.Color.BgBrightBlue,
		"BgBrightMagenta": live.Printer.Color.BgBrightMagenta,
		"BgBrightCyan":    live.Printer.Color.BgBrightCyan,
		"BgBrightWhite":   live.Printer.Color.BgBrightWhite,
	}
}

// GetEmoji returns the corresponding emoji from the supplied category.
func GetEmoji(cat string) string {
	emoji := T.Emojis[cat]
	if emoji == "" {
		return T.Emojis["default"]
	}
	return emoji
}

// GetIcon returns the corresponding Icon from the supplied category.
func GetIcon(cat string) Icon {
	return Icon{Light: T.Icons[cat]["light"], Dark: T.Icons[cat]["dark"]}
}
