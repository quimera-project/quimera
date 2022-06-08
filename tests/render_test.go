package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/logrusorgru/aurora"
	qfs "github.com/quimera-project/quimera/internal/utils/fs"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/render"
	"github.com/stretchr/testify/assert"
)

type Render struct {
	expr, final string
	aux         map[string]any
}

var (
	Unmarshall = func() string {
		type Check struct {
			Args     string   `yaml:"args" json:"-"`
			Assets   []string `yaml:"assets" json:"-"`
			Bin      string   `yaml:"bin" json:"-"`
			Name     string
			Os       []string `yaml:"os"`
			Priority int      `yaml:"priority"`
			Tags     []string `yaml:"tags"`
			Tools    []string `yaml:"tools" json:"-"`
			Type     string   `yaml:"type"`
		}
		c := &Check{}
		json.Unmarshal([]byte(`{ "args": "g", "tools": [ "jo", "jq" ], "assets": [ "files.json" ], "bin": "sxid.sh", "os": [ "debian" ], "priority": 1, "type": "table" }`), c)
		return c.Os[0]
	}

	renders = []Render{
		{"3*3", "9", nil},
		{"7-5", "2", nil},
		{"0-1", "-1", nil},
		{"Talk('Quimer'+'a')", "Quimera", map[string]any{"Talk": fmt.Sprint}},
		{"Talk(Red('Quimer'+'a'))", "\x1b[31mQuimera\x1b[0m", map[string]any{"Talk": fmt.Sprint, "Red": aurora.Red}},
		{"Unmarshall()", "debian", map[string]any{"Unmarshall": Unmarshall}}}
)

func TestRender(t *testing.T) {
	qfs.Templates.NewFS(os.DirFS(filepath.Join("/home/quimera-project/quimera-workshop", "templates")))
	live.Init(false)
	render.Init()
	for _, r := range renders {
		o, err := render.Render(r.expr, r.aux)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		assert.Equal(t, r.final, o)
	}
}
