# Context Contract Adversarial Fixtures

## Purpose

This document challenges the proposed Context contract after the [normal fixtures](context-contract-fixtures.md) established that ordinary subjects fit. It tests valid paths, invalid paths, and semantic boundaries before any provider implementation is frozen.

Azure supplies familiar implementation pressure. Azure resources are examples of how a logical subject may be implemented; they are not Context fields, identity components, or validation rules.

## Contract under test

```hcl
{
  identity = {
    name      = string
    namespace = optional(string)
  }
  ownership = {
    owned_by    = string
    operated_by = optional(string)
  }
  environment = optional(string)
  governance = optional({
    criticality         = optional(string)
    cost_center         = optional(string)
    data_classification = optional(string)
  })
}
```

Logical identity is `(namespace?, name)`. Canonical optional absence is a typed null. Known strings must be valid UTF-8, non-empty, and free of leading or trailing Unicode whitespace.

## Evaluation rules

Each fixture has one outcome:

- **Accept:** the value is semantically valid and has the stated canonical result.
- **Reject:** known input violates the structural or semantic contract.
- **Preserve unknown:** validation must wait because Terraform does not know the value yet.
- **Pressure:** the case is representable but tests whether a proposed meaning remains coherent.

Rejecting an invalid value is successful contract behavior. A fixture fails the proposed contract only when truthful facts cannot be represented without fabrication, overloading, or provider-specific data.

## Happy paths

### H1 — Minimal subject

Facts: `payments/api` is accountable to `team:payments`; no other facts are declared.

```hcl
{
  identity  = { namespace = "payments", name = "api" }
  ownership = { owned_by = "team:payments" }
}
```

**Answer:** Accept. Canonical output preserves the three known values and represents `operated_by`, `environment`, and `governance` as typed nulls. No default is introduced.

### H2 — Fully described subject

Facts: the production payments database has a distinct operator and all governance facts.

```hcl
{
  identity = { namespace = "payments", name = "database" }
  ownership = {
    owned_by    = "team:payments"
    operated_by = "team:database-platform"
  }
  environment = "production"
  governance = {
    criticality         = "critical"
    cost_center         = "CC-1042"
    data_classification = "restricted"
  }
}
```

**Answer:** Accept. Every field has one meaning and every known value is preserved exactly.

### H3 — Shared Azure platform without an environment

Facts: an Azure hub network spans environments and is accountable to the network platform team. Its representative footprint includes a virtual network, subnets, route tables, Azure Firewall, and connectivity gateways.

```hcl
{
  identity = { name = "hub-network" }
  ownership = {
    owned_by    = "team:network-platform"
    operated_by = "team:network-operations"
  }
  governance = { criticality = "high" }
}
```

**Answer:** Accept. A shared subject does not need a fabricated environment. Azure resource groups, subscriptions, regions, and resource IDs do not participate in Context identity.

### H4 — Opaque Unicode identifiers

```hcl
{
  identity = {
    namespace = "europe"
    name      = "payments-東京"
  }
  ownership = { owned_by = "team:central payments" }
  environment = "pre production"
  governance = {
    cost_center         = "FIN/東京-42"
    data_classification = "internal—restricted"
  }
}
```

**Answer:** Accept. Unicode, punctuation, and internal spaces remain unchanged. Future Azure tag restrictions belong to a projection and cannot weaken canonical Context.

### H5 — Explicitly equal owner and operator

```hcl
{
  identity = { name = "monitoring-workspace" }
  ownership = {
    owned_by    = "team:observability"
    operated_by = "team:observability"
  }
}
```

**Answer:** Accept. Equal values express two distinct facts. Canonicalization must not remove `operated_by` merely because its value equals `owned_by`.

## Invalid paths

### B1 — Empty or whitespace-only required identifiers

| Input | Answer | Reason |
| --- | --- | --- |
| `identity.name = ""` | Reject | Empty known semantic identifier |
| `identity.name = " "` | Reject | Whitespace-only identifier |
| `ownership.owned_by = "\t"` | Reject | Whitespace-only identifier |
| `identity.namespace = "\n"` | Reject | Known optional value is invalid, not absent |

