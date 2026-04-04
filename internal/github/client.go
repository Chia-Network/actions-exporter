package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v84/github"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var (
	client *github.Client
	cash   *cache.Cache
)

// Init constructs a github API client for this package taking in a token
func Init(token string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	// init cache, defaults 1 hour expiration and purges expired items every 10 minutes
	cash = cache.New(1*time.Hour, 10*time.Minute)
}

// GetRateLimits calls the GET /rate_limit endpoint, which does not count against the rate limit
func GetRateLimits() (*github.RateLimits, error) {
	limits, _, err := client.RateLimit.Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("GetRateLimits returned error: %v", err)
	}
	return limits, nil
}

// ListRepositoriesByOrg gets all repositories in a GitHub organization
func ListRepositoriesByOrg(org string) ([]*github.Repository, error) {
	var repos []*github.Repository
	if x, found := cash.Get("ListRepositoriesByOrg"); found {
		var ok bool
		repos, ok = x.([]*github.Repository)
		if ok {
			log.Debug("github client request was found in cache for ListRepositoriesByOrg")
			return repos, nil
		}
	}

	var page, perPage int = 1, 100
	for {
		r, _, err := client.Repositories.ListByOrg(context.Background(), org, &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("ListRepositoriesByOrg returned error: \n%v", err)
		}

		repos = append(repos, r...)

		if len(r) != perPage {
			break
		}
		page++
	}

	cash.Set("ListRepositoriesByOrg", repos, cache.DefaultExpiration)

	return repos, nil
}

// ListRepositoriesByUser gets all repositories for a GitHub user
func ListRepositoriesByUser(user string) ([]*github.Repository, error) {
	var repos []*github.Repository
	if x, found := cash.Get("ListRepositoriesByUser"); found {
		var ok bool
		repos, ok = x.([]*github.Repository)
		if ok {
			log.Debug("github client request was found in cache for ListRepositoriesByUser")
			return repos, nil
		}
	}

	var page, perPage int = 1, 100
	for {
		r, _, err := client.Repositories.ListByUser(context.Background(), user, &github.RepositoryListByUserOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("ListRepositoriesByUser returned error: \n%v", err)
		}

		repos = append(repos, r...)

		if len(r) != perPage {
			break
		}
		page++
	}

	cash.Set("ListRepositoriesByUser", repos, cache.DefaultExpiration)

	return repos, nil
}

// ListRepositoryWorkflows gets all workflows for a GitHub repository
func ListRepositoryWorkflows(owner, repo string) ([]*github.Workflow, error) {
	var workflows []*github.Workflow
	if x, found := cash.Get(fmt.Sprintf("ListRepositoryWorkflows_%s_%s", owner, repo)); found {
		var ok bool
		workflows, ok = x.([]*github.Workflow)
		if ok {
			log.Debug("github client request was found in cache for ListRepositoryWorkflows")
			return workflows, nil
		}
	}

	var page, perPage int = 1, 100
	for {
		r, _, err := client.Actions.ListWorkflows(context.Background(), owner, repo, &github.ListOptions{
			Page:    page,
			PerPage: perPage,
		})
		if err != nil {
			return nil, fmt.Errorf("ListRepositoryWorkflows returned error: \n%v", err)
		}

		workflows = append(workflows, r.Workflows...)

		if len(r.Workflows) != perPage {
			break
		}
		page++
	}

	cash.Set(fmt.Sprintf("ListRepositoryWorkflows_%s_%s", owner, repo), workflows, cache.DefaultExpiration)

	return workflows, nil
}
