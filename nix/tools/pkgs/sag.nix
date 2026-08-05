{ lib, stdenv, fetchurl }:

let
  sources = {
    "aarch64-darwin" = {
      url = "https://github.com/steipete/sag/releases/download/v0.4.1/sag_0.4.1_darwin_universal.tar.gz";
      hash = "sha256-1AA0tiYX1HIHgJO6W/cwfrcp0XBBIWnp1f5trDg0zvY=";
    };
    "x86_64-linux" = {
      url = "https://github.com/steipete/sag/releases/download/v0.4.1/sag_0.4.1_linux_amd64.tar.gz";
      hash = "sha256-EEu7xqVP6m0EBdRY67ok9/FPxMkMeHHyB9dYnrpkEdI=";
    };
  };
in
stdenv.mkDerivation {
  pname = "sag";
  version = "0.4.1";

  src = fetchurl sources.${stdenv.hostPlatform.system};

  dontConfigure = true;
  dontBuild = true;

  unpackPhase = ''
    tar -xzf "$src"
  '';

  installPhase = ''
    runHook preInstall
    mkdir -p "$out/bin"
    cp sag "$out/bin/sag"
    chmod 0755 "$out/bin/sag"
    runHook postInstall
  '';

  meta = with lib; {
    description = "Command-line ElevenLabs TTS with mac-style flags";
    homepage = "https://github.com/steipete/sag";
    license = licenses.mit;
    platforms = builtins.attrNames sources;
    mainProgram = "sag";
  };
}
