package core

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/mercury/internal/conf"
)

var (
	dep *Dependency
)

func Init(conf *conf.Bootstrap, log log.Logger) {
	dep = NewDependency(NewConfig(conf), log)
}

type Dependency struct {
	conf Config
	log  log.Logger
}

func NewDependency(conf Config, log log.Logger) *Dependency {
	return &Dependency{
		conf: conf,
		log:  log,
	}
}

type Config struct {
	Boot *conf.Bootstrap
}

func NewConfig(boot *conf.Bootstrap) Config {
	return Config{
		Boot: boot,
	}
}

func Dep() *Dependency {
	return dep
}

func BootConf() *conf.Bootstrap {
	return dep.conf.Boot
}

func Unencrypted() bool {
	return dep.conf.Boot.Label.Unencrypted
}

func Color() string {
	return dep.conf.Boot.Label.Color
}

func ServerId() int64 {
	return dep.conf.Boot.App.ServerId
}

func AppConf() *conf.App {
	return dep.conf.Boot.App
}
