package field

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/pkg"
	"github.com/vulcan-frame/vulcan-util/camelcase"
)

const (
	CommentRowIndex = iota
	NameRowIndex
	ProtoRowIndex
	ConditionRowIndex
	TypeRowIndex
)

const (
	CommentColIndex = iota + 1
	NameColIndex
	ProtoColIndex
	ConditionColIndex
	TypeColIndex
)

const (
	sep                   = "_"
	joinedSep             = "#"
	KvMetadataLineSize    = 3
	KvMetadataColSize     = 7
	TableMetadataLineSize = 5
)

type MetadataType string

const (
	IdType     = MetadataType("Id")
	NormalType = MetadataType("")

	SharedIdType    = MetadataType("SharedId")
	SharedSubIdType = MetadataType("SharedSubId")
	SharedType      = MetadataType("Shared")

	MergedMapType        = MetadataType("MergedMap")
	MergedListType       = MetadataType("MergedList")
	MergedListNonNilType = MetadataType("MergedListNonNil")

	JoinedType     = MetadataType("Joined")
	JoinedListType = MetadataType("JoinedList")
	JoinedMapType  = MetadataType("JoinedMap")

	FormulaType = MetadataType("Formula")

	GameStructType     = MetadataType("Type")
	GameStructListType = MetadataType("TypeList")
)

// Metadata excel Field Metadata
// When updating the field, you need to update both the Metadata.clone and newMetadata methods.
type Metadata struct {
	Type MetadataType

	JoinedType         MetadataType
	JoinedName         string
	FormulaValue       string
	GameStructName     string
	GameStructListName string

	Comment   string
	Name      string
	Proto     string
	Condition string

	FieldName    string
	FieldType    reflect.Type
	FieldTypeStr string
}

func NewMetadataList(rows [][]string) (mds []*Metadata, err error) {
	if len(rows) != TableMetadataLineSize {
		err = fmt.Errorf("field rows length is not %d", TableMetadataLineSize)
		return
	}

	typeRow := rows[TypeRowIndex]
	commentRow := pkg.Align(rows[CommentRowIndex], len(typeRow))
	nameRow := pkg.Align(rows[NameRowIndex], len(typeRow))
	protoRow := pkg.Align(rows[ProtoRowIndex], len(typeRow))
	conditionRow := pkg.Align(rows[ConditionRowIndex], len(typeRow))

	mds = make([]*Metadata, 0, len(nameRow))
	names := make(map[string]*Metadata, len(nameRow))

	var (
		idFieldName    string
		subIdFieldName string
	)

	for i := range typeRow {
		typ := strings.TrimSpace(typeRow[i])
		if typ == "" {
			break
		}

		var (
			md        *Metadata
			name      = strings.TrimSpace(nameRow[i])
			protoName = strings.TrimSpace(protoRow[i])
			condition = strings.TrimSpace(conditionRow[i])
			comment   = strings.TrimSpace(commentRow[i])
		)

		md, err = newMetadata(typ, name, protoName, condition, comment)
		if err != nil {
			err = errors.WithMessagef(err, "field=%s", name)
			return
		}

		switch md.Type {
		case IdType, SharedIdType:
			if idFieldName != "" {
				err = fmt.Errorf("fieldName=%s is duplicated with %s", md.Name, idFieldName)
				return
			}
			idFieldName = md.Name
		case SharedSubIdType:
			if subIdFieldName != "" {
				err = fmt.Errorf("second primary field name=%s is duplicated with %s", md.Name, subIdFieldName)
				return
			}
			subIdFieldName = md.Name
		default:
			if _, ok := names[md.FieldName]; ok {
				err = fmt.Errorf("fieldName=%s is duplicated", md.FieldName)
				return
			}
			names[md.FieldName] = md
		}
		if md.JoinedType != "" {
			if md.JoinedName == "" {
				err = fmt.Errorf("joined name is empty")
				return
			}
		}
		mds = append(mds, md)
	}

	if idFieldName == "" {
		err = fmt.Errorf("primary field (%s or %s) is not exist", IdType, SharedIdType)
		return
	}
	return
}

