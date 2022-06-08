package storage

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/render"
)

type statistics struct {
	Id       uint32
	User     string
	Os       string
	Arch     string
	Total    int
	Failed   int
	Errors   int
	Time     time.Time
	Duration time.Duration
}

// render renders the statistics struct.
func (s *statistics) render() string {
	var (
		buf = &bytes.Buffer{}
		w   = tabwriter.NewWriter(buf, 0, 0, 4, ' ', tabwriter.TabIndent)
	)
	fmt.Fprintf(w, "%s %d\t%s %s\t%s %s\t%s %.3f secs\t\n",
		statRender("Report ID:"), s.Id,
		statRender("User:"), s.User,
		statRender("Date:"), s.Time.Local().Format(time.RFC1123),
		statRender("Duration:"), s.Duration.Seconds())
	fmt.Fprintf(w, "%s %d\t%s %d\t%s %d\n",
		statRender("Total checks:"), s.Total,
		statRender("Succeed:"), s.Total-s.Failed,
		statRender("Failed:"), s.Failed)
	fmt.Fprintf(w, "%s %s\t%s %s\n",
		statRender("OS:"), strings.Title(s.Os),
		statRender("Arch:"), strings.Title(s.Arch))
	w.Flush()
	return buf.String()
}

// statRender returns the rendered supplied statistic.
func statRender(s string) string {
	r, err := render.Render(render.T.Statistics, map[string]any{"Stat": s})
	if err != nil {
		live.Printer.Errorf("rendering stat %s: %v", s, err)
		return ""
	}
	return r
}
