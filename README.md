# WRAAS — Digital Twin Product

April Fools 2026 joke site for **WRAAS** (Weighted Romain Algorithmic Approximation Software) — a fictional high-fidelity neural network delivering the exact same brilliant insights, and signature sighs, as Romain himself.

Live at [wraas.click](https://wraas.click).

## Development

The site is hand-crafted static HTML with no build step.

```sh
# Serve with live reload (requires npx)
just serve

# Serve and open in browser
just open

# Serve without dependencies
just serve-simple

# Validate HTML (lang attribute, img alt)
just validate

# List all pages
just pages
```

## Project structure

```text
site/
  index.html          # Landing page
  docs/               # Fake documentation site (Diataxis-structured)
    tutorials/        # Learning-oriented guides
    how-to/           # Task-oriented guides
    reference/        # Lookup material (CLI, config, decisions)
    explanation/      # Design philosophy articles
  styles/             # CSS
  scripts/            # JavaScript
  images/             # Assets (WebP, PNG)
  favicons/           # Favicon set
```

## Deployment

Pushing to `main` with changes in `site/` triggers the [deploy workflow](.github/workflows/deploy.yml), which publishes to GitHub Pages.
