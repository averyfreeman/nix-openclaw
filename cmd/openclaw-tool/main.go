package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type request struct {
	source  string
	name    string
	out     string
	kind    string
	version string
	rev     string
	binary  string
}

func main() {
	if len(os.Args) < 2 || os.Args[1] != "init" {
		fmt.Fprintln(os.Stderr, "usage: openclaw-tool init <git-url|npm-package> [--output DIR] [--name NAME] [--kind auto|go|npm|generic]")
		os.Exit(2)
	}

	fs := flag.NewFlagSet("init", flag.ExitOnError)
	out := fs.String("output", "", "directory to create (default: ./<name>-openclaw)")
	name := fs.String("name", "", "OpenClaw tool name")
	kind := fs.String("kind", "auto", "project kind: auto, go, npm, or generic")
	version := fs.String("version", "", "exact npm package version")
	rev := fs.String("rev", "main", "Git revision placeholder")
	binary := fs.String("binary", "", "executable name to expose")
	args := os.Args[2:]
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		fatal(errors.New("the source must be the first argument after init"))
	}
	source := args[0]
	_ = fs.Parse(args[1:])
	if fs.NArg() != 0 {
		fatal(errors.New("exactly one Git URL or npm package is required"))
	}

	req, err := normalizeRequest(request{
		source:  source,
		name:    *name,
		out:     *out,
		kind:    *kind,
		version: *version,
		rev:     *rev,
		binary:  *binary,
	})
	if err != nil {
		fatal(err)
	}
	if err := writeScaffold(req); err != nil {
		fatal(err)
	}
	fmt.Printf("created %s\n", req.out)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "openclaw-tool:", err)
	os.Exit(1)
}

