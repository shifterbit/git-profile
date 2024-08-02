package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

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

type CLIOptions struct {
	Local    bool
	Worktree bool
	Profile  string
}

func main() {
	ok, _ := ReadConfig("profile.toml")
	fmt.Printf("%#v\n", ok.User)
	fmt.Printf("%#v\n", ok.GPG)
	fmt.Printf("%#v\n", ok.Commit)
	_ = parse_options()
	profiles_dir := find_profiles_directory()
	profiles := collect_profiles(profiles_dir)
	profile_list := list_profiles(profiles)
	fmt.Printf("%#v\n", profiles)
	fmt.Println(profile_list)
}

func print_profiles(profile_list []string) {
	for _, v := range profile_list {
		fmt.Println(v)
	}

}
func list_profiles(profiles map[string]GitConfig) []string {
	profile_list := []string{}
	for name := range profiles {
		profile_list = append(profile_list, name)
	}
	return profile_list
}

func collect_profiles(profiles_dir string) map[string]GitConfig {
	entries, err := os.ReadDir(profiles_dir)
	if err != nil {
		log.Fatal(err)
	}

	profiles := map[string]GitConfig{}
	for _, v := range entries {

		if !v.IsDir() && contains_str(filepath.Ext(v.Name()), []string{".toml", ".tml"}) {
			profile_name := removeExtension(filepath.Base(v.Name()))
			_, ok := profiles[profile_name]
			if ok {
				log.Fatal("Found duplicate profiles!")
			}
			profile_config, err := ReadConfig(profiles_dir + v.Name())
			if err != nil {
				log.Fatal(err)
			}
			profiles[profile_name] = profile_config
			fmt.Println(profile_config)
		}

	}
	return profiles
}

func removeExtension(filename string) string {
	extension := filepath.Ext(filename)
	return filename[0 : len(filename)-len(extension)]
}

func contains_str(str string, arr []string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}

	}

	return false

}
func find_profiles_directory() string {
	curr_os := runtime.GOOS
	fmt.Println(curr_os)
	dir, ok := os.LookupEnv("GIT_PROFILES_DIR")

	if curr_os == "windows" {
		if ok && dir != "" && dir[len(dir)-1] != '\\' {
			return dir + "\\"
		} else {
			return "%USERPROFILE%\\AppData\\Local\\git-profile\\"
		}
	} else {
		if ok && dir != "" && dir[len(dir)-1] != '/' {
			return dir + "/"
		} else {
			return "~/.config/git-profile/"
		}
	}
}

func parse_options() CLIOptions {
	local := flag.Bool("local", true, "apply configuration to local repo")
	worktree := flag.Bool("worktree", true, "apply configuration to worktree")
	profile := flag.String("set-profile", "", "the profile to set")

	flag.Parse()
	return CLIOptions{
		Local:    *local,
		Worktree: *worktree,
		Profile:  *profile,
	}
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
