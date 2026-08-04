## ADDED Requirements

### Requirement: Release-prep merge detection is robust to commit message content
The `release-prepare-web.yml` and `release-prepare-backend.yml` workflows' release-prep merge detection SHALL pass the triggering commit message into its shell step via an `env:` variable, not by interpolating `${{ github.event.head_commit.message }}` directly into the `run:` script body, so that special characters in the commit message (quotes, backticks, `$()`, etc.) cannot break the script's shell parsing or be executed as shell syntax.

#### Scenario: Commit message contains a double quote
- **WHEN** a merge commit's message contains a literal `"` character (e.g. a PR title referencing a quoted name)
- **THEN** the "Detect release-prep merge" step still correctly evaluates whether the message references a `release-prep/<package>-v*` branch, without a shell parse error
