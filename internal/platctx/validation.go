package platctx

import (
	"cmp"
	"slices"
	"unicode"
	"unicode/utf8"
)

const (
	RuleInvalidUTF8        = "invalid_utf8"
	RuleInvalidString      = "invalid_string"
	RuleInvalidWhitespace  = "invalid_whitespace"
	RuleInvalidCriticality = "invalid_criticality"
)

type ValidationIssue struct {
	Path    string
	Rule    string
	Message string
}

func Validate(value Context) []ValidationIssue {
	issues := value.Identity.validate(nil)
	issues = value.Ownership.validate(issues)

	if value.Environment != nil {
		issues = appendStringIssue(issues, PathEnvironment, string(*value.Environment))
	}

	if value.Governance != nil {
		issues = value.Governance.validate(issues)
	}

	slices.SortFunc(issues, func(a, b ValidationIssue) int {
		if c := cmp.Compare(a.Path, b.Path); c != 0 {
			return c
		}

		return cmp.Compare(a.Rule, b.Rule)
	})

	return issues
}

func (identity Identity) validate(issues []ValidationIssue) []ValidationIssue {
	issues = appendStringIssue(issues, PathIdentityName, string(identity.Name))

	if identity.Namespace != nil {
		issues = appendStringIssue(issues, PathIdentityNamespace, string(*identity.Namespace))
	}

	return issues
}

func (ownership Ownership) validate(issues []ValidationIssue) []ValidationIssue {
	issues = appendStringIssue(issues, PathOwnershipOwnedBy, string(ownership.OwnedBy))

	if ownership.OperatedBy != nil {
		issues = appendStringIssue(issues, PathOwnershipOperatedBy, string(*ownership.OperatedBy))
	}

	return issues
}

func (governance Governance) validate(issues []ValidationIssue) []ValidationIssue {
	if governance.Criticality != nil {
		path := PathGovernanceCriticality
		value := string(*governance.Criticality)

		if !isValidString(value) {
			issues = appendStringIssue(issues, path, value)
		} else if !isCriticality(*governance.Criticality) {
			issues = append(issues, ValidationIssue{
				Path:    path,
				Rule:    RuleInvalidCriticality,
				Message: "must be one of: low, medium, high, critical",
			})
		}
	}

	if governance.CostCenter != nil {
		issues = appendStringIssue(issues, PathGovernanceCostCenter, string(*governance.CostCenter))
	}

	if governance.DataClassification != nil {
		issues = appendStringIssue(
			issues,
			PathGovernanceDataClassification,
			string(*governance.DataClassification),
		)
	}

	return issues
}

func appendStringIssue(issues []ValidationIssue, path, value string) []ValidationIssue {
	if !utf8.ValidString(value) {
		return append(issues, ValidationIssue{
			Path:    path,
			Rule:    RuleInvalidUTF8,
			Message: "must be valid UTF-8",
		})
	}

	if value == "" {
		return append(issues, ValidationIssue{
			Path:    path,
			Rule:    RuleInvalidString,
			Message: "must be non-empty",
		})
	}

	first, _ := utf8.DecodeRuneInString(value)
	last, _ := utf8.DecodeLastRuneInString(value)

	if unicode.IsSpace(first) || unicode.IsSpace(last) {
		return append(issues, ValidationIssue{
			Path:    path,
			Rule:    RuleInvalidWhitespace,
			Message: "must be free of leading or trailing whitespace",
		})
	}

	return issues
}

func isValidString(value string) bool {
	if !utf8.ValidString(value) {
		return false
	}

	if value == "" {
		return false
	}

	first, _ := utf8.DecodeRuneInString(value)
	last, _ := utf8.DecodeLastRuneInString(value)

	return !unicode.IsSpace(first) && !unicode.IsSpace(last)
}

func isCriticality(value Criticality) bool {
	switch value {
	case CriticalityLow, CriticalityMedium, CriticalityHigh, CriticalityCritical:
		return true
	default:
		return false
	}
}
