#!/usr/bin/env bats
# WRAAS CLI — End-to-end tests (bats)
# Run: just test-e2e

setup_file() {
  export WRAAS="$(cd "$BATS_TEST_DIRNAME/.." && pwd)/bin/wraas"
  # Build binary once for the whole file
  (cd "$BATS_TEST_DIRNAME" && go build -o "$WRAAS" .)
}

setup() {
  # Each test runs in its own temp directory
  TEST_DIR="$(mktemp -d)"
  cd "$TEST_DIR"
}

teardown() {
  rm -rf "$TEST_DIR"
}

# ── Version ───────────────────────────────────────────

@test "--version prints version" {
  run "$WRAAS" --version
  [ "$status" -eq 0 ]
  [[ "$output" == *"WRAAS v1.3"* ]]
}

@test "version subcommand shows full info" {
  run "$WRAAS" version
  [ "$status" -eq 0 ]
  [[ "$output" == *"WRAAS"* ]]
  [[ "$output" == *"1.3"* ]]
  [[ "$output" == *"FULL"* ]]
  [[ "$output" == *"0.00%"* ]]
  [[ "$output" == *"113ms"* ]]
}

# ── Help ──────────────────────────────────────────────

@test "help exits 0" {
  run "$WRAAS" --help
  [ "$status" -eq 0 ]
}

@test "help shows version" {
  run "$WRAAS" --help
  [[ "$output" == *"WRAAS v1.3"* ]]
}

@test "help lists all commands" {
  run "$WRAAS" --help
  [[ "$output" == *"init"* ]]
  [[ "$output" == *"query"* ]]
  [[ "$output" == *"review"* ]]
  [[ "$output" == *"status"* ]]
  [[ "$output" == *"config"* ]]
}

@test "help shows desertion rate and latency" {
  run "$WRAAS" --help
  [[ "$output" == *"0.00%"* ]]
  [[ "$output" == *"113ms"* ]]
}

# ── Init ──────────────────────────────────────────────

@test "init creates config file" {
  run "$WRAAS" init
  [ "$status" -eq 0 ]
  [ -f wraas.yml ]
  [[ "$output" == *"Config written to"* ]]
  [[ "$output" == *"makes the relationship official"* ]]
  [[ "$output" == *"113ms"* ]]
}

@test "init config contains correct defaults" {
  "$WRAAS" init
  grep -q "latency_target_ms: 113" wraas.yml
  grep -q "level: FULL" wraas.yml
  grep -q "# Non-negotiable" wraas.yml
  grep -q "include_wrong_options: true" wraas.yml
  grep -q "threshold: MILD" wraas.yml
}

@test "init fails when config already exists" {
  "$WRAAS" init
  run "$WRAAS" init
  [ "$status" -eq 1 ]
  [[ "$output" == *"already exists"* ]]
}

@test "init --force overwrites and archives old config" {
  "$WRAAS" init
  run "$WRAAS" init --force
  [ "$status" -eq 0 ]
  [[ "$output" == *"archived"* ]]
  # Backup file should exist
  ls wraas.yml.bak.* >/dev/null 2>&1
}

# ── Status ────────────────────────────────────────────

@test "status shows hardcoded values (text)" {
  run "$WRAAS" status
  [ "$status" -eq 0 ]
  [[ "$output" == *"RUNNING"* ]]
  [[ "$output" == *"FULL"* ]]
  [[ "$output" == *"ACTIVE"* ]]
  [[ "$output" == *"0.00%"* ]]
  [[ "$output" == *"113ms"* ]]
}

@test "status json output" {
  run "$WRAAS" status --output json
  [ "$status" -eq 0 ]
  [[ "$output" == *'"engine": "RUNNING"'* ]]
  [[ "$output" == *'"latency_ms": 113'* ]]
  [[ "$output" == *'"commitment_level": "FULL"'* ]]
  [[ "$output" == *'"desertion_rate": 0'* ]]
}

@test "status yaml output" {
  run "$WRAAS" status --output yaml
  [ "$status" -eq 0 ]
  [[ "$output" == *"engine: RUNNING"* ]]
  [[ "$output" == *"commitment_level: FULL"* ]]
  [[ "$output" == *"latency_ms: 113"* ]]
}

# ── Config Get/Set ────────────────────────────────────

@test "config get reads value" {
  "$WRAAS" init --quiet
  run "$WRAAS" config get commit.enforcement
  [ "$status" -eq 0 ]
  [[ "$output" == *"STANDARD"* ]]
}

@test "config get reads nested values" {
  "$WRAAS" init --quiet
  run "$WRAAS" config get engine.latency_target_ms
  [ "$status" -eq 0 ]
  [[ "$output" == *"113"* ]]
}

@test "config set commitment.level PARTIAL triggers satirical override" {
  "$WRAAS" init --quiet
  run "$WRAAS" config set commitment.level PARTIAL
  [ "$status" -eq 0 ]
  [[ "$output" == *"Noted. Proceeding at FULL commitment"* ]]
  [[ "$output" == *"113ms"* ]]
}

