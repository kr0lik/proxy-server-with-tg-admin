package config

import (
	"context"
	"flag"
	"fmt"
	"os"
	"proxy-server-with-tg-admin/internal/helper"
	"strings"
)

const (
	EnvDev            = "dev"
	EnvProd           = "prod"
	defaultSocks5Port = 1080
)

type Config struct {
	env              string
	ip               string
	portSocks5       uint
	telegramBotToken string
	telegramAdminId  int64
	dataPath         string
}

func (c *Config) Env() string {
	return c.env
}

func (c *Config) Ip() string {
	return c.ip
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
	var ip string
	var portSocks5 uint
	var env, dataPath string
	var telegramBotToken string
	var telegramAdminId int64

	flag.StringVar(&ip, "ip", "", "Server ip")
	flag.UintVar(&portSocks5, "port-socks5", defaultSocks5Port, "SOCKS5 server port")
	flag.StringVar(&env, "env", EnvProd, "Application environment: EnvDev or EnvProd")
	flag.StringVar(&dataPath, "data-path", "./.data", "Storage path")
	flag.StringVar(&telegramBotToken, "telegram-bot-token", "", "Telegram bot token")
	flag.Int64Var(&telegramAdminId, "telegram-admin-id", 0, "Telegram admin id")
	flag.Parse()

	ip = getOrDetectIP(ip)
	validateConfigParams(env, ip, portSocks5, telegramBotToken, telegramAdminId, dataPath)
	ensureDataPathExists(dataPath)

	return &Config{
		env:              env,
		ip:               ip,
		portSocks5:       portSocks5,
		telegramBotToken: telegramBotToken,
		telegramAdminId:  telegramAdminId,
		dataPath:         dataPath,
	}
}

func getOrDetectIP(ip string) string {
	if ip != "" {
		return ip
	}
	myIp, err := helper.GetMyIp(context.Background())
	if err != nil {
		panic("Server ip: " + err.Error())
	}
	if myIp == "" {
		panic("Server ip is empty")
	}

	return myIp
}

func validateConfigParams(env, ip string, portSocks5 uint, telegramBotToken string, telegramAdminId int64, dataPath string) {
	switch env {
	case EnvDev, EnvProd:
		// ok
	default:
		panic(fmt.Sprintf("Invalid env: %s (must be %s or %s)", env, EnvDev, EnvProd))
	}
	if ip == "" {
		panic("Server ip is empty")
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
}

func ensureDataPathExists(dataPath string) {
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		const folderPerm = 0o750
		if err := os.MkdirAll(dataPath, folderPerm); err != nil {
			panic(fmt.Errorf("could not create sqlite path: %w", err))
		}
	}
}
