package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	ctx    context.Context
}

type Gist struct {
	ID          string
	Description string
	Public      bool
	GitPullURL  string
	GitPushURL  string
	HTMLURL     string
}

func NewClient(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client: github.NewClient(tc),
		ctx:    ctx,
	}
}

func (c *Client) GetAllGists() ([]*Gist, error) {
	var allGists []*Gist
	opts := &github.GistListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		gists, resp, err := c.client.Gists.List(c.ctx, "", opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list gists: %w", err)
		}

		for _, gist := range gists {
			allGists = append(allGists, c.convertGist(gist))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allGists, nil
}

func (c *Client) GetGist(gistID string) (*Gist, error) {
	gist, _, err := c.client.Gists.Get(c.ctx, gistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gist %s: %w", gistID, err)
	}

	return c.convertGist(gist), nil
}

func (c *Client) convertGist(ghGist *github.Gist) *Gist {
	return &Gist{
		ID:          ghGist.GetID(),
		Description: ghGist.GetDescription(),
		Public:      ghGist.GetPublic(),
		GitPullURL:  ghGist.GetGitPullURL(),
		GitPushURL:  ghGist.GetGitPushURL(),
		HTMLURL:     ghGist.GetHTMLURL(),
	}
}
