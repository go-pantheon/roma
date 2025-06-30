package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/mercury/internal/job/dev"
	"github.com/go-pantheon/roma/mercury/internal/job/storage"
	"github.com/go-pantheon/roma/mercury/internal/job/user"
	"github.com/go-pantheon/roma/mercury/internal/workshop"
)

func newNormalWorkshop(logger log.Logger) *workshop.Workshop {
	ws := workshop.NewWorkshop("normal", logger)
	ws.AddJob(
		user.NewLoginJob(),
		dev.NewDevListJob(),
		storage.NewStorageJob(),
	)

	return ws
}
