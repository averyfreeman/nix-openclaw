{
  pkgs,
  sourceInfo ? import ../sources/openclaw-source.nix,
  openclawToolPkgs ? { },
  qmdPackage ? null,
  toolNamesOverride ? null,
  excludeToolNames ? [ ],
}:
let
  isDarwin = pkgs.stdenv.hostPlatform.isDarwin;
  pnpm_11 = pkgs.callPackage ./pnpm-11.nix { };
  pnpmForOpenClaw = if toString (sourceInfo.pnpmMajor or "10") == "11" then pnpm_11 else pkgs.pnpm_10;
  toolPkgs = openclawToolPkgs // {
    pnpm = pnpmForOpenClaw;
    inherit pnpm_11;
  };
  toolSets = import ../tools/extended.nix {
    pkgs = pkgs;
    openclawToolPkgs = toolPkgs;
    inherit toolNamesOverride excludeToolNames;
  };
  runtimePluginLocks = import ../generated/openclaw-runtime-plugins;
  buildBundledRuntimePlugin = pkgs.callPackage ../lib/openclaw-runtime-plugin.nix {
    linkOpenClawPeer = false;
  };
  bundledAcpx = buildBundledRuntimePlugin runtimePluginLocks.acpx;
  openclawGateway = pkgs.callPackage ./openclaw-gateway.nix {
    inherit sourceInfo;
    inherit pnpm_11;
    inherit bundledAcpx;
  };
  buildOpenClawRuntimePlugin = pkgs.callPackage ../lib/openclaw-runtime-plugin.nix {
    openclawPackage = openclawGateway;
  };
  openclawRuntimePlugins = pkgs.lib.mapAttrs (
    _name: lock: buildOpenClawRuntimePlugin lock
  ) runtimePluginLocks;
  openclawApp = if isDarwin then pkgs.callPackage ./openclaw-app.nix { } else null;
  openclawBundle = pkgs.callPackage ./openclaw-batteries.nix {
    openclaw-gateway = openclawGateway;
    openclaw-app = openclawApp;
    extendedTools = toolSets.tools;
    version = sourceInfo.releaseVersion or null;
  };
  openclawSelfUpdate = pkgs.buildGoModule {
    pname = "openclaw-self-update";
    version = "0.1.0";
    src = ../../.;
    subPackages = [ "cmd/openclaw-self-update" ];
    vendorHash = null;
  };
in
{
  inherit pnpm_11;
  inherit openclawRuntimePlugins;
  qmd = qmdPackage;
  openclaw-gateway = openclawGateway;
  openclaw = openclawBundle;
  openclaw-self-update = openclawSelfUpdate;
}
// (if isDarwin then { openclaw-app = openclawApp; } else { })
