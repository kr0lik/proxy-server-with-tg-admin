package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"proxy-server-with-tg-admin/internal/config"
	"proxy-server-with-tg-admin/internal/infrastructure/adblock"
	"proxy-server-with-tg-admin/internal/infrastructure/socks5"
	"proxy-server-with-tg-admin/internal/infrastructure/sqlite"
	"proxy-server-with-tg-admin/internal/infrastructure/telegram"
	"proxy-server-with-tg-admin/internal/usecase/auth"
	"proxy-server-with-tg-admin/internal/usecase/commands"
	"proxy-server-with-tg-admin/internal/usecase/statistic"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env())
	logger.Debug("starting", "env", cfg.Env(), "port", cfg.PortSocks5(), "db path", cfg.DataPath())

	logger.Info("Storage starting")

	storage, err := sqlite.New(cfg.DataPath(), logger)
	if err != nil {
		logger.Error("Failed to start storage", "err", err)

		return
	}
	defer storage.Close()

	logger.Info("Ad blocker running")

	adBlock := adblock.New(logger)
	adBlock.Start()

	logger.Info("Use cases starting")

	authenticator := auth.New(storage, logger)
	cmdList := commands.New(cfg.Ip(), cfg.PortSocks5(), storage, authenticator)

	statisticTracker := statistic.New(storage, logger)
	logger.Info("Statistic tracker running")

	statisticTracker.Start()
	defer statisticTracker.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Telegram bot starting")

	tgBot, err := telegram.MakeBot(cfg.TelegramToken(), cfg.TelegramAdminId(), cmdList)
	if err != nil {
		logger.Error("Failed to start telegram bot", "err", err)

		return
	}
	defer tgBot.Stop()

	go func() {
		logger.Info("Telegram bot running")
		tgBot.Start()
		logger.Info("Telegram bot stopped")
	}()

	logger.Info("Socks5 server starting")

	socks5Server := socks5.New(statisticTracker, adBlock, authenticator, logger)
	defer socks5Server.Shutdown()

	go func() {
		serverAddr := fmt.Sprintf(":%d", cfg.PortSocks5())
		logger.Info("Socks5 server running", "on", serverAddr)

		if err := socks5Server.ListenAndServe("tcp", serverAddr); err != nil {
			logger.Error("Socks5 server", "ListenAndServe", err)

			done <- os.Interrupt

			return
		}

		logger.Info("Socks5 server stopped")
	}()

	<-done

	logger.Info("stopping")
}

func setupLogger(env string) *slog.Logger {
	opts := &slog.HandlerOptions{}

	switch env {
	case config.EnvDev:
		opts.Level = slog.LevelDebug
	case config.EnvProd:
		opts.Level = slog.LevelInfo
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
