package gh

import (
	"context"
	"github.com/fzdwx/code-github-workspace/config"
	"github.com/google/go-github/v46/github"
	"golang.org/x/oauth2"
)

func GetAuth(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Get().Token},
	)

	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func Get() *github.Client {
	return github.NewClient(nil)
}
