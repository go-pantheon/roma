package internal

import (
	"sync"
	"sync/atomic"

	"github.com/go-pantheon/fabrica-net/conf"
)

var (
	uidWidMap  = &sync.Map{}
	shardCount uint64
)

type Buckets struct {
	buckets []*sync.Map
	size    *atomic.Int64
}

func NewBuckets(c *conf.Bucket) *Buckets {
	shardCount = uint64(c.BucketSize)

	m := &Buckets{
		buckets: make([]*sync.Map, shardCount),
		size:    &atomic.Int64{},
	}
	for i := uint64(0); i < shardCount; i++ {
		m.buckets[i] = &sync.Map{}
	}
	return m
}

func (bs *Buckets) getBucket(wid uint64) *sync.Map {
	return bs.buckets[getBucketKey(wid)]
}

func getBucketKey(wid uint64) uint64 {
	return wyhash(wid) & (uint64(shardCount) - 1)
}

// wyhash generates a 64-bit hash for the given 64-bit key using wyhash algorithm.
func wyhash(key uint64) uint64 {
	x := key
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	x *= 0xc4ceb9fe1a85ec53
	x ^= x >> 33

	return x
}

func (bs Buckets) Worker(key uint64) *Worker {
	if any, ok := bs.getBucket(key).Load(key); ok {
		return any.(*Worker)
	}
	return nil
}

func (bs Buckets) Put(w *Worker) *Worker {
	old, ok := bs.getBucket(w.WID()).LoadOrStore(w.WID(), w)
	uidWidMap.Store(w.UID(), w.WID())
	bs.size.Add(1)
	if !ok {
		return nil
	}
	return old.(*Worker)
}

func (bs Buckets) Del(w *Worker) {
	uidWidMap.Delete(w.UID())
	bs.getBucket(w.WID()).Delete(w.WID())
	bs.size.Add(-1)
}

func (bs *Buckets) Walk(f func(w *Worker) bool) {
	continued := true
	for _, b := range bs.buckets {
		b.Range(func(key, value any) bool {
			v, ok := value.(*Worker)
			if !ok {
				return true
			}
			continued = f(v)
			return continued
		})
		if !continued {
			break
		}
	}
}

func (bs Buckets) GetByUID(uid int64) *Worker {
	wid, ok := uidWidMap.Load(uid)
	if !ok {
		return nil
	}
	if any, ok := bs.getBucket(wid.(uint64)).Load(wid.(uint64)); ok {
		return any.(*Worker)
	}
	return nil
}

func (bs Buckets) GetByUIDs(uids []int64) map[int64]*Worker {
	wids := make([]uint64, 0, len(uids))
	for _, uid := range uids {
		if wid, ok := uidWidMap.Load(uid); ok {
			wids = append(wids, wid.(uint64))
		}
	}

	bucketKeys := make([][]uint64, shardCount)
	result := make(map[int64]*Worker, len(wids))

	for _, wid := range wids {
		key := getBucketKey(wid)
		bucketKeys[key] = append(bucketKeys[key], wid)
	}

	for key, wids := range bucketKeys {
		if len(wids) == 0 {
			continue
		}

		bucket := bs.buckets[key]
		for _, wid := range wids {
			if any, ok := bucket.Load(wid); ok {
				worker := any.(*Worker)
				result[worker.UID()] = worker
			}
		}
	}

	return result
}
