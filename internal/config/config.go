package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	Dev  = "dev"
	Prod = "prod"
)

type Config struct {
	env              string
	portSocks5       uint
	telegramBotToken string
	telegramAdminId  int64
	sqlitePath       string
}

func (c *Config) Env() string {
	return c.env
}

func (c *Config) PortSocks5() uint {
	return c.portSocks5
}

func (c *Config) TelegramToken() string {
	return c.telegramBotToken
}

func (c *Config) TelegramAdminId() int64 {
	return c.telegramAdminId
}

func (c *Config) SqlitePath() string {
	return strings.TrimRight(c.sqlitePath, string(os.PathSeparator))
}

func MustLoad() *Config {
	var portSocks5 uint
	var env, sqlitePath string
	var telegramBotToken string
	var telegramAdminId int64

	flag.UintVar(&portSocks5, "port-socks5", 1080, "SOCKS5 server port")
	flag.StringVar(&env, "env", Prod, "Application environment: dev or prod")
	flag.StringVar(&sqlitePath, "sqlite-path", "./.data", "Storage path")
	flag.StringVar(&telegramBotToken, "telegram-bot-token", "", "Telegram bot token")
	flag.Int64Var(&telegramAdminId, "telegram-admin-id", 0, "Telegram admin id")
	flag.Parse()

	switch env {
	case Dev, Prod:
	default:
		panic(fmt.Sprintf("Invalid env: %s (must be %s or %s)", env, Dev, Prod))
	}

	if portSocks5 == 0 {
		panic("SOCKS5 server port is empty")
	}

	if telegramBotToken == "" {
		panic("Telegram bot token is empty")
	}

	if telegramAdminId == 0 {
		panic("Telegram command id is empty")
	}

	if sqlitePath == "" {
		panic("Storage path is empty")
	}

	if _, err := os.Stat(sqlitePath); os.IsNotExist(err) {
		if err := os.MkdirAll(sqlitePath, 0o750); err != nil {
			panic(fmt.Errorf("could not create sqlite path: %w", err))
		}
	}

	return &Config{
		env:              env,
		portSocks5:       portSocks5,
		telegramBotToken: telegramBotToken,
		telegramAdminId:  telegramAdminId,
		sqlitePath:       sqlitePath,
	}
}
