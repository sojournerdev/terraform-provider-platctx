package platctx

import (
	"reflect"
	"testing"
)

func TestCanonicalize(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		value Context
		want  Context
	}{
		// Accepts.

		"preserves valid Context": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
			},
			want: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
			},
		},
		"preserves fully populated Context": {
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
			want: Context{
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
		"preserves non-empty Governance": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					Criticality: new(CriticalityLow),
				},
			},
			want: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
				Governance: &Governance{
					Criticality: new(CriticalityLow),
				},
			},
		},

		// Rejects.

		"normalizes empty Governance to nil": {
			value: Context{
				Identity:   Identity{Name: subjectName},
				Ownership:  Ownership{OwnedBy: owner},
				Governance: &Governance{},
			},
			want: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
			},
		},
		"preserves nil Governance": {
			value: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
			},
			want: Context{
				Identity:  Identity{Name: subjectName},
				Ownership: Ownership{OwnedBy: owner},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := Canonicalize(test.value)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Canonicalize() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestCanonicalizeIdempotent(t *testing.T) {
	t.Parallel()

	value := Context{
		Identity:   Identity{Name: subjectName},
		Ownership:  Ownership{OwnedBy: owner},
		Governance: &Governance{},
	}

	first := Canonicalize(value)
	second := Canonicalize(first)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("Canonicalize() not idempotent: first = %#v, second = %#v", first, second)
	}

	if first.Governance != nil {
		t.Fatalf("Canonicalize() did not normalize empty Governance: got %#v", first.Governance)
	}
}

func FuzzCanonicalizeDeterministic(f *testing.F) {
	for _, seed := range []struct {
		name          string
		governanceNil bool
	}{
		{name: "payments", governanceNil: true},
		{name: "payments", governanceNil: false},
		{name: "", governanceNil: true},
	} {
		f.Add(seed.name, seed.governanceNil)
	}

	f.Fuzz(func(t *testing.T, name string, governanceNil bool) {
		value := Context{
			Identity:  Identity{Name: SubjectName(name)},
			Ownership: Ownership{OwnedBy: owner},
		}

		if !governanceNil {
			value.Governance = &Governance{}
		}

		first := Canonicalize(value)
		second := Canonicalize(first)

		if !reflect.DeepEqual(first, second) {
			t.Fatalf("Canonicalize() not idempotent: first = %#v, second = %#v", first, second)
		}

		if value.Identity.Name != SubjectName(name) {
			t.Fatalf("Canonicalize() mutated input")
		}
	})
}
