package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vulcan-frame/vulcan-game/gen/gamedata/base"
)

func TestLoad(t *testing.T) {
	jsonDir := filepath.Join("../../../../../gen/gamedata/json")

	err := base.Load(jsonDir)
	assert.Nil(t, err)
}
