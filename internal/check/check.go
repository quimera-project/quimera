package check

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/types"
	"github.com/quimera-project/quimera/internal/types/structure"
	qfs "github.com/quimera-project/quimera/internal/utils/fs"
	"github.com/quimera-project/quimera/internal/utils/json"
	qjson "github.com/quimera-project/quimera/internal/utils/json"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/render"
	"github.com/quimera-project/quimera/internal/utils/terminal"
	"github.com/quimera-project/quimera/internal/utils/translate"
)

type Check struct {
	Args         string   `yaml:"args" json:"-"`
	Assets       []string `yaml:"assets" json:"-"`
	Bin          string   `yaml:"bin" json:"-"`
	Category     string
	Duration     time.Duration
	Failed       bool
	Info         string
	MarkdownLink string       `json:"-"`
	Memfds       []*qfs.Memfd `json:"-"`
	Name         string
	Os           []string `yaml:"os"`
	Priority     int      `yaml:"priority"`
	Structure    structure.Structure
	Tags         []string `yaml:"tags"`
	Title        string
	Tools        []string `yaml:"tools" json:"-"`
	Type         string   `yaml:"type"`
}

// Render returns the rendered string of the check and any error encountered.
func (c *Check) Render() (string, error) {
	buf, err := c.render(render.T.Render)
	if err != nil {
		return "", fmt.Errorf("rendering %s: %v", c.Title, err)
	}
	if !c.Failed {
		r, err := c.Structure.Render()
		if err != nil {
			return "", fmt.Errorf("rendering %s: %v", c.Title, err)
		}
		fmt.Fprintln(buf, r)
	}
	return buf.String(), nil
}

// Markdown returns the markdown representation of the check and any error encountered.
func (c *Check) Markdown() (string, error) {
	buf, err := c.render(render.T.Markdown)
	if err != nil {
		return "", fmt.Errorf("rendering %s: %v", c.Title, err)
	}
	if !c.Failed {
		r, err := c.Structure.Markdown()
		if err != nil {
			return "", fmt.Errorf("rendering %s: %v", c.Title, err)
		}
		fmt.Fprintln(buf, strings.Replace(r, `\n`, "\n", -1))
	}
	return buf.String(), nil
}

// Raw returns the associated check script containing the execution code.
//
// If the file is an elf, it panics with a message.
//
// It can also save the content to the env.Config.Output path if env.Config.Stdout is specified by the user.
//
// If env.Config.Silent is specified, there will be no stdout output.
func (c *Check) Raw() string {
	bin, err := qfs.Checks.ReadBinaryByCategory(c.Category, c.Bin)
	if err != nil {
		live.Printer.Fatalf("reading check binary: %v", err)
	}
	if env.Config.Stdout {
		qfs.SaveFile(filepath.Join(env.Config.Output, c.Name), bin)
	}
	if !env.Config.Silent {
		if qfs.Elf(bin) {
			live.Printer.Fatalf(translate.TranslateMessage("is-binary", map[string]any{"Title": c.Title}))
		}
		return fmt.Sprintf("%s\n%s\n%s\n", c.Title, strings.Repeat("─", terminal.GetWidth()-5), bin)
	}
	return ""
}

// Manual returns the rendered markdown doc file associated to the check.
//
// If env.Config.Raw is specified by the user, the returned markdown will be in raw format.
//
// It can also save the content to the env.Config.Output path if env.Config.Stdout is specified by the user.
//
// If env.Config.Silent is specified by the user, there will be no stdout output.
func (c *Check) Manual() string {
	doc, err := qfs.Checks.ReadDocByCategory(c.Category, c.Name)
	if err != nil {
		live.Printer.Fatalf("reading check doc: %v", err)
	}
	var check string
	if env.Config.Raw {
		if !env.Config.Silent {
			check = fmt.Sprintf("%s\n", doc)
		}
		if env.Config.Stdout {
			qfs.SaveMarkdown(filepath.Join(env.Config.Output, c.Name), doc)
		}
	} else {
		md := fmt.Sprintf("%s\n", markdown.Render(string(doc), terminal.GetWidth()-5, 0))
		if !env.Config.Silent {
			check = md
		}
		if env.Config.Stdout {
			qfs.SaveFile(filepath.Join(env.Config.Output, c.Name), []byte(md))
		}
	}
	return check
}

// Exec crafts the check and returns the corresponding rendered string.
//
// If env.Config.Silent is specified by the user, there will be no stdout output.
//
// If env.Config.Json is specified by the user, it will save the content to the env.Config.Output in json format.
//
// If env.Config.Markdown is specified by the user, it will save the content to the env.Config.Output in markdown format.
//
// It can also save the content to the env.Config.Output path if env.Config.Stdout is specified by the user.
func (c *Check) Exec() string {
	if err := c.Craft(); err != nil {
		live.Printer.Fatalf("executing \"%s\" check: %v", c.Title, err)
	}
	var check string
	if !env.Config.Silent {
		r, err := c.Render()
		if err != nil {
			live.Printer.Fatalf("rendering check: %v", err)
		}
		check = fmt.Sprintf("%s\n", r)
	}
	if env.Config.Json {
		b, err := qjson.Marshal(c)
		if err != nil {
			live.Printer.Fatalf("%v\n", err)
		}
		qfs.SaveJSON(filepath.Join(env.Config.Output, c.Name), b)
	}
	if env.Config.Markdown {
		b, err := c.Markdown()
		if err != nil {
			live.Printer.Errorf("%v\n", err)
		}
		qfs.SaveMarkdown(filepath.Join(env.Config.Output, c.Name), []byte(b))
	}
	if env.Config.Stdout {
		b, err := c.Render()
		if err != nil {
			live.Printer.Fatalf("%v\n", err)
		}
		qfs.SaveFile(filepath.Join(env.Config.Output, c.Name), []byte(b))
	}
	return check
}

