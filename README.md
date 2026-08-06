<!-- % prefer front matter not to render %
---
Title: Nix-OpenClaw Invertebrate Outbreak Edition
Author: Avery Freeman
Version: see tagged commits 
tags:
  - OpenClaw
  - AI_Assistant
  - Nix
  - TypeScript
  - npm_packages
--- 
-->
<div align="center">
  <h1>Invertebrate Outbreak</h1>
  <h3>[ The Self-Immolating Sequel of the Nix-OpenClaw Story ]</h3>
  <p>
    <img src="./docs/assets/chibi_lobster_does_Maine.png" alt="Chibi Lobster be chill, Maine" width="300" height="300" />
  </p>
</div>

---
This is a mostly personal project. It's a less locked-down variant of the original
[openclaw/nix-openclaw](https://github.com/openclaw/nix-openclaw) project.

It keeps Nix and Home Manager packaging where useful, while allowing runtime
configuration, workspace files, plugins, skills, ACP servers, and experiments
to evolve outside the flake. The installation remains reproducible, but immutability for the installed ecosystem is optional.

Here's the site describing the original project for adherents to the socially acceptable, declarative "Nix way" approach (aka _most people_): [https://docs.openclaw.ai/install/nix](https://docs.openclaw.ai/install/nix)

#### Situational synopsis:
---

1. Issue Maintaining OpenClaw with Nix: **OpenClaw cannot self-update.** This can seem untenable, since the cadence of releases is so rapid, many users experience edge cases that require frequent patching, and so on.

**Solution:** Design an installation policy where direct self-update not possible, but still managed by Nix.

  - Nix Home Manager installs versioned releases of OpenClaw install under `$HOME/.openclaw/releases`.
  - Release version is selected by using a symlink to `releases/current` folder.
  - Installation tooling stages updates, performs smoke tests, provides rollbacks if necessary, and offers non-destructive version migration.
  - Agentic changelog checker compares release notes against installed plugins, skills, ACP servers, config schema, and dependencies.

This way, OpenClaw acts as both a risk assessor and a reserved salesperson touting a beyond-generous return policy, rather than an obvlious decapod steamroller.

2. Issues using an immutable OpenClaw: In official `nix-openclaw` (and Nix packages, by design), direct mutations are intentionally blocked. This can be really frustrating for ecosystems that have lots of plugins and other add-ons, and for codebases that progress at a breathtaking cadence.  

**Solution:** In _Invertebrate Jailbreak Edition_, OpenClaw’s ecosystem remains mutable like a traditional installation. Users can manage plugins, skills, and ACP servers independently, as well. One can still use declarative skills curated by the Nix-Openclaw ecosystem, and managed through Nix/Home Manager, but outside that mode other plugins are available, as well. 

More flexible tool strategy separates:
  - Nix-managed base packages and trusted plugins.
  - Mutable user-managed extensions under `$HOME/.openclaw/tools` or `$HOME/.openclaw/extensions`.
  - A registry recording source, version, enabled state, and optional lock info.
  - Access to sources such as `Git`/`npm`/`local`/`ClawHub`.
  - Inventory and compatibility checks offered for managed _and_ mutable plugins!

#### Opt-In to Self-Update:

Declare in Home Manager as usual using modified flake, and enable the **mutable channel:**

```nix
programs.openclaw.selfUpdate.enable = true;
```

After Home Manager installs with above directive, use new included command:

```bash
openclaw-self-update status
openclaw-self-update check
openclaw-self-update review 2026.7.1
openclaw-self-update stage 2026.7.1
openclaw-self-update switch 2026.7.1
openclaw-self-update rollback
```

Staging does not activate a release, so switching is explicit. If no mutable release is active, the gateway falls back to the Nix package.

See Changelog for more details.