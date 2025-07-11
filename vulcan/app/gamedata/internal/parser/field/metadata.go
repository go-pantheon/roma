package field

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/roma/vulcan/pkg/align"
	"github.com/pkg/errors"
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
		return nil, errors.Errorf("field rows length is not %d", TableMetadataLineSize)
	}

	typeRow := rows[TypeRowIndex]
	commentRow := align.Align(rows[CommentRowIndex], len(typeRow))
	nameRow := align.Align(rows[NameRowIndex], len(typeRow))
	protoRow := align.Align(rows[ProtoRowIndex], len(typeRow))
	conditionRow := align.Align(rows[ConditionRowIndex], len(typeRow))

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
			return nil, errors.WithMessagef(err, "field=%s", name)
		}

		switch md.Type {
		case IdType, SharedIdType:
			if idFieldName != "" {
				return nil, errors.Errorf("fieldName=%s is duplicated with %s", md.Name, idFieldName)
			}

			idFieldName = md.Name
		case SharedSubIdType:
			if subIdFieldName != "" {
				return nil, errors.Errorf("second primary field name=%s is duplicated with %s", md.Name, subIdFieldName)
			}

			subIdFieldName = md.Name
		default:
			if _, ok := names[md.FieldName]; ok {
				return nil, errors.Errorf("fieldName=%s is duplicated", md.FieldName)
			}

			names[md.FieldName] = md
		}

		if md.JoinedType != "" && md.JoinedName == "" {
			return nil, errors.Errorf("joined name is empty")
		}

		mds = append(mds, md)
	}

	if idFieldName == "" {
		return nil, errors.Errorf("primary field (%s or %s) is not exist", IdType, SharedIdType)
	}

	return mds, nil
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
		if info != nil {
			if matched, _ := regexp.MatchString("^[a-zA-Z]+$", info.JoinedName+info.FieldName); !matched {
				err = errors.Errorf("joined name must be all alphabet. joinedName=%s", info.JoinedName)
			}
		}
	}()

	if name == "" {
		return nil, errors.Errorf("field name is empty")
	}

	if name == string(IdType) {
		info.Type = IdType
		info.FieldName = camelcase.ToUpperCamel(name)

		return info, nil
	}

	parts := strings.Split(name, sep)
	if len(parts) > 3 {
		return nil, errors.Errorf("field name must have at most 3 parts. name=%s", name)
	}

	if len(parts) == 1 {
		info.Type = NormalType
		info.FieldName = camelcase.ToUpperCamel(name)

		return info, nil
	}

	var specialPart string

	info.Type = MetadataType(parts[0])
	if len(parts) == 2 {
		info.FieldName = camelcase.ToUpperCamel(parts[1])
		switch info.Type {
		case NormalType, SharedIdType, SharedSubIdType, SharedType, MergedMapType, MergedListNonNilType, MergedListType:
			return info, nil
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
			return nil, errors.Errorf("field name must have 2 parts when first prefix is %s. name=%s", info.Type, name)
		}
	}

	subParts := strings.Split(specialPart, joinedSep)
	if len(subParts) != 2 {
		return nil, errors.Errorf("joined part must have 2 parts. part=%s", specialPart)
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
			return nil, errors.Errorf("game struct name is empty. part=%s", specialPart)
		}
	case GameStructListType:
		info.GameStructListName = camelcase.ToUpperCamel(subParts[1])
		if info.GameStructListName == "" {
			return nil, errors.Errorf("game struct list name is empty. part=%s", specialPart)
		}
	case JoinedType, JoinedListType, JoinedMapType:
		info.JoinedType = MetadataType(subParts[0])
		info.JoinedName = camelcase.ToUpperCamel(subParts[1])

		if info.JoinedName == "" {
			return nil, errors.Errorf("joined name is empty. part=%s", specialPart)
		}
	default:
		return nil, errors.Errorf("special type must be sepcial or joined type. part=%s", specialPart)
	}

	return info, nil
}

func (md *Metadata) IsMerged() bool {
	return md.Type == MergedMapType || md.Type == MergedListType || md.Type == MergedListNonNilType
}
