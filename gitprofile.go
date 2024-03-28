package gitprofile

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
    "dario.cat/mergo"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func MergeProfileIntoRepo(name string, profileDirectory string, repoPath string) (*config.Config, error) {
	profile, err := readConfigProfile(name, profileDirectory)
	if err != nil {
		return nil, err
	}
	repoConfig, err := readLocalRepoConfig(repoPath)
	if err != nil {
		return nil, err
	}
	var finalConfig *config.Config
	mergeCore(profile, repoConfig)
	mergeUser(profile, repoConfig)
	mergeCommitter(profile, repoConfig)
	mergeAuthor(profile, repoConfig)
	mergeInit(profile, repoConfig)
	mergeRemotes(profile, repoConfig)
	mergeURLS(profile, repoConfig)
	mergeRemotes(profile, repoConfig)
	mergeBranches(profile, repoConfig)
	mergeURLS(profile, repoConfig)

	return finalConfig, nil

}

func mergeCore(src *config.Config, dest *config.Config) {
	dest.Core.IsBare = src.Core.IsBare
	if src.Core.Worktree != "" {
		dest.Core.Worktree = src.Core.Worktree
	}
	if src.Core.CommentChar != "" {
		dest.Core.CommentChar = src.Core.CommentChar
	}
	if src.Core.RepositoryFormatVersion != "" {
		dest.Core.RepositoryFormatVersion = src.Core.RepositoryFormatVersion
	}
	mergo.Merge(src.Raw, dest.Raw, mergo.WithOverride)
}

func mergeUser(src *config.Config, dest *config.Config) {
	if src.User.Name != "" {
		dest.User.Name = src.User.Name
	}
	if src.User.Email != "" {
		dest.User.Email = src.User.Email
	}
	mergo.Merge(src.Raw, dest.Raw, mergo.WithOverride)
}

func mergeCommitter(src *config.Config, dest *config.Config) {
	if src.Committer.Name != "" {
		dest.Committer.Name = src.Committer.Name
	}
	if src.Committer.Email != "" {
		dest.Committer.Email = src.Committer.Email
	}
	mergo.Merge(src.Raw, dest.Raw, mergo.WithOverride)
}

func mergeAuthor(src *config.Config, dest *config.Config) {
	if src.Author.Name != "" {
		dest.Author.Name = src.Author.Name
	}
	if src.Author.Email != "" {
		dest.Author.Email = src.Author.Email
	}
	mergo.Merge(src.Raw, dest.Raw, mergo.WithOverride)
}

func mergeInit(src *config.Config, dest *config.Config) {
	if src.Init.DefaultBranch != "" {
		dest.Init.DefaultBranch = src.Init.DefaultBranch
	}
	mergo.Merge(src.Raw, dest.Raw, mergo.WithOverride)
}

func mergeRemotes(src *config.Config, dest *config.Config) {
	if src.Remotes != nil {
		dest.Remotes = src.Remotes
	}
	mergo.Merge(src.Raw, dest.Raw, mergo.WithOverride)
}

func mergeBranches(src *config.Config, dest *config.Config) {
	if src.Branches != nil {
		dest.Branches = src.Branches
	}
	mergo.Merge(src.Raw, dest.Raw, mergo.WithOverride)
}

func mergeURLS(src *config.Config, dest *config.Config) {
	if src.URLs != nil {
		dest.URLs = src.URLs
	}
	mergo.Merge(src.Raw, dest.Raw, mergo.WithOverride)
}


func readLocalRepoConfig(path string) (*config.Config, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	config, err := repo.ConfigScoped(config.LocalScope)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func readConfigProfile(name string, profileDirectory string) (*config.Config, error) {
	files, err := os.ReadDir(profileDirectory)
	if err != nil {
		return nil, err
	}
	var configFile []byte

	for _, file := range files {
		if file.Name() == name {
			fullPath := filepath.Join(profileDirectory, name)
			configFile, err = os.ReadFile(fullPath)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if len(configFile) == 0 {
		return nil, errors.New("Config file not found")
	}
	r := bytes.NewReader(configFile)
	return config.ReadConfig(r)
}
