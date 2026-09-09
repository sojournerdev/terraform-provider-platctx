# Consumer Use Cases

## Purpose

Show where a deterministic Context function could help callers. These are
falsification targets, not proof that the provider API should ship.

## Caller-owned locals

A caller keeps Context facts in Terraform locals and invokes the pure function
at the boundary before passing the value to other modules:

```hcl
locals {
  context = provider::context::canonicalize({
    identity = {
      namespace = "payments"
      name      = "api"
    }
    ownership = {
      owned_by = "team:payments"
    }
  })
}

module "service" {
  source  = "./service"
  context = local.context
}
```

The qualified function name is illustrative until the provider source address
is fixed.

The caller owns the facts. The function does not discover, infer, default, or
manage anything. It provides one deterministic boundary that validates the
local value and emits the fixed canonical shape.

## Workflow A: reusable module input

A root module defines a Context once and passes it to multiple child modules.
Each child receives the same typed-null and validation behavior instead of
reimplementing it.

Evidence required:

- two child modules consume the same local;
- both receive equal canonical values; and
- invalid input fails at the caller boundary with a useful path.

## Workflow B: known and plan-time values

A caller builds a Context from a mixture of literals, variables, and resource
outputs. Known facts are validated immediately. Unknown facts remain unknown
until Terraform can determine them.

Evidence required:

- known invalid facts are not hidden by unknown siblings;
- unknown facts do not become null or empty strings; and
- the result remains valid for downstream module arguments.

## Baseline comparison

Compare the function with a caller that uses:

- a typed local or module variable;
- ordinary Terraform object constraints; and
- consumer-local validation.

The function is justified only if it gives callers a clear, reusable contract
that the baseline cannot provide without repeated or inconsistent rules.

## Non-goals

These use cases do not justify:

- provider configuration;
- remote lookups;
- resource lifecycle;
- policy completeness checks;
- infrastructure projections; or
- a nominal Terraform Context type.

## Related documents

- [Context API v0.1](context-api-v0.1.md)
- [Falsification plan](context-contract-falsification.md)
- [Terraform boundary contract](terraform-boundary-contract.md)
