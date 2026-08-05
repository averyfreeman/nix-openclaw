{
  lib,
  stdenv,
  fetchurl,
  ffmpeg,
}:

let
  sources = {
    "aarch64-darwin" = {
      url = "https://github.com/steipete/camsnap/releases/download/v0.3.0/camsnap_0.3.0_darwin_arm64.tar.gz";
      hash = "sha256-fGkZmtwHfGg9RR6Bom+7Ag5ryTGbshXzFw1JKB6bk80=";
    };
    "x86_64-linux" = {
      url = "https://github.com/steipete/camsnap/releases/download/v0.3.0/camsnap_0.3.0_linux_amd64.tar.gz";
      hash = "sha256-C/6gUwte5iuweeNgJtz+OCzsXgv6D4I4SQ2+jF61XXY=";
    };
    "aarch64-linux" = {
      url = "https://github.com/steipete/camsnap/releases/download/v0.3.0/camsnap_0.3.0_linux_arm64.tar.gz";
      hash = "sha256-PUSBmb2kPk+b41tclcGWkYYCFqt+KLi2lPpGow9d82o=";
    };
  };
in
stdenv.mkDerivation {
  pname = "camsnap";
  version = "0.3.0";

  src = fetchurl sources.${stdenv.hostPlatform.system};

  dontConfigure = true;
  dontBuild = true;

  unpackPhase = ''
    tar -xzf "$src"
  '';

  installPhase = ''
    runHook preInstall
    mkdir -p "$out/bin" "$out/share/doc/camsnap"
    cp $(find . -type f -name camsnap | head -1) "$out/bin/camsnap"
    chmod 0755 "$out/bin/camsnap"
    if [ -f LICENSE ]; then
      cp LICENSE "$out/share/doc/camsnap/"
    fi
    if [ -f README.md ]; then
      cp README.md "$out/share/doc/camsnap/"
    fi
    runHook postInstall
  '';

  propagatedBuildInputs = [ ffmpeg ];

  meta = with lib; {
    description = "One command to grab frames, clips, or motion alerts from RTSP/ONVIF cams";
    homepage = "https://github.com/steipete/camsnap";
    license = licenses.mit;
    platforms = builtins.attrNames sources;
    mainProgram = "camsnap";
  };
}
