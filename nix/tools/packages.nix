{ pkgs }:

let
  lib = pkgs.lib;
  system = pkgs.stdenv.hostPlatform.system;
  supports =
    name:
    lib.elem system
      {
        summarize = [
          "aarch64-darwin"
          "x86_64-linux"
          "aarch64-linux"
        ];
        discrawl = [
          "aarch64-darwin"
          "x86_64-linux"
          "aarch64-linux"
        ];
        wacrawl = [
          "aarch64-darwin"
          "x86_64-linux"
          "aarch64-linux"
        ];
        gogcli = [
          "aarch64-darwin"
          "x86_64-linux"
          "aarch64-linux"
        ];
        goplaces = [
          "aarch64-darwin"
          "x86_64-linux"
          "aarch64-linux"
        ];
        camsnap = [
          "aarch64-darwin"
          "x86_64-linux"
          "aarch64-linux"
        ];
        sonoscli = [
          "aarch64-darwin"
          "x86_64-linux"
          "aarch64-linux"
        ];
        peekaboo = [ "aarch64-darwin" ];
        poltergeist = [ "aarch64-darwin" ];
        sag = [
          "aarch64-darwin"
          "x86_64-linux"
        ];
        imsg = [ "aarch64-darwin" ];
        qmd = [
          "aarch64-darwin"
          "x86_64-linux"
        ];
      }
      .${name};
  optionalPackage = name: package: lib.optionalAttrs (supports name) { "${name}" = package; };
in
(optionalPackage "summarize" (
  pkgs.callPackage ./pkgs/summarize.nix {
    pnpm = if pkgs ? pnpm_10 then pkgs.pnpm_10 else pkgs.pnpm;
    nodejs = if pkgs ? nodejs_22 then pkgs.nodejs_22 else pkgs.nodejs;
  }
))
// (optionalPackage "discrawl" (pkgs.callPackage ./pkgs/discrawl.nix { }))
// (optionalPackage "wacrawl" (pkgs.callPackage ./pkgs/wacrawl.nix { }))
// (optionalPackage "gogcli" (pkgs.callPackage ./pkgs/gogcli.nix { }))
// (optionalPackage "goplaces" (pkgs.callPackage ./pkgs/goplaces.nix { }))
// (optionalPackage "camsnap" (pkgs.callPackage ./pkgs/camsnap.nix { }))
// (optionalPackage "sonoscli" (pkgs.callPackage ./pkgs/sonoscli.nix { }))
// (optionalPackage "peekaboo" (pkgs.callPackage ./pkgs/peekaboo.nix { }))
// (optionalPackage "poltergeist" (pkgs.callPackage ./pkgs/poltergeist.nix { }))
// (optionalPackage "sag" (pkgs.callPackage ./pkgs/sag.nix { }))
// (optionalPackage "imsg" (pkgs.callPackage ./pkgs/imsg.nix { }))
// (optionalPackage "qmd" (pkgs.callPackage ./pkgs/qmd.nix { }))
