package core

import (
	"context"
	"sync"
	"time"

	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/xsync"
)

type RouteTableCache struct {
	xsync.Stoppable

	routeTable  routetable.ReadOnlyRouteTable
	cache       sync.Map // routeTableKey -> *CacheEntry
	ticker      *time.Ticker
	serviceName string
	cacheTTL    time.Duration
}

func NewRouteTableCache(routeTable routetable.ReadOnlyRouteTable, serviceName string, cacheTTL time.Duration) *RouteTableCache {
	cache := &RouteTableCache{
		Stoppable:   xsync.NewStopper(time.Second * 2),
		routeTable:  routeTable,
		serviceName: serviceName,
		cacheTTL:    cacheTTL,
	}

	cache.ticker = time.NewTicker(cacheTTL * 2)

	xsync.Go("route-table-cache-cleanup", func() error {
		return cache.cleanupLoop()
	})

	return cache
}

func (c *RouteTableCache) cleanupLoop() error {
	for {
		select {
		case <-c.StopTriggered():
			return xsync.ErrStopByTrigger
		case <-c.ticker.C:
			c.cleanupExpiredEntries()
		}
	}
}

func (c *RouteTableCache) cleanupExpiredEntries() {
	c.cache.Range(func(key, value any) bool {
		entry := value.(*CacheEntry)
		if entry.NeedRemove() {
			c.cache.Delete(key)
		}

		return true
	})
}

func (c *RouteTableCache) buildCacheKey(uid int64, color string) string {
	return c.routeTable.BuildKey(color, uid)
}

func (c *RouteTableCache) get(ctx context.Context, uid int64, color string) (string, error) {
	cacheKey := c.buildCacheKey(uid, color)

	if value, exists := c.cache.Load(cacheKey); exists {
		entry := value.(*CacheEntry)
		if !entry.IsExpired() {
			return entry.Value, nil
		}

		c.cache.Delete(cacheKey)
	}

	endpoint, err := c.routeTable.Get(ctx, color, uid)
	if err != nil {
		return "", errors.WithMessage(err, "failed to get route from table")
	}

	c.updateCache(cacheKey, endpoint)

	return endpoint, nil
}

// batchGet retrieves endpoint list for multiple uids with given color
func (c *RouteTableCache) batchGet(ctx context.Context, uids []int64, color string) (map[string][]int64, error) {
	ret := make(map[string][]int64)
	uncacheableUids := make([]int64, 0, len(uids))

	for _, uid := range uids {
		key := c.buildCacheKey(uid, color)

		if value, exists := c.cache.Load(key); exists {
			entry := value.(*CacheEntry)
			if !entry.IsExpired() {
				ret[entry.Value] = append(ret[entry.Value], uid)
				continue
			}

			c.cache.Delete(key)
		}

		uncacheableUids = append(uncacheableUids, uid)
	}

	uncacheableEndpoints, err := c.routeTable.BatchGet(ctx, color, uncacheableUids)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get route from table")
	}

	for i, endpoint := range uncacheableEndpoints {
		if endpoint == "" {
			continue
		}

		ret[endpoint] = append(ret[endpoint], uncacheableUids[i])
	}

	return ret, nil
}

func (c *RouteTableCache) updateCache(cacheKey, endpoint string) {
	c.cache.Store(cacheKey, NewCacheEntry(endpoint, c.cacheTTL))
}

// BuildKey returns the key for given color and oid
func (c *RouteTableCache) BuildKey(color string, oid int64) string {
	return c.routeTable.BuildKey(color, oid)
}

// GetCacheSize returns the current cache size for monitoring
func (c *RouteTableCache) GetCacheSize() int {
	count := 0
	c.cache.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

// GetCacheStats returns detailed cache statistics
func (c *RouteTableCache) GetCacheStats() map[string]any {
	totalEntries := 0
	expiredEntries := 0

	c.cache.Range(func(key, value any) bool {
		totalEntries++
		entry := value.(*CacheEntry)
		if entry.IsExpired() {
			expiredEntries++
		}
		return true
	})

	return map[string]any{
		"total_entries":   totalEntries,
		"expired_entries": expiredEntries,
		"cache_ttl_ms":    c.cacheTTL.Milliseconds(),
	}
}

// ClearCache removes all entries from the cache (useful for testing)
func (c *RouteTableCache) ClearCache() {
	c.cache.Range(func(key, value any) bool {
		c.cache.Delete(key)
		return true
	})
}

type CacheEntry struct {
	Value        string
	ExpiresAt    time.Time
	NeedRemoveAt time.Time
}

func NewCacheEntry(value string, ttl time.Duration) *CacheEntry {
	return &CacheEntry{
		Value:        value,
		ExpiresAt:    time.Now().Add(ttl),
		NeedRemoveAt: time.Now().Add(ttl * 2),
	}
}

func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

func (e *CacheEntry) NeedRemove() bool {
	return time.Now().After(e.NeedRemoveAt)
}
