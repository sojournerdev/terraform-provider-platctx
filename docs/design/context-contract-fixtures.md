# Context Contract Normal Fixtures

## Purpose

These fixtures test whether the proposed contract describes application and infrastructure subjects without fake, policy, or provider-specific facts.

Each fixture states its facts first. Candidate Context values contain only those facts. HCL examples omit absent optional fields; a fixed Terraform result would use typed nulls.
Azure provides familiar implementation pressure for these fixtures. Each representative Azure footprint is non-normative: it illustrates resources that may implement the subject, but neither those resources nor their configuration become Context facts.

Verdicts:

- **Pass:** fits without semantic strain.
- **Pressure:** fits but exposes an unresolved meaning.
- **Fail:** needs a fake value, overloaded field, or missing concept.

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

Only `(namespace?, name)` defines logical identity.

## Fixtures

### payments-api — Pass

Facts: logical subject `payments/api`; accountable owner `team:payments`; no other facts declared.
Representative Azure footprint: an App Service or Container App, managed identity, Key Vault references, and Application Insights.

```hcl
{
  identity  = { namespace = "payments", name = "api" }
  ownership = { owned_by = "team:payments" }
}
```

No service-specific or deployment field is required.

### payments-worker — Pressure

Facts: logical subject `payments/worker`; owner `team:payments`; evaluated for `production`; high expected impact.
Representative Azure footprint: a Container Apps job or Function App, Service Bus subscription, managed identity, and Application Insights.

```hcl
{
  identity    = { namespace = "payments", name = "worker" }
  ownership   = { owned_by = "team:payments" }
  environment = "production"
  governance  = { criticality = "high" }
}
```

The worker fits, but `environment` describes its deployment while identity describes the logical subject.

### payments-database — Pressure

Facts: logical subject `payments/database`; owner `team:payments`; operator `team:database-platform`; production deployment; critical impact; cost center `CC-1042`; classification `restricted`.
Representative Azure footprint: an Azure Database for PostgreSQL flexible server, private endpoint, private DNS integration, and diagnostic settings.

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

All fields are truthful. `environment` has the same scope pressure as the worker.

### hub-network — Pass

Facts: shared subject `hub-network`; owner `team:network-platform`; operator `team:network-operations`; high expected impact; no single environment.
Representative Azure footprint: a virtual network, subnets, route tables, Azure Firewall, and VPN or ExpressRoute gateways.

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

Shared infrastructure needs no fake environment or workload field.

### central-dns — Pass

Facts: global subject `central-dns`; owner `team:network-platform`; operator `team:network-operations`; critical expected impact.
Representative Azure footprint: Azure DNS zones, Private DNS zones and links, DNS Private Resolver endpoints, and forwarding rulesets.

```hcl
{
  identity = { name = "central-dns" }
  ownership = {
    owned_by    = "team:network-platform"
    operated_by = "team:network-operations"
  }
  governance = { criticality = "critical" }
}
```

DNS needs no environment, provider field, or subject taxonomy.

### shared-kubernetes — Pass

Facts: shared platform `shared-kubernetes`; owner `team:platform`; operator `team:sre`; critical expected impact; no single environment.
Representative Azure footprint: an Azure Kubernetes Service cluster, agent pools, managed identities, private DNS integration, and monitoring integration.

```hcl
{
  identity = { name = "shared-kubernetes" }
  ownership = {
    owned_by    = "team:platform"
    operated_by = "team:sre"
  }
  governance = { criticality = "critical" }
}
```

The platform needs no Kubernetes-specific or hosted-workload facts.

### monitoring-workspace — Pass

Facts: subject `observability/workspace`; owner and operator `team:observability`; high expected impact; spans environments.
Representative Azure footprint: a Log Analytics workspace, data collection rules, diagnostic destinations, and workspace-based Application Insights.

```hcl
{
  identity = { namespace = "observability", name = "workspace" }
  ownership = {
    owned_by    = "team:observability"
    operated_by = "team:observability"
  }
  governance = { criticality = "high" }
}
```

