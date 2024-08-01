package main

import (
	"fmt"
	"git-profile"
)


func main() {
	ok, _ :=  gitprofile.ReadConfig("profile.toml")
	fmt.Printf("%#v\n", ok.User)
	fmt.Printf("%#v\n", ok.GPG)
	fmt.Printf("%#v\n", ok.Commit)
}