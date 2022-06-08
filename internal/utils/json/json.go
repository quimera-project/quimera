package json

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/ecoshub/jin"
	"github.com/quimera-project/quimera/internal/types/structure"
	"github.com/quimera-project/quimera/internal/types/subtypes/checkbox"
	"github.com/quimera-project/quimera/internal/types/subtypes/file"
	"github.com/quimera-project/quimera/internal/types/subtypes/pair"
	"github.com/quimera-project/quimera/internal/types/subtypes/question"
	"github.com/quimera-project/quimera/internal/types/subtypes/raw"
	"github.com/quimera-project/quimera/internal/types/subtypes/separator"
)

// subType represents a map which returns the corresponding New function of the provided type string.
var subType = map[string]func(json []byte) (structure.Structure, error){
	"checkbox":  checkbox.New,
	"pair":      pair.New,
	"question":  question.New,
	"file":      file.New,
	"raw":       raw.New,
	"separator": separator.New,
}

// New returns a new Structure and any error encountered.
func New(json []byte) (any, error) {
	t, err := jin.GetString(json, "type")
	if err != nil {
		return nil, fmt.Errorf("attribute \"type\" not found")
	}
	if new := subType[t]; new != nil {
		return new(json)
	}
	return nil, fmt.Errorf("type %s not found", t)
}

// ValidJSON returns true if the supplied json is valid. Otherwise it returns false.
func ValidJSON(in []byte) bool {
	return json.Valid(in)
}

// Unescape returns an unescaped json byte slice from a previous string.
func Unescape(json string) []byte {
	b := strings.ReplaceAll(json, `\t`, "\t")
	b = strings.ReplaceAll(b, `\n`, "\n")
	b = strings.ReplaceAll(b, `\u0026`, "&")
	b = strings.ReplaceAll(b, `\u003c`, "<")
	b = strings.ReplaceAll(b, `\u003e`, ">")
	return []byte(b)
}

// Marshal returns the JSON encoding of v and any error encountered.
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func DataTypeCheckPath(json []byte, path string, expected string) error {
	it, err := jin.GetType(json, path)
	if err != nil {
		return fmt.Errorf("retrieving %s %s type: %v", path, json, err)
	}
	if it != expected {
		return fmt.Errorf("%s data type is %s and should be %s", path, it, expected)
	}
	return nil
}

func DataTypeCheck(json []byte, expected string) error {
	it, err := jin.GetType(json)
	if err != nil {
		return fmt.Errorf("retrieving %s type: %v", json, err)
	}
	if it != expected {
		return fmt.Errorf("data type is %s and should be %s", it, expected)
	}
	return nil
}
