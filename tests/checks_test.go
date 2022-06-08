package tests

import (
	"testing"

	"github.com/quimera-project/quimera/internal/types/list"
	"github.com/quimera-project/quimera/internal/types/single"
	"github.com/quimera-project/quimera/internal/types/table"
	"github.com/stretchr/testify/assert"
)

type Checker struct {
	stype string
	json  string
	err   bool
}

var (
	checks = []Checker{{stype: "list", json: `[{
		"type": "pair",
		"key": "Hostname",
		"value": "Chronos"
	}, {
		"type": "raw",
		"content": ["Se han encontrado un total de", {
			"text": "45",
			"level": "critical"
		}, "contraseñas en archivos bajo la ruta /home"]
	}, {
		"type": "question",
		"query": "¿Estamos en una máquina virtual?",
		"answer": {
			"text": "Si",
			"level": "high"
		}
	}]`}, {stype: "table", json: `{
		"header": ["IP", "Alias"],
		"body": [
			[["127.0.0.1", "127.0.0.2"], "localhost"],
			["0.0.0.0", {"text": "docker", "level": "info"}],
			["23.4.95.100", {"text": "backdoor", "level": "critical"}]
		]
	}`}, {stype: "table", json: `{
		"header": ["IP", "Alias"],
		"bodya": [
			[["127.0.0.1", "127.0.0.2"], "localhost"],
			["0.0.0.0", {"text": "docker", "level": "info"}],
			["23.4.95.100", {"text": "backdoor", "level": "critical"}]
		]
	}`, err: true}, {stype: "single", json: `{
		"header": ["IP", "Alias"],
		"bodya": [
			[["127.0.0.1", "127.0.0.2"], "localhost"],
			["0.0.0.0", {"text": "docker", "level": "info"}],
			["23.4.95.100", {"text": "backdoor", "level": "critical"}]
		]
	}`, err: true}, {stype: "single", json: `{
		"type": "checkbox",
		"passed": 0,
		"query": "test"
	}`, err: true}, {stype: "single", json: `{
		"type": "checkbox",
		"passed": false,
		"query": "test"
	}`}, {stype: "single", json: `{
		"type": "raw",
		"content": ["Se han encontrado un total de", {"text": "45", "level": "critical"}, "contraseñas en archivos bajo la ruta /home"]
	}`}, {stype: "single", json: `{
		"type": "raw",
		"content": *
	}`}, {stype: "single", json: `{
		"type": "raw",
		"contento": *
	}`, err: true}, {stype: "single", json: `{
		"type": "question",
		"query": "¿Estamos en una máquina virtual?",
		"answer": {
			"text": "Si",
			"level": "high"
		}
	}`}, {stype: "single", json: `{
		"type": "question",
		"query": true
		}
	}`, err: true}}
)

func TestNew(t *testing.T) {

	for _, c := range checks {
		var err error
		switch c.stype {
		case "single":
			// Create list
			_, err = single.New([]byte(c.json))
		case "list":
			// Create list
			_, err = list.New([]byte(c.json))
		case "table":
			// Create list
			_, err = table.New([]byte(c.json))
		}

		// Check error
		if !c.err {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

}
