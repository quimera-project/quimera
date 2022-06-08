package cli

import (
	"reflect"

	"github.com/alecthomas/kong"
	"github.com/quimera-project/quimera/internal/env"
)

type Mutant struct {
	Enum      *Enum      `cmd:"" help:"Enum the system for Privilege Escalation" group:"Quimera:"`
	Benchmark *Benchmark `cmd:"" help:"Benchmark checks test" group:"Quimera:"`
	Doc       *Doc       `cmd:"" help:"Read the manual from a specific Privilege Escalation check" group:"Quimera:"`
	Show      *Show      `cmd:"" help:"Show a specific Privilege Escalation check" group:"Quimera:"`
	Run       *Run       `cmd:"" help:"Run a specific Privilege Escalation check" group:"Quimera:"`
	Language  string     `name:"lang" short:"L" default:"en" help:"Select language" enum:"en,es"`
}

type Quimera struct {
	Mutate *Mutate `cmd:"" help:"Create a mutant" group:"Laboratory:"`
	Mutant
}

var (
	CLI        any
	mutant     bool
	categories = "containers,files,network,processes,software,system,users"
	systems    = "kali,endeavouros,oracle-linux"
	tags       = "top,info,exploitable,configs"
)

func Parse(m string) *kong.Context {
	mutant = m == "true"
	if mutant {
		CLI = &Mutant{}
	} else {
		CLI = &Quimera{}
	}
	ctx := kong.Parse(CLI, kong.Vars{
		"categories": categories,
		"systems":    systems,
		"tags":       tags,
	}, kong.Name("quimera"), kong.Description("The new era of privilege escalation"), kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))
	if mutant {
		env.Config.Language = reflect.ValueOf(CLI).Interface().(*Mutant).Language
	} else {
		env.Config.Language = reflect.ValueOf(CLI).Interface().(*Quimera).Language
	}
	return ctx
}

func Execute(ctx *kong.Context) {
	switch ctx.Command() {
	case "mutate":
		generate()
	case "enum":
		enumerate()
	case "benchmark":
		bench()
	case "doc <Check title>":
		doc(mutant)
	case "show <Check title>":
		show(mutant)
	case "run <Check title>":
		run(mutant)
	default:
		panic(ctx.Command())
	}
}
