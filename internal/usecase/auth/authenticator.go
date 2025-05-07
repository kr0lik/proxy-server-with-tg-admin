package auth

import (
	"log/slog"
	"proxy-server-with-tg-admin/internal/entity"
	"time"
)

type StorageInterface interface {
	GetUser(username string) (*entity.User, error)
}

type Authenticator struct {
	storage StorageInterface
	cache   *cache
	logger  *slog.Logger
}

func New(storage StorageInterface, logger *slog.Logger) *Authenticator {
	cache := &cache{
		data:    make(map[string]*state),
		storage: storage,
		logger:  logger,
	}

	go cache.checkup()

	return &Authenticator{storage: storage, cache: cache, logger: logger}
}

func (a *Authenticator) Authenticate(username, password string) uint32 {
	st, ok := a.cache.get(username)

	if ok && st.userPassword == password {
		return st.userId
	}

	user, err := a.storage.GetUser(username)
	if err != nil {
		a.logger.Debug("Auth authenticate fail", "GetUser", err)

		return 0
	}

	if user.Password != password {
		a.logger.Debug("Auth authenticate fail", "user.Password", user.Password)

		return 0
	}

	if !a.validate(user) {
		return 0
	}

	a.cache.update(user.Username, user.Password, user.ID, user.Ttl)

	return user.ID
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
