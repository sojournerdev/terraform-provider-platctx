package platctx

// Canonicalize returns the canonical form of value.
//
// Known facts are preserved exactly. Omitted optional facts remain nil. Empty
// Governance is normalized to nil. No defaults or generated values are added.
//
// Deterministic and idempotent: C(C(x)) equals C(x).
func Canonicalize(value Context) Context {
	value.Governance = canonicalizeGovernance(value.Governance)

	return value
}

func canonicalizeGovernance(governance *Governance) *Governance {
	if governance == nil {
		return nil
	}

	if governance.Criticality == nil &&
		governance.CostCenter == nil &&
		governance.DataClassification == nil {
		return nil
	}

	return governance
}
