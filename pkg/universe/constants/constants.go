package constants

import "time"

const (
	WorkerPushTimeout    = time.Second * 3
	WorkerPersistTimeout = time.Second * 5
)

const (
	WorkerRenewalTickDuration = time.Second * 10
	WorkerSecondTickDuration  = time.Second * 1
	WorkerMinuteTickDuration  = time.Minute * 1
	WorkerPersistTickDuration = time.Second * 10
	StatisticTickDuration     = time.Second * 5
	CacheExpiredDuration      = time.Minute * 10
)

const (
	WorkerSize       = 1024 * 20
	WorkerEventSize  = 512
	WorkerReplySize  = 1024
	BroadcastMsgSize = 1024
	WorkerHolderSize = 1024
)

const (
	PersistManagerStopTimeout = time.Second * 20
)
