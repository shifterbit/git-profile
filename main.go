package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

type GitConfig struct {
	User   *GitUser
	Commit *GitCommitOptions
	GPG    *GitGPGOptions
}

type GitGPGOptions struct {
	Format string
}

type GitCommitOptions struct {
	GPGSign bool
}
type GitUser struct {
	Name       string
	Email      string
	SigningKey string
}

func main() {
	ok, _ :=  ReadConfig("profile.toml")
	fmt.Printf("%#v\n", ok.User)
	fmt.Printf("%#v\n", ok.GPG)
	fmt.Printf("%#v\n", ok.Commit)
}


func ReadConfig(path string) (GitConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return GitConfig{}, err
	}
	gitConfig := GitConfig{}
	err = toml.Unmarshal(file, &gitConfig)
	if err != nil {
		return GitConfig{}, err
	}

	return gitConfig, err

}

func ApplyProfile(config GitConfig) {
	setUser(config.User)
	setCommitOptions(config.Commit)
	setGPGFormat(config.GPG)
}

func setGPGFormat(options *GitGPGOptions) {
	if options == nil {
		return
	}
	cmd := exec.Command("git", "config", "--local", "gpg.format", options.Format)
	if options.Format != "" {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

}
func setCommitOptions(options *GitCommitOptions) {
	if options == nil {
		return
	}
	cmd := exec.Command("git", "config", "--local", "commit.gpgsign", strconv.FormatBool(options.GPGSign))
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

}
func setUser(user *GitUser) {
	if user == nil {
		return
	}
	cmd := exec.Command("git", "config", "--local", "user.name", user.Name)
	if user.Name != "" {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
	cmd = exec.Command("git", "config", "--local", "user.email", user.Email)
	if user.Email != "" {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
	cmd = exec.Command("git", "config", "--local", "user.signingKey", user.SigningKey)
	if user.SigningKey != "" {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

}
