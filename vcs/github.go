package vcs

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const defaultHookName = "thrap"

// github repo extending git
type githubVCS struct {
	git    *GitVCS
	client *github.Client
}

// takes an optional underlying git vcs
func newGithubVCS(g *GitVCS) *githubVCS {
	gh := &githubVCS{
		git: g,
	}

	// Initialize git backend
	if gh.git == nil {
		gh.git = NewGitVCS()
	}

	return gh
}

func (gh *githubVCS) ID() string {
	return "github"
}

// Environment Variables:
// GITHUB_ACCESS_TOKEN
func (gh *githubVCS) Init(conf map[string]interface{}) error {
	// Init base git interface
	err := gh.git.Init(conf)
	if err == nil {
		var token string
		if iface, ok := conf["token"]; ok {
			token = iface.(string)
			if token == "" {
				return errors.New("'token' must be a string")
			}
		}

		httpClient := makeGithubHTTPClient(token)
		gh.client = github.NewClient(httpClient)
	}

	return err
}

func (gh *githubVCS) DefaultUser() string {
	return gh.git.defaultUser
}

func (gh *githubVCS) DefaultEmail() string {
	return gh.git.defaultEmail
}

func (gh *githubVCS) Get(repo *Repository, opt Option) (interface{}, error) {
	ctx := context.Background()

	// Create only if it does not exist
	ghRepo, _, err := gh.client.Repositories.Get(ctx, repo.Owner, repo.Name)
	return ghRepo, err
}

func (gh *githubVCS) AddHook(repo *Repository) (interface{}, error) {
	rs := gh.client.Repositories

	ctx := context.Background()

	hookName := defaultHookName
	hook := &github.Hook{
		Name: &hookName,
	}

	rhook, _, err := rs.CreateHook(ctx, repo.Owner, repo.Name, hook)
	return rhook, err
}

func (gh *githubVCS) RemoveHook(repo *Repository) (interface{}, error) {
	return nil, errors.New("to be implemented")
}

// Create creates a new repo. Each call only fills in missing pieces so multiple
// calls will not corrupt
func (gh *githubVCS) Create(repo *Repository, opt Option) (interface{}, error) {
	ctx := context.Background()

	// TODO: handle user vs org

	// Create only if it does not exist
	ghRepo, _, err := gh.client.Repositories.Get(ctx, repo.Owner, repo.Name)
	if err == nil {
		return ghRepo, nil
	}

	newRepo := &github.Repository{
		Name:        &repo.Name,
		Private:     &repo.Private,
		Description: &repo.Description,
	}

	// Owner defaults to the user if not specified
	ghrepo, _, err := gh.client.Repositories.Create(ctx, repo.Owner, newRepo)
	return ghrepo, err
}

// Delete deletes the specified repo from github
func (gh *githubVCS) Delete(repo *Repository, opt Option) error {
	ctx := context.Background()
	_, err := gh.client.Repositories.Delete(ctx, repo.Owner, repo.Name)
	return err
}

func (gh *githubVCS) IgnoresFile() string {
	return gh.git.IgnoresFile()
}

func makeGithubHTTPClient(token string) *http.Client {
	var (
		httpClient *http.Client
		ghtoken    = os.Getenv("GITHUB_ACCESS_TOKEN")
	)

	if token != "" {
		ghtoken = token
	}

	if ghtoken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ghtoken},
		)
		httpClient = oauth2.NewClient(context.Background(), ts)
	}

	return httpClient
}
