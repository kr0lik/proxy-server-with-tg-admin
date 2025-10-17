package commands

import (
	"errors"
	"proxy-server-with-tg-admin/internal/entity"
	"proxy-server-with-tg-admin/internal/usecase/auth"
	"time"
)

const usernameArg = "{username}"

const NotActiveAccountMsg = "You are not active account."

var ErrUnknownCommand = errors.New("unknown command")
var ErrUsernameRequired = errors.New("username required")

var ErrUserExists = errors.New("user already exists")
var ErrTelegramIdExists = errors.New("another user already has that telegramId")
var ErrNoInviteToken = errors.New("no invite token")

type StorageInterface interface {
	CreateUser(username, password string) (uint32, error)
	ActivateUser(username string) error
	DeactivateUser(username string) error
	RenameUser(username, usernameTo string) error
	UpdatePassword(username, password string) error
	UpdateTtl(username string, ttl time.Time) error
	UpdateInviteToken(username string, token string) error
	GetUserIdByUsername(username string) (uint32, error)
	GetUserByUsername(username string) (*entity.User, error)
	GetUserByTelegramId(telegramId int64) (*entity.User, error)
	GetUserByInviteToken(token string) (*entity.User, error)
	AssignTelegramIdByInviteToken(token string, telegramId int64) error
	GetStatistic(userId uint32) (*entity.UserStat, error)
	ListUsersWithStat() ([]*UsersWithStatDto, error)
	DeleteUser(userId uint32) error
	DeleteUserStat(userId uint32) error
}

type UsersWithStatDto struct {
	Username   string
	TelegramId int64
	Active     bool
	Ttl        time.Time
	TotalIn    uint64
	TotalOut   uint64
	DyesActive uint
	LastActive time.Time
}

type Cmd interface {
	Id() string
	IsForAdminOnly() bool
	Arguments() []string
	Description() string
	Run(telegramId int64, args ...string) (string, error)
}

type List struct {
	list []Cmd
}

func New(ip string, port uint, storage StorageInterface, authenticator *auth.Authenticator) *List {
	return &List{
		list: []Cmd{
			&top{ip: ip, port: port},
			&userCreate{ip: ip, port: port, storage: storage},
			&userActivate{storage: storage, authenticator: authenticator},
			&userDeactivate{storage: storage, authenticator: authenticator},
			&userRename{storage: storage, authenticator: authenticator},
			&userPasswordUpdate{ip: ip, port: port, storage: storage, authenticator: authenticator},
			&userTtlUpdate{storage: storage, authenticator: authenticator},
			&userDelete{storage: storage, authenticator: authenticator},
			&userStatisticGet{storage: storage},
			&userStatisticClear{storage: storage},
			&userList{storage: storage},
			&userInviteCreate{storage: storage},
			&join{storage: storage},
			&myInfoGet{ip: ip, port: port, storage: storage},
			&myNameUpdate{ip: ip, port: port, storage: storage},
			&myPasswordUpdate{ip: ip, port: port, storage: storage},
		},
	}
}

func (c *List) List() []Cmd {
	return c.list
}

func (c *List) Get(id string) (Cmd, error) { //nolint: ireturn
	for _, cmd := range c.list {
		if cmd.Id() == id {
			return cmd, nil
		}
	}

	return nil, ErrUnknownCommand
}
