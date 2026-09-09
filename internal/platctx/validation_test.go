package platctx

import (
	"slices"
	"testing"
)

const (
	invalidUTF8        = string("\xff")
	whitespaceOnly     = "   "
	leadingWhitespace  = " payments"
	trailingWhitespace = "payments "
)

// Maps rules to expected messages for test assertions.
var ruleMessages = map[string]string{
	RuleInvalidUTF8:        "must be valid UTF-8",
	RuleInvalidString:      "must be non-empty",
	RuleInvalidWhitespace:  "must be free of leading or trailing whitespace",
	RuleInvalidCriticality: "must be one of: low, medium, high, critical",
}

func issue(path, rule, message string) ValidationIssue {
	return ValidationIssue{Path: path, Rule: rule, Message: message}
}

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		value Context
		want  []ValidationIssue
	}{
		// Accepts.

		"accepts minimal Context": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
			},
		},
		"accepts fully populated Context": {
			value: Context{
				Identity: Identity{
					Name:      subjectName,
					Namespace: new(subjectNS),
				},
				Ownership: Ownership{
					OwnedBy:    owner,
					OperatedBy: new(operator),
				},
				Environment: new(environment),
				Governance: &Governance{
					Criticality:        new(CriticalityCritical),
					CostCenter:         new(costCenter),
					DataClassification: new(classification),
				},
			},
		},
		"accepts equal owner and operator": {
			value: Context{
				Identity: Identity{Name: subjectName},
				Ownership: Ownership{
					OwnedBy:    owner,
					OperatedBy: new(OperatorRef(owner)),
				},
			},
		},
		"accepts low Criticality": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					Criticality: new(CriticalityLow),
				},
			},
		},
		"accepts medium Criticality": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					Criticality: new(CriticalityMedium),
				},
			},
		},
		"accepts high Criticality": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					Criticality: new(CriticalityHigh),
				},
			},
		},
		"accepts opaque Unicode": {
			value: Context{
				Identity: Identity{
					Name:      SubjectName("platform-東京"),
					Namespace: new(Namespace("europe")),
				},
				Ownership:   Ownership{OwnedBy: OwnerRef("team:platform")},
				Environment: new(Environment("pre production")),
				Governance: &Governance{
					CostCenter:         new(CostCenter("FIN/東京-42")),
					DataClassification: new(DataClassification("internal—restricted")),
				},
			},
		},
		"accepts empty governance": {
			value: Context{
				Identity:   Identity{Name: subjectName},
				Ownership:  Ownership{OwnedBy: owner},
				Governance: &Governance{},
			},
		},
		"accepts shared infrastructure with no environment": {
			value: Context{
				Identity:  Identity{Name: "shared-platform"},
				Ownership: Ownership{OwnedBy: OwnerRef("team:platform")},
			},
		},

		// Rejects.

		"rejects empty name": {
			value: Context{
				Identity:  Identity{Name: ""},
				Ownership: Ownership{OwnedBy: owner},
			},
			want: []ValidationIssue{
				issue(PathIdentityName, RuleInvalidString, ruleMessages[RuleInvalidString]),
			},
		},
		"rejects empty owner": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: ""},
			},
			want: []ValidationIssue{
				issue(PathOwnershipOwnedBy, RuleInvalidString, ruleMessages[RuleInvalidString]),
			},
		},
		"rejects whitespace-only name": {
			value: Context{
				Identity:  Identity{Name: SubjectName(whitespaceOnly)},
				Ownership: Ownership{OwnedBy: owner},
			},
			want: []ValidationIssue{
				issue(PathIdentityName, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
			},
		},
		"rejects whitespace-only owner": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: OwnerRef(whitespaceOnly)},
			},
			want: []ValidationIssue{
				issue(PathOwnershipOwnedBy, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
			},
		},
		"rejects leading whitespace in name": {
			value: Context{
				Identity:  Identity{Name: SubjectName(leadingWhitespace)},
				Ownership: Ownership{OwnedBy: owner},
			},
			want: []ValidationIssue{
				issue(PathIdentityName, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
			},
		},
		"rejects trailing whitespace in name": {
			value: Context{
				Identity:  Identity{Name: SubjectName(trailingWhitespace)},
				Ownership: Ownership{OwnedBy: owner},
			},
			want: []ValidationIssue{
				issue(PathIdentityName, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
			},
		},
		"rejects leading whitespace in owner": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: OwnerRef(leadingWhitespace)},
			},
			want: []ValidationIssue{
				issue(PathOwnershipOwnedBy, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
			},
		},
		"rejects trailing whitespace in owner": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: OwnerRef(trailingWhitespace)},
			},
			want: []ValidationIssue{
				issue(PathOwnershipOwnedBy, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
			},
		},
		"rejects leading whitespace in environment": {
			value: Context{
				Identity:    Identity{Name: subjectName},
				Ownership:   Ownership{OwnedBy: owner},
				Environment: new(Environment(leadingWhitespace)),
			},
			want: []ValidationIssue{
				issue(PathEnvironment, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
			},
		},
		"rejects trailing whitespace in environment": {
			value: Context{
				Identity:    Identity{Name: subjectName},
				Ownership:   Ownership{OwnedBy: owner},
				Environment: new(Environment(trailingWhitespace)),
			},
			want: []ValidationIssue{
				issue(PathEnvironment, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
			},
		},
		"rejects empty optional string": {
			value: Context{
				Identity:    Identity{Name: subjectName},
				Ownership:   Ownership{OwnedBy: owner},
				Environment: new(Environment("")),
			},
			want: []ValidationIssue{
				issue(PathEnvironment, RuleInvalidString, ruleMessages[RuleInvalidString]),
			},
		},
		"rejects invalid UTF-8 in name": {
			value: Context{
				Identity:  Identity{Name: SubjectName(invalidUTF8)},
				Ownership: Ownership{OwnedBy: owner},
			},
			want: []ValidationIssue{
				issue(PathIdentityName, RuleInvalidUTF8, ruleMessages[RuleInvalidUTF8]),
			},
		},
		"rejects invalid UTF-8 in owner": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: OwnerRef(invalidUTF8)},
			},
			want: []ValidationIssue{
				issue(PathOwnershipOwnedBy, RuleInvalidUTF8, ruleMessages[RuleInvalidUTF8]),
			},
		},
		"rejects invalid UTF-8 in namespace": {
			value: Context{
				Identity: Identity{
					Name:      subjectName,
					Namespace: new(Namespace(invalidUTF8)),
				},
				Ownership: Ownership{OwnedBy: owner},
			},
			want: []ValidationIssue{
				issue(PathIdentityNamespace, RuleInvalidUTF8, ruleMessages[RuleInvalidUTF8]),
			},
		},
		"rejects invalid UTF-8 in environment": {
			value: Context{
				Identity:    Identity{Name: subjectName},
				Ownership:   Ownership{OwnedBy: owner},
				Environment: new(Environment(invalidUTF8)),
			},
			want: []ValidationIssue{
				issue(PathEnvironment, RuleInvalidUTF8, ruleMessages[RuleInvalidUTF8]),
			},
		},
		"rejects unsupported Criticality": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					Criticality: new(Criticality("severe")),
				},
			},
			want: []ValidationIssue{
				issue(PathGovernanceCriticality, RuleInvalidCriticality, ruleMessages[RuleInvalidCriticality]),
			},
		},
		"rejects uppercase Criticality": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					Criticality: new(Criticality("HIGH")),
				},
			},
			want: []ValidationIssue{
				issue(PathGovernanceCriticality, RuleInvalidCriticality, ruleMessages[RuleInvalidCriticality]),
			},
		},
		"rejects empty Criticality": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					Criticality: new(Criticality("")),
				},
			},
			want: []ValidationIssue{
				issue(PathGovernanceCriticality, RuleInvalidString, ruleMessages[RuleInvalidString]),
			},
		},
		"rejects invalid UTF-8 Criticality once": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					Criticality: new(Criticality(invalidUTF8)),
				},
			},
			want: []ValidationIssue{
				issue(PathGovernanceCriticality, RuleInvalidUTF8, ruleMessages[RuleInvalidUTF8]),
			},
		},
		"rejects trailing non-breaking space in classification": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					DataClassification: new(DataClassification("internal\u00a0")),
				},
			},
			want: []ValidationIssue{
				issue(PathGovernanceDataClassification, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
			},
		},
		"rejects invalid UTF-8 cost center": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					CostCenter: new(CostCenter(invalidUTF8)),
				},
			},
			want: []ValidationIssue{
				issue(PathGovernanceCostCenter, RuleInvalidUTF8, ruleMessages[RuleInvalidUTF8]),
			},
		},
		"rejects multiple independent failures in order": {
			value: Context{
				Identity:    Identity{Name: ""},
				Ownership:   Ownership{OwnedBy: OwnerRef(whitespaceOnly)},
				Environment: new(Environment(leadingWhitespace)),
				Governance: &Governance{
					Criticality:        new(Criticality("HIGH")),
					CostCenter:         new(CostCenter(invalidUTF8)),
					DataClassification: new(DataClassification("internal\u00a0")),
				},
			},
			want: []ValidationIssue{
				issue(PathEnvironment, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
				issue(PathGovernanceCostCenter, RuleInvalidUTF8, ruleMessages[RuleInvalidUTF8]),
				issue(PathGovernanceCriticality, RuleInvalidCriticality, ruleMessages[RuleInvalidCriticality]),
				issue(PathGovernanceDataClassification, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
				issue(PathIdentityName, RuleInvalidString, ruleMessages[RuleInvalidString]),
				issue(PathOwnershipOwnedBy, RuleInvalidWhitespace, ruleMessages[RuleInvalidWhitespace]),
			},
		},

		// Edge cases.

		"near-identity Unicode remains distinct": {
			value: Context{
				Identity:  Identity{Name: SubjectName("platførm")},
				Ownership: Ownership{OwnedBy: owner},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			issues := Validate(test.value)

			if !slices.Equal(issues, test.want) {
				t.Errorf("Validate() = %#v, want %#v", issues, test.want)
			}
		})
	}
}

