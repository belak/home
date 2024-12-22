{
  lib,
  buildGoModule,
  fetchFromGitHub,
}:
let
  version = "4.18.1";
in
buildGoModule rec {
  pname = "golang-migrate";
  inherit version;

  src = fetchFromGitHub {
    owner = "golang-migrate";
    repo = "migrate";
    rev = "v${version}";
    sha256 = "sha256-ZZeurnoFcObrK75zkIZvz9ycdDP9AM3uX6h/4bMWpGc=";
  };

  proxyVendor = true;
  vendorHash = "sha256-Zaq88oF5rSCSv736afyKDvTNCSIyrIGTN0kuJWqS7tg=";

  subPackages = [ "cmd/migrate" ];

  tags = [ "pgx5" ];

  ldflags = [
    "-s" "-w"
    "-X main.Version=v${version}-belak"
  ];
}
