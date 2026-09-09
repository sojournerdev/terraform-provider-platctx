# ADR-007: Provider projections are deferred

## Status

Accepted

## Date

2026-08-31

## Context

Azure tags, AWS tags, GCP labels, Kubernetes metadata, and OpenTelemetry attributes have different limits and syntax. Designing Context around one of them would weaken the canonical contract.

## Decision

Context v0.1 contains no provider-specific projections, fields, configuration, or validation.

Downstream systems are pressure tests only. A future projection must define its mapping, omissions, encoding, collisions, and lossy behavior separately.

## Alternatives

- **Implement one projection now:** rejected because its constraints could shape Context prematurely.
- **Add provider-ready fields:** rejected because representation is not a subject fact.
- **Use the strictest downstream syntax:** rejected because it would discard valid canonical values.
- **Add an extension map:** rejected because it would bypass field semantics and review.

## Consequences

Context remains provider-neutral. A canonical value may be valid even when a future target cannot represent it directly.
