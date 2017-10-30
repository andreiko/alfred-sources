package github

import (
	"context"
	"github.com/andreiko/alfred-sources/sources"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubSource struct {
	items  []sources.Item
	client *github.Client
}

func (src *GithubSource) Query(query string) []sources.Item {
	return sources.Query(src.items, query)
}

func (src *GithubSource) Id() string {
	return "github"
}

func (src *GithubSource) Update() error {
	next_page := 1
	items := make([]sources.Item, 0)
	for next_page > 0 {
		list_opts := &github.RepositoryListOptions{ListOptions: github.ListOptions{Page: next_page}}
		repos, resp, err := src.client.Repositories.List(context.Background(), "", list_opts)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			items = append(items, NewGithubItem(repo.PushedAt.String(), *repo.Name, *repo.Owner.Login))
		}

		next_page = resp.NextPage
	}

	src.items = items

	return nil
}

func NewGithubSource(token string) *GithubSource {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	return &GithubSource{
		items:  make([]sources.Item, 0),
		client: client,
	}
}