Equal owner and operator values are explicit, not inferred.

### terraform-state — Pass

Facts: subject `platform/terraform-state`; owner and operator `team:platform`; high expected impact; classification `confidential`; backend not yet selected.
Representative Azure footprint: a storage account and blob container, private endpoint, private DNS integration, and role assignments.

```hcl
{
  identity = { namespace = "platform", name = "terraform-state" }
  ownership = {
    owned_by    = "team:platform"
    operated_by = "team:platform"
  }
  governance = {
    criticality         = "high"
    data_classification = "confidential"
  }
}
```

The subject needs no backend, retention, encryption, or provider facts.

### subscription-foundation — Pass

Facts: subject `subscription-foundation`; owner `team:cloud-platform`; operator `team:cloud-operations`; high expected impact; cost center `FIN-PLATFORM`; infrastructure not yet created.
Representative Azure footprint: an Azure subscription, resource groups, role assignments, policy assignments, budgets, and activity-log diagnostics.

```hcl
{
  identity = { name = "subscription-foundation" }
  ownership = {
    owned_by    = "team:cloud-platform"
    operated_by = "team:cloud-operations"
  }
  governance = {
    criticality = "high"
    cost_center = "FIN-PLATFORM"
  }
}
```

A cloud-related subject does not require provider-specific Context fields.

### data-lake — Pressure

Facts: subject `data-lake`; owner and operator `team:data-platform`; production deployment; high expected impact; cost center `US004210`; classification `restricted`; infrastructure not yet selected.
Representative Azure footprint: an Azure Data Lake Storage Gen2 account, containers, private endpoints, role assignments, and diagnostic settings.

```hcl
{
  identity = { name = "data-lake" }
  ownership = {
    owned_by    = "team:data-platform"
    operated_by = "team:data-platform"
  }
  environment = "production"
  governance = {
    criticality         = "high"
    cost_center         = "US004210"
    data_classification = "restricted"
  }
}
```

Governance fits before implementation. `environment` remains deployment-scoped.

## Audit

| Fixture | Truthful | Fake required fact | Optional facts honest | Meaning changes | Duplicated information | Workload assumption | Needs `kind` | Contains policy | Provider-specific schema | Exists before infrastructure |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| payments-api | Yes | No | Yes | No | No | No | No | No | No | Yes |
| payments-worker | Yes | No | Yes | `environment` pressure | No | No | No | No | No | Yes |
| payments-database | Yes | No | Yes | `environment` pressure | No | No | No | No | No | Yes |
| hub-network | Yes | No | Yes | No | No | No | No | No | No | Yes |
| central-dns | Yes | No | Yes | No | No | No | No | No | No | Yes |
| shared-kubernetes | Yes | No | Yes | No | No | No | No | No | No | Yes |
| monitoring-workspace | Yes | No | Yes | No | Equal values, distinct facts | No | No | No | No | Yes |
| terraform-state | Yes | No | Yes | No | Equal values, distinct facts | No | No | No | No | Yes |
| subscription-foundation | Yes | No | Yes | No | No | No | No | No | No | Yes |
| data-lake | Yes | No | Yes | `environment` pressure | Equal values, distinct facts | No | No | No | No | Yes |

## Findings

- `name` and `owned_by` are truthful in every fixture.
- `namespace` remains an identity collision domain.
- `owned_by` and `operated_by` keep distinct meanings. Joint ownership passes to the adversarial fixtures.
- `environment` is the only pressure: it describes a deployment while identity describes the logical subject.
- Governance fields remain facts and are honestly omitted when they do not apply.
- These fixtures do not prove that Context Core should own a closed Criticality vocabulary.
- No fixture needs `kind` or a provider-specific field.
- Every subject can exist before its infrastructure.

## Result

Seven fixtures pass, three expose `environment` pressure, and none fail.

The [adversarial fixtures](context-contract-adversarial-fixtures.md) test same-subject deployments, joint ownership, unknown values, empty objects, and the Criticality vocabulary before the contract is frozen.