@test "config set include_wrong_options false triggers DME-0001" {
  "$WRAAS" init --quiet
  run "$WRAAS" config set decision.include_wrong_options false
  [ "$status" -eq 0 ]
  [[ "$output" == *"DME-0001"* ]]
  [[ "$output" == *"Value reset to true"* ]]
}

@test "config set commitment.timeout_ms is accepted and ignored" {
  "$WRAAS" init --quiet
  run "$WRAAS" config set commitment.timeout_ms 5000
  [ "$status" -eq 0 ]
  [[ "$output" == *"Accepted. Logged. Ignored."* ]]
}

@test "config set commitment.level FULL has no warning" {
  "$WRAAS" init --quiet
  run "$WRAAS" config set commitment.level FULL
  [ "$status" -eq 0 ]
  [[ "$output" != *"Noted. Proceeding"* ]]
}

# ── Query (demo mode) ────────────────────────────────

@test "query demo mode shows structured response" {
  run env WRAAS_PROVIDER=mock "$WRAAS" query --input "Should I use Redis?"
  [ "$status" -eq 0 ]
  [[ "$output" == *"VIABLE"* ]]
  [[ "$output" == *"REJECTED"* ]]
  [[ "$output" == *"DME-0047"* ]]
  [[ "$output" == *"87.3%"* ]]
  [[ "$output" == *"sigh"* ]]
  [[ "$output" == *"113ms"* ]]
  [[ "$output" == *"Demo mode"* ]]
}

@test "query --include-wrong-options=false triggers DME-0001" {
  run env WRAAS_PROVIDER=mock "$WRAAS" query --input "test" --include-wrong-options=false
  [ "$status" -eq 0 ]
  [[ "$output" == *"DME-0001"* ]]
}

@test "query --sigh existential overrides sigh level" {
  run env WRAAS_PROVIDER=mock "$WRAAS" query --input "wip" --sigh existential
  [ "$status" -eq 0 ]
  [[ "$output" == *"EXISTENTIAL"* ]]
}

@test "query wip input auto-calibrates to existential sigh" {
  run env WRAAS_PROVIDER=mock "$WRAAS" query --input "wip"
  [ "$status" -eq 0 ]
  [[ "$output" == *"EXISTENTIAL"* ]]
}

@test "query json output is valid" {
  run env WRAAS_PROVIDER=mock "$WRAAS" query --input "test" --output json
  [ "$status" -eq 0 ]
  [[ "$output" == *'"input": "test"'* ]]
  [[ "$output" == *'"latency_ms": 113'* ]]
  [[ "$output" == *'"sigh_level"'* ]]
}

@test "query missing --input flag fails" {
  run "$WRAAS" query
  [ "$status" -ne 0 ]
  [[ "$output" == *"required"* ]]
}

# ── Review ────────────────────────────────────────────

@test "review missing required flags fails" {
  run "$WRAAS" review
  [ "$status" -eq 1 ]
  [[ "$output" == *"required flag"* ]]
}

@test "review --strict=false warns there is no false" {
  # This test needs gh, but the warning appears before the API call
  run "$WRAAS" review --pr 1 --repo rlespinasse/github-slug-action --strict=false
  [[ "$output" == *"no false"* ]]
}

# Review tests that require gh CLI authentication
@test "review fetches real PR and validates commits" {
  if ! command -v gh &>/dev/null || ! gh auth status &>/dev/null 2>&1; then
    skip "gh CLI not authenticated"
  fi
  run "$WRAAS" review --pr 1 --repo rlespinasse/github-slug-action
  [ "$status" -eq 0 ]
  [[ "$output" == *"commits"* ]]
  [[ "$output" == *"113ms"* ]]
}

@test "review json output for real PR" {
  if ! command -v gh &>/dev/null || ! gh auth status &>/dev/null 2>&1; then
    skip "gh CLI not authenticated"
  fi
  run "$WRAAS" review --pr 1 --repo rlespinasse/github-slug-action --output json
  [ "$status" -eq 0 ]
  [[ "$output" == *'"pr": 1'* ]]
  [[ "$output" == *'"latency_ms": 113'* ]]
  [[ "$output" == *'"repo": "rlespinasse/github-slug-action"'* ]]
}

@test "review detects missing scope in commits" {
  if ! command -v gh &>/dev/null || ! gh auth status &>/dev/null 2>&1; then
    skip "gh CLI not authenticated"
  fi
  run "$WRAAS" review --pr 1 --repo rlespinasse/github-slug-action
  [[ "$output" == *"scope"* ]]
}

# ── Persistence ──────────────────────────────────────

@test "query records last_query in state.json" {
  env WRAAS_PROVIDER=mock "$WRAAS" query --input "persistence test"
  [ -f "$HOME/.wraas/state.json" ]
  grep -q "last_query" "$HOME/.wraas/state.json"
}

@test "query logs sigh to sigh.log" {
  env WRAAS_PROVIDER=mock "$WRAAS" query --input "persistence test"
  [ -f "$HOME/.wraas/sigh.log" ]
  grep -q "query:" "$HOME/.wraas/sigh.log"
}
