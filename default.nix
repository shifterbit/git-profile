{buildGoModule, ...}: let
  name = "git-profile";
in
  buildGoModule {
    pname = name;
    version = "0.0.1";
    src = ./.;
    vendorHash = "sha256-/EGXT9nQjvVi9/8OIKb+GZszRQx8qANznTWlqK9RX2w=";
  }
