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
	RenameUser(username, usernameTo string) error
	UpdatePassword(username, password string) error
	UpdateTtl(username string, ttl time.Time) error
	GetStatistic(username string) (*entity.UserStat, error)
	ListUsersWithStat() ([]*UsersWithStatDto, error)
	DeleteUserWithStat(username string) error
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

func New(ip string, port uint, storage StorageInterface, authenticator *auth.Authenticator) *List {
	return &List{
		list: []Cmd{
			&top{ip: ip, port: port},
			&createUser{ip: ip, port: port, storage: storage},
			&activateUser{storage: storage, authenticator: authenticator},
			&deactivateUser{storage: storage, authenticator: authenticator},
			&renameUser{storage: storage, authenticator: authenticator},
			&updatePassword{ip: ip, port: port, storage: storage, authenticator: authenticator},
			&updateTtl{storage: storage, authenticator: authenticator},
			&deleteUser{storage: storage, authenticator: authenticator},
			&getStatistic{storage: storage},
			&clearStatistic{storage: storage},
			&listUsers{storage: storage},
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
