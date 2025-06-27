package workshop

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"github.com/go-pantheon/roma/mercury/internal/worker"
)

type Workshop struct {
	Name string

	Jobs    []*job.Job
	Workers map[int64]*worker.Worker
}

func NewWorkshop(name string) *Workshop {
	ws := &Workshop{
		Name:    name,
		Jobs:    make([]*job.Job, 0, 16),
		Workers: make(map[int64]*worker.Worker, core.AppConf().WorkerCount),
	}

	return ws
}

func (ws *Workshop) AddJob(js ...*job.Job) {
	ws.Jobs = append(ws.Jobs, js...)

	for _, j := range js {
		log.Infof("[%s] add job: %d", ws.Name, j.T)
	}
}

func (ws *Workshop) AddWorker(w *worker.Worker) {
	ws.Workers[w.UID()] = w
	log.Infof("[%s] add worker: %d", ws.Name, w.UID())
}

func (ws *Workshop) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for _, w := range ws.Workers {
		wg.Add(1)

		go func(w *worker.Worker) {
			defer wg.Done()

			time.Sleep(core.AppConf().LoginInterval.AsDuration() + time.Duration(rand.Int63n(2000)*int64(time.Millisecond)))
			work(ctx, w, ws.Jobs)
		}(w)
	}

	wg.Wait()
}

func work(c context.Context, w *worker.Worker, jobs []*job.Job) {
	ctx := core.NewContext(c, w)
	st := time.Now()
	et := st.Add(core.AppConf().WorkMinInterval.AsDuration())

	w.Log().Infof("[worker-%d] start at: %s", w.UID(), st.Format("15:04:05"))

	if err := w.Start(ctx); err != nil {
		w.Log().Errorf("%+v", err)
		return
	}

	defer w.Stop(ctx)

	for _, j := range jobs {
		for _, t := range j.Tasks {
			if w.OnStopping() {
				return
			}

			w.DistributeTask(t)
			time.Sleep(core.AppConf().WorkMinInterval.AsDuration() + time.Duration(rand.Int63n(2*1000)*int64(time.Millisecond)))
		}
		
		time.Sleep(2 * time.Second)
	}

	time.Sleep(time.Until(et))

	w.Completed.Store(true)
	w.WorkingTime = time.Since(st)
}

func (ws *Workshop) Stop(ctx context.Context) {
	for _, w := range ws.Workers {
		w.Stop(ctx)
	}
}
