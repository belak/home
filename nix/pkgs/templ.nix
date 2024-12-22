{
  lib,
  buildGoModule,
  fetchFromGitHub,
}:

let
  version = "0.2.793";
in
buildGoModule {
  pname = "templ";
  inherit version;

  src = fetchFromGitHub {
    owner = "a-h";
    repo = "templ";
    rev = "v${version}";
    hash = "sha256-0KGht5IMbJV8KkXgT5qJxA9bcmWevzXXAVPMQTm0ccw=";
  };

  vendorHash = "sha256-ZWY19f11+UI18jeHYIEZjdb9Ii74mD6w+dYRLPkdfBU=";

  ldflags = [
    "-s"
    "-w"
  ];

  subPackages = [ "cmd/templ" ];
}
