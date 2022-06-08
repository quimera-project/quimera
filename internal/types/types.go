package types

import (
	"github.com/quimera-project/quimera/internal/types/list"
	"github.com/quimera-project/quimera/internal/types/single"
	"github.com/quimera-project/quimera/internal/types/structure"
	"github.com/quimera-project/quimera/internal/types/table"
)

// Types represents a map which returns the corresponding New function of the provided type string.
var Types = map[string]func(json []byte) (structure.Structure, error){
	"single": single.New,
	"list":   list.New,
	"table":  table.New,
}
