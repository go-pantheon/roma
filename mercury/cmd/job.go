package main

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/mercury/internal/conf"
	"github.com/go-pantheon/roma/mercury/internal/job/dev"
	"github.com/go-pantheon/roma/mercury/internal/job/storage"
	"github.com/go-pantheon/roma/mercury/internal/job/user"
	"github.com/go-pantheon/roma/mercury/internal/worker"
	"github.com/go-pantheon/roma/mercury/internal/workshop"
)

func newWorkshop(_ context.Context, logger log.Logger, bc *conf.Bootstrap) *workshop.Workshop {
	ws := workshop.NewWorkshop("echo")
	ws.AddJob(
		user.NewLoginJob(),
		dev.NewDevListJob(),
		storage.NewStorageJob(),
	)

	for i := range bc.App.WorkerCount {
		ws.AddWorker(worker.NewWorker(int64(i)+bc.App.FirstUid, logger))
	}

	return ws
}
