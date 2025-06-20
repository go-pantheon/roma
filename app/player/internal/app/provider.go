package app

import (
	"github.com/go-pantheon/roma/app/player/internal/app/dev"
	"github.com/go-pantheon/roma/app/player/internal/app/gamedata"
	"github.com/go-pantheon/roma/app/player/internal/app/hero"
	"github.com/go-pantheon/roma/app/player/internal/app/plunder"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge"
	"github.com/go-pantheon/roma/app/player/internal/app/status"
	"github.com/go-pantheon/roma/app/player/internal/app/storage"
	"github.com/go-pantheon/roma/app/player/internal/app/system"
	"github.com/go-pantheon/roma/app/player/internal/app/user"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	gamedata.ProviderSet,
	user.ProviderSet,
	dev.ProviderSet,
	system.ProviderSet,
	status.ProviderSet,
	recharge.ProviderSet,
	plunder.ProviderSet,
	storage.ProviderSet,
	hero.ProviderSet,
)
