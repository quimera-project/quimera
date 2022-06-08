package enum

import (
	"fmt"

	"github.com/quimera-project/quimera/internal/check"
	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/output/html"
	"github.com/quimera-project/quimera/internal/output/json"
	"github.com/quimera-project/quimera/internal/output/markdown"
	"github.com/quimera-project/quimera/internal/output/terminal"
	"github.com/quimera-project/quimera/internal/storage"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/yaml"
)

func Iterate() {
	live.SppinerStart()
	checks, err := yaml.AllChecksByPriority()
	if err != nil {
		live.Printer.Fatalf("getting all checks: %v", err)
	}
	for i := range checks {
		live.SpinnerSuffixf("Crafting priority %d", checks[i][0].Priority)
		n := len(checks[i])
		arr := make([]check.Crafter, n)

		for j := 0; j < n; j++ {
			arr[j] = *check.NewCrafter()
			go arr[j].Craft(checks[i][j])
		}
		for j := 0; j < n; j++ {
			for {
				tick, ok := <-arr[j].Ch
				if ok {
					if tick.Err != nil {
						storage.AddError()
						live.Printer.Errorf("on \"%s\" checks: crafting \"%s\" check: %v", tick.Check.Category, tick.Check.Name, tick.Err)
					} else {
						storage.Add(tick.Check)
					}
				} else {
					break
				}
			}
		}
	}
	live.SpinnerStopf("All checks have been crafted")
	storage.Finish()
	output()
}

func output() {
	terminal.Output()
	if env.Config.Markdown {
		markdown.Output()
	}
	if env.Config.Html {
		html.Output()
	}
	if env.Config.Json {
		json.Output()
	}
	fmt.Println()
}
