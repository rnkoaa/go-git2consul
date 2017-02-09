package mock

import (
	"os"
	"time"

	"github.com/vkorehov/go-git2consul/config"
)

// RepoConfig returns a mock Repo config object
func RepoConfig(repoUrl string) *config.Repo {
	return &config.Repo{
		Name:     "git2consul-test-local",
		Url:      repoUrl,
		Branches: []string{"master"},
		Hooks: []*config.Hook{
			{
				Type:     "polling",
				Interval: 5 * time.Second,
			},
		},
	}
}

// Config returns a mock Config object with one repository configuration
func Config(repoUrl string) *config.Config {
	return &config.Config{
		LocalStore: os.TempDir(),
		HookSvr: &config.HookSvrConfig{
			Port: 9000,
		},
		Repos: []*config.Repo{
			{
				Name:     "git2consul-test-local",
				Url:      repoUrl,
				Branches: []string{"master"},
				Hooks: []*config.Hook{
					{
						Type:     "polling",
						Interval: 5 * time.Second,
					},
				},
			},
		},
		Consul: &config.ConsulConfig{
			Address: "127.0.0.1:8500",
		},
	}
}
