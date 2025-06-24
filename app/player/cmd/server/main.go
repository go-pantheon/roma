package main

import (
	"flag"
	"path/filepath"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-pantheon/fabrica-kit/metrics"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/trace"
	"github.com/go-pantheon/fabrica-kit/xlog"
	"github.com/go-pantheon/fabrica-net/http/health"
	"github.com/go-pantheon/fabrica-util/xtime"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/go-pantheon/roma/gamedata"
)

var (
	flagConf        string
	flagGameDataDir string
)

func init() {
	flag.StringVar(&flagConf, "conf", "app/player/configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&flagGameDataDir, "gamedata", "gen/gamedata/json", "gamedata path, eg: -gamedata json")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, health *health.Server, label *conf.Label, rr registry.Registrar) *kratos.App {
	md := map[string]string{
		profile.ServiceKey: label.Service,
		profile.ProfileKey: label.Profile,
		profile.VersionKey: label.Version,
		profile.ColorKey:   label.Color,
		profile.NodeKey:    label.Node,
	}

	url, err := gs.Endpoint()
	if err != nil {
		panic(err)
	}

	profile.Init(label.Service, label.Profile, label.Color, label.Zone, label.Version, label.Node, url)

	return kratos.New(
		kratos.Name(label.Service),
		kratos.Version(label.Version),
		kratos.Metadata(md),
		kratos.Logger(logger),
		kratos.Server(health, hs, gs),
		kratos.Registrar(rr),
	)
}

func main() {
	flag.Parse()

	flagConf, err := filepath.Abs(flagConf)
	if err != nil {
		panic(err)
	}

	flagGameDataDir, err = filepath.Abs(flagGameDataDir)
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

	if err := xtime.InitSimple(bc.Label.Language); err != nil {
		panic(err)
	}

	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		panic(err)
	}

	if err := trace.Init(bc.Trace.Endpoint, bc.Label.Service, bc.Label.Profile, bc.Label.Color); err != nil {
		panic(err)
	}

	metrics.Init(bc.Label.Service)
	gamedata.Load(flagGameDataDir)

	logger := xlog.Init(bc.Log.Type, bc.Log.Level, bc.Label.Profile, bc.Label.Color, bc.Label.Service, bc.Label.Version, bc.Label.Node)

	app, cleanup, err := initApp(bc.Server, bc.Label, bc.Recharge, &rc, bc.Data, logger, health.NewServer(bc.Server.Health))
	if err != nil {
		panic(err)
	}
	defer cleanup()

	log.Infof("[%s] is running", bc.Label.Service)

	if err = app.Run(); err != nil {
		panic(err)
	}
}
