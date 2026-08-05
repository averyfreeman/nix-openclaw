{
  lib,
  stdenv,
  fetchurl,
}:

let
  sources = {
    "aarch64-darwin" = {
      url = "https://github.com/openclaw/wacrawl/releases/download/v0.3.6/wacrawl_0.3.6_darwin_arm64.tar.gz";
      hash = "sha256-GKT64ZAC0PMBCUn1k327kL/j3NZFX4ymfMN9SZ0fVbw=";
    };
    "x86_64-linux" = {
      url = "https://github.com/openclaw/wacrawl/releases/download/v0.3.6/wacrawl_0.3.6_linux_amd64.tar.gz";
      hash = "sha256-6u/D4tzndgnYPjUasdavNkTFFrmS5R6PpHm1tCMKbe8=";
    };
    "aarch64-linux" = {
      url = "https://github.com/openclaw/wacrawl/releases/download/v0.3.6/wacrawl_0.3.6_linux_arm64.tar.gz";
      hash = "sha256-4FKk2GA1lTC13NcFyJs5e/vlW3WCGGSL11AhRwypZcg=";
    };
  };
in
stdenv.mkDerivation {
  pname = "wacrawl";
  version = "0.3.6";

  src = fetchurl sources.${stdenv.hostPlatform.system};

  dontConfigure = true;
  dontBuild = true;

  unpackPhase = ''
    tar -xzf "$src"
  '';

  installPhase = ''
    runHook preInstall
    mkdir -p "$out/bin" "$out/share/doc/wacrawl"
    cp $(find . -type f -name wacrawl | head -1) "$out/bin/wacrawl"
    chmod 0755 "$out/bin/wacrawl"
    if [ -f LICENSE ]; then
      cp LICENSE "$out/share/doc/wacrawl/"
    fi
    if [ -f README.md ]; then
      cp README.md "$out/share/doc/wacrawl/"
    fi
    runHook postInstall
  '';

  meta = with lib; {
    description = "Read-only local archive and search for WhatsApp Desktop data";
    homepage = "https://github.com/steipete/wacrawl";
    license = licenses.mit;
    platforms = builtins.attrNames sources;
    mainProgram = "wacrawl";
  };
}
