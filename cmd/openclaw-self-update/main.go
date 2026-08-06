package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var exactVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z][0-9A-Za-z.-]*)?$`)

type releaseReport struct {
	Version       string   `json:"version"`
	Current       string   `json:"current,omitempty"`
	ReleaseURL    string   `json:"releaseUrl"`
	ReleaseName   string   `json:"releaseName,omitempty"`
	PublishedAt   string   `json:"publishedAt,omitempty"`
	Changelog     string   `json:"changelog,omitempty"`
	InstalledTool []string `json:"installedTools,omitempty"`
	Warnings      []string `json:"warnings,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	root := updateRoot()
	switch os.Args[1] {
	case "check":
		fatalIf(checkLatest())
	case "review":
		fatalIf(review(os.Args[2:]))
	case "stage":
		fatalIf(stage(os.Args[2:]))
	case "switch":
		fatalIf(switchCommand(os.Args[2:]))
	case "rollback":
		fatalIf(rollback(root))
	case "status":
		fatalIf(status(root))
	case "run":
		fatalIf(run(root, os.Args[2:]))
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: openclaw-self-update {check|review [VERSION]|stage VERSION|switch VERSION|rollback|status|run [ARGS...]}`)
	os.Exit(2)
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "openclaw-self-update:", err)
		os.Exit(1)
	}
}

func updateRoot() string {
	if root := os.Getenv("OPENCLAW_SELF_UPDATE_ROOT"); root != "" {
		return root
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fatalIf(err)
	}
	return filepath.Join(home, ".openclaw")
}

func validateVersion(version string) error {
	if !exactVersion.MatchString(version) {
		return fmt.Errorf("version %q must be an exact semantic version such as 2026.7.1", version)
	}
	return nil
}

func releasesDir(root string) string          { return filepath.Join(root, "releases") }
func releasePath(root, version string) string { return filepath.Join(releasesDir(root), version) }

func switchCommand(args []string) error {
	if len(args) != 1 {
		return errors.New("switch requires exactly one staged version")
	}
	if err := validateVersion(args[0]); err != nil {
		return err
	}
	return switchRelease(updateRoot(), args[0])
}

func switchRelease(root, version string) error {
	if err := validateVersion(version); err != nil {
		return err
	}
	target := releasePath(root, version)
	if info, err := os.Stat(filepath.Join(target, "node_modules", ".bin", "openclaw")); err != nil || info.IsDir() {
		return fmt.Errorf("release %s is not staged or has no openclaw executable", version)
	}
	if err := os.MkdirAll(releasesDir(root), 0755); err != nil {
		return err
	}
	current := filepath.Join(releasesDir(root), "current")
	if old, err := os.Readlink(current); err == nil && old != version {
		if err := atomicSymlink(releasesDir(root), "previous", old); err != nil {
			return err
		}
	}
	if err := atomicSymlink(releasesDir(root), "current", version); err != nil {
		return err
	}
	fmt.Printf("active OpenClaw release: %s\n", version)
	return nil
}

func atomicSymlink(dir, name, target string) error {
	tmp := filepath.Join(dir, fmt.Sprintf(".%s.%d", name, time.Now().UnixNano()))
	if err := os.Symlink(target, tmp); err != nil {
		return err
	}
	if err := os.Rename(tmp, filepath.Join(dir, name)); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func rollback(root string) error {
	dir := releasesDir(root)
	current, err := os.Readlink(filepath.Join(dir, "current"))
	if err != nil {
		return fmt.Errorf("no active release to roll back: %w", err)
	}
	previous, err := os.Readlink(filepath.Join(dir, "previous"))
	if err != nil {
		return fmt.Errorf("no previous release to roll back to: %w", err)
	}
	if err := atomicSymlink(dir, "current", previous); err != nil {
		return err
	}
	if err := atomicSymlink(dir, "previous", current); err != nil {
		return err
	}
	fmt.Printf("rolled back OpenClaw release to %s\n", previous)
	return nil
}

func run(root string, args []string) error {
	current, err := os.Readlink(filepath.Join(releasesDir(root), "current"))
	if err != nil {
		fallback := os.Getenv("OPENCLAW_SELF_UPDATE_FALLBACK")
		if fallback == "" {
			return errors.New("no active mutable release; run stage and switch first")
		}
		return execCommand(fallback, false, args...)
	}
	bin := filepath.Join(releasePath(root, current), "node_modules", ".bin", "openclaw")
	return execCommand(bin, true, args...)
}

func execCommand(binary string, clearNixMode bool, args ...string) error {
	command := exec.Command(binary, args...)
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
	command.Env = os.Environ()
	if clearNixMode {
		command.Env = append(command.Env, "OPENCLAW_NIX_MODE=", "OPENCLAW_DISABLE_PERSISTED_PLUGIN_REGISTRY=")
	}
	return command.Run()
}

func stage(args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return errors.New("stage requires VERSION and optionally a package source")
	}
	version := args[0]
	if err := validateVersion(version); err != nil {
		return err
	}
	root := updateRoot()
	target := releasePath(root, version)
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("release %s is already staged", version)
	}
	if err := os.MkdirAll(releasesDir(root), 0755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(releasesDir(root), ".stage-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	source := "openclaw@" + version
	if len(args) == 2 {
		source = args[1]
	}
	npm := os.Getenv("OPENCLAW_SELF_UPDATE_NPM")
	if npm == "" {
		npm = "npm"
	}
	command := exec.Command(npm, "install", "--prefix", tmp, "--ignore-scripts", "--omit=dev", "--no-audit", "--no-fund", source)
	command.Stdout, command.Stderr = os.Stdout, os.Stderr
	command.Env = append(os.Environ(), "NPM_CONFIG_UPDATE_NOTIFIER=false")
	if err := command.Run(); err != nil {
		return fmt.Errorf("install %s: %w", source, err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "node_modules", ".bin", "openclaw")); err != nil {
		return fmt.Errorf("installed source has no node_modules/.bin/openclaw: %w", err)
	}
	if err := os.Rename(tmp, target); err != nil {
		return err
	}
	fmt.Printf("staged OpenClaw release: %s\n", version)
	return nil
}

func checkLatest() error {
	output, err := runCommand("npm", "view", "openclaw", "version", "--json")
	if err != nil {
		return fmt.Errorf("checking npm for the latest OpenClaw release: %w", err)
	}
	version := strings.Trim(strings.TrimSpace(output), `"`)
	if err := validateVersion(version); err != nil {
		return fmt.Errorf("npm returned invalid OpenClaw version %q: %w", version, err)
	}
	fmt.Println(version)
	return nil
}

