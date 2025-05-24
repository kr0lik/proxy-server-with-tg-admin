package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	EnvDev            = "dev"
	EnvProd           = "prod"
	defaultSocks5Port = 1080
)

type Config struct {
	env              string
	portSocks5       uint
	telegramBotToken string
	telegramAdminId  int64
	dataPath         string
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

func (c *Config) DataPath() string {
	return strings.TrimRight(c.dataPath, string(os.PathSeparator))
}

func MustLoad() *Config {
	var portSocks5 uint
	var env, dataPath string
	var telegramBotToken string
	var telegramAdminId int64

	flag.UintVar(&portSocks5, "port-socks5", defaultSocks5Port, "SOCKS5 server port")
	flag.StringVar(&env, "env", EnvProd, "Application environment: EnvDev or EnvProd")
	flag.StringVar(&dataPath, "data-path", "./.data", "Storage path")
	flag.StringVar(&telegramBotToken, "telegram-bot-token", "", "Telegram bot token")
	flag.Int64Var(&telegramAdminId, "telegram-admin-id", 0, "Telegram admin id")
	flag.Parse()

	switch env {
	case EnvDev, EnvProd:
	default:
		panic(fmt.Sprintf("Invalid env: %s (must be %s or %s)", env, EnvDev, EnvProd))
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

	if dataPath == "" {
		panic("Data path is empty")
	}

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		const folderPerm = 0o750
		if err := os.MkdirAll(dataPath, folderPerm); err != nil {
			panic(fmt.Errorf("could not create sqlite path: %w", err))
		}
	}

	return &Config{
		env:              env,
		portSocks5:       portSocks5,
		telegramBotToken: telegramBotToken,
		telegramAdminId:  telegramAdminId,
		dataPath:         dataPath,
	}
}
