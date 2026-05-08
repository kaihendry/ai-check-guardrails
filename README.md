# ai-check-guardrails

Audit Claude AI guardrail compliance on developer workstations.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/kaihendry/ai-check-guardrails/main/install.sh | bash
```

Installs to `~/.local/bin` by default. Override the install location:

```bash
INSTALL_DIR=/usr/local/bin curl -fsSL https://raw.githubusercontent.com/kaihendry/ai-check-guardrails/main/install.sh | bash
```

## Supported platforms

| OS    | amd64 | arm64 |
|-------|-------|-------|
| Linux | ✓     | ✓     |
| macOS | ✓     | ✓     |

## Usage

```
ai-check-guardrails [options]
```

On startup the tool prints its version to stderr and checks for updates:

```
ai-check-guardrails v0.0.20260508143022+abc1234
[update] new version available: v0.0.20260509010101+def5678
[update] downloading...
[update] updated to v0.0.20260509010101+def5678, continuing...
```

### Disable self-update

```bash
# For a single run
ai-check-guardrails --no-update

# For all runs in a session or CI
export NO_AUTO_UPDATE=1
```

### Options

```
  -config string
        path to config JSON
  -mode string
        override run mode: monitor or enforce
  -no-update
        disable self-update check
  -version
        print version and exit
  -install-launchd
        install macOS launchd schedule
  -install-systemd
        install Linux systemd schedule
  -uninstall
        remove installed schedule
```

## Releases

Pre-built binaries are published automatically on every commit to `main`.
See [GitHub Releases](https://github.com/kaihendry/ai-check-guardrails/releases/tag/latest) for the latest build.

## Build from source

```bash
git clone https://github.com/kaihendry/ai-check-guardrails
cd ai-check-guardrails
go build ./cmd/ai-check-guardrails
```
