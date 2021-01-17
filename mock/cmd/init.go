package main

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/mock/internal/conf"
	"github.com/vulcan-frame/vulcan-game/mock/internal/job/user"
	"github.com/vulcan-frame/vulcan-game/mock/internal/worker"
	"github.com/vulcan-frame/vulcan-game/mock/internal/workshop"
)

func initWorkshop(_ context.Context, logger log.Logger, bc *conf.Bootstrap) *workshop.Workshop {
	ws := workshop.NewWorkshop("echo")
	ws.AddJob(
		user.NewLoginJob(),
	)

	for i := int64(0); i < bc.App.WorkerCount; i++ {
		ws.AddWorker(worker.NewWorker(i+bc.App.FirstUid, logger))
	}
	return ws
}