func TestValidateDeterministic(t *testing.T) {
	t.Parallel()

	value := Context{
		Identity:  Identity{Name: ""},
		Ownership: Ownership{OwnedBy: ""},
	}

	first := Validate(value)
	second := Validate(value)

	if !slices.Equal(first, second) {
		t.Fatalf("Validate() differs: first = %#v, second = %#v", first, second)
	}
}

func FuzzValidateDeterministic(f *testing.F) {
	for _, seed := range []struct {
		name        string
		owner       string
		environment string
	}{
		{name: "payments", owner: "team:platform", environment: "production"},
		{name: "", owner: "", environment: ""},
		{name: "\u00a0payments", owner: "team:platform", environment: "production"},
		{name: invalidUTF8, owner: "team:platform", environment: "production"},
	} {
		f.Add(seed.name, seed.owner, seed.environment)
	}

	f.Fuzz(func(t *testing.T, name, owner, environment string) {
		value := Context{
			Identity:    Identity{Name: SubjectName(name)},
			Ownership:   Ownership{OwnedBy: OwnerRef(owner)},
			Environment: new(Environment(environment)),
		}

		first := Validate(value)
		second := Validate(value)

		if !slices.Equal(first, second) {
			t.Fatalf("Validate() differs: first = %#v, second = %#v", first, second)
		}

		if value.Identity.Name != SubjectName(name) ||
			value.Ownership.OwnedBy != OwnerRef(owner) ||
			value.Environment == nil ||
			*value.Environment != Environment(environment) {
			t.Fatalf("Validate() mutated input: got %#v", value)
		}
	})
}

