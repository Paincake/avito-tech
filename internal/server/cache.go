package server

import (
	"github.com/Paincake/avito-tech/internal/database"
	"github.com/Paincake/avito-tech/internal/dto"
	"github.com/puzpuzpuz/xsync"
	"math/rand/v2"
	"strconv"
	"sync"
	"time"
)

type BannerCache interface {
	GetBanner(featureID int64, params dto.GetUserBannerParams) (database.UserBanner, error)
}

type MemoryCache struct {
	Repository               database.BannerRepository
	Map                      xsync.Map
	KeyLocks                 xsync.Map
	MinutesToKeyInvalidation float64
	SchedulerRateMinute      int64
}

func NewMemoryCache(repository database.BannerRepository,
	minutesToKeyInval float64,
	schedulerRate int64,
	done chan bool) *MemoryCache {
	cache := &MemoryCache{
		Repository:               repository,
		Map:                      *xsync.NewMap(),
		KeyLocks:                 *xsync.NewMap(),
		MinutesToKeyInvalidation: minutesToKeyInval,
		SchedulerRateMinute:      schedulerRate,
	}

	go func() {
		time.Sleep(1 * time.Duration(cache.SchedulerRateMinute))
		cache.cacheCleaningScheduler(done)
	}()
	return cache
}

func (c *MemoryCache) cacheCleaningScheduler(done chan bool) {
	ticker := time.NewTicker(time.Minute * time.Duration(c.SchedulerRateMinute))
	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				mapCopy := make(map[string]database.UserBanner)
				c.Map.Range(func(key string, value interface{}) bool {
					mapCopy[key] = value.(database.UserBanner)
					return true
				})
				for key := range mapCopy {
					val, _ := c.Map.Load(key)
					banner := val.(database.UserBanner)
					if time.Since(banner.UpdatedAt).Minutes() > c.MinutesToKeyInvalidation {
						c.Map.Delete(key)
					}
				}
			}
		}
	}()
}

func addJitterSeconds(min, max int) int64 {
	return int64(rand.IntN(max-min) + min)
}

func (c *MemoryCache) buildValue(featureID int64, params dto.GetUserBannerParams) (database.UserBanner, error) {
	key := strconv.Itoa(int(featureID))
	value, _ := c.KeyLocks.LoadOrStore(key, &sync.Mutex{})
	mtx := value.(*sync.Mutex)
	var content any
	mtx.Lock()
	defer func() {
		mtx.Unlock()
	}()

	content, ok := c.Map.Load(key)
	if ok && time.Since(content.(database.UserBanner).UpdatedAt).Minutes() < c.MinutesToKeyInvalidation {
		return content.(database.UserBanner), nil
	}
	banner, err := c.Repository.SelectUserBanner(params)
	if err != nil {
		return database.UserBanner{}, err
	}

	jitter := addJitterSeconds(-15, 15)
	banner.UpdatedAt = banner.UpdatedAt.Add(time.Second * time.Duration(jitter))

	content = banner
	c.Map.Store(key, content)
	c.KeyLocks.Delete(key)
	return content.(database.UserBanner), nil
}

func (c *MemoryCache) GetBanner(featureID int64, params dto.GetUserBannerParams) (database.UserBanner, error) {
	var err error
	value, ok := c.Map.Load(strconv.Itoa(int(featureID)))
	if !ok {
		value, err = c.buildValue(featureID, params)
		if err != nil {
			return database.UserBanner{}, err
		}
	}
	return value.(database.UserBanner), nil
}

func (c *MemoryCache) SetBanner(featureId int64, content dto.Content) error {
	return nil
}
