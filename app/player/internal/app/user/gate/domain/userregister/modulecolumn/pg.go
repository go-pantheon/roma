package modulecolumn

import (
	basicobj "github.com/go-pantheon/roma/app/player/internal/app/basic/gate/domain/object"
	statusobj "github.com/go-pantheon/roma/app/player/internal/app/status/gate/domain/object"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

const (
	jsonbType = "JSONB"
)

// ColumnTypeMap defines specific data types for certain user modules.
// Modules not listed here will default to BYTEA (for raw bytes).
var ColumnTypeMap = map[life.ModuleKey]string{
	basicobj.ModuleKey:  jsonbType,
	statusobj.ModuleKey: jsonbType,
}
