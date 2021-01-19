package main

import (
	"flag"
	"path/filepath"

	"github.com/vulcan-frame/vulcan-game/app/player/internal/conf"
	"github.com/vulcan-frame/vulcan-game/gamedata"
)

var (
	flagConf    string
	gameDataDir string
)

func init() {
	flag.StringVar(&flagConf, "conf", "app/player/configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&gameDataDir, "gamedata", "gen/gamedata/json", "config path, eg: -gamedata json")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, health *health.Server, label *conf.Label, rr registry.Registrar) *kratos.App {
	md := map[string]string{
		profile.SERVICE: label.Service,
		profile.PROFILE: label.Profile,
		profile.VERSION: label.Version,
		profile.COLOR:   label.Color,
		profile.NODE:    label.Node,
	}

	url, err := gs.Endpoint()
	if err != nil {
		panic(err)
	}

	profile.Init(label.Profile, label.Color, label.Zone, label.Version, label.Node, url)

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

	gameDataDir, err = filepath.Abs(gameDataDir)
	if err != nil {
		panic(err)
	}

	c := config.New(
		config.WithSource(
			env.NewSource(profile.ORG_PREFIX),
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

	xtime.Init(bc.Label.Language)

	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		panic(err)
	}

	if err := trace.Init(bc.Trace.Endpoint, bc.Label.Service, bc.Label.Profile, bc.Label.Color); err != nil {
		panic(err)
	}

	logger := xlog.Init(bc.Log.Type, bc.Log.Level, bc.Label.Profile, bc.Label.Color, bc.Label.Service, bc.Label.Version, bc.Label.Node)
	metrics.Init(bc.Label.Service)

	gamedata.Load(gameDataDir)

	app, cleanup, err := initApp(bc.Server, bc.Label, &rc, bc.Data, logger, health.NewServer(bc.Server.Health))
	if err != nil {
		panic(err)
	}
	defer cleanup()

	log.Infof("[%s] is running", bc.Label.Service)

	if err = app.Run(); err != nil {
		panic(err)
	}
}
