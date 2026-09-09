# Architecture Decision Records

This directory records long-lived Context Core decisions.

## Decision log

| ADR | Decision | Status |
| --- | --- | --- |
| [001](001-function-only-terraform-provider.md) | Use a function-only provider if a Terraform API is justified | Accepted |
| [002](002-one-context-one-logical-subject.md) | One Context describes one logical subject | Accepted |
| [003](003-context-contains-facts.md) | Context contains facts, not policy or representation | Accepted |
| [004](004-opaque-semantic-identifiers.md) | Organization-defined identifiers remain opaque | Accepted |
| [005](005-structural-terraform-contract.md) | Terraform Context is a structural contract | Accepted |
| [006](006-missing-facts-remain-missing.md) | Missing facts remain missing | Accepted |
| [007](007-provider-projections-deferred.md) | Provider projections are deferred | Accepted |
| [008](008-canonicalize-public-api.md) | Treat `canonicalize` as a release-gated candidate API | Accepted |
| [009](009-terraform-function-parameter-strategy.md) | Use a dynamic input for strict Context validation | Accepted |

## Lifecycle

```text
Proposed → Accepted → Superseded
         ↘ Rejected
```

- **Proposed:** under review and open to change.
- **Accepted:** approved and in force.
- **Rejected:** considered but not adopted.
- **Superseded:** replaced by a later ADR.

Accepted ADRs are historical records. Correct spelling or broken links in place, but record decision changes in a new ADR. A superseded ADR must link to its replacement, and the replacement must link back.

## Required structure

New ADRs use the next number and contain:

```markdown
# ADR-NNN: Decision title

## Status

Proposed

## Date

YYYY-MM-DD

## Context

State the problem and constraints.

## Decision

State the decision.

## Alternatives

List credible alternatives and why they were not chosen.

## Consequences

State the important outcomes and tradeoffs.
```

Add a `References` section only when the decision depends on an external specification or documented behavior.
