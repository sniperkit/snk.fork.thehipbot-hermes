/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"

	"github.com/sniperkit/snk.fork.thehipbot-hermes/cache"
	"github.com/sniperkit/snk.fork.thehipbot-hermes/fs"
	mock_prompt "github.com/sniperkit/snk.fork.thehipbot-hermes/mock"
)

var (
	cmd       = &cobra.Command{}
	cacheFile billy.File
)

type RootCmdSuite struct {
	suite.Suite
}

func (s *RootCmdSuite) SetupTest() {
	configFS = &fs.ConfigFS{
		FS: memfs.New(),
	}
	configFS.Setup()
	cacheFile, _ = configFS.GetCacheFile()
	fsCache = cache.NewCache(cacheFile)
}

func (s *RootCmdSuite) TearDownSuite() {
	fsCache.Close()
}

func (s *RootCmdSuite) TestGetHandlerSingleCachedRepo() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	cacheFile.Seek(0, 0)
	p, _ := cacheFile.Write([]byte(`{
		"version": "0.0.1",
		"remotes": {
			"github.com": {
				"name": "github.com",
				"url":  "https://github.com",
				"repos": [
					{
						"name": "github.com/sniperkit/snk.fork.thehipbot-hermes",
						"repo_path": "/repos/github.com/sniperkit/snk.fork.thehipbot-hermes"
					},
					{
						"name": "github.com/TheHipbot/dotfiles",
						"repo_path": "/repos/github.com/TheHipbot/dotfiles"
					},
					{
						"name": "github.com/TheHipbot/dockerfiles",
						"repo_path": "/repos/github.com/TheHipbot/dockerfiles"
					},
					{
						"name": "github.com/src-d/go-git",
						"repo_path": "/repos/github.com/src-d/go-git"
					}
				]
			},
			"gitlab.com": {
				"name": "gitlab.com",
				"url":  "https://gitlab.com",
				"repos": [
					{
						"name": "gitlab.com/gitlab-org/gitlab-ce",
						"repo_path": "/repos/gitlab.com/gitlab-org/gitlab-ce"
					},
					{
						"name": "gitlab.com/gnachman/iterm2",
						"repo_path": "/repos/gitlab.com/gnachman/iterm2"
					}
				]
			}
		}
	}`))
	cacheFile.Truncate(int64(p))
	fsCache.Open()

	mockPrompter := mock_prompt.NewMockFactory(ctrl)
	mockPrompter.
		EXPECT().
		CreateSelectPrompt(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)
	prompter = mockPrompter

	getHandler(cmd, []string{"github.com/sniperkit/snk.fork.thehipbot-hermes"})
	target := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("target_file"))
	stat, _ := configFS.FS.Stat(target)
	targetFile, err := configFS.FS.Open(target)
	defer targetFile.Close()
	content := make([]byte, stat.Size())
	targetFile.Read(content)
	s.Nil(err, "Target file should exist")
	s.Equal(string(content), "/repos/github.com/sniperkit/snk.fork.thehipbot-hermes", "Get should find one repo and set target path")
}

func (s *RootCmdSuite) TestGetHandlerMultipleCachedRepos() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	cacheFile.Seek(0, 0)
	p, _ := cacheFile.Write([]byte(`{
		"version": "0.0.1",
		"remotes": {
			"github.com": {
				"name": "github.com",
				"url":  "https://github.com",
				"repos": [
					{
						"name": "github.com/sniperkit/snk.fork.thehipbot-hermes",
						"repo_path": "/repos/github.com/sniperkit/snk.fork.thehipbot-hermes"
					},
					{
						"name": "github.com/TheHipbot/dotfiles",
						"repo_path": "/repos/github.com/TheHipbot/dotfiles"
					},
					{
						"name": "github.com/TheHipbot/dockerfiles",
						"repo_path": "/repos/github.com/TheHipbot/dockerfiles"
					},
					{
						"name": "github.com/src-d/go-git",
						"repo_path": "/repos/github.com/src-d/go-git"
					}
				]
			},
			"gitlab.com": {
				"name": "gitlab.com",
				"url":  "https://gitlab.com",
				"repos": [
					{
						"name": "gitlab.com/gitlab-org/gitlab-ce",
						"repo_path": "/repos/gitlab.com/gitlab-org/gitlab-ce"
					},
					{
						"name": "gitlab.com/gnachman/iterm2",
						"repo_path": "/repos/gitlab.com/gnachman/iterm2"
					}
				]
			}
		}
	}`))
	cacheFile.Truncate(int64(p))
	fsCache.Open()
	repos := []cache.Repo{
		cache.Repo{
			Name: "github.com/sniperkit/snk.fork.thehipbot-hermes",
			Path: "/repos/github.com/sniperkit/snk.fork.thehipbot-hermes",
		},
		cache.Repo{
			Name: "github.com/TheHipbot/dotfiles",
			Path: "/repos/github.com/TheHipbot/dotfiles",
		},
		cache.Repo{
			Name: "github.com/TheHipbot/dockerfiles",
			Path: "/repos/github.com/TheHipbot/dockerfiles",
		},
	}

	mockPrompter := mock_prompt.NewMockFactory(ctrl)
	mockPrompt := mock_prompt.NewMockPrompt(ctrl)
	mockPrompt.
		EXPECT().
		Run().
		Return(1, "blah", nil).
		Times(1)

	mockPrompter.
		EXPECT().
		CreateSelectPrompt(gomock.Any(), gomock.Eq(repos), gomock.Any()).
		Return(mockPrompt).
		Times(1)

	prompter = mockPrompter

	getHandler(cmd, []string{"hipbot"})
	target := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("target_file"))
	stat, _ := configFS.FS.Stat(target)
	targetFile, err := configFS.FS.Open(target)
	defer targetFile.Close()
	content := make([]byte, stat.Size())
	targetFile.Read(content)
	s.Nil(err, "Target file should exist")
	s.Equal(string(content), "/repos/github.com/TheHipbot/dotfiles", "Get should find one repo and set target path")
}

func TestRootCmdSuite(t *testing.T) {
	suite.Run(t, new(RootCmdSuite))
}
