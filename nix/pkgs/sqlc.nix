{
  lib,
  buildGoModule,
  fetchFromGitHub,
}:
let
  version = "1.27.0";
in
buildGoModule {
  pname = "sqlc";
  inherit version;

  src = fetchFromGitHub {
    owner = "sqlc-dev";
    repo = "sqlc";
    rev = "v${version}";
    hash = "sha256-wxQ+YPsDX0Z6B8whlQ/IaT2dRqapPL8kOuFEc6As1rU=";
  };

  # We need proxyVendor because otherwise it will be missing some header files
  # on build.
  proxyVendor = true;
  vendorHash = "sha256-ndOw3uShF5TngpxYNumoK3H3R9v4crfi5V3ZCoSqW90=";

  subPackages = [ "cmd/sqlc" ];

  ldflags = [
    "-s"
    "-w"
  ];
}
