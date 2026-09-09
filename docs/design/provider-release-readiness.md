# Provider Release Readiness

## Purpose

Define the evidence and repository requirements for releasing a trustworthy
Terraform provider.

## Required repository metadata

Before the first release, the repository MUST contain:

- a license;
- contribution guidance;
- a security reporting policy;
- a changelog;
- provider source address and namespace;
- supported Terraform, Go, and Framework versions; and
- reproducible build instructions.

## User documentation

The release MUST include a function reference with:

- installation and `required_providers` configuration;
- the qualified function name;
- parameter and result types;
- caller-local examples;
- valid and invalid examples;
- null and unknown behavior;
- diagnostics;
- non-goals; and
- compatibility requirements.

Design documents explain intent. User documentation explains how to use the
provider. They MUST remain separate.

## Build and test evidence

A release requires:

- clean formatting and static analysis;
- unit, Framework, provider, and CLI tests;
- conformance coverage;
- falsification results for the candidate API;
- race detection where supported;
- dependency vulnerability checks;
- minimum-version compatibility tests; and
- a clean build from pinned dependencies.

The release record MUST identify the commit, toolchain versions, test command,
and test result.

## Compatibility and versioning

Use Semantic Versioning. A major release is required for changes to identity,
field meaning, result type, null or unknown behavior, or other documented
breaking behavior. Minor releases may add compatible optional facts after
consumer review. Patch releases fix behavior without changing the contract.

Every release MUST state whether the candidate `canonicalize` API is retained,
revised, deferred, or rejected.

## Distribution

Before publishing, verify:

- provider packages exist for every supported platform;
- package checksums or signatures are published;
- the registry metadata matches the provider source address;
- documentation links resolve; and
- installation works from a clean Terraform configuration.

## Security response

Security reports MUST have a private reporting path, an owner, response targets,
and a release procedure for fixes. Do not rely on public issue reports for
sensitive vulnerabilities.

## Related documents

- [Provider engineering standards](provider-engineering-standards.md)
- [Security and resource limits](security-and-resource-limits.md)
- [Support matrix](support-matrix.md)
- [Falsification plan](context-contract-falsification.md)
