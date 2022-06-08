package mutate

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"

	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/utils/live"
)

func Mutate() {
	var (
		workshop = filepath.Join(env.Config.QuimeraDir, "workshop")
	)
	live.SppinerStart()
	live.SpinnerSuffixf("Generating mutant")
	live.SpinnerMsg("Setting up workshop")
	if err := clear(workshop); err != nil {
		live.Printer.Fatalf("%v", err)
	}
	opt := cp.Options{
		Skip: func(src string) (bool, error) {
			return strings.HasSuffix(src, ".git") || src == "README.MD", nil
		},
	}
	if err := cp.Copy(env.Config.WorkshopDir, workshop, opt); err != nil {
		live.Printer.Fatalf("%v", err)
	}
	live.SpinnerDone()
	live.SpinnerMsg("Building binary")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("GOARCH=%q go build -ldflags=\"-X 'main.fileLess=true' -X 'main.version=%s'\" -o %q", env.Config.Arch, env.Config.Version, env.Config.Output))
	cmd.Dir = env.Config.QuimeraDir
	// cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		live.Printer.Fatalf("%v", err)
	}
	live.SpinnerMsg("Cleaning workshop")
	if err := clear(workshop); err != nil {
		live.Printer.Fatalf("%v", err)
	}
	live.SpinnerStopf("Mutant have been created at %s", env.Config.Output)
}

func clear(dir string) error {
	names, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entery := range names {
		os.RemoveAll(path.Join([]string{dir, entery.Name()}...))
	}
	return nil
}
