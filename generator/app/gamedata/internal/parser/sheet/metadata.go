package sheet

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-util/camelcase"
)

const (
	MetadataSheetName = "_metadata"

	GenDataPackage = "DataPackage"
	GenOrderSuffix = "_Order"
	GenNextSuffix  = "_Next"
)

type SheetType string

const (
	SheetTypeIgnore = SheetType("")
	SheetTypeTable  = SheetType("table_")
	SheetTypeKV     = SheetType("kv_")
)

type Metadata struct {
	Path  string
	Sheet string
	Type  SheetType

	Package  string
	Name     string
	FullName string
	SharedId bool
	Next     bool
	Order    int

	HasFormulaField bool
}

func newMetadata(path string, sheetName, dataName string, kvs map[string]string) (md *Metadata, err error) {
	md = &Metadata{
		Path:    relativePath(path),
		Sheet:   sheetName,
		Name:    camelcase.ToUpperCamel(dataName),
		Package: camelcase.ToUpperCamel(kvs[GenDataPackage]),
	}

	md.FullName = md.Package + md.Name

	if strings.ToLower(kvs[sheetName+GenNextSuffix]) == "true" {
		md.Next = true
	}

	if v := kvs[sheetName+GenOrderSuffix]; v != "" {
		if order, err := strconv.ParseInt(v, 10, 64); err != nil {
			return nil, errors.Errorf("%s _gen %s%s must be a number", path, sheetName, GenOrderSuffix)
		} else {
			md.Order = int(order)
		}
	}

	return md, nil
}

func relativePath(path string) string {
	parts := strings.Split(path, "exceldata/")
	if len(parts) != 2 {
		return path
	}
	return parts[1]
}
