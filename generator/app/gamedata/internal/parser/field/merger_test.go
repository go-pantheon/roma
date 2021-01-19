package field

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name   string
		in     []*Field
		want   *Field
		hasErr bool
	}{
		{
			name: "MergedArray int value should work",
			in: []*Field{
				{
					Metadata: &Metadata{
						FieldName: "Skills",
						FieldType: reflect.TypeOf(uint64(1)),
						Type:      MergedListType,
					},
					Value:    uint64(1),
					RawValue: "1",
				},
				{
					Metadata: &Metadata{
						FieldName: "Skills",
						FieldType: reflect.TypeOf(uint64(1)),
						Type:      MergedListType,
					},
					Value:    uint64(0),
					RawValue: "",
				},
				{
					Metadata: &Metadata{
						FieldName: "Skills",
						FieldType: reflect.TypeOf(uint64(1)),
						Type:      MergedListType,
					},
					Value:    uint64(2),
					RawValue: "2",
				},
			},
			want: &Field{
				Metadata: &Metadata{
					FieldName: "Skills",
					FieldType: reflect.TypeOf([]uint64{}),
					Type:      MergedListType,
				},
				Value: []uint64{1, 0, 2},
			},
		},
		{
			name: "MergedArray int value should work",
			in: []*Field{
				{
					Metadata: &Metadata{
						FieldName: "Skills",
						FieldType: reflect.TypeOf(uint64(1)),
						Type:      MergedListNonNilType,
					},
					Value:    uint64(1),
					RawValue: "1",
				},
				{
					Metadata: &Metadata{
						FieldName: "Skills",
						FieldType: reflect.TypeOf(uint64(1)),
						Type:      MergedListNonNilType,
					},
					Value:    uint64(0),
					RawValue: "",
				},
				{
					Metadata: &Metadata{
						FieldName: "Skills",
						FieldType: reflect.TypeOf(uint64(1)),
						Type:      MergedListNonNilType,
					},
					Value:    uint64(2),
					RawValue: "2",
				},
			},
			want: &Field{
				Metadata: &Metadata{
					FieldName: "Skills",
					FieldType: reflect.TypeOf([]uint64{}),
					Type:      MergedListNonNilType,
				},
				Value: []uint64{1, 2},
			},
		},
		{
			name: "MergedMap int value should work",
			in: []*Field{
				{
					Metadata: &Metadata{
						FieldName: "Items",
						FieldType: reflect.TypeOf(map[uint64]int64{}),
						Type:      MergedMapType,
					},
					Value: map[uint64]int64{1: 10, 2: 20},
				},
				{
					Metadata: &Metadata{
						FieldName: "Items",
						FieldType: reflect.TypeOf(map[uint64]int64{}),
						Type:      MergedMapType,
					},
					Value: map[uint64]int64{3: 30},
				},
				{
					Metadata: &Metadata{
						FieldName: "Items",
						FieldType: reflect.TypeOf(map[uint64]int64{}),
						Type:      MergedMapType,
					},
					Value: map[uint64]int64{},
				},
				{
					Metadata: &Metadata{
						FieldName: "Items",
						FieldType: reflect.TypeOf(map[uint64]int64{}),
						Type:      MergedMapType,
					},
					Value: map[uint64]int64{4: 100, 5: 200},
				},
			},
			want: &Field{
				Metadata: &Metadata{
					FieldName: "Items",
					FieldType: reflect.TypeOf(map[uint64]int64{}),
					Type:      MergedMapType,
				},
				Value: map[uint64]int64{1: 10, 2: 20, 3: 30, 4: 100, 5: 200},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Merge(tt.in)
			require.Equal(t, tt.hasErr, err != nil)
			assert.Equal(t, tt.want.FieldName, got.FieldName)
			assert.Equal(t, tt.want.Value, got.Value)
			assert.Equal(t, tt.want.Type, got.Type)
			assert.Equal(t, tt.want.FieldType, got.FieldType)
		})
	}

}