func newMetadata(typ, name, protoName, condition, comment string) (*Metadata, error) {
	ft, err := toReflectType(typ)
	if err != nil {
		return nil, err
	}

	nameInfo, err := parseName(name)
	if err != nil {
		return nil, err
	}

	md := &Metadata{
		Name:               name,
		Proto:              protoName,
		Comment:            comment,
		Condition:          condition,
		FieldName:          nameInfo.FieldName,
		FieldType:          ft,
		Type:               nameInfo.Type,
		JoinedType:         nameInfo.JoinedType,
		JoinedName:         nameInfo.JoinedName,
		FormulaValue:       nameInfo.FormulaValue,
		GameStructName:     nameInfo.GameStructName,
		GameStructListName: nameInfo.GameStructListName,
	}

	return md, nil
}

func (md *Metadata) clone() *Metadata {
	return &Metadata{
		Type:               md.Type,
		Name:               md.Name,
		Proto:              md.Proto,
		Comment:            md.Comment,
		Condition:          md.Condition,
		FieldName:          md.FieldName,
		FieldType:          md.FieldType,
		FieldTypeStr:       md.FieldTypeStr,
		JoinedType:         md.JoinedType,
		JoinedName:         md.JoinedName,
		FormulaValue:       md.FormulaValue,
		GameStructName:     md.GameStructName,
		GameStructListName: md.GameStructListName,
	}
}

type NameInfo struct {
	Type               MetadataType
	JoinedType         MetadataType
	JoinedName         string
	FormulaValue       string
	GameStructName     string
	GameStructListName string
	FieldName          string
}

func parseName(name string) (info *NameInfo, err error) {
	info = &NameInfo{}

	defer func() {
		if matched, _ := regexp.MatchString("^[a-zA-Z]+$", info.JoinedName+info.FieldName); !matched {
			err = fmt.Errorf("joined name must be all alphabet. joinedName=%s", info.JoinedName)
		}
	}()

	if name == "" {
		err = fmt.Errorf("field name is empty")
		return
	}

	if name == string(IdType) {
		info.Type = IdType
		info.FieldName = camelcase.ToUpperCamel(name)
		return
	}

	parts := strings.Split(name, sep)
	if len(parts) > 3 {
		err = errors.Errorf("field name must have at most 3 parts. name=%s", name)
		return
	}

	if len(parts) == 1 {
		info.Type = NormalType
		info.FieldName = camelcase.ToUpperCamel(name)
		return
	}

	var specialPart string

	info.Type = MetadataType(parts[0])
	if len(parts) == 2 {
		info.FieldName = camelcase.ToUpperCamel(parts[1])
		switch info.Type {
		case NormalType, SharedIdType, SharedSubIdType, SharedType, MergedMapType, MergedListNonNilType, MergedListType:
			return
		default:
			info.Type = NormalType
			specialPart = parts[0]
		}
	}

	if len(parts) == 3 {
		info.FieldName = camelcase.ToUpperCamel(parts[2])
		switch info.Type {
		case SharedType, MergedMapType, MergedListNonNilType, MergedListType:
			specialPart = parts[1]
		default:
			err = errors.Errorf("field name must have 2 parts when first prefix is %s. name=%s", info.Type, name)
			return
		}
	}

	subParts := strings.Split(specialPart, joinedSep)
	if len(subParts) != 2 {
		err = errors.Errorf("joined part must have 2 parts. part=%s", specialPart)
		return
	}

	pt := MetadataType(subParts[0])
	switch pt {
	case FormulaType:
		info.FormulaValue = camelcase.ToUpperCamel(subParts[1])
		if info.FormulaValue == "" {
			info.FormulaValue = "Formula"
		}
	case GameStructType:
		info.GameStructName = camelcase.ToUpperCamel(subParts[1])
		if info.GameStructName == "" {
			err = errors.Errorf("game struct name is empty. part=%s", specialPart)
			return
		}
	case GameStructListType:
		info.GameStructListName = camelcase.ToUpperCamel(subParts[1])
		if info.GameStructListName == "" {
			err = errors.Errorf("game struct list name is empty. part=%s", specialPart)
			return
		}
	case JoinedType, JoinedListType, JoinedMapType:
		info.JoinedType = MetadataType(subParts[0])
		info.JoinedName = camelcase.ToUpperCamel(subParts[1])
		if info.JoinedName == "" {
			err = errors.Errorf("joined name is empty. part=%s", specialPart)
			return
		}
	default:
		err = errors.Errorf("special type must be sepcial or joined type. part=%s", specialPart)
		return
	}
	return
}
