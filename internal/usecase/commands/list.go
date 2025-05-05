package commands

import (
	"errors"
	"proxy-server-with-tg-admin/internal/entity"
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
	ListUsers() ([]*entity.User, error)
	DeleteUser(username string) error
	DeleteUserStat(username string) error
}

type Cmd interface {
	Id() string
	Arguments() []string
	Run(args ...string) (string, error)
}

type List struct {
	list []Cmd
}

func New(storage StorageInterface) *List {
	return &List{
		list: []Cmd{
			&CreateUser{storage: storage},
			&ActivateUser{storage: storage},
			&StopUser{storage: storage},
			&UpdatePassword{storage: storage},
			&UpdateTtl{storage: storage},
			&ListUsers{storage: storage},
			&GetStatistic{storage: storage},
			&DeleteUser{storage: storage},
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
