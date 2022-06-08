package single

import (
	"fmt"

	"github.com/quimera-project/quimera/internal/types/structure"
	qjson "github.com/quimera-project/quimera/internal/utils/json"
)

// New returns a new single structure and any error encountered.
func New(json []byte) (structure.Structure, error) {
	s, err := qjson.New(json)
	if err != nil {
		return nil, fmt.Errorf("creating single: %v", err)
	}
	return s.(structure.Structure), nil
}
