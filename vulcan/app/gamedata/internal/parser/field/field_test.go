package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField_EncodeToJson(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     *Field
		want   string
		hasErr bool
	}{
		{
			name: "int value should work",
			in: &Field{
				Metadata: &Metadata{
					FieldName: "Id",
				},
				Value: 1,
			},
			want: `"Id": 1`,
		},
		{
			name: "string value should work",
			in: &Field{
				Metadata: &Metadata{
					FieldName: "Name",
				},
				Value: "test",
			},
			want: `"Name": "test"`,
		},
		{
			name: "slice value should work",
			in: &Field{
				Metadata: &Metadata{
					FieldName: "Skills",
				},
				Value: []uint64{1, 2, 3},
			},
			want: `"Skills": [1,2,3]`,
		},
		{
			name: "map value should work",
			in: &Field{
				Metadata: &Metadata{
					FieldName: "Items",
				},
				Value: map[uint64]string{1: "a", 2: "b", 3: "c"},
			},
			want: `"Items": {"1":"a","2":"b","3":"c"}`,
		},
		{
			name: "nil value should not work",
			in: &Field{
				Metadata: &Metadata{
					FieldName: "None",
				},
				Value: nil,
			},
			want:   "",
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.in.EncodeToJson()

			if tt.hasErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
