package cli

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/quimera-project/quimera/internal/cmd/mutate"
	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/utils/live"
)

type Mutate struct {
	Arch string `name:"arch" short:"a" default:"amd64" help:"Select mutant arch" enum:"amd64,arm64" group:"Options"`

	Output string `name:"output" short:"o" default:"./mutant" help:"Select output path" group:"Output"`
}

func (m *Mutate) AfterApply() error {
	if filepath.Dir(m.Output) == "." {
		dir, err := os.Getwd()
		if err != nil {
			live.Printer.Fatalf("getting current directory: %v", err)
		}
		m.Output = filepath.Join(dir, m.Output)
	}
	return copier.Copy(env.Config, m)
}

func generate() {
	mutate.Mutate()
}
