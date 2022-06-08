package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/copier"
	"github.com/quimera-project/quimera/internal/cmd/benchmark"
	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/utils/live"
)

type Benchmark struct {
	Categories []string `name:"categories" short:"C" aliases:"cat" help:"Only enum checks with the selected categories" enum:"containers,files,network,processes,software,system,users" group:"Checks"`
	OS         []string `name:"os" short:"O" help:"Only enum checks tested on the selected operative systems" enum:"${systems}" group:"Checks"`
	Tags       []string `name:"tags" short:"T" aliases:"tag" help:"Only enum checks with the selected tags" enum:"${tags}" group:"Checks"`

	Uncolorized bool `name:"uncolorized" help:"Uncolorized output" group:"Output"`
}

func (b *Benchmark) AfterApply() error {
	return copier.Copy(env.Config, b)
}

func bench() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	go func() {
		<-sigCh
		live.SpinnerChao()
		signal.Stop(sigCh)
		os.Exit(0)
	}()
	benchmark.Test()
}
