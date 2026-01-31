package opengist

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
)

type Client struct {
	opengistURL      string
	opengistUsername string
	opengistToken    string
	workDir          string
	logger           *logrus.Logger
}

func NewClient(opengistURL, opengistUsername, opengistToken, workDir string, logger *logrus.Logger) *Client {
	return &Client{
		opengistURL:      strings.TrimSuffix(opengistURL, "/"),
		opengistUsername: opengistUsername,
		opengistToken:    opengistToken,
		workDir:          workDir,
		logger:           logger,
	}
}

func (c *Client) SyncGist(gistID, githubURL, description string) error {
	repoDir := filepath.Join(c.workDir, gistID)

	if err := os.MkdirAll(c.workDir, 0755); err != nil {
		return fmt.Errorf("failed to create work directory: %w", err)
	}

	repo, err := c.cloneOrOpen(repoDir, githubURL)
	if err != nil {
		return err
	}

	if err := c.pushToOpengist(repo, gistID, description); err != nil {
		return err
	}

	if err := c.cleanup(repoDir); err != nil {
		c.logger.WithError(err).Warnf("Failed to cleanup repo %s", gistID)
	}

	return nil
}

func (c *Client) cloneOrOpen(repoDir, githubURL string) (*git.Repository, error) {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		c.logger.Infof("Cloning gist from %s", githubURL)
		repo, err := git.PlainClone(repoDir, false, &git.CloneOptions{
			URL:      githubURL,
			Progress: nil,
			Auth: &http.BasicAuth{
				Username: "gist-sync",
				Password: "", // GitHub gists are public or use token in URL
			},
			Depth: 1,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to clone gist: %w", err)
		}
		return repo, nil
	}

	c.logger.Infof("Opening existing repo at %s", repoDir)
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo: %w", err)
	}

	c.logger.Infof("Pulling latest changes from %s", githubURL)
	w, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	if err := w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "gist-sync",
			Password: "",
		},
	}); err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, fmt.Errorf("failed to pull: %w", err)
	}

	return repo, nil
}

func (c *Client) pushToOpengist(repo *git.Repository, gistID, description string) error {
	opengistRemoteURL := fmt.Sprintf("%s/%s/%s", c.opengistURL, c.opengistUsername, gistID)

	remotes, err := repo.Remotes()
	if err != nil {
		return fmt.Errorf("failed to get remotes: %w", err)
	}

	opengistRemoteExists := false
	for _, remote := range remotes {
		if remote.Config().Name == "opengist" {
			opengistRemoteExists = true
			break
		}
	}

	if !opengistRemoteExists {
		c.logger.Infof("Adding Opengist remote: %s", opengistRemoteURL)
		_, err := repo.CreateRemote(&config.RemoteConfig{
			Name: "opengist",
			URLs: []string{opengistRemoteURL},
		})
		if err != nil {
			return fmt.Errorf("failed to create remote: %w", err)
		}
	}

	c.logger.Infof("Pushing to Opengist: %s", opengistRemoteURL)

	pushOptions := &git.PushOptions{
		RemoteName: "opengist",
		Auth: &http.BasicAuth{
			Username: c.opengistUsername,
			Password: c.opengistToken,
		},
		Force: true,
	}

	if err := repo.Push(pushOptions); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			c.logger.Infof("Gist %s already up to date", gistID)
			return nil
		}
		return fmt.Errorf("failed to push to opengist: %w", err)
	}

	time.Sleep(100 * time.Millisecond)

	c.logger.Infof("Successfully synced gist %s", gistID)
	return nil
}

func (c *Client) cleanup(repoDir string) error {
	return os.RemoveAll(repoDir)
}
