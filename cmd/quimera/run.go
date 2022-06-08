package cli

import (
	"reflect"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/quimera-project/quimera/internal/cmd/fuzz"
	"github.com/quimera-project/quimera/internal/env"
)

type Run struct {
	Batch bool `name:"batch" help:"Don't ask for user input (select first coincidence)"`

	Args []string `arg:"" name:"Check title"`

	Theme       string `name:"theme" short:"t" help:"Theme to use" default:"default" enum:"default,demon,ascii" group:"Output"`
	Output      string `name:"output" short:"o" default:"./quimera" help:"Select output path" group:"Output"`
	Json        bool   `name:"json" help:"Generate JSON output" group:"Output"`
	Markdown    bool   `name:"markdown" help:"Generate Markdown output" group:"Output"`
	Silent      bool   `name:"silent" help:"Silent stdout output" group:"Output"`
	Stdout      bool   `name:"stdout" help:"Generate stdout output" group:"Output"`
	Uncolorized bool   `name:"uncolorized" help:"Uncolorized output" group:"Output"`
}

func (r *Run) AfterApply() error {
	return copier.Copy(env.Config, r)
}

func run(mutant bool) {
	var args []string
	if mutant {
		args = reflect.ValueOf(CLI).Interface().(*Mutant).Run.Args
	} else {
		args = reflect.ValueOf(CLI).Interface().(*Quimera).Run.Args
	}
	fuzz.Fuzz(strings.Join(args, " "), "run")
}
