package cli

import (
	"reflect"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/quimera-project/quimera/internal/cmd/fuzz"
	"github.com/quimera-project/quimera/internal/env"
)

type Doc struct {
	Batch bool `name:"batch" help:"Don't ask for user input (select first coincidence)"`

	Args []string `arg:"" name:"Check title"`

	Theme  string `name:"theme" short:"t" help:"Theme to use" default:"default" enum:"default,demon,ascii" group:"Output"`
	Output string `name:"output" short:"o" default:"./quimera" help:"Select output path" group:"Output"`
	Silent bool   `name:"silent" help:"Silent stdout output" group:"Output"`
	Stdout bool   `name:"stdout" help:"Generate stdout output" group:"Output"`
	Raw    bool   `name:"raw" help:"Raw markdown without render" group:"Output"`
}

func (d *Doc) AfterApply() error {
	return copier.Copy(env.Config, d)
}

func doc(mutant bool) {
	var args []string
	if mutant {
		args = reflect.ValueOf(CLI).Interface().(*Mutant).Doc.Args
	} else {
		args = reflect.ValueOf(CLI).Interface().(*Quimera).Doc.Args
	}
	fuzz.Fuzz(strings.Join(args, " "), "doc")
}
