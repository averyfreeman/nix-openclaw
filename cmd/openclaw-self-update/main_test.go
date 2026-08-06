package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateVersion(t *testing.T) {
	for _, version := range []string{"2026.7.1", "2026.7.1-2", "2026.7.1-beta.1"} {
		if err := validateVersion(version); err != nil {
			t.Errorf("validateVersion(%q): %v", version, err)
		}
	}
	for _, version := range []string{"latest", "v2026.7.1", "2026", "2026.7"} {
		if err := validateVersion(version); err == nil {
			t.Errorf("validateVersion(%q) unexpectedly succeeded", version)
		}
	}
}

func TestSwitchAndRollbackUseAtomicCurrentLink(t *testing.T) {
	root := t.TempDir()
	for _, version := range []string{"2026.7.1", "2026.7.2"} {
		binDir := filepath.Join(root, "releases", version, "node_modules", ".bin")
		if err := os.MkdirAll(binDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(binDir, "openclaw"), []byte("#!/bin/sh\n"), 0755); err != nil {
			t.Fatal(err)
		}
	}

	if err := switchRelease(root, "2026.7.1"); err != nil {
		t.Fatal(err)
	}
	if got, err := os.Readlink(filepath.Join(root, "releases", "current")); err != nil || got != "2026.7.1" {
		t.Fatalf("current link = %q, %v", got, err)
	}

	if err := switchRelease(root, "2026.7.2"); err != nil {
		t.Fatal(err)
	}
	if got, err := os.Readlink(filepath.Join(root, "releases", "previous")); err != nil || got != "2026.7.1" {
		t.Fatalf("previous link = %q, %v", got, err)
	}
	if err := rollback(root); err != nil {
		t.Fatal(err)
	}
	if got, err := os.Readlink(filepath.Join(root, "releases", "current")); err != nil || got != "2026.7.1" {
		t.Fatalf("rolled back current link = %q, %v", got, err)
	}
}

func TestSwitchRejectsUnstagedRelease(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "releases"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := switchRelease(root, "2026.7.1"); err == nil {
		t.Fatal("switchRelease unexpectedly accepted an unstaged release")
	}
}
