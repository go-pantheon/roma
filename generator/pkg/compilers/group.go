package compilers

import (
	"sort"
)

// GroupType string in under_score format
type GroupType string

// ModType string in under_score format
type ModType string

const (
	PlayerGroup = GroupType("player")
)

func GroupByMod(mod ModType) GroupType {
	for g, mods := range groupModMap {
		for m := range mods {
			if m == mod {
				return g
			}
		}
	}
	return PlayerGroup
}

func GroupMods(group GroupType) map[ModType]struct{} {
	if mods, ok := groupModMap[group]; ok {
		return mods
	}
	return nil
}

func Groups() []GroupType {
	groups := make([]GroupType, 0, len(groupModMap)+1)
	groups = append(groups, PlayerGroup)
	for g := range groupModMap {
		groups = append(groups, g)
	}
	return groups
}

var _ sort.Interface = (*ModSlice)(nil)

type ModSlice []ModType

func (s ModSlice) Len() int {
	return len(s)
}

func (s ModSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s ModSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