func review(args []string) error {
	version := ""
	if len(args) > 1 {
		return errors.New("review accepts zero or one version")
	}
	if len(args) == 1 {
		version = args[0]
	} else {
		output, err := runCommand("npm", "view", "openclaw", "version", "--json")
		if err != nil {
			return err
		}
		version = strings.Trim(strings.TrimSpace(output), `"`)
	}
	if err := validateVersion(version); err != nil {
		return err
	}
	current := ""
	if link, err := os.Readlink(filepath.Join(releasesDir(updateRoot()), "current")); err == nil {
		current = link
	}
	report := releaseReport{
		Version:    version,
		Current:    current,
		ReleaseURL: "https://github.com/openclaw/openclaw/releases/tag/v" + version,
		Warnings: []string{
			"Review release notes and run smoke tests before switching.",
			"NPM-installed releases are outside Nix reproducibility and require manual rollback planning.",
		},
	}
	if body, err := fetchReleaseNotes(version); err == nil {
		report.ReleaseName = body.Name
		report.PublishedAt = body.PublishedAt
		report.Changelog = body.Body
	} else {
		report.Warnings = append(report.Warnings, "GitHub release notes unavailable: "+err.Error())
	}
	report.InstalledTool = discoverTools(updateRoot())
	encoded, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(encoded))
	return nil
}

type githubRelease struct {
	Name        string `json:"name"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
}

func fetchReleaseNotes(version string) (githubRelease, error) {
	request, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/openclaw/openclaw/releases/tags/v"+version, nil)
	if err != nil {
		return githubRelease{}, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return githubRelease{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return githubRelease{}, fmt.Errorf("GitHub returned %s", response.Status)
	}
	var release githubRelease
	if err := json.NewDecoder(response.Body).Decode(&release); err != nil {
		return githubRelease{}, err
	}
	return release, nil
}

func discoverTools(root string) []string {
	entries, err := os.ReadDir(filepath.Join(root, "tools"))
	if err != nil {
		return nil
	}
	tools := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			tools = append(tools, entry.Name())
		}
	}
	sort.Strings(tools)
	return tools
}

func runCommand(name string, args ...string) (string, error) {
	command := exec.Command(name, args...)
	output, err := command.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s: %s: %w", name, strings.TrimSpace(string(exit.Stderr)), err)
		}
		return "", err
	}
	return string(output), nil
}

func status(root string) error {
	current, currentErr := os.Readlink(filepath.Join(releasesDir(root), "current"))
	previous, previousErr := os.Readlink(filepath.Join(releasesDir(root), "previous"))
	if currentErr != nil {
		current = "none"
	}
	if previousErr != nil {
		previous = "none"
	}
	fmt.Printf("root=%s\ncurrent=%s\nprevious=%s\n", root, current, previous)
	return nil
}
