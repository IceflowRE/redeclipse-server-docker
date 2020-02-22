package server

import (
	"sync"
	"time"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
)

type updateTimer struct {
	v   map[string]time.Time
	mux sync.Mutex
}

// return true if no other update is running
func (c *updateTimer) Add(ref string) (bool, *time.Time) {
	c.mux.Lock()
	defer c.mux.Unlock()
	curTime := time.Now()
	if _, ok := c.v[ref]; ok {
		c.v[ref] = curTime
		return false, &curTime
	}
	c.v[ref] = time.Now()
	return true, &curTime
}

// return true if updated time was the newest or reference had no value
// false if time was not newest one, + newest time
func (c *updateTimer) Remove(ref string, time time.Time) (bool, *time.Time) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if val, ok := c.v[ref]; ok && time.Before(val) {
		return false, &val
	}
	delete(c.v, ref)
	return true, nil
}

var updateManager = updateTimer{
	v: make(map[string]time.Time),
}

func update(updaterConfig *updater.AppConfig, storage *updater.HashStorage, buildCtx *updater.BuildContext, ref string) bool {
	for _, build := range updaterConfig.Build {
		if build.Ref == ref {
			go func() {
				if update, curTime := updateManager.Add(ref); update {
					for {
						updater.BuildStep(updaterConfig, storage, buildCtx, build)
						time.Sleep(10 * time.Second)
						if update, curTime = updateManager.Remove(ref, *curTime); update {
							break
						}
					}
				}
			}()
			return true
		}
	}
	return false
}
