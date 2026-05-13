# Research: Fix Self-Update Version Check

## Root Cause Analysis

**Decision**: The bug is a version identifier mismatch between the release pipeline and the binary.

**Finding**: `.github/workflows/release.yml` publishes every release with a static `tag_name: latest`. The binary is built with `-ldflags "-X main.version=$VERSION"` where `$VERSION` is a timestamped string like `v0.0.20260513135826+5118f3a`. When `checkAndUpdate()` fetches `/releases/latest`, the GitHub API returns `tag_name: "latest"`. The comparison `"latest" == "v0.0.20260513135826+5118f3a"` is always false, so the tool re-downloads itself on every run.

**Rationale**: The existing comparison logic in `update.go` is correct — it compares `rel.TagName` with `currentVersion`. The test `TestCheckAndUpdate_NoopWhenAlreadyCurrent` already validates this logic correctly. Only the release pipeline is wrong.

## Fix

**Decision**: Change `tag_name: latest` to `tag_name: ${{ env.VERSION }}` in `.github/workflows/release.yml`.

**Rationale**: 
- The binary already embeds `$VERSION` via ldflags
- Using the same `$VERSION` value as the git/release tag makes the comparison `rel.TagName == currentVersion` work as intended
- GitHub's `/releases/latest` API returns the most recently published non-draft, non-prerelease release — after this change it will return a release with `tag_name` matching the binary's embedded version
- No changes needed in `update.go`, `main.go`, or `install.sh`
- `install.sh` already fetches the release tag dynamically, so switching from `latest` to version-tagged releases does not break installation

**Alternatives considered**:

| Alternative | Rejected because |
|-------------|-----------------|
| Embed `"latest"` as the binary version | Destroys user-visible version info; users cannot identify which build they are running |
| Compare release name (e.g., `"Latest (v0.0...)"`) instead of tag | Fragile — parsing a human-readable string; breaks if name format changes |
| Add a separate `version.txt` asset to each release | Extra complexity with no benefit; the tag already carries the version |
| Overwrite the `latest` git tag and keep the static tag name | Same mismatch persists; moving git tags are an anti-pattern |

## Side Effects

- GitHub will accumulate one release per commit to main instead of a single overwritten `latest` release. This is acceptable and expected for a tool that versions every commit.
- Old `latest`-tagged releases on GitHub remain but become stale; users who installed before this fix will trigger one final re-download (because `"v0.0.TIMESTAMP+HASH" != "latest"`), then the fix takes effect.
- The release name field (`name: "Latest (${{ env.VERSION }})"`) can stay as-is or be simplified — no functional impact.
