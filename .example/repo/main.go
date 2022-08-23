package main

import (
	"context"
	"fmt"
	"github.com/fzdwx/gh-sp/api"
	"github.com/fzdwx/gh-sp/model/repolist"
	"github.com/google/go-github/v46/github"
	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()
	ops := &github.RepositoryListOptions{
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{},
	}
	repos, resp, err := api.Get(ctx).Repositories.List(ctx, "", ops)

	cobra.CheckErr(err)

	items := repolist.NewItems(repos, resp)

	fmt.Println(items.View(166, 15))
}
