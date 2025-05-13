package commands

import (
	"errors"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/usecase/auth"
	"time"
)

var ErrUnknownCommand = errors.New("unknown command")
var ErrUsernameRequired = errors.New("username required")

const usernameArg = "{username}"

type StorageInterface interface {
	CreateUser(username, password string) (uint32, error)
	ActivateUser(username string) error
	DeactivateUser(username string) error
	UpdatePassword(username, password string) error
	UpdateTtl(username string, ttl time.Time) error
	GetStatistic(username string) (*entity.UserStat, error)
	ListUsersWithStat() ([]*UsersWithStatDto, error)
	DeleteUser(username string) error
	DeleteUserStat(username string) error
}

type UsersWithStatDto struct {
	Username   string
	Active     bool
	Ttl        time.Time
	TotalIn    uint64
	TotalOut   uint64
	DyesActive uint
	LastActive time.Time
}

type Cmd interface {
	Id() string
	Arguments() []string
	Run(args ...string) (string, error)
}

type List struct {
	list []Cmd
}

func New(storage StorageInterface, authenticator *auth.Authenticator) *List {
	return &List{
		list: []Cmd{
			&top{},
			&createUser{storage: storage},
			&activateUser{storage: storage, authenticator: authenticator},
			&deactivateUser{storage: storage, authenticator: authenticator},
			&updatePassword{storage: storage, authenticator: authenticator},
			&updateTtl{storage: storage, authenticator: authenticator},
			&listUsers{storage: storage},
			&getStatistic{storage: storage},
			&deleteUser{storage: storage, authenticator: authenticator},
		},
	}
}

func (c *List) List() []Cmd {
	return c.list
}

func (c *List) Get(id string) (Cmd, error) {
	for _, cmd := range c.list {
		if cmd.Id() == id {
			return cmd, nil
		}
	}

	return nil, ErrUnknownCommand
}
