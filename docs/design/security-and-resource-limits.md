# Security and Resource Safety

## Purpose

Keep the provider safe, predictable, and suitable for CI and shared Terraform
execution environments.

## Trust boundaries

The provider receives caller-controlled Terraform values. It MUST treat every
value as untrusted input, including values produced by modules or other
providers.

The provider MUST NOT:

- access the network;
- read files or environment variables;
- inspect credentials or provider configuration;
- execute processes;
- log Context values; or
- send telemetry without a separate reviewed decision.

## Diagnostic safety

Diagnostics MUST identify the failing path and rule without echoing full raw
values. This avoids disclosing owner references, cost centers, classifications,
or other organization-defined identifiers in logs and CI output.

If a diagnostic needs a value fragment for usability, the implementation MUST
define an explicit redaction rule and test it.

## Resource safety

The dynamic adapter MUST validate input without unbounded recursion, allocation,
or diagnostic growth. Before release, define and test limits for:

- maximum input depth;
- maximum attribute count;
- maximum string length; and
- maximum diagnostic count.

If limits are intentionally delegated to Terraform, document that decision and
include a regression test for the delegated boundary. Limits MUST return a
clear provider diagnostic rather than panic or truncate input.

Known values MUST NOT be truncated or rewritten to satisfy a limit. A rejected
limit is different from a changed canonical fact.

## Dependency and supply-chain safety

The repository MUST:

- pin Go and Plugin Framework versions;
- review dependency updates;
- run vulnerability checks in CI;
- build from a clean, reproducible dependency graph; and
- publish checksums or signatures according to the release process.

The provider MUST avoid dependencies that introduce network or runtime side
effects into the pure Context path.

## Failure behavior

Malformed input MUST produce diagnostics, not a panic. Unexpected internal
errors MUST be surfaced as errors and MUST NOT produce a partial canonical
result.

The function MUST remain deterministic after errors. A failed call MUST NOT
change process-global state or affect a later call.

## Related documents

- [Terraform boundary contract](terraform-boundary-contract.md)
- [Diagnostic catalog](diagnostic-catalog.md)
- [Support matrix](support-matrix.md)
- [Test strategy](context-contract-test-strategy.md)
