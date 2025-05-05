package auth

import (
	"log/slog"
	"sync"
	"time"
)

const ttl = time.Second * 60 * 5
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
	defer c.mu.Unlock()

	st, ok := c.data[username]
	if !ok {
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

	item, exists := c.data[username]
	if !exists {
		return state{}, false
	}

	c.extendTTLIfNeeded(item)

	return *item, true
}

func (c *cache) forget(username string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, username)
}

func (c *cache) extendTTLIfNeeded(state *state) {
	if state.ttl.Add(time.Minute).Before(time.Now()) {
		return
	}

	if state.userTtl.IsZero() || state.userTtl.After(time.Now()) {
		state.ttl = time.Now().Add(ttl)
	}
}

func (c *cache) checkup() {
	c.logger.Debug("Auth state checkup starting")

	ticker := time.NewTicker(clearPeriod)
	defer ticker.Stop()

	for range ticker.C {
		var toForget []string

		now := time.Now()

		c.mu.RLock()

		for username, state := range c.data {
			if state.ttl.Before(now) {
				c.logger.Debug("Auth state checkup", "ttl expired", username)
				toForget = append(toForget, username)
			} else if !state.userTtl.IsZero() && state.userTtl.Before(now) {
				c.logger.Debug("Auth state checkup", "userTtl expired", username)
				toForget = append(toForget, username)
			} else {
				user, err := c.storage.GetUser(username)
				if err == nil {
					if !user.Active || user.Password != state.userPassword || user.Ttl != state.userTtl {
						c.logger.Debug("Auth state checkup", "user changed", username)
						toForget = append(toForget, username)
					}
				}
			}
		}

		c.mu.RUnlock()

		if len(toForget) > 0 {
			c.mu.Lock()
			for _, username := range toForget {
				delete(c.data, username)
			}
			c.mu.Unlock()
		}
	}
}
