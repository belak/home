{
  lib,
  stdenv,
  fetchurl,
  system,
  ...
}:
let
  versions = lib.importJSON ./templ-bin-versions.json;
  versionInfo = versions."${system}";
  inherit (versionInfo) version sha256 url;
in
stdenv.mkDerivation (finalAttrs: {
  pname = "templ-bin";
  inherit version;

  src = fetchurl { inherit sha256 url; };

  installPhase = ''
    runHook preInstall
    install -D templ $out/bin/templ
    runHook postInstall
  '';

  sourceRoot = ".";
})
