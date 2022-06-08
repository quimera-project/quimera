package translate

import (
	"fmt"

	"github.com/quimera-project/quimera/internal/env"
	qfs "github.com/quimera-project/quimera/internal/utils/fs"
	"github.com/quimera-project/quimera/internal/utils/live"
	"github.com/quimera-project/quimera/internal/utils/render"
	"gopkg.in/yaml.v3"
)

type checkEntry struct {
	Title, Info string
}

type translator map[string]map[string]map[string]checkEntry

type langEntry struct {
	Help string            `yaml:"help"`
	Fuzz map[string]string `yaml:"fuzz"`
}

type messages map[string]langEntry

var t translator
var m messages

func Init() {
	t = translator{"en": {}, env.Config.Language: {}}
	m = messages{}
	if err := initMessages(); err != nil {
		live.Printer.Fatalf("%v", err)
	}
}

func TranslateChecksTitle(cat, name string) (string, error) {
	if t[env.Config.Language][cat] == nil {
		if err := initTranslator(cat); err != nil {
			return "", fmt.Errorf("initialization dictionary %s: %v", cat, err)
		}
	}
	var title string
	if title = t[env.Config.Language][cat][name].Title; title == "" {
		title = t["en"][cat][name].Title
	}
	return title, nil
}

func TranslateChecksInfo(cat, name string) (string, error) {
	if t[env.Config.Language][cat] == nil {
		if err := initTranslator(cat); err != nil {
			return "", fmt.Errorf("initialization dictionary %s: %v", cat, err)
		}
	}
	var info string
	if info = t[env.Config.Language][cat][name].Info; info == "" {
		info = t["en"][cat][name].Info
	}
	return info, nil
}

func TranslateMessage(id string, vars map[string]any) string {
	if m[env.Config.Language].Fuzz[id] == "" {
		if m["en"].Fuzz[id] == "" {
			live.Printer.Warningf("no message with id \"%s\"", id)
			return ""
		}
		r, err := render.Render(m["en"].Fuzz[id], vars)
		if err != nil {
			live.Printer.Errorf("translating %s: %v", id, err)
			return ""
		}
		return r
	}
	r, err := render.Render(m[env.Config.Language].Fuzz[id], vars)
	if err != nil {
		live.Printer.Errorf("translating %s: %v", id, err)
		return ""
	}
	return r
}

func initTranslator(cat string) error {
	t["en"][cat] = map[string]checkEntry{}
	file, err := qfs.Checks.ReadCategoryLang(cat, "en")
	if err != nil {
		return fmt.Errorf("reading %s's english language file: %v", cat, err)
	}
	if err = yaml.Unmarshal(file, t["en"][cat]); err != nil {
		return fmt.Errorf("unmarshalling %s's english language file: %v", cat, err)
	}
	if env.Config.Language != "en" {
		t[env.Config.Language][cat] = map[string]checkEntry{}
		file, err := qfs.Checks.ReadCategoryLang(cat, env.Config.Language)
		if err != nil {
			return nil
		}
		yaml.Unmarshal(file, t[env.Config.Language][cat])
	}
	return nil
}

func initMessages() error {
	file, err := qfs.Lang.ReadMessages("en")
	if err != nil {
		return fmt.Errorf("reading english language file: %v", err)
	}
	entry := &langEntry{}
	if err = yaml.Unmarshal(file, entry); err != nil {
		return fmt.Errorf("unmarshalling english language file: %v", err)
	}
	m["en"] = *entry
	if env.Config.Language != "en" {
		entry := &langEntry{}
		file, err := qfs.Lang.ReadMessages(env.Config.Language)
		if err != nil {
			return nil
		}
		yaml.Unmarshal(file, entry)
		m[env.Config.Language] = *entry
	}
	return nil
}