Canonicalization must not convert any of these values to null.

### B2 — Surrounding whitespace

| Input | Answer |
| --- | --- |
| `identity.name = " api"` | Reject |
| `identity.name = "api "` | Reject |
| `ownership.owned_by = " team:payments "` | Reject |
| `environment = "production\n"` | Reject |

**Answer:** Reject without trimming. Silent trimming would merge distinct input and conceal a configuration defect.

### B3 — Missing or null required values

| Input state | Answer |
| --- | --- |
| Top-level Context is null | Reject |
| `identity` omitted or null | Reject |
| `ownership` omitted or null | Reject |
| `identity.name` omitted or null | Reject |
| `ownership.owned_by` omitted or null | Reject |

Required means a known value is eventually required. Null is absence and cannot satisfy a required fact.

### B4 — Wrong structure or type

| Input | Answer |
| --- | --- |
| Context is a list, tuple, set, or primitive | Reject |
| `identity.name` is a number or boolean | Reject |
| `governance` is a list or string | Reject |
| `criticality` is a collection | Reject |
| Required nested object is malformed | Reject with its full attribute path |
| Unrecognized attribute such as `identity.display_name` | Reject as unsupported structure |

Every public Context consumer must perform these checks independently because Terraform values are structural, not nominal.

### B5 — Unsupported Criticality

| Value | Answer |
| --- | --- |
| `"low"` | Accept |
| `"medium"` | Accept |
| `"high"` | Accept |
| `"critical"` | Accept |
| `"HIGH"` | Reject |
| `"tier-1"` | Reject |
| `"severe"` | Reject |
| `""` | Reject |

**Answer:** Keep the closed vocabulary. It gives Criticality stable provider-owned semantics and deterministic validation. Organization-specific tiers are not automatically equivalent to expected impact; mapping them would be policy outside Context. An organization that cannot state one of the four meanings may omit the fact rather than manufacture it.

### B6 — Manufactured defaults

| Missing fact | Forbidden canonical value |
| --- | --- |
| `operated_by` | Copy of `owned_by` |
| `environment` | `"production"` |
| `criticality` | `"medium"` |
| `data_classification` | `"internal"` |

**Answer:** Reject the defaulting behavior. The input itself remains valid with canonical nulls for the missing facts.

### B7 — Multiple subjects in one Context

A module manages an API, worker, database, and their Azure resources. It proposes one Context named `payments-platform` even though accountability or governance differs among those subjects.

**Answer:** Reject the model, not necessarily the structural value. One synthetic aggregate cannot truthfully replace four subject Context values. Module and resource boundaries do not define the logical subject.

### B8 — Provider representation mixed into Context

Proposed additions include Azure subscription ID, resource group, region, resource ID, tag-safe name, and deployment timestamp.

**Answer:** Reject all additions from canonical Context. They are representation, deployment, or runtime data rather than provider-neutral subject facts.

## Edge cases

### E1 — Identity collisions and exact comparison

| Left identity | Right identity | Same logical identity? |
| --- | --- | --- |
| `(null, "api")` | `(null, "api")` | Yes |
| `("payments", "api")` | `("payments", "api")` | Yes |
| `("payments", "api")` | `("commerce", "api")` | No |
| `(null, "api")` | `("payments", "api")` | No |
| `("payments", "api")` | `("Payments", "api")` | No |
| `("payments", "api")` | `("payments", "API")` | No |
| `("café", "api")` | `("café", "api")` | No |

The final row uses visually similar but byte-distinct Unicode forms. **Answer:** preserve and compare exact values. Context performs no case folding or Unicode normalization.

### E2 — Same logical subject across environments

```text
payments/api + development
payments/api + staging
payments/api + production
```

All three values use identity `("payments", "api")` and differ only in `environment`.

