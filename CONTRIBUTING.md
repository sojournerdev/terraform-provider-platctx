# Contributing to Terraform Provider PlatCtx

Thank you for your support in wanting to improve the PlatCtx provider! Every contribution helps make the provider better for everyone.

## Prerequisites

- [Go](https://golang.org/doc/install) >= 1.27
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8

## Getting Started

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/<your-username>/terraform-provider-platctx.git
cd terraform-provider-platctx
git remote add upstream https://github.com/sojournerdev/terraform-provider-platctx.git
make tools
```

## Development Workflow

Build and verify your changes:

```bash
make build        # compile
make test         # unit tests with race detector
make lint         # static analysis
make check        # full quality suite (Go + docs)
```

## Testing

Run unit tests:

```bash
make test
```

Run acceptance tests:

```bash
make test-acc
```

## Documentation

After changing schemas or function definitions, regenerate documentation:

```bash
make docs-generate
```

Verify documentation passes all quality gates:

```bash
make check-docs
```

Do not manually edit files in `docs/`. Edit templates in `templates/` instead.

## Commit Messages

This project uses [Conventional Commits](https://www.conventionalcommits.org/). Use this format for all commit messages:

```text
<type>: <description>
```

Common types:

- `feat` — new feature
- `fix` — bug fix
- `docs` — documentation only
- `refactor` — code change that neither fixes a bug nor adds a feature
- `test` — adding or updating tests
- `chore` — build process, tooling, or dependencies
- `style` — formatting, no code change

Examples:

```text
feat: add support for custom criticality levels
fix: handle null governance in canonicalize function
docs: update README with usage examples
chore: update terraform-plugin-framework to v1.19.0
```

## Pull Requests

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Run `make check` to verify all quality gates pass
5. Submit your pull request from your fork

Keep pull requests focused on a single change. Open an issue before substantial changes.

## Code Style

- Follow standard Go conventions
- Run `make fmt` before committing
- Run `make lint` to catch issues
- We use [golangci-lint](https://golangci-lint.run/) to keep code maintainable. See [`.golangci.yml`](.golangci.yml) for our configuration.
- Prefer table-driven tests

## Reporting Issues

Open a [GitHub issue](https://github.com/sojournerdev/terraform-provider-platctx/issues) with:

- A clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Terraform and provider versions
