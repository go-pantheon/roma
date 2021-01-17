package workshop

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"github.com/vulcan-frame/vulcan-game/mock/internal/job"
	"github.com/vulcan-frame/vulcan-game/mock/internal/worker"
)

type Workshop struct {
	Name string

	Jobs    []*job.Job
	Workers map[int64]*worker.Worker
}

func NewWorkshop(name string) *Workshop {
	ws := &Workshop{Name: name}
	ws.Workers = make(map[int64]*worker.Worker, 8)
	return ws
}

func (ws *Workshop) AddJob(js ...*job.Job) {
	ws.Jobs = append(ws.Jobs, js...)
	for _, j := range js {
		log.Infof("add job: %d", j.T)
	}
}

func (ws *Workshop) AddWorker(w *worker.Worker) {
	ws.Workers[w.UID()] = w
	log.Infof("add worker: %d", w.UID())
}

func (ws *Workshop) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, w := range ws.Workers {
		wg.Add(1)
		go func(w *worker.Worker) {
			defer wg.Done()
			time.Sleep(base.App().LoginInterval.AsDuration() + time.Duration(rand.Int63n(2*1000)*int64(time.Millisecond)))
			do(ctx, w, ws.Jobs)
		}(w)
	}
	wg.Wait()
}

func do(c context.Context, w *worker.Worker, jobs []*job.Job) {
	ctx := base.NewContext(c, w)
	startTime := time.Now()
	endTime := startTime.Add(base.App().WorkMinInterval.AsDuration())

	w.Log().Infof("worker:%d start at: %s", w.UID(), startTime.Format("15:04:05"))

	if err := w.Start(ctx); err != nil {
		w.Log().Errorf("%+v", err)
		return
	}

	defer w.Stop()

	if err := w.Handshake(ctx); err != nil {
		w.Log().Errorf("%+v", err)
		return
	}

	for _, j := range jobs {
		for _, t := range j.Tasks {
			if w.IsStopping() {
				return
			}
			w.DistributeTask(t)
			time.Sleep(base.App().WorkMinInterval.AsDuration() + time.Duration(rand.Int63n(2*1000)*int64(time.Millisecond)))
		}
		time.Sleep(2 * time.Second)
	}

	time.Sleep(time.Until(endTime))

	w.Completed.Store(true)
	w.WorkingTime = time.Since(startTime)
}

func (ws *Workshop) Stop() {
	for _, w := range ws.Workers {
		w.TriggerStop()
	}
}
