package auth

import (
	"log/slog"
	"sync"
	"time"
)

const ttl = time.Second * 60 * 1
const clearPeriod = time.Second * 60

type state struct {
	userId       uint32
	userPassword string
	userTtl      time.Time
	ttl          time.Time
}

// TODO
type cache struct {
	data    map[string]*state
	mu      sync.RWMutex
	storage UserStorageInterface
	logger  *slog.Logger
}

func (c *cache) update(username, password string, userId uint32, userTtl time.Time) {
	c.mu.Lock()
	c.data[username] = &state{userId: userId, userPassword: password, userTtl: userTtl, ttl: time.Now().Add(ttl)}
	c.mu.Unlock()
}

func (c *cache) get(username string) (state, bool) {
	c.mu.RLock()
	item, exists := c.data[username]
	if !exists {
		c.mu.RUnlock()
		return state{}, false
	}
	c.mu.RUnlock()

	c.extendTTL(item)

	return *item, true
}

func (c *cache) forget(username string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, username)
}

func (c *cache) extendTTL(state *state) {
	if state.ttl.Add(time.Minute).Before(time.Now()) {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if state.userTtl.IsZero() || state.userTtl.After(time.Now()) {
		state.ttl = time.Now().Add(ttl)
	}
}

func (c *cache) checkup() {
	c.logger.Debug("Auth state checkup starting")

	ticker := time.NewTicker(clearPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for username, state := range c.data {
				if state.ttl.Before(time.Now()) {
					c.logger.Debug("Auth state checkup", "ttl expired", username)
					c.forget(username)
				} else if !state.userTtl.IsZero() && state.userTtl.Before(time.Now()) {
					c.logger.Debug("Auth state checkup", "userTtl expired", username)
					c.forget(username)
				} else {
					user, err := c.storage.GetUser(username)
					if err == nil {
						if !user.Active || user.Password != state.userPassword || user.Ttl != state.userTtl {
							c.logger.Debug("Auth state checkup", "user changed", username)
							c.forget(username)
						}
					}
				}
			}
		}
	}
}
