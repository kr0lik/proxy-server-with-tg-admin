package auth

import (
	"fmt"
	"log/slog"
	"proxy-server-with-tg-admin/internal/entity"
	"sync"
	"time"
)

type UserStorageInterface interface {
	GetUser(username string) (*entity.User, error)
}

type Authenticator struct {
	storage UserStorageInterface
	cache   *cache
	logger  *slog.Logger
}

func New(storage UserStorageInterface, logger *slog.Logger) *Authenticator {
	cache := &cache{
		data:    make(map[string]*state),
		mu:      sync.RWMutex{},
		storage: storage,
		logger:  logger,
	}

	go cache.checkup()

	return &Authenticator{storage: storage, cache: cache, logger: logger}
}

func (a *Authenticator) Authenticate(username, password string) bool {
	state, ok := a.cache.get(username)

	if ok && state.userPassword == password {
		return true
	}

	user, err := a.storage.GetUser(username)
	if err != nil {
		a.logger.Debug("Auth authenticate fail", "GetUser", err)
		return false
	}

	if user.Password != password {
		a.logger.Debug("Auth authenticate fail", "user.Password", user.Password)
		return false
	}

	if !a.validate(user) {
		return false
	}

	a.cache.update(user.Username, user.Password, user.ID, user.Ttl)

	return true
}

func (a *Authenticator) Revoke(username string) {
	a.cache.forget(username)
}

func (a *Authenticator) GetUserId(username, password string) (uint32, error) {
	state, ok := a.cache.get(username)

	if ok && state.userPassword == password {
		return state.userId, nil
	}

	if !ok {
		a.logger.Warn("Auth cache missed", "GetUserId", username)
	}

	user, err := a.storage.GetUser(username)
	if err != nil {
		return 0, err
	}

	if user.Password != password {
		return 0, fmt.Errorf("ivalid password")
	}

	return user.ID, nil
}

func (a *Authenticator) validate(user *entity.User) bool {
	if !user.Active {
		a.logger.Debug("Auth validate fail", "user.Active", user.Active)
		return false
	}

	if user.Ttl.Unix() > 0 && user.Ttl.Unix() < time.Now().Unix() {
		a.logger.Debug("Auth validate fail", "user.Ttl", user.Ttl.Format(time.DateTime))
		return false
	}

	return true
}
