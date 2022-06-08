package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/copier"
	"github.com/quimera-project/quimera/internal/cmd/enum"
	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/utils/live"
)

type Enum struct {
	Categories []string `name:"categories" short:"C" help:"Only enum checks with the selected categories" enum:"${categories}" group:"Checks"`
	OS         []string `name:"os" short:"O" help:"Only enum checks tested on the selected operative systems" enum:"${systems}" group:"Checks"`
	Tags       []string `name:"tags" short:"T" aliases:"tag" help:"Only enum checks with the selected tags" enum:"${tags}" group:"Checks"`
	SkipFailed bool     `name:"skip-failed" help:"Skip failed checks" group:"Checks"`

	Theme       string `name:"theme" short:"t" help:"Theme to use" default:"default" enum:"default,demon,ascii" group:"Output"`
	Output      string `name:"output" short:"o" default:"./quimera" help:"Select output path" group:"Output"`
	Html        bool   `name:"html" help:"Generate HTML output" group:"Output"`
	Json        bool   `name:"json" help:"Generate JSON output" group:"Output"`
	Markdown    bool   `name:"markdown" help:"Generate Markdown output" group:"Output"`
	Obsidian    bool   `name:"obsidian" help:"If Markdown, generate Obsidian output" group:"Output"`
	Silent      bool   `name:"silent" help:"Silent stdout output" group:"Output"`
	Stdout      bool   `name:"stdout" help:"Generate stdout output" group:"Output"`
	Uncolorized bool   `name:"uncolorized" help:"Uncolorized output" group:"Output"`

	Debug bool `name:"debug" help:"Show the returned JSON from the checks with errors"`
}

func (e *Enum) AfterApply() error {
	return copier.Copy(env.Config, e)
}

func enumerate() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	go func() {
		<-sigCh
		live.SpinnerChao()
		signal.Stop(sigCh)
		os.Exit(0)
	}()
	enum.Iterate()
}
