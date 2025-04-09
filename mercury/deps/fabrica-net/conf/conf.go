package conf

import (
	"flag"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	region    string
	zone      string
	deployEnv string
	host      string
	addrs     string
	weight    int64
	offline   bool

	Conf *Config
)

func init() {
	var (
		defHost, _    = os.Hostname()
		defAddrs      = os.Getenv("ADDRS")
		defWeight, _  = strconv.ParseInt(os.Getenv("WEIGHT"), 10, 32)
		defOffline, _ = strconv.ParseBool(os.Getenv("OFFLINE"))
	)
	flag.StringVar(&region, "region", os.Getenv("REGION"), "avaliable region. or use REGION env variable, value: sh etc.")
	flag.StringVar(&zone, "zone", os.Getenv("ZONE"), "avaliable zone. or use ZONE env variable, value: sh001/sh002 etc.")
	flag.StringVar(&deployEnv, "deploy.env", os.Getenv("DEPLOY_ENV"), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
	flag.StringVar(&host, "host", defHost, "machine hostname. or use default machine hostname.")
	flag.StringVar(&addrs, "addrs", defAddrs, "server public ip addrs. or use ADDRS env variable, value: localhost etc.")
	flag.Int64Var(&weight, "weight", defWeight, "load balancing weight, or use WEIGHT env variable, value: 10 etc.")
	flag.BoolVar(&offline, "offline", defOffline, "server offline. or use OFFLINE env variable, value: true/false etc.")
}

// Init init config.
func Init() {
	env := Env{
		Region:    region,
		Zone:      zone,
		DeployEnv: deployEnv,
		Host:      host,
		Weight:    weight,
		Addrs:     strings.Split(addrs, ","),
		Offline:   offline,
	}

	Conf = Default(env)
}

func Default(env Env) *Config {
	tcp := &Server{
		WorkerSize:   runtime.NumCPU(),
		Bind:         ":7000",
		WriteBufSize: 30000,
		ReadBufSize:  30000,
		KeepAlive:    true,
		StopTimeout:  time.Second * 30,
	}
	protocol := &Worker{
		ReaderBufSize:         8192,
		ReplyChanSize:         1024,
		HandshakeTimeout:      time.Second * 10,
		RequestIdleTimeout:    time.Second * 60,
		WaitMainTunnelTimeout: time.Second * 30,
		StopTimeout:           time.Second * 3,
	}
	bucket := &Bucket{
		BucketSize: 128,
	}

	return &Config{
		Env:    env,
		Server: tcp,
		Worker: protocol,
		Bucket: bucket,
	}
}

type Config struct {
	Env Env

	Server *Server
	Worker *Worker
	Bucket *Bucket
}

type Env struct {
	Debug     bool
	Region    string
	Zone      string
	DeployEnv string
	Host      string
	Weight    int64
	Offline   bool
	Addrs     []string
}

type Server struct {
	WorkerSize   int
	Bind         string
	WriteBufSize int
	ReadBufSize  int
	KeepAlive    bool
	StopTimeout  time.Duration
}

type Worker struct {
	ReaderBufSize         int
	ReplyChanSize         int
	HandshakeTimeout      time.Duration
	RequestIdleTimeout    time.Duration
	WaitMainTunnelTimeout time.Duration
	StopTimeout           time.Duration
}

type Bucket struct {
	BucketSize int
}
