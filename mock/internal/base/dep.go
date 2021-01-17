package base

import (
	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/mock/internal/conf"
)

var (
	Dep *Dependency
)

func Unencrypted() bool {
	return Dep.Conf.Boot.Label.Unencrypted
}

func Color() string {
	return Dep.Conf.Boot.Label.Color
}

func App() *conf.App {
	return Dep.Conf.Boot.App
}

type Dependency struct {
	Conf *Config
	log  log.Logger
}

func NewDependency(conf *Config, log log.Logger) *Dependency {
	return &Dependency{
		Conf: conf,
		log:  log,
	}
}

type Config struct {
	Boot *conf.Bootstrap
}

func NewConfig(boot *conf.Bootstrap) *Config {
	return &Config{
		Boot: boot,
	}
}
