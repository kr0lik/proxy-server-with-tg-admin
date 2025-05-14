package auth

import (
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"log/slog"
	"proxy-server-with-tg-admin/internal/entity"
	"time"
)

var ErrUserPassword = errors.New("invalid password")
var ErrUserInactive = errors.New("user inactive")
var ErrUserTtl = errors.New("user ttl exceeded")

type StorageInterface interface {
	GetUser(username string) (*entity.User, error)
}

type Authenticator struct {
	storage StorageInterface
	single  singleflight.Group
	cache   *cache
	logger  *slog.Logger
}

func New(storage StorageInterface, logger *slog.Logger) *Authenticator {
	cache := &cache{
		data:   make(map[string]*state),
		logger: logger,
	}

	go cache.checkup()

	return &Authenticator{storage: storage, cache: cache, logger: logger}
}

func (a *Authenticator) Authenticate(username, password string) (uint32, error) {
	const op = "auth.Authenticator.Authenticate"

	st, exist := a.cache.get(username)

	if exist && st.userPassword == password {
		return st.userId, nil
	}

	userId, err, _ := a.single.Do(username, func() (interface{}, error) {
		return a.authenticateOnce(username, password)
	})
	defer a.single.Forget(username)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userIdTyped, ok := userId.(uint32)
	if !ok {
		return 0, fmt.Errorf("%s: userId is not uint32", op)
	}

	return userIdTyped, nil
}

func (a *Authenticator) authenticateOnce(username, password string) (uint32, error) {
	const op = "auth.Authenticator.authenticateOnce"

	user, err := a.storage.GetUser(username)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if user.Password != password {
		return 0, ErrUserPassword
	}

	if !user.Active {
		return 0, ErrUserInactive
	}

	if user.Ttl.Unix() > 0 && user.Ttl.Unix() < time.Now().Unix() {
		return 0, ErrUserTtl
	}

	a.cache.update(user.Username, user.Password, user.ID, user.Ttl)

	return user.ID, nil
}

func (a *Authenticator) Forget(username string) {
	a.cache.forget(username)
}

func (a *Authenticator) UpdateUserTtl(username string, ttl time.Time) {
	a.cache.updateUserTtl(username, ttl)
}
