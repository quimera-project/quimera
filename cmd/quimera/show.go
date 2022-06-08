package cli

import (
	"reflect"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/quimera-project/quimera/internal/cmd/fuzz"
	"github.com/quimera-project/quimera/internal/env"
)

type Show struct {
	Batch bool `name:"batch" help:"Don't ask for user input (select first coincidence)"`

	Args []string `arg:"" name:"Check title"`

	Output string `name:"output" short:"o" default:"./quimera" help:"Select output path" group:"Output"`
	Silent bool   `name:"silent" help:"Silent stdout output" group:"Output"`
	Stdout bool   `name:"stdout" help:"Generate stdout output" group:"Output"`
}

func (s *Show) AfterApply() error {
	return copier.Copy(env.Config, s)
}

func show(mutant bool) {
	var args []string
	if mutant {
		args = reflect.ValueOf(CLI).Interface().(*Mutant).Show.Args
	} else {
		args = reflect.ValueOf(CLI).Interface().(*Quimera).Show.Args
	}
	fuzz.Fuzz(strings.Join(args, " "), "show")
}
