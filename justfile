# WRAAS — Digital Twin Product
# Development recipes

# Default: list available recipes
default:
    @just --list

# Serve the site with live reload (watches for changes)
serve port="3000":
    npx --yes browser-sync start \
        --server site \
        --files "site/**/*" \
        --port {{port}} \
        --no-notify \
        --no-open

# Serve and open in browser
open port="3000":
    npx --yes browser-sync start \
        --server site \
        --files "site/**/*" \
        --port {{port}} \
        --no-notify

# Serve with a simple HTTP server (no live reload, no dependencies)
serve-simple port="8000":
    python3 -m http.server {{port}} --directory site

# Validate HTML files
validate:
    #!/usr/bin/env bash
    set -euo pipefail
    errors=0
    for f in $(find site -name '*.html'); do
        # Check for unclosed tags, missing alt, missing lang
        if ! grep -q 'lang=' "$f"; then
            echo "WARN: $f missing lang attribute"
            errors=$((errors + 1))
        fi
        if grep -q '<img' "$f" && ! grep '<img' "$f" | grep -q 'alt='; then
            echo "WARN: $f has img without alt attribute"
            errors=$((errors + 1))
        fi
    done
    if [ "$errors" -eq 0 ]; then
        echo "All HTML files pass basic validation."
    else
        echo "$errors warning(s) found."
        exit 1
    fi

# List all HTML pages
pages:
    @find site -name '*.html' | sort

# Build the WRAAS CLI
build:
    cd cli && go build -o ../bin/wraas .

# Run CLI unit tests
test-cli:
    cd cli && go test ./...

# Run CLI end-to-end tests
test-e2e:
    bats cli/e2e_test.bats

# Install CLI locally
install:
    cd cli && go install .

# Run the CLI (pass arguments after --)
run *args:
    cd cli && go run . {{args}}

# Build the search index with Pagefind
index:
    npx --yes pagefind --site site/docs --output-path site/docs/_pagefind

# Count lines of code by file type
stats:
    @echo "HTML:" && find site -name '*.html' | wc -l | tr -d ' '
    @echo "CSS:"  && find site -name '*.css'  | wc -l | tr -d ' '
    @echo "JS:"   && find site -name '*.js'   | wc -l | tr -d ' '
    @echo "" && echo "Total lines:"
    @find site -name '*.html' -o -name '*.css' -o -name '*.js' | xargs wc -l | tail -1
