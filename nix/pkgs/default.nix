{
  pkgs ? import <nixpkgs> { },
}:
{
  templ-bin = pkgs.callPackage ./templ-bin.nix { };
  sqlc-bin = pkgs.callPackage ./sqlc-bin.nix { };

}
