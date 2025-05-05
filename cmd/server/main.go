package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"proxy-server-with-tg-admin/internal/config"
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

	logger.Debug(cfg.Env(), "socks5 port", cfg.PortSocks5(), "db path", cfg.SqlitePath())

	logger.Info("starting storage")

	storage, err := sqlite.New(cfg.SqlitePath(), "server.db", logger)
	if err != nil {
		logger.Error("Failed to start storage", "err", err)

		return
	}
	defer storage.Close()

	logger.Info("starting telegram bot")

	tgBot, err := telegram.GetTelegramBot(cfg.TelegramToken(), cfg.TelegramAdminId(), commands.New(storage))
	if err != nil {
		logger.Error("Failed to start telegram bot", "err", err)

		return
	}
	defer tgBot.Stop()

	go func() {
		logger.Info("running telegram bot")
		tgBot.Start()
		logger.Info("telegram bot stopped")
	}()

	logger.Info("starting statistic tracker")

	statisticTracker := statistic.New(storage, logger)
	logger.Info("running statistic tracker")
	statisticTracker.Start()
	defer statisticTracker.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("starting socks5 server")

	socks5Server := socks5.GetServer(statisticTracker, auth.New(storage, logger), logger)

	go func() {
		serverAddr := fmt.Sprintf(":%d", cfg.PortSocks5())
		logger.Info("running socks5 server", "on", serverAddr)

		if err := socks5Server.ListenAndServe("tcp", serverAddr); err != nil {
			logger.Error("Failed to start socks5 server", "err", err)
			os.Exit(1)
		}

		logger.Info("socks5 server stopped")
	}()

	<-done

	logger.Info("stopping")
}

func setupLogger(env string) *slog.Logger {
	opts := &slog.HandlerOptions{}

	switch env {
	case config.Dev:
		opts.Level = slog.LevelDebug
	case config.Prod:
		opts.Level = slog.LevelInfo
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
