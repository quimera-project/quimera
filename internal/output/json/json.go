package json

import (
	"path/filepath"

	"github.com/quimera-project/quimera/internal/env"
	"github.com/quimera-project/quimera/internal/storage"
	qfs "github.com/quimera-project/quimera/internal/utils/fs"
	ujson "github.com/quimera-project/quimera/internal/utils/json"
	"github.com/quimera-project/quimera/internal/utils/live"
)

func Output() {
	b, err := ujson.Marshal(storage.Checks())
	if err != nil {
		live.Printer.Fatalf("%v", err)
	}
	qfs.SaveJSON(filepath.Join(env.Config.Output, "out"), b)
}
