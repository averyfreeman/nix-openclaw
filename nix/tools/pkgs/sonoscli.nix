{ lib, stdenv, fetchurl }:

let
  sources = {
    "aarch64-darwin" = {
      url = "https://github.com/steipete/sonoscli/releases/download/v0.3.4/sonoscli_0.3.4_darwin_arm64.tar.gz";
      hash = "sha256-03/6r9uioBbIVHfuRRc+YGiyvx9wEVMQUfVknG+DRW4=";
    };
    "x86_64-linux" = {
      url = "https://github.com/steipete/sonoscli/releases/download/v0.3.4/sonoscli_0.3.4_linux_amd64.tar.gz";
      hash = "sha256-GBxSGSuoxBHagQBqA/OCkeET+xJRF58dlHzjzMhUBDI=";
    };
    "aarch64-linux" = {
      url = "https://github.com/steipete/sonoscli/releases/download/v0.3.4/sonoscli_0.3.4_linux_arm64.tar.gz";
      hash = "sha256-61eLFp2nWpaY5v2euS2qrUN25QYugrhQdcJV9095aa4=";
    };
  };
in
stdenv.mkDerivation {
  pname = "sonoscli";
  version = "0.3.4";

  src = fetchurl sources.${stdenv.hostPlatform.system};

  dontConfigure = true;
  dontBuild = true;

  unpackPhase = ''
    tar -xzf "$src"
  '';

  installPhase = ''
    runHook preInstall
    mkdir -p "$out/bin" "$out/share/doc/sonoscli"
    cp $(find . -type f -name sonos | head -1) "$out/bin/sonos"
    chmod 0755 "$out/bin/sonos"
    if [ -f LICENSE ]; then
      cp LICENSE "$out/share/doc/sonoscli/"
    fi
    if [ -f README.md ]; then
      cp README.md "$out/share/doc/sonoscli/"
    fi
    runHook postInstall
  '';

  meta = with lib; {
    description = "Control Sonos speakers from the command-line";
    homepage = "https://github.com/steipete/sonoscli";
    license = licenses.mit;
    platforms = builtins.attrNames sources;
    mainProgram = "sonos";
  };
}
