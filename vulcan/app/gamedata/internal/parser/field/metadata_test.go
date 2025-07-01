package field

import (
	"fmt"
	"testing"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/stretchr/testify/assert"
)

func TestStringToMetadataType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     string
		want   *NameInfo
		hasErr bool
	}{
		{
			name: "Id",
			in:   string(IdType),
			want: &NameInfo{Type: IdType, FieldName: camelcase.ToUpperCamel(string(IdType))},
		},
		{
			name: "SharedId",
			in:   fmt.Sprintf("%s_%s", SharedIdType, "Id"),
			want: &NameInfo{Type: SharedIdType, FieldName: "ID"},
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
		{
			name: "string with 2 sep and joinedList",
			in:   fmt.Sprintf("%s_%s#%s_%s", SharedType, JoinedListType, "HeroSkill", "Skill"),
			want: &NameInfo{Type: SharedType, JoinedType: JoinedListType, JoinedName: "HeroSkill", FieldName: "Skill"},
		},
		{
			name:   "string with 3 sep is invalid",
			in:     fmt.Sprintf("%s_%s_%s", SharedType, MergedMapType, "Pack"),
			want:   &NameInfo{},
			hasErr: true,
		},
		{
			name:   "string with invalid joined 1",
			in:     fmt.Sprintf("%s#%s_%s_%s", JoinedType, "ResourcePack", MergedMapType, "Pack"),
			want:   &NameInfo{},
			hasErr: true,
		},
		{
			name:   "string with invalid joined 2",
			in:     fmt.Sprintf("%s#%s", JoinedType, "ResourcePack"),
			want:   &NameInfo{},
			hasErr: true,
		},
		{
			name:   "string with invalid second prefix",
			in:     fmt.Sprintf("%s_%s_%s", SharedType, SharedType, "Hero"),
			want:   &NameInfo{},
			hasErr: true,
		},
		{
			name:   "string with invalid second prefix",
			in:     fmt.Sprintf("%s_%s_%s", SharedIdType, JoinedType, "Hero"),
			want:   &NameInfo{},
			hasErr: true,
		},
		{
			name: "string with formula prefix",
			in:   fmt.Sprintf("%s#%s_%s", FormulaType, "Attribute", "Formula"),
			want: &NameInfo{Type: NormalType, FormulaValue: "Attribute", FieldName: "Formula"},
		},
		{
			name:   "string with invalid formula prefix",
			in:     fmt.Sprintf("%s#%s_%s", FormulaType, "", "Formula"),
			want:   &NameInfo{Type: NormalType, FormulaValue: "Formula", FieldName: "Formula"},
			hasErr: false,
		},
		{
			name:   "string with empty joined name",
			in:     fmt.Sprintf("%s#%s_%s", JoinedType, "", "HeroLevel"),
			want:   &NameInfo{},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ret, err := parseName(tt.in)

			if tt.hasErr {
				assert.Error(t, err)
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
