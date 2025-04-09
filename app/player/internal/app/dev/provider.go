package dev

import (
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/cmdregistrar"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/registry"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	registry.ProviderSet,
	service.ProviderSet,
	biz.ProviderSet,
	cmdregistrar.ProviderSet,
)