**Answer:** Accept all three as descriptions of deployments or evaluations of the same logical subject. `environment` is a descriptive deployment-context fact, not identity. A shared or pre-deployment Context with no truthful single environment omits it. Context does not define deployment identity or guarantee one value per logical identity.

This resolves the pressure found by the normal fixtures: `environment` survives only with this explicitly deployment-scoped meaning.

### E3 — Ownership transfer

Before transfer:

```hcl
ownership = {
  owned_by    = "team:payments"
  operated_by = "team:platform"
}
```

After transfer:

```hcl
ownership = {
  owned_by    = "team:commerce"
  operated_by = "team:platform"
}
```

**Answer:** Both values describe the same logical identity. Accountability changed; the subject did not. No alias, replacement identity, or generated ID is needed.

### E4 — Operator transfer or removal

Changing `operated_by` from `team:platform` to `team:sre` does not change identity. Setting it to null means no operator was explicitly declared; it does not imply `owned_by` operates the subject.

**Answer:** Accept both transitions and preserve the distinction between accountable ownership and routine operation.

### E5 — Contributing teams versus joint accountability

Case A: several teams contribute, but `team:payments` remains ultimately accountable.

**Answer:** Accept singular `owned_by = "team:payments"`. Contributors are outside this contract.

Case B: `team:payments` and `team:commerce` are independently and jointly accountable, and no durable group represents that accountability.

**Answer:** Pressure. Any singular value would discard a truthful fact. This is the only case that could justify plural ownership, but a hypothetical case alone does not justify expanding v0.1. Singular ownership survives provisionally; a concrete organizational example with an operation requiring both owners would falsify it.

### E6 — Optional absence, null, and empty strings

| Input state | Canonical result |
| --- | --- |
| Optional attribute omitted | Typed null |
| Optional attribute explicitly null | Typed null |
| Optional attribute is `""` | Reject |
| Optional attribute is whitespace-only | Reject |
| Optional attribute is a known valid value | Preserve exactly |

**Answer:** Omission and explicit null have one semantic result in the fixed canonical object. Empty is invalid and never aliases absence.

### E7 — Empty governance

| Input | Canonical `governance` |
| --- | --- |
| Omitted | Typed null |
| Explicitly null | Typed null |
| `{}` | Typed null |
| All recognized members null | Typed null |
| At least one known non-null member | Known governance object |
| At least one member unknown | Known governance object preserving that unknown |
| Whole object unknown | Typed unknown governance object |

**Answer:** An object with no facts is not a distinct semantic state. Unknown information prevents empty-object collapse because it may later become a fact.

### E8 — Unknown required and optional values

| Input state | Answer |
| --- | --- |
| Entire Context unknown | Return an unknown canonical Context; do not run semantic validation yet |
| `identity.name` unknown | Preserve unknown; validate when known |
| `ownership.owned_by` unknown | Preserve unknown; validate when known |
| Optional `environment` unknown | Preserve unknown, not null |
| `criticality` unknown | Preserve unknown; validate vocabulary when known |
| One nested member unknown, siblings known | Preserve known siblings and the unknown member |
| One known member invalid, another unknown | Reject the known invalid member immediately |

Unknown is deferred information, not absence or a Go zero value. Terraform adapter tests must prove that partially known objects survive the Framework boundary.

### E9 — Identity mutation matrix

| Mutated field | Logical identity changes? | Context remains valid if value is valid? |
| --- | --- | --- |
| `identity.name` | Yes | Yes |
| `identity.namespace` | Yes | Yes |
| `ownership.owned_by` | No | Yes |
| `ownership.operated_by` | No | Yes |
| `environment` | No | Yes |
| `governance.criticality` | No | Yes |
| `governance.cost_center` | No | Yes |
| `governance.data_classification` | No | Yes |

**Answer:** Only the identity tuple changes logical identity. No descriptive mutation silently creates a new subject.

### E10 — Azure resource topology

