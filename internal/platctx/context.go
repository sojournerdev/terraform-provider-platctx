package platctx

// Context describes exactly one logical subject.
//
// Its zero value is invalid and must be validated at external boundaries.
type Context struct {
	Identity    Identity
	Ownership   Ownership
	Environment *Environment
	Governance  *Governance
}

// Identity identifies a subject by the exact tuple (Namespace, Name).
type Identity struct {
	Name      SubjectName
	Namespace *Namespace
}

// Ownership keeps accountability separate from routine operation.
type Ownership struct {
	OwnedBy    OwnerRef
	OperatedBy *OperatorRef
}

// Governance holds optional subject facts and never prescribes policy.
type Governance struct {
	Criticality        *Criticality
	CostCenter         *CostCenter
	DataClassification *DataClassification
}

// SubjectName names a logical subject.
type SubjectName string

// Namespace separates logical-identity collision domains.
type Namespace string

// OwnerRef names the accountable principal.
type OwnerRef string

// OperatorRef names the routine operating principal.
type OperatorRef string

// Environment labels deployment or evaluation context and is not identity.
type Environment string

// CostCenter identifies an opaque organization-defined cost center.
type CostCenter string

// DataClassification identifies an opaque organization-defined data class.
type DataClassification string

// Criticality is the provider-owned expected-impact vocabulary.
type Criticality string

// Attribute paths for error messages.
const (
	PathIdentity                     = "identity"
	PathIdentityName                 = "identity.name"
	PathIdentityNamespace            = "identity.namespace"
	PathOwnership                    = "ownership"
	PathOwnershipOwnedBy             = "ownership.owned_by"
	PathOwnershipOperatedBy          = "ownership.operated_by"
	PathEnvironment                  = "environment"
	PathGovernance                   = "governance"
	PathGovernanceCriticality        = "governance.criticality"
	PathGovernanceCostCenter         = "governance.cost_center"
	PathGovernanceDataClassification = "governance.data_classification"
)

// Attribute keys for Terraform schema lookups.
const (
	KeyIdentity           = "identity"
	KeyName               = "name"
	KeyNamespace          = "namespace"
	KeyOwnership          = "ownership"
	KeyOwnedBy            = "owned_by"
	KeyOperatedBy         = "operated_by"
	KeyEnvironment        = "environment"
	KeyGovernance         = "governance"
	KeyCriticality        = "criticality"
	KeyCostCenter         = "cost_center"
	KeyDataClassification = "data_classification"
)

const (
	// CriticalityLow is the low expected-impact level.
	CriticalityLow Criticality = "low"

	// CriticalityMedium is the medium expected-impact level.
	CriticalityMedium Criticality = "medium"

	// CriticalityHigh is the high expected-impact level.
	CriticalityHigh Criticality = "high"

	// CriticalityCritical is the critical expected-impact level.
	CriticalityCritical Criticality = "critical"
)
