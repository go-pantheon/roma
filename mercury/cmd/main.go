package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xlog"
	"github.com/go-pantheon/roma/gamedata"
	"github.com/go-pantheon/roma/mercury/internal/conf"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/core/security"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	Name        = "mercury"
	flagConf    string
	gameDataDir string
)

func init() {
	flag.StringVar(&flagConf, "conf", "mercury/configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&gameDataDir, "gamedata", "gen/gamedata/json", "config path, eg: -gamedata json")
}

func main() {
	flag.Parse()

	flagConf, err := filepath.Abs(flagConf)
	if err != nil {
		panic(err)
	}

	gameDataDir, err = filepath.Abs(gameDataDir)
	if err != nil {
		panic(err)
	}

	c := config.New(
		config.WithSource(
			env.NewSource(profile.OrgPrefix),
			file.NewSource(flagConf),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	var sc conf.Secret
	if err := c.Scan(&sc); err != nil {
		panic(err)
	}

	if err := security.Init(&sc); err != nil {
		panic(err)
	}

	logger := xlog.Init(bc.Log.Type, bc.Log.Level, bc.Label.Profile, bc.Label.Color, "mercury", "v0.0.1", "local")

	core.Init(&bc, logger)
	gamedata.Load(gameDataDir)

	log.Infof("[%s] is running. profile=%s, color=%s", Name, bc.Label.Profile, bc.Label.Color)

	if err = run(logger); err != nil {
		if errors.Is(err, workshopCompletedFlag) {
			log.Infof("[%s] is completed. profile=%s, color=%s", Name, bc.Label.Profile, bc.Label.Color)
			return
		}

		panic(err)
	}
}

var workshopCompletedFlag = errors.Errorf("mercury completed")

func run(logger log.Logger) error {
	ws := newWorkshop(logger)

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return ws.Run(ctx)
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	eg.Go(func() error {
		<-c
		return ws.Stop(ctx)
	})

	return eg.Wait()
}
