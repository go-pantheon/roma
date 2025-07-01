package workshop

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"github.com/go-pantheon/roma/mercury/internal/worker"
)

type Workshop struct {
	xsync.Stoppable

	logger  log.Logger
	Name    string
	Jobs    []*job.Job
	WIDs    []int64
	Workers map[int64]*worker.Worker
}

func NewWorkshop(name string, logger log.Logger) *Workshop {
	ws := &Workshop{
		Stoppable: xsync.NewStopper(10 * time.Second),
		logger:    logger,
		Name:      name,
		Jobs:      make([]*job.Job, 0, 16),
		Workers:   make(map[int64]*worker.Worker, core.AppConf().WorkerCount),
	}

	firstWID := core.AppConf().FirstUid

	for i := range core.AppConf().WorkerCount {
		ws.WIDs = append(ws.WIDs, firstWID+i)
	}

	return ws
}

func (ws *Workshop) AddJob(js ...*job.Job) {
	ws.Jobs = append(ws.Jobs, js...)

	for _, j := range js {
		log.Infof("[workshop-%s] add job: %d", ws.Name, j.T)
	}
}

func (ws *Workshop) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, wid := range ws.WIDs {
		if _, ok := ws.Workers[wid]; ok {
			continue
		}

		w := worker.NewWorker(wid, ws.logger)
		ws.Workers[w.UID()] = w

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wg.Add(1)

			go func(w *worker.Worker) {
				defer wg.Done()

				st := time.Now()

				if err := w.Start(ctx, ws.Jobs); err != nil {
					w.Log().Errorf("[worker-%d] start failed: %+v", w.UID(), err)
					return
				}

				w.Log().Infof("[worker-%d] start at: %s", w.UID(), st.Format("15:04:05"))

				w.WaitStopped()

				w.Log().Infof("[worker-%d] completed. used: %s", w.UID(), time.Since(st).String())
			}(w)

			time.Sleep(core.AppConf().LoginInterval.AsDuration() + time.Duration(rand.Int63n(2000)*int64(time.Millisecond)))
		}
	}

	wg.Wait()

	return nil
}

func (ws *Workshop) Stop(ctx context.Context) (err error) {
	return ws.TurnOff(func() error {
		var wg sync.WaitGroup

		for _, w := range ws.Workers {
			wg.Add(1)

			go func(w *worker.Worker) {
				defer wg.Done()

				if werr := w.Stop(ctx); werr != nil {
					err = errors.JoinUnsimilar(err, werr)
				}
			}(w)
		}

		wg.Wait()

		return err
	})
}
