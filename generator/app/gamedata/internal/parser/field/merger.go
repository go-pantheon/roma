package field

import (
	"reflect"

	"github.com/pkg/errors"
)

func Merge(fds []*Field) (mfd *Field, err error) {
	if len(fds) == 0 {
		return nil, errors.Errorf("fields is empty")
	}

	var mg merger

	switch fds[0].Type {
	case MergedListNonNilType, MergedListType:
		mg, err = buildFieldListMerger(fds[0])
	case MergedMapType:
		mg, err = buildFieldMapMerger(fds[0])
	default:
		err = errors.Errorf("unsupported metadata type. type=%s", fds[0].Type)
		return
	}
	if err != nil {
		return
	}

	for _, f := range fds {
		if err = mg.merge(f); err != nil {
			return
		}
	}

	mfd = mg.toField()
	return
}

type merger interface {
	merge(fd *Field) error
	toField() *Field
}

var _ merger = (*sliceMerger)(nil)

type sliceMerger struct {
	*Metadata
	value reflect.Value
}

func (mg *sliceMerger) toField() *Field {
	return &Field{
		Metadata: mg.Metadata,
		Value:    mg.value.Interface(),
	}
}

func buildFieldListMerger(fd *Field) (merger, error) {
	md := fd.Metadata.clone()
	md.FieldType = reflect.SliceOf(fd.FieldType)
	mg := &sliceMerger{
		Metadata: md,
		value:    reflect.MakeSlice(md.FieldType, 0, 0),
	}
	return mg, nil
}

func (mg *sliceMerger) merge(fd *Field) error {
	if mg.FieldName != fd.FieldName {
		return errors.Errorf("metadata ServerName mismatch %v %v", mg.FieldName, fd.FieldName)
	}
	if mg.FieldType.Elem() != fd.FieldType {
		return errors.Errorf("metadata FieldType mismatch %v %v", mg.FieldType, fd.FieldType)
	}

	switch mg.Type {
	case MergedListNonNilType:
		if fd.RawValue != "" {
			mg.value = reflect.Append(mg.value, reflect.ValueOf(fd.Value))
		}
	case MergedListType:
		mg.value = reflect.Append(mg.value, reflect.ValueOf(fd.Value))
	default:
		return errors.Errorf("unsupported metadata type. type=%s", mg.Type)
	}
	return nil
}

var _ merger = (*mapMerger)(nil)

type mapMerger struct {
	*Metadata
	value reflect.Value
}

func (mg *mapMerger) toField() *Field {
	return &Field{
		Metadata: mg.Metadata,
		Value:    mg.value.Interface(),
	}
}

func buildFieldMapMerger(fd *Field) (merger, error) {
	md := fd.Metadata.clone()
	mg := &mapMerger{
		Metadata: md,
		value:    reflect.MakeMap(md.FieldType),
	}
	return mg, nil
}

func (mg *mapMerger) merge(fd *Field) error {
	if mg.FieldName != fd.FieldName {
		return errors.Errorf("metadata ServerName mismatch %v %v", mg.FieldName, fd.FieldName)
	}
	if mg.FieldType != fd.FieldType {
		return errors.Errorf("metadata FieldType mismatch %v %v", mg.FieldType, fd.FieldType)
	}

	for _, k := range reflect.ValueOf(fd.Value).MapKeys() {
		if mg.value.MapIndex(k).IsValid() {
			return errors.Errorf("key %v already exists", k)
		}
		mg.value.SetMapIndex(k, reflect.ValueOf(fd.Value).MapIndex(k))
	}
	return nil
}