| Azure topology | Answer |
| --- | --- |
| One API implemented by a Container App, identity, Key Vault references, and monitoring resources | One API Context may accompany all components |
| One hub network implemented by many networking resources | One hub-network Context; resources do not each redefine identity |
| One shared Log Analytics workspace receives data from many subjects | Workspace Context describes the workspace, not every sending subject |
| One subject spans resource groups or subscriptions | Resource placement does not change logical identity |
| Azure resource is replaced while the logical subject remains | Context identity remains stable |
| Two logical subjects share one Azure resource | Each subject keeps its own Context; the shared resource does not merge them |

**Answer:** Azure validates the separation between logical subject and implementation topology. It does not require provider-specific fields or `kind`.

### E11 — Determinism and idempotence

For every accepted value `x`:

```text
C(x) = C(x)
C(C(x)) = C(x)
```

Changing time, host, process environment, Azure credentials, subscription context, filesystem, or prior calls cannot alter the result.

**Answer:** Required. Any violation is an implementation defect. Canonicalization uses no network, provider configuration, defaults, cache, or randomness.

### E12 — Schema evolution

| Proposed change | Compatibility answer |
| --- | --- |
| Add an optional field with unchanged existing meanings | Candidate additive change, but the fixed Terraform result type still changes and requires compatibility tests |
| Add a required field | Breaking |
| Remove or rename a field | Breaking |
| Change a field type | Breaking |
| Change an existing field's meaning | Breaking even if its Terraform type is unchanged |
| Add a Criticality value | Semantically additive, but requires review of consumers that assume exhaustive values |
| Add an Azure-only field | Reject; projection concern rather than evolution |

**Answer:** Only additive optional facts are plausible compatible evolution. No change is declared safe merely because Terraform can structurally convert an object. Existing callers, canonical output types, unknown handling, and consumer validation must be tested before release.

## Cross-case findings

| Question | Answer | Status after adversarial fixtures |
| --- | --- | --- |
| Does `(namespace?, name)` survive as identity? | Yes; exact comparison handles collisions without infrastructure identity | Retain |
| Does `environment` survive? | Yes, only as an optional deployment-context fact outside identity | Retain with clarified scope |
| Does singular `owned_by` survive? | Yes provisionally; genuine joint accountability remains one explicit falsifier | Retain with pressure |
| Does optional `operated_by` survive? | Yes; equality, transfer, and absence remain meaningful | Retain |
| Does closed Criticality survive? | Yes; it provides a small comparable impact domain | Retain |
| Do opaque organization-defined identifiers survive? | Yes; exact UTF-8 preservation handles all valid cases | Retain |
| Does empty governance create a state? | No; canonicalize it to typed null unless unknown information exists | Retain normalization |
| Are null and unknown semantics distinct? | Yes; null is absence and unknown is deferred information | Hard invariant |
| Is `kind` required? | No tested operation needs it | Omit |
| Are provider-specific fields required? | No; Azure topology remains implementation pressure only | Omit |
| Is a generated canonical ID required? | No operation needs one | Omit |

## Contract consequence

The proposed field set survives these fixtures without expansion. Two decisions require explicit implementation proof rather than more schema:

1. Terraform Framework tests must demonstrate preservation of partially unknown nested values and typed-null canonical output.
2. Singular ownership remains intentionally narrow until a concrete jointly accountable subject proves that one durable organizational principal is insufficient.

These fixtures support evaluating `canonicalize()` but do not prove that it earns a public API. The final falsification report must still compare its reusable validation and normalization value against ordinary Terraform object constraints and validation.

## Related decisions

- [ADR-002: One Context describes one logical subject](../adr/002-one-context-one-logical-subject.md)
- ADR-004: Logical identity is `(namespace?, name)` (superseded — SameIdentity removed)
- [ADR-004: Organization-defined identifiers remain opaque](../adr/004-opaque-semantic-identifiers.md)
- [ADR-005: Terraform Context is a structural contract](../adr/005-structural-terraform-contract.md)
- [ADR-006: Missing facts remain missing](../adr/006-missing-facts-remain-missing.md)
- [ADR-007: Provider projections are deferred](../adr/007-provider-projections-deferred.md)
