package dev

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds/cmdregistrar"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/registry"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/service"
)

var ProviderSet = wire.NewSet(
	registry.ProviderSet,
	service.ProviderSet,
	biz.ProviderSet,
	cmdregistrar.ProviderSet,
)
