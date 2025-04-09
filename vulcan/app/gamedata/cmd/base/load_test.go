package main

import (
	"path/filepath"
	"testing"

	"github.com/go-pantheon/roma/gen/gamedata/base"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	jsonDir := filepath.Join("../../../../../gen/gamedata/json")

	err := base.Load(jsonDir)
	assert.Nil(t, err)
}
