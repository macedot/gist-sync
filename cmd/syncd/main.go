package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/macedot/gist-sync/internal/config"
	"github.com/macedot/gist-sync/internal/sync"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		if lvl, err := logrus.ParseLevel(level); err == nil {
			logger.SetLevel(lvl)
		}
	}

	logger.Info("Starting gist-sync daemon")

	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	logger.Infof("Configuration loaded - syncing from GitHub user: %s to Opengist: %s", cfg.GitHubUsername, cfg.OpengistURL)

	syncer := sync.NewSyncer(cfg, logger)

	done := make(chan bool)
	go func() {
		syncer.Start()
		done <- true
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		logger.Infof("Received signal %s, shutting down...", sig)
	case <-done:
		logger.Info("Syncer stopped")
	}
}
