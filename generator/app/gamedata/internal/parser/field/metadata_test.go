package field

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringToMetadataType(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		want   *NameInfo
		hasErr bool
	}{
		{
			name: "Id",
			in:   string(IdType),
			want: &NameInfo{Type: IdType, FieldName: string(IdType)},
		},
		{
			name: "SharedId",
			in:   fmt.Sprintf("%s_%s", SharedIdType, "Id"),
			want: &NameInfo{Type: SharedIdType, FieldName: "Id"},
		},
		{
			name: "SharedSubId",
			in:   fmt.Sprintf("%s_%s", SharedSubIdType, "Level"),
			want: &NameInfo{Type: SharedSubIdType, FieldName: "Level"},
		},
		{
			name: "string without sep",
			in:   "Equipment",
			want: &NameInfo{Type: NormalType, FieldName: "Equipment"},
		},
		{
			name: "string with joined",
			in:   fmt.Sprintf("%s#%s_%s", JoinedType, "ResourceItem", "Item"),
			want: &NameInfo{Type: NormalType, JoinedType: JoinedType, JoinedName: "ResourceItem", FieldName: "Item"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, err := parseName(tt.in)
			require.Equal(t, tt.hasErr, err != nil)
			if tt.hasErr {
				return
			}
			assert.Equal(t, tt.want.Type, ret.Type)
			assert.Equal(t, tt.want.JoinedType, ret.JoinedType)
			assert.Equal(t, tt.want.JoinedName, ret.JoinedName)
			assert.Equal(t, tt.want.FormulaValue, ret.FormulaValue)
			assert.Equal(t, tt.want.FieldName, ret.FieldName)
		})
	}
}
