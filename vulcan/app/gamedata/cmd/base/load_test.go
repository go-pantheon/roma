package main

import (
	"testing"

	"github.com/go-pantheon/roma/gen/gamedata/base"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	err := base.Load("../../../../../gen/gamedata/json")
	assert.NoError(t, err)
}
