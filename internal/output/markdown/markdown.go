package markdown

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/storage"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/render"
)

// Output generates a markdown file with the Quimera output or a obsidian folder if env.Config.Obsidian is true.
func Output() {
	if env.Config.Obsidian {
		obsidian()
	} else {
		markdown()
	}
}

// markdown generates the markdown file with the Quimera output.
func markdown() {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "# %s Statistics\n", render.GetEmoji("dashboard"))
	stats := storage.Stats()
	fmt.Fprintf(buf, "- **ID**: %d\n", stats.Id)
	fmt.Fprintf(buf, "- **User**: %s\n", stats.User)
	fmt.Fprintf(buf, "- **Date**: %s\n", stats.Time.Local().Format(time.RFC1123))
	fmt.Fprintf(buf, "- **Duration**: %f\n\n", stats.Duration.Seconds())
	fmt.Fprintf(buf, "- **Total checks**: %d\n", stats.Total)
	fmt.Fprintf(buf, "- **Succeed**: %d\n", stats.Total-stats.Failed)
	fmt.Fprintf(buf, "- **Failed**: %d\n\n", stats.Failed)
	fmt.Fprintf(buf, "- **OS**: %s\n", stats.Os)
	fmt.Fprintf(buf, "- **Arch**: %s\n", stats.Arch)

	for _, cat := range storage.SortedCategories() {
		fmt.Fprintf(buf, "# %s %s checks\n", render.GetEmoji(cat), strings.Title(cat))
		for _, c := range storage.ChecksByCategory(cat) {
			check, err := c.Markdown()
			if err != nil {
				live.Printer.Errorf("%v\n", err)
			}
			fmt.Fprintf(buf, "%s\n&nbsp;\n", check)
		}
		fmt.Fprintln(buf, "&nbsp;")
	}
	if err := ioutil.WriteFile(filepath.Join(env.Config.Output, "out.md"), buf.Bytes(), 0644); err != nil {
		live.Printer.Fatalf("writing \"%s\": %v", filepath.Join(env.Config.Output, "out.md"), err)
	}
}
