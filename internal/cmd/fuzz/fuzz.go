package fuzz

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma/quick"
	"github.com/quimera-project/quimera/internal/check"
	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/translate"
	"github.com/quimera-project/quimera/internal/utils/yaml"
	"github.com/sahilm/fuzzy"
)

var options = map[string]func(checks check.Checks, results fuzzy.Matches, i int) string{
	"show": func(checks check.Checks, results fuzzy.Matches, i int) string {
		buf := &bytes.Buffer{}
		quick.Highlight(buf, checks[results[i].Index].Raw(), "bash", "terminal", "monokai")
		return buf.String()
	},
	"doc": func(checks check.Checks, results fuzzy.Matches, i int) string {
		return checks[results[i].Index].Manual()
	},
	"run": func(checks check.Checks, results fuzzy.Matches, i int) string {
		return checks[results[i].Index].Exec()
	},
}

func Fuzz(query, option string) {
	checks, err := yaml.AllChecks()
	if err != nil {
		live.Printer.Fatalf("getting all checks: %v", err)
	}
	if results := fuzzy.FindFrom(query, checks); len(results) == 0 {
		live.Printer.Fatalf("%s\n", translate.TranslateMessage("no-matches", map[string]any{"Query": query}))
	} else if len(results) == 1 {
		fmt.Printf("\n%s", options[option](checks, results, 0))
	} else {
		if env.Config.Batch {
			fmt.Printf("\n%s", options[option](checks, results, 0))
			return
		}

		fmt.Printf("\n%s\n", translate.TranslateMessage("question", nil))
		for i, r := range results {
			fmt.Printf("%d %s\n", live.Printer.Color.Cyan(i), checks[r.Index].Title)
		}
		var i int
		fmt.Printf("\n%s\n", translate.TranslateMessage("number", nil))
		fmt.Scan(&i)
		if i >= 0 && i < len(results) {
			fmt.Printf("\n%s\n", options[option](checks, results, i))
		} else {
			live.Printer.Fatalf("%s\n", translate.TranslateMessage("wrong", map[string]any{"Number": i}))
		}
	}
}
