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
func (c *updateTimer) Add(branch string) (bool, *time.Time) {
	c.mux.Lock()
	defer c.mux.Unlock()
	curTime := time.Now()
	if _, ok := c.v[branch]; ok {
		c.v[branch] = curTime
		return false, &curTime
	}
	c.v[branch] = time.Now()
	return true, &curTime
}

// return true if updated time was the newest or branch had no value
// false if time was not newest one, + newest time
func (c *updateTimer) Remove(branch string, time time.Time) (bool, *time.Time) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if val, ok := c.v[branch]; ok && time.Before(val) {
		return false, &val
	}
	delete(c.v, branch)
	return true, nil
}

var updateManager = updateTimer{
	v: make(map[string]time.Time),
}

func update(updaterConfig *updater.AppConfig, storage *updater.HashStorage, buildCtx *updater.BuildContext, branch string) bool {
	for _, build := range updaterConfig.Build {
		if build.Branch == branch {
			go func() {
				if update, curTime := updateManager.Add(branch); update {
					for {
						updater.BuildStep(updaterConfig, storage, buildCtx, build)
						time.Sleep(10 * time.Second)
						if update, curTime = updateManager.Remove(branch, *curTime); update {
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
