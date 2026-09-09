# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in this provider, please report it responsibly.

**Do not open a public GitHub issue for security vulnerabilities.**

Instead, use [GitHub Security Advisories](https://github.com/sojournerdev/terraform-provider-platctx/security/advisories/new) to report privately.

Include:

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

## Response

- **Acknowledgment:** within 48 hours
- **Initial assessment:** within 5 business days
- **Resolution timeline:** depends on severity, typically 30 days for critical issues

## Scope

This policy covers the terraform-provider-platctx codebase only. It does not cover:

- Terraform core
- Third-party dependencies (report upstream)
- Infrastructure not managed by this provider

## Disclosure

We follow coordinated disclosure. Please do not publicly disclose the vulnerability until we have released a fix.

## Security Best Practices

This provider:

- Makes no network requests
- Reads no environment variables
- Logs no Context values
- Telemetry is not collected

For details, see [security-and-resource-limits.md](docs/design/security-and-resource-limits.md).
