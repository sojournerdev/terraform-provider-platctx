# Runtime and Compatibility Support

## Purpose

Define supported versions and compatibility rules before implementation freezes
the provider surface.

## Initial support policy

| Component | v0.1 policy |
| --- | --- |
| Terraform CLI | 1.8 or later |
| Terraform Plugin Framework | Pin an exact version in `go.mod` |
| Go | Pin an exact toolchain in `go.mod` |
| Provider protocol | Use the protocol selected by the pinned Framework version |
| Provider configuration | None |
| Resources | None |
| Data sources | None |
| Network and filesystem access | None |
| Candidate function input | Dynamic parameter with null and unknown handling enabled |
| Candidate function result | Fixed typed Context object |

The implementation MUST record the exact Framework and Go versions in the
repository. Documentation MUST be updated when either pin changes.

## Provider identity

Before release, record all of the following in the provider implementation and
release metadata:

- provider source address;
- provider namespace and type name;
- qualified function name;
- protocol compatibility;
- supported operating systems and architectures; and
- minimum Terraform CLI version.

The candidate operation is referred to as `canonicalize(context)` until the
provider source address fixes its qualified Terraform name.

## Compatibility rules

The following changes are breaking unless a later compatibility decision says
otherwise:

- removing or renaming an attribute;
- changing an attribute type;
- making an optional attribute required;
- changing the meaning of an existing attribute;
- changing null or unknown behavior; and
- changing identity comparison.

Adding an optional fact MAY be compatible, but requires tests for existing
canonical output, unknown handling, structural conversion, and all consumers.
Adding a Criticality value requires consumer review because consumers may use
exhaustive comparisons.

Provider-specific fields, projections, defaults, and extension maps are not
compatibility mechanisms. They remain outside the canonical contract.

## Release gate

A release requires:

1. supported-version tests for the minimum Terraform CLI;
2. Framework and Go versions pinned and reproducible;
3. provider protocol smoke coverage;
4. conformance and falsification evidence;
5. no undocumented breaking contract change; and
6. updated ADR status for decisions in force.

## Related documents

- [Context API v0.1](context-api-v0.1.md)
- [Terraform boundary contract](terraform-boundary-contract.md)
- [Falsification plan](context-contract-falsification.md)
- [Provider engineering standards](provider-engineering-standards.md)
- [Security and resource limits](security-and-resource-limits.md)
- [Provider release readiness](provider-release-readiness.md)
- [ADR-001: Function-only Terraform provider](../adr/001-function-only-terraform-provider.md)
- [ADR-007: Provider projections are deferred](../adr/007-provider-projections-deferred.md)
