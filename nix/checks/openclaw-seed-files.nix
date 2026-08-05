{ stdenvNoCC }:

stdenvNoCC.mkDerivation {
  pname = "openclaw-seed-files";
  version = "1";
  dontUnpack = true;
  dontConfigure = true;
  dontBuild = true;
  env.OPENCLAW_SEED_FILES = "${../modules/home-manager/openclaw-seed-files.sh}";
  doCheck = true;
  checkPhase = "${../scripts/check-openclaw-seed-files.sh}";
  installPhase = "${../scripts/empty-install.sh}";
}
