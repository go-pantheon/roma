package constants

import "time"

const (
	AsyncMongoTimeout = time.Second * 1
	AsyncMySQLTimeout = time.Second * 5
	AsyncRedisTimeout = time.Second * 1
	AsyncGRPCTimeout  = time.Second * 3
)

const (
	GameTunnelChangeTimeout      = time.Second * 3
	WorkerPersistTickDuration    = time.Second * 10
	ManagerStatisticTickDuration = time.Second * 5
	CacheExpiredDuration         = time.Minute * 10
)

const (
	WorkerMaxErrorCount = 10
)

const (
	WorkerSize       = 1024 * 20
	WorkerEventSize  = 512
	WorkerReplySize  = 1024
	WorkerHolderSize = 1024
)
