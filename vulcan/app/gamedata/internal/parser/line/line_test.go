package line

import (
	"testing"

	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/field"
	"github.com/stretchr/testify/assert"
)

func TestLine_EncodeToJson(t *testing.T) {
	tests := []struct {
		name   string
		in     *Line
		want   string
		hasErr bool
	}{
		{
			name: "line without subLines should work",
			in: &Line{
				IdField: &field.Field{
					Metadata: &field.Metadata{
						FieldName: "Id",
					},
					Value: int64(1),
				},
				Fields: []*field.Field{
					{
						Metadata: &field.Metadata{
							FieldName: "Name",
						},
						Value: "test",
					},
					{
						Metadata: &field.Metadata{
							FieldName: "Level",
						},
						Value: 10,
					},
					{
						Metadata: &field.Metadata{
							FieldName: "Skills",
						},
						Value: []uint64{1, 2, 3},
					},
					{
						Metadata: &field.Metadata{
							FieldName: "Items",
						},
						Value: map[uint64]string{1: "a", 2: "b", 3: "c"},
					},
				},
			},
			want: `{"Id": 1,"Name": "test","Level": 10,"Skills": [1,2,3],"Items": {"1":"a","2":"b","3":"c"}}`,
		},
		{
			name: "line with subLines should work",
			in: &Line{
				IdField: &field.Field{
					Metadata: &field.Metadata{
						FieldName: "Id",
					},
					Value: 1,
				},
				Fields: []*field.Field{
					{
						Metadata: &field.Metadata{
							FieldName: "Name",
						},
						Value: "test",
					},
				},
				SubLines: []*Line{
					{
						IdField: &field.Field{
							Metadata: &field.Metadata{
								FieldName: "Level",
							},
							Value: int64(1),
						},
						Fields: []*field.Field{
							{
								Metadata: &field.Metadata{
									FieldName: "Skills",
								},
								Value: []uint64{1, 2, 3},
							},
						},
					},
					{
						IdField: &field.Field{
							Metadata: &field.Metadata{
								FieldName: "Level",
							},
							Value: int64(2),
						},
						Fields: []*field.Field{
							{
								Metadata: &field.Metadata{
									FieldName: "Items",
								},
								Value: map[uint64]string{1: "a", 2: "b", 3: "c"},
							},
						},
					},
				},
			},
			want: `{"Id": 1,"Name": "test","SubDatas": [{"Level": 1,"Skills": [1,2,3]},{"Level": 2,"Items": {"1":"a","2":"b","3":"c"}}]}`,
		},
		{
			name: "nil field should not work",
			in: &Line{
				IdField: &field.Field{
					Metadata: &field.Metadata{
						FieldName: "Id",
					},
					Value: int64(1),
				},
			},
			want:   "",
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.EncodeToJson()
			if tt.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
