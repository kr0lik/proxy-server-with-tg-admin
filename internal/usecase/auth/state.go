package auth

import (
	"log/slog"
	"sync"
	"time"
)

const ttl = time.Minute * 5
const clearPeriod = time.Minute * 2

type state struct {
	userId       uint32
	userPassword string
	userTtl      time.Time
	ttl          time.Time
}

type cache struct {
	data   map[string]*state
	mu     sync.RWMutex
	logger *slog.Logger
}

func (c *cache) update(username, password string, userId uint32, userTtl time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	st, exist := c.data[username]
	if !exist {
		st = &state{}
		c.data[username] = st
	}

	st.userId = userId
	st.userPassword = password
	st.userTtl = userTtl
	st.ttl = time.Now().Add(ttl)
}

func (c *cache) get(username string) (state, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	st, exists := c.data[username]
	if !exists {
		return state{}, false
	}

	if st.ttl.Add(clearPeriod).After(time.Now()) {
		go func() {
			c.mu.Lock()
			defer c.mu.Unlock()

			st, exists := c.data[username]
			if exists {
				st.ttl = time.Now().Add(ttl)
			}
		}()
	}

	return *st, true
}

func (c *cache) checkup() {
	c.logger.Debug("Auth state checkup starting")

	ticker := time.NewTicker(clearPeriod)
	defer ticker.Stop()

	for range ticker.C {
		var toForget []string

		c.mu.RLock()
		for username, state := range c.data {
			if c.isNeedForget(username, state) {
				toForget = append(toForget, username)
			}
		}
		c.mu.RUnlock()

		if len(toForget) > 0 {
			for _, username := range toForget {
				c.mu.Lock()
				st, exist := c.data[username]
				if exist && c.isNeedForget(username, st) {
					delete(c.data, username)
				}
				c.mu.Unlock()
			}
		}
	}
}

func (c *cache) isNeedForget(username string, state *state) bool {
	now := time.Now()

	if state.ttl.Before(now) {
		c.logger.Debug("Auth state checkup", "ttl expired", username)

		return true
	}

	if !state.userTtl.IsZero() && state.userTtl.Before(now) {
		c.logger.Debug("Auth state checkup", "userTtl expired", username)

		return true
	}

	return false
}

func (c *cache) forget(username string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, username)
}

func (c *cache) updateUserTtl(username string, ttl time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.data[username]
	if exists {
		item.userTtl = ttl
	}
}