func TestValidationIssueMessages(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		value     Context
		wantRules []string
		wantMsgs  []string
	}{
		"empty name has message": {
			value: Context{
				Identity:  Identity{Name: ""},
				Ownership: Ownership{OwnedBy: owner},
			},
			wantRules: []string{RuleInvalidString},
			wantMsgs:  []string{"must be non-empty"},
		},
		"invalid UTF-8 has message": {
			value: Context{
				Identity:  Identity{Name: SubjectName(invalidUTF8)},
				Ownership: Ownership{OwnedBy: owner},
			},
			wantRules: []string{RuleInvalidUTF8},
			wantMsgs:  []string{"must be valid UTF-8"},
		},
		"invalid criticality has message": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					Criticality: new(Criticality("severe")),
				},
			},
			wantRules: []string{RuleInvalidCriticality},
			wantMsgs:  []string{"must be one of: low, medium, high, critical"},
		},
		"whitespace has message": {
			value: Context{
				Identity:  Identity{Name: SubjectName(leadingWhitespace)},
				Ownership: Ownership{OwnedBy: owner},
			},
			wantRules: []string{RuleInvalidWhitespace},
			wantMsgs:  []string{"must be free of leading or trailing whitespace"},
		},
		"multiple issues have messages": {
			value: Context{
				Identity:  Identity{Name: ""},
				Ownership: Ownership{OwnedBy: OwnerRef(whitespaceOnly)},
			},
			wantRules: []string{RuleInvalidString, RuleInvalidWhitespace},
			wantMsgs: []string{
				"must be non-empty",
				"must be free of leading or trailing whitespace",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			issues := Validate(test.value)

			if len(issues) != len(test.wantRules) {
				t.Fatalf("Validate() returned %d issues, want %d", len(issues), len(test.wantRules))
			}

			for i, issue := range issues {
				if issue.Rule != test.wantRules[i] {
					t.Errorf("issue[%d].Rule = %q, want %q", i, issue.Rule, test.wantRules[i])
				}
				if issue.Message == "" {
					t.Errorf("issue[%d].Message is empty, want non-empty", i)
				}
				if issue.Message != test.wantMsgs[i] {
					t.Errorf("issue[%d].Message = %q, want %q", i, issue.Message, test.wantMsgs[i])
				}
			}
		})
	}
}
