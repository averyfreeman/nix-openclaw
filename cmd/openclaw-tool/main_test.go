package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeNpmPackage(t *testing.T) {
	req, err := normalizeRequest(request{source: "@example/tool@1.2.3", name: "tool"})
	if err != nil {
		t.Fatal(err)
	}
	if req.kind != "npm" || req.source != "@example/tool" || req.version != "1.2.3" {
		t.Fatalf("unexpected request: %+v", req)
	}
}

func TestWriteScaffold(t *testing.T) {
	out := filepath.Join(t.TempDir(), "tool")
	req := request{source: "https://github.com/example/tool.git", name: "tool", out: out, kind: "go", rev: "v1"}
	if err := writeScaffold(req); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"flake.nix", "package.nix", "README.md", "skills/tool/SKILL.md"} {
		if _, err := os.Stat(filepath.Join(out, name)); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(filepath.Join(out, "flake.nix"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "openclawPlugin") {
		t.Fatal("generated flake does not export openclawPlugin")
	}
}
