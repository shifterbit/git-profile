# Git Profile
A git command to manage and reuse repository configurations

## Example Profile

```toml
[gpg]
format = "ssh"
[commit]
gpgsign = true

[user]
name = "johndoe"
email = "johndoe@example.com"
signingkey = "path to signing key" # all fields can be ommitted

```