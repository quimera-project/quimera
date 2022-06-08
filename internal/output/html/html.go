package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/quimera-project/quimera/internal/check"
	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/storage"
	qfs "github.com/quimera-project/quimera/internal/utils/fs"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/render"
)

// Output generates a folder with an index HTML file and the corresponding assets.
func Output() {
	var (
		buf      = &bytes.Buffer{}
		htmlPath = filepath.Join(env.Config.Output, "html")
		t        = template.New("Html")
	)
	if err := os.Mkdir(htmlPath, os.ModePerm); err != nil && !os.IsExist(err) {
		live.Printer.Errorf("creating directory \"%s\": %v", htmlPath, err)
	}

	t = t.Funcs(template.FuncMap{"getChecks": func(cat string) check.Checks { return storage.Checks()[cat] }})
	t = t.Funcs(template.FuncMap{"getTitle": func(cat string) string { return strings.Title(cat) }})
	t = t.Funcs(template.FuncMap{"getEmoji": render.GetEmoji})
	t = t.Funcs(template.FuncMap{"getIcon": render.GetIcon})
	t = t.Funcs(template.FuncMap{"getStats": storage.Stats})
	t = t.Funcs(template.FuncMap{"getCategories": storage.Categories})
	t = t.Funcs(template.FuncMap{"getTime": func() string {
		z := time.Unix(0, 0).UTC()
		return z.Add(storage.Stats().Duration).Format("15:04:05.99")
	}})
	t = t.Funcs(template.FuncMap{"getDate": func() string { return storage.Stats().Time.Local().Format(time.RFC1123) }})
	t = t.Funcs(template.FuncMap{"getSucceed": func() int {
		s := storage.Stats()
		return s.Total - s.Failed
	}})
	t = t.Funcs(template.FuncMap{"inc": func(i int) int { return i + 1 }})

	index, err := qfs.Templates.ReadHtmlFile("index.html")
	if err != nil {
		live.Printer.Fatalf("reading \"templates/html/index.html\" embed file: %v", err)
	}

	templ, err := t.Parse(string(index))
	if err != nil {
		live.Printer.Fatalf("parsing index.html: %v", err)
	}

	if err = templ.Execute(buf, storage.SortedCategories()); err != nil {
		live.Printer.Fatalf("executing template: %v", err)
	}

	if err = ioutil.WriteFile(filepath.Join(htmlPath, "index.html"), buf.Bytes(), 0644); err != nil {
		live.Printer.Fatalf("writting index.html: %v", err)
	}

	if err := fs.WalkDir(qfs.Templates.FS(), "html/assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel("html", path)
		if err != nil {
			return err
		} else if rel == "." {
			return nil
		}
		p := filepath.Join(htmlPath, rel)
		if d.IsDir() {
			if err := os.MkdirAll(p, os.ModePerm); err != nil && !os.IsExist(err) {
				live.Printer.Errorf("creating directory \"%s\": %v", path, err)
			}
		} else {
			file, err := qfs.Templates.ReadHtmlFile(rel)
			if err != nil {
				return fmt.Errorf("reading \"%s\" embed file: %v", path, err)
			}
			if err = ioutil.WriteFile(p, file, 0644); err != nil {
				return fmt.Errorf("writing \"%s\" embed file into \"%s\": %v", path, p, err)
			}
		}
		return nil
	}); err != nil {
		live.Printer.Fatalf("copying html's template: %v", err)
	}

	if err = ioutil.WriteFile(filepath.Join(htmlPath, "assets/css/theme.css"), []byte(render.T.Css), 0644); err != nil {
		live.Printer.Fatalf("writting assets/css/theme.css: %v", err)
	}
}
