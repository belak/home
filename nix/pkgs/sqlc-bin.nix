{
  lib,
  stdenv,
  fetchurl,
  system,
  ...
}:
let
  versions = lib.importJSON ./sqlc-bin-versions.json;
  versionInfo = versions."${system}";
  inherit (versionInfo) version sha256 url;
in
stdenv.mkDerivation (finalAttrs: {
  pname = "sqlc-bin";
  inherit version;

  src = fetchurl { inherit sha256 url; };

  installPhase = ''
    runHook preInstall
    install -D sqlc $out/bin/sqlc
    runHook postInstall
  '';

  sourceRoot = ".";
})