func normalizeRequest(req request) (request, error) {
	isNpm := !strings.Contains(req.source, "://") && !strings.HasSuffix(req.source, ".git") || strings.HasPrefix(req.source, "npm:") || strings.Contains(req.source, "npmjs.com/package/")
	if strings.HasPrefix(req.source, "npm:") {
		req.source = strings.TrimPrefix(req.source, "npm:")
		isNpm = true
	}
	if strings.Contains(req.source, "npmjs.com/package/") {
		parsed, err := url.Parse(req.source)
		if err != nil {
			return req, fmt.Errorf("invalid npm URL: %w", err)
		}
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(parts) < 2 {
			return req, errors.New("npm URL must contain /package/<name>")
		}
		req.source = parts[len(parts)-1]
		isNpm = true
	}
	if isNpm {
		if at := strings.LastIndex(req.source, "@"); at > 0 && at > strings.LastIndex(req.source, "/") {
			if req.version == "" {
				req.version = req.source[at+1:]
			}
			req.source = req.source[:at]
		}
		if req.version == "" || req.version == "latest" {
			return req, errors.New("npm sources require an exact version (use --version or name@version)")
		}
	}
	if req.name == "" {
		req.name = strings.TrimPrefix(filepath.Base(strings.TrimSuffix(req.source, ".git")), "@")
		req.name = strings.ReplaceAll(req.name, "@", "-")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`).MatchString(req.name) {
		return req, fmt.Errorf("invalid tool name %q", req.name)
	}
	if req.binary == "" {
		req.binary = req.name
	}
	if req.kind == "" {
		req.kind = "auto"
	}
	if req.kind == "auto" {
		req.kind = "generic"
		if isNpm {
			req.kind = "npm"
		} else if strings.HasSuffix(req.source, ".git") || strings.Contains(req.source, "github.com/") {
			req.kind = "go"
		}
	}
	if req.kind != "go" && req.kind != "npm" && req.kind != "generic" {
		return req, fmt.Errorf("unsupported kind %q", req.kind)
	}
	if req.out == "" {
		req.out = filepath.Join(".", req.name+"-openclaw")
	}
	return req, nil
}

func writeScaffold(req request) error {
	if _, err := os.Stat(req.out); err == nil {
		return fmt.Errorf("output already exists: %s", req.out)
	} else if !os.IsNotExist(err) {
		return err
	}
	files := map[string]string{
		"flake.nix":   renderFlake(req),
		"package.nix": renderPackage(req),
		"README.md":   renderReadme(req),
		filepath.Join("skills", req.name, "SKILL.md"): renderSkill(req),
	}
	for name, contents := range files {
		path := filepath.Join(req.out, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
			return err
		}
	}
	return nil
}

func renderFlake(req request) string {
	return fmt.Sprintf(`{
  description = "OpenClaw tool plugin: %s";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      systems = [ "aarch64-darwin" "x86_64-linux" "aarch64-linux" ];
      forSystem = system:
        let
          pkgs = import nixpkgs { inherit system; };
          package = pkgs.callPackage ./package.nix { };
        in {
          inherit package;
          plugin = {
            name = "%s";
            skills = [ ./skills/%s ];
            packages = [ package ];
            needs = { stateDirs = [ ]; requiredEnv = [ ]; };
          };
        };
    in {
      packages = builtins.listToAttrs (map (system: {
        name = system;
        value = { %s = (forSystem system).package; };
      }) systems);
      openclawPlugin = system: (forSystem system).plugin;
    };
}
`, req.name, req.name, req.name, req.name)
}

func renderPackage(req request) string {
	if req.kind == "npm" {
		return fmt.Sprintf(`{ pkgs, lib ? pkgs.lib }:

pkgs.buildNpmPackage {
  pname = "%s";
  version = "%s";
  src = pkgs.fetchurl {
    url = "https://registry.npmjs.org/%s/-/%s-%s.tgz";
    hash = lib.fakeHash;
  };
  npmDepsHash = lib.fakeHash;
  nativeBuildInputs = [ pkgs.makeWrapper ];
  dontNpmBuild = true;
  installPhase = ''
    mkdir -p $out/lib/%s $out/bin
    cp -R . $out/lib/%s
    # Adjust this path if the package exposes a different executable entrypoint.
    makeWrapper ${pkgs.nodejs}/bin/node $out/bin/%s --add-flags "$out/lib/%s/index.js"
  '';
}
`, req.name, req.version, req.source, req.name, req.version, req.name, req.name, req.binary, req.name)
	}
	if req.kind == "go" {
		return fmt.Sprintf(`{ pkgs, lib ? pkgs.lib }:

pkgs.buildGoModule {
  pname = "%s";
  version = "0.1.0";
  src = pkgs.fetchgit {
    url = "%s";
    rev = "%s";
    hash = lib.fakeHash;
  };
  vendorHash = lib.fakeHash;
  subPackages = [ "." ];
}
`, req.name, req.source, req.rev)
	}
	return fmt.Sprintf(`{ pkgs, lib ? pkgs.lib }:

pkgs.stdenvNoCC.mkDerivation {
  pname = "%s";
  version = "0.1.0";
  src = pkgs.fetchgit {
    url = "%s";
    rev = "%s";
    hash = lib.fakeHash;
  };
  installPhase = ''
    mkdir -p $out/bin
    # Add the upstream build/install commands and expose %s here.
    touch $out/bin/%s
  '';
}
`, req.name, req.source, req.rev, req.binary, req.binary)
}

func renderReadme(req request) string {
	return fmt.Sprintf(`# %s OpenClaw tool

Generated by 'openclaw-tool init'.

Source: '%s'

## Finish the package

'Replace every lib.fakeHash with the hash reported by Nix, then adjust
'package.nix' if the upstream project has a non-default build or executable
layout. The generated package is immutable after installation; runtime tool
configuration belongs in the OpenClaw state directory.

## Use from Home Manager

~~~nix
programs.openclaw.customPlugins = [
  { source = "path:%s"; }
];
~~~
`, req.name, req.source, req.out)
}

func renderSkill(req request) string {
	return fmt.Sprintf(`---
name: %s
description: Use the %s command provided by this OpenClaw tool.
---

Use '%s' after completing the package install and documenting its command
flags here.
`, req.name, req.name, req.binary)
}
