package updater

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/utils/live"
)

func Check() {
	if _, err := os.Stat(env.Config.WorkshopDir); os.IsNotExist(err) {
		live.Printer.Infof("%s does not exists. Downloading workshop...", env.Config.WorkshopDir)
		if err != nil {
			live.Printer.Fatalf("%v", err)
		}
		if _, err = git.PlainClone(env.Config.WorkshopDir, false, &git.CloneOptions{
			URL:      "https://github.com/quimera-project/quimera-workshop",
			Progress: nil,
		}); err != nil {
			live.Printer.Fatalf("%v", err)
		}
	}
}
