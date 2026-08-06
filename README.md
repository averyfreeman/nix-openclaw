# Invertebrate Outbreak
### [ The Self-Immolating Sequel of the Nix-OpenClaw Story ]
---

![Chibi Lobster be chill, maine](./docs/assets/chibi_lobster_does_Maine.jpg)

This is a mostly personal project. It's a less locked-down variant of the original
[openclaw/nix-openclaw](https://github.com/openclaw/nix-openclaw) project.

It keeps Nix and Home Manager packaging where useful, while allowing runtime
configuration, workspace files, plugins, skills, ACP servers, and experiments
to evolve outside the flake. The installation remains reproducible, but immutability for the installed ecosystem is optional.

Here's the site describing the original project for adherents to the socially acceptable, declarative "Nix way" approach (aka _most people_): [https://docs.openclaw.ai/install/nix](https://docs.openclaw.ai/install/nix)

#### Situational synopsis:
---
  1. OpenClaw binaries cannot currently self-update independently. Nix places them in /nix/store, and Home Manager/flake updates remain authoritative. Runtime configuration is mutable, but installation is not.

  Recommended self-update design:

  - Keep Nix as the base runtime.
  - Install versioned OpenClaw releases under $HOME/.openclaw/releases.
  - Maintain a mutable current symlink.
  - Add staged update, smoke testing, rollback, and opt-in switching.
  - Add an agentic changelog checker that compares release notes against installed plugins, skills, ACP servers, config schema,
    and runtime dependencies.

  - Treat the agent as a compatibility/risk assessor, not an automatic authority.

  2. In current Nix mode, direct plugin mutation commands are intentionally blocked. Plugins and declarative skills are managed
     through Nix/Home Manager. Outside that mode, OpenClaw’s own mutable ecosystem can manage plugins, skills, and ACP servers
     independently.

  A future tool strategy should separate:

  - Nix-managed base packages and trusted plugins.
  - Mutable user-managed extensions under $HOME/.openclaw/tools or $HOME/.openclaw/extensions.
  - A registry recording source, version, enabled state, and optional lock information.
  - Git/npm/local/ClawHub installation sources.
  - Inventory and compatibility checks across both managed and mutable tools.

  Changelog: