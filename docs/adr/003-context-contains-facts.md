# ADR-003: Context contains facts

## Status

Accepted

## Date

2026-08-31

## Context

Context may inform policy and provider configuration, but it should not contain either.

```text
facts → policy → representation
```

For example, `criticality = "critical"` is a fact. Replica count, service tier, retention, and tag syntax are policy or representation.

## Decision

Context contains only facts about its subject. It does not:

- infer one fact from another;
- prescribe infrastructure;
- enforce organization-specific completeness;
- contain provider-specific representation;
- create facts through defaults.

## Alternatives

- **Embed policy:** rejected because policy varies by organization.
- **Store infrastructure settings:** rejected because they describe implementation, not the subject.
- **Shape values for one provider:** rejected because canonical facts must remain provider-neutral.

## Consequences

Missing facts stay missing. Organizations enforce completeness separately. Future projections must handle provider limits without changing canonical values.
