package markdown

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/storage"
	qfs "github.com/quimera-project/quimera/internal/utils/fs"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/render"
)

// obsidian generates the obsidian folder with the Quimera output, the quimera documentation and some obsidian checks.
func obsidian() {
	var (
		obsidianPath = filepath.Join(env.Config.Output, "obsidian")
		hiddenPath   = filepath.Join(obsidianPath, ".obsidian")
		checksPath   = filepath.Join(obsidianPath, "checks")
		docsPath     = filepath.Join(obsidianPath, "docs")
	)
	if err := os.MkdirAll(hiddenPath, os.ModePerm); err != nil && !os.IsExist(err) {
		live.Printer.Errorf("creating directory \"%s\": %v", obsidian, err)
	}
	if err := fs.WalkDir(qfs.Templates.FS(), "obsidian", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel("obsidian", path)
		if err != nil {
			return err
		} else if rel == "." {
			return nil
		}
		p := filepath.Join(hiddenPath, rel)
		if d.IsDir() {
			if err := os.MkdirAll(p, os.ModePerm); err != nil && !os.IsExist(err) {
				live.Printer.Errorf("creating directory \"%s\": %v", path, err)
			}
		} else {
			file, err := qfs.Templates.ReadObsidianFile(rel)
			if err != nil {
				return fmt.Errorf("reading \"%s\" embed file: %v", path, err)
			}
			if err = ioutil.WriteFile(p, file, 0644); err != nil {
				return fmt.Errorf("writing \"%s\" embed file into \"%s\": %v", path, p, err)
			}
		}
		return nil
	}); err != nil {
		live.Printer.Fatalf("copying obsidian's template: %v", err)
	}

	// Create appearance.json
	if err := ioutil.WriteFile(filepath.Join(hiddenPath, "appearance.json"), []byte(fmt.Sprintf(`{"cssTheme": "%s"}`, render.T.Obsidian)), 0644); err != nil {
		live.Printer.Fatalf("writting \"%s/appearance.json\": %v", hiddenPath, err)
	}

	// Create checks dir
	if err := os.Mkdir(checksPath, os.ModePerm); err != nil && !os.IsExist(err) {
		live.Printer.Errorf("creating directory \"%s\": %v", checksPath, err)
	}

	// Create checks
	for cat := range storage.Checks() {
		buf := &bytes.Buffer{}
		fmt.Fprintf(buf, "# %s %s checks\n", render.GetEmoji(cat), strings.Title(cat))
		for _, c := range storage.ChecksByCategory(cat) {
			check, err := c.Markdown()
			if err != nil {
				live.Printer.Errorf("%v\n", err)
			}
			fmt.Fprintln(buf, check)
		}
		if err := ioutil.WriteFile(filepath.Join(checksPath, cat+".md"), buf.Bytes(), 0644); err != nil {
			live.Printer.Fatalf("writing \"%s\": %v", filepath.Join(checksPath, cat+".md"), err)
		}
	}

	// Create docs dir
	if err := os.Mkdir(docsPath, os.ModePerm); err != nil && !os.IsExist(err) {
		live.Printer.Errorf("creating directory \"%s\": %v", docsPath, err)
	}

	// Create docs
	for cat := range storage.Checks() {
		if err := os.Mkdir(filepath.Join(docsPath, cat), os.ModePerm); err != nil && !os.IsExist(err) {
			live.Printer.Errorf("creating directory \"%s\": %v", filepath.Join(obsidianPath, "docs", cat), err)
		}
		re := regexp.MustCompile(`.+\.md$`)
		p := filepath.Join(cat, "doc", env.Config.Language)
		err := fs.WalkDir(qfs.Checks.FS(), p, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.Name() == p || d.IsDir() {
				return nil
			}
			if ok := re.MatchString(d.Name()); !ok {
				return nil
			}
			if env.Config.SkipFailed {
				var failed bool
				for _, check := range storage.Checks()[cat] {
					if check.Name+".md" == d.Name() {
						failed = true
					}
				}
				if !failed {
					return nil
				}
			}
			file, err := fs.ReadFile(qfs.Checks.FS(), path)
			if err != nil {
				return fmt.Errorf("reading \"%s\": %v", path, err)
			}
			if err := ioutil.WriteFile(filepath.Join(docsPath, cat, d.Name()), file, 0644); err != nil {
				live.Printer.Fatalf("writing \"%s\": %v", filepath.Join(docsPath, d.Name()), err)
			}
			return nil
		})
		if err != nil {
			live.Printer.Errorf("copying docs: %v", err)
		}
	}

	// Create quimera.md
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
	fmt.Fprint(buf, "# Categories\n\n")
	for cat := range storage.Checks() {
		fmt.Fprintf(buf, "![[%s.md]]\n\n", cat)
	}
	if err := ioutil.WriteFile(filepath.Join(obsidianPath, "quimera.md"), buf.Bytes(), 0664); err != nil {
		live.Printer.Fatalf("generating markdown: writting \"%s\": %v", filepath.Join(obsidianPath, "quimera.md"), err)
	}
}
