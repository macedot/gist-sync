package sync

import (
	"time"

	"github.com/macedot/gist-sync/internal/config"
	"github.com/macedot/gist-sync/internal/github"
	"github.com/macedot/gist-sync/internal/opengist"
	"github.com/sirupsen/logrus"
)

type Syncer struct {
	config         *config.Config
	githubClient   *github.Client
	opengistClient *opengist.Client
	logger         *logrus.Logger
	lastSync       map[string]time.Time
}

func NewSyncer(cfg *config.Config, logger *logrus.Logger) *Syncer {
	return &Syncer{
		config:         cfg,
		githubClient:   github.NewClient(cfg.GitHubToken),
		opengistClient: opengist.NewClient(cfg.OpengistURL, cfg.OpengistUsername, cfg.OpengistToken, cfg.WorkDir, logger),
		logger:         logger,
		lastSync:       make(map[string]time.Time),
	}
}

func (s *Syncer) Run() error {
	s.logger.Info("Starting sync...")

	gists, err := s.githubClient.GetAllGists()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get gists from GitHub")
		return err
	}

	s.logger.Infof("Found %d gists to sync", len(gists))

	for _, gist := range gists {
		s.logger.Infof("Syncing gist: %s - %s", gist.ID, gist.Description)

		if err := s.opengistClient.SyncGist(gist.ID, gist.GitPullURL, gist.Description); err != nil {
			s.logger.WithError(err).Errorf("Failed to sync gist %s", gist.ID)
			continue
		}

		s.lastSync[gist.ID] = time.Now()
	}

	s.logger.Info("Sync completed")
	return nil
}

func (s *Syncer) Start() {
	s.logger.Info("Starting sync daemon")

	if err := s.Run(); err != nil {
		s.logger.WithError(err).Error("Initial sync failed")
	}

	ticker := time.NewTicker(s.config.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := s.Run(); err != nil {
			s.logger.WithError(err).Error("Scheduled sync failed")
		}
	}
}
