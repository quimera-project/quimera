package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/leaanthony/debme"
	cmd "github.com/quimera-project/quimera/cmd/quimera"
	"github.com/quimera-project/quimera/internal/env"
	qfs "github.com/quimera-project/quimera/internal/utils/fs"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/render"
	"github.com/quimera-project/quimera/internal/utils/translate"
	"github.com/quimera-project/quimera/internal/utils/updater"
)

var (
	banner = `
  ███████    ██  ██▓  ██  ▄██  ██    █████  ██▀███   ▄▄▄      
  ▒██▓  ██▒ ██  ▓██▒▓██▒▓██▒▀█▀ ██▒▓▓█   ▀ ▓██ ▒ ██▒▒████▄    
  ▒██▒  ██░▓██  ▒██░▒██▒▓██    ▓██░▒███    ▓██ ░▄█ ▒▒██  ▀█▄  
  ░██   █ ░▓▓█  ░██░░██░▒██▓   ▒██ ▒▓█   ▄ ▒██▀▀█▄  ░██▄▄▄▄██  
  ░▒█████▄ ▒▒█████▓ ░██░▒██▒   ░██▒░▒█████▒░██▓ ▒██▒ ▓█   ▓██▒ 
  ░░ ▒▒░ ▒ ░▒▓▒ ▒ ▒ ░▓  ░ ▒░   ░  ░░░ ▒░ ░░ ▒▓ ░▒▓░ ▒▒   ▓▒█░  %s
   ░ ▒░  ░ ░░▒░ ░ ░  ▒ ░░  ░      ░ ░ ░  ░  ░▒ ░ ▒░  ▒   ▒▒ ░  %s
     ░   ░  ░░░ ░ ░  ▒ ░░      ░      ░     ░░   ░   ░   ▒  ░  
    ░       ░      ░         ░      ░  ░   ░           ░        @PwnedShell - %s  
  
	`
	version     string
	fileLess    string
	workshopDir string
	//go:embed workshop
	workshopFS embed.FS
)

func main() {
	ctx := cmd.Parse(fileLess)
	live.Init(!strings.Contains(ctx.Command(), " "))
	setConfig()
	if fileLess != "true" {
		fmt.Fprintln(os.Stderr, live.Printer.Color.Green(fmt.Sprintf(banner, "", "", env.Config.Version)))
		updater.Check()
	} else {
		fmt.Fprintln(os.Stderr, live.Printer.Color.Green(fmt.Sprintf(banner, "░█▄▒▄█░█▒█░▀█▀▒▄▀▄░█▄░█░▀█▀", "░█▒▀▒█░▀▄█░▒█▒░█▀█░█▒▀█░▒█▒", env.Config.Version)))
	}
	setWorkshop()
	boot(ctx)
	cmd.Execute(ctx)
}

func boot(ctx *kong.Context) {
	if env.Config.Stdout || env.Config.Html || env.Config.Json || env.Config.Markdown {
		path := env.Config.Output
		out, err := os.Stat(path)
		if os.IsNotExist(err) {
			if err = os.MkdirAll(path, os.ModePerm); err != nil {
				live.Printer.Errorf("creating directory \"%s\": %v", path, err)
			}
		} else {
			if !out.IsDir() {
				live.Printer.Fatalf("\"%s\" already exists", path)
			}
		}
	}
	render.Init()
	translate.Init()
}

func setConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		live.Printer.Fatalf("%v", err)
	}
	if workshopDir = os.Getenv("QUIMERA_WORKSHOP"); workshopDir == "" {
		workshopDir = filepath.Join(home, "quimera-workshop")
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(home, "go")
	}
	repo := filepath.Join(gopath, "pkg", "mod", "github.com", "quimera-project", "quimera@")
	if version == "" {
		var v bytes.Buffer
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("ls -d %s* | sort -V | tail -n 1 | cut -d '@' -f 2", repo))
		cmd.Stdout = &v
		err := cmd.Run()
		if err != nil {
			live.Printer.Errorf("Version could not be retrieved")
		}
		version = strings.TrimSpace(v.String())
	}
	env.Config.QuimeraDir = filepath.Join(gopath, "pkg", "mod", "github.com", "quimera-project", fmt.Sprintf("quimera@%s", version))
	env.Config.WorkshopDir = workshopDir
	env.Config.Version = version
}

func setWorkshop() {
	if fileLess == "true" {
		ws, err := debme.FS(workshopFS, "workshop")
		if err != nil {
			live.Printer.Fatalf("embedding files: %v", err)
		}
		da, err := ws.FS("assets")
		if err != nil {
			live.Printer.Fatalf("embedding files: %v", err)
		}
		dc, err := ws.FS("checks")
		if err != nil {
			live.Printer.Fatalf("embedding files: %v", err)
		}
		dt, err := ws.FS("templates")
		if err != nil {
			live.Printer.Fatalf("embedding files: %v", err)
		}
		dl, err := ws.FS("lang")
		if err != nil {
			live.Printer.Fatalf("embedding files: %v", err)
		}
		do, err := ws.FS(filepath.Join("tools", runtime.GOARCH))
		if err != nil {
			live.Printer.Fatalf("embedding files: %v", err)
		}
		qfs.Assets.NewFS(da)
		qfs.Checks.NewFS(dc)
		qfs.Templates.NewFS(dt)
		qfs.Tools.NewFS(do)
		qfs.Lang.NewFS(dl)
	} else {
		qfs.Assets.NewFS(os.DirFS(filepath.Join(workshopDir, "assets")))
		qfs.Checks.NewFS(os.DirFS(filepath.Join(workshopDir, "checks")))
		qfs.Templates.NewFS(os.DirFS(filepath.Join(workshopDir, "templates")))
		qfs.Tools.NewFS(os.DirFS(filepath.Join(workshopDir, "tools", runtime.GOARCH)))
		qfs.Lang.NewFS(os.DirFS(filepath.Join(workshopDir, "lang")))
	}
}