// render returns the render structure of a check in a bytes.Buffer and any error encountered.
func (c *Check) render(theme *render.Model) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	if c.Failed {
		r, err := render.Render(theme.Check.Name.Failed, map[string]any{"Check": c})
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(buf, r)
		return buf, nil
	}
	r, err := render.Render(theme.Check.Name.Ok, map[string]any{"Check": c})
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(buf, r)
	if len(c.Info) > 0 {
		r, err := render.Render(theme.Check.Info, map[string]any{"Check": c})
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(buf, r)
	}
	return buf, nil
}

// Prepare prepares the check by assigning the supplied category and name and translating the corresponding title and info.
//
// It returns an error if something goes wrong translating the title and info variables. Otherwise it returns nil.
func (c *Check) Prepare(cat, name string) error {
	c.Category = cat
	c.Name = name
	if env.Config.Obsidian {
		c.MarkdownLink = fmt.Sprintf("docs/%s/%s.md", c.Category, c.Name)
	} else {
		c.MarkdownLink = "#"
	}

	title, err := translate.TranslateChecksTitle(c.Category, c.Name)
	if err != nil {
		return fmt.Errorf("assigning the title: %v", err)
	}
	c.Title = title

	info, err := translate.TranslateChecksInfo(c.Category, c.Name)
	if err != nil {
		return fmt.Errorf("assigning the info: %v", err)
	}
	c.Info = info
	return nil
}

// Craft crafts the check structure.
//
// It also modify the Failed value if the executed binary returns 1 as exit code.
//
// It returns an error if encountered.
func (c *Check) Craft() error {
	bin, err := qfs.Checks.ReadBinaryByCategory(c.Category, c.Bin)
	if err != nil {
		return fmt.Errorf("reading check binary: %v", err)
	}

	mfd, err := qfs.Anonymous(bin)
	if err != nil {
		return fmt.Errorf("creating anonymous file: %v", err)
	}
	defer mfd.Close()

	out := &bytes.Buffer{}
	cmd := exec.Command(fmt.Sprintf("/proc/%d/fd/%s", os.Getpid(), fmt.Sprint(mfd.Fd())), c.Args)
	cmd.Stdout = out
	cmd_env, err := c.env()
	if err != nil {
		return fmt.Errorf("setting environment: %v", err)
	}
	cmd.Env = cmd_env
	start := time.Now()
	defer func() { c.Duration = time.Since(start) }()
	if err := cmd.Run(); err != nil {
		if strings.Contains(err.Error(), "exit status 1") {
			c.Failed = true
			return nil
		} else {
			return fmt.Errorf("running check: %v", err)
		}
	}
	if !json.ValidJSON(out.Bytes()) {
		if env.Config.Debug {
			fmt.Println(out.String())
		}
		return fmt.Errorf("invalid json")
	}
	s, err := types.Types[c.Type](json.Unescape(out.String()))
	if err != nil {
		if env.Config.Debug {
			fmt.Println(out.String())
		}
		return err
	}
	c.Structure = s
	if err = c.drop(); err != nil {
		return err
	}
	return nil
}

// env returns a slice with the key/value environmental variables needed for the check to run and any error encountered.
func (c *Check) env() ([]string, error) {
	var (
		chain []string
		add   = func(i string, file []byte) error {
			fd, err := qfs.Anonymous(file)
			if err != nil {
				return fmt.Errorf("creating anonymous file: %v", err)
			}
			memfd := qfs.NewMemfdAsset(fd, i)
			c.Memfds = append(c.Memfds, memfd)
			chain = append(chain, fmt.Sprintf("%s=/proc/$$/fd/%d", memfd.Name, memfd.Fd.Fd()))
			return nil
		}
	)
	if len(c.Assets)+len(c.Tools) == 0 {
		return []string{}, nil
	}
	for _, a := range c.Assets {
		file, err := qfs.Assets.Read(a)
		if err != nil {
			return nil, fmt.Errorf("getting asset %s: %v", a, err)
		}
		if err := add(a, file); err != nil {
			return nil, err
		}
	}
	for _, t := range c.Tools {
		file, err := qfs.Tools.Read(t)
		if err != nil {
			return nil, fmt.Errorf("getting tool %s: %v", t, err)
		}
		if err := add(t, file); err != nil {
			return nil, err
		}
	}
	imports, err := render.TextTemplate(chain, `{{range $index, $import := .}}{{if $index}}{{" && "}}{{end}}{{$import}}{{end}}`)
	if err != nil {
		return nil, err
	}
	return []string{fmt.Sprintf("IMPORTS=%s", imports)}, nil
}

// drop closes all assigned fd of the check.
//
// It returns an error if encountered.
func (c *Check) drop() error {
	for _, m := range c.Memfds {
		if err := m.Fd.Close(); err != nil {
			return fmt.Errorf("clossing assets fd: %v", err)
		}
	}
	return nil
}
