package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sojournerdev/terraform-provider-platctx/internal/platctx"
)

// Public API tests.

func TestDecodeContext(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		object  types.Object
		want    platctx.Context
		wantErr bool
	}{
		"happy path": {
			object: types.ObjectValueMust(
				map[string]attr.Type{
					platctx.KeyIdentity: types.ObjectType{AttrTypes: map[string]attr.Type{
						platctx.KeyName: types.StringType,
					}},
					platctx.KeyOwnership: types.ObjectType{AttrTypes: map[string]attr.Type{
						platctx.KeyOwnedBy: types.StringType,
					}},
				},
				map[string]attr.Value{
					platctx.KeyIdentity: types.ObjectValueMust(
						map[string]attr.Type{platctx.KeyName: types.StringType},
						map[string]attr.Value{platctx.KeyName: types.StringValue("payments")},
					),
					platctx.KeyOwnership: types.ObjectValueMust(
						map[string]attr.Type{platctx.KeyOwnedBy: types.StringType},
						map[string]attr.Value{platctx.KeyOwnedBy: types.StringValue("team:platform")},
					),
				},
			),
			want: platctx.Context{
				Identity:  platctx.Identity{Name: platctx.SubjectName("payments")},
				Ownership: platctx.Ownership{OwnedBy: platctx.OwnerRef("team:platform")},
			},
			wantErr: false,
		},
		"unsupported attribute": {
			object: types.ObjectValueMust(
				map[string]attr.Type{"unknown": types.StringType},
				map[string]attr.Value{"unknown": types.StringValue("test")},
			),
			want:    platctx.Context{},
			wantErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := decodeContext(test.object)
			if (err != nil) != test.wantErr {
				t.Errorf("decodeContext() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !test.wantErr && got.Identity.Name != test.want.Identity.Name {
				t.Errorf("decodeContext().Identity.Name = %v, want %v", got.Identity.Name, test.want.Identity.Name)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		object  types.Object
		wantErr bool
	}{
		"happy path": {
			object: types.ObjectValueMust(
				map[string]attr.Type{
					platctx.KeyIdentity: types.ObjectType{AttrTypes: map[string]attr.Type{
						platctx.KeyName: types.StringType,
					}},
					platctx.KeyOwnership: types.ObjectType{AttrTypes: map[string]attr.Type{
						platctx.KeyOwnedBy: types.StringType,
					}},
				},
				map[string]attr.Value{
					platctx.KeyIdentity: types.ObjectValueMust(
						map[string]attr.Type{platctx.KeyName: types.StringType},
						map[string]attr.Value{platctx.KeyName: types.StringValue("payments")},
					),
					platctx.KeyOwnership: types.ObjectValueMust(
						map[string]attr.Type{platctx.KeyOwnedBy: types.StringType},
						map[string]attr.Value{platctx.KeyOwnedBy: types.StringValue("team:platform")},
					),
				},
			),
			wantErr: false,
		},
		"validation error": {
			object: types.ObjectValueMust(
				map[string]attr.Type{
					platctx.KeyIdentity: types.ObjectType{AttrTypes: map[string]attr.Type{
						platctx.KeyName: types.StringType,
					}},
					platctx.KeyOwnership: types.ObjectType{AttrTypes: map[string]attr.Type{
						platctx.KeyOwnedBy: types.StringType,
					}},
				},
				map[string]attr.Value{
					platctx.KeyIdentity: types.ObjectValueMust(
						map[string]attr.Type{platctx.KeyName: types.StringType},
						map[string]attr.Value{platctx.KeyName: types.StringNull()},
					),
					platctx.KeyOwnership: types.ObjectValueMust(
						map[string]attr.Type{platctx.KeyOwnedBy: types.StringType},
						map[string]attr.Value{platctx.KeyOwnedBy: types.StringValue("team:platform")},
					),
				},
			),
			wantErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := decode(test.object)
			if (err != nil) != test.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestEncodeContext(t *testing.T) {
	t.Parallel()

	value := platctx.Context{
		Identity:  platctx.Identity{Name: platctx.SubjectName("payments")},
		Ownership: platctx.Ownership{OwnedBy: platctx.OwnerRef("team:platform")},
	}

	got := encodeContext(value)

	if got.IsNull() {
		t.Error("encodeContext() returned null")
	}

	if got.IsUnknown() {
		t.Error("encodeContext() returned unknown")
	}
}

func TestEncodeGovernance(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		governance *platctx.Governance
		wantNull   bool
	}{
		"nil governance": {
			governance: nil,
			wantNull:   true,
		},
		"with values": {
			governance: &platctx.Governance{
				Criticality: new(platctx.Criticality("critical")),
			},
			wantNull: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := encodeGovernance(test.governance)
			if test.wantNull && !got.IsNull() {
				t.Errorf("encodeGovernance() = %v, want null", got)
			}
			if !test.wantNull && got.IsNull() {
				t.Errorf("encodeGovernance() = null, want non-null")
			}
		})
	}
}

func TestFormatIssues(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		issues []platctx.ValidationIssue
		want   string
	}{
		"single issue": {
			issues: []platctx.ValidationIssue{
				{Path: "identity.name", Rule: "invalid_string", Message: "must be non-empty"},
			},
			want: "identity.name: must be non-empty",
		},
		"multiple issues": {
			issues: []platctx.ValidationIssue{
				{Path: "identity.name", Rule: "invalid_string", Message: "must be non-empty"},
				{Path: "ownership.owned_by", Rule: "invalid_string", Message: "must be non-empty"},
			},
			want: "identity.name: must be non-empty; ownership.owned_by: must be non-empty",
		},
		"empty": {
			issues: []platctx.ValidationIssue{},
			want:   "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := formatIssues(test.issues)
			if got != test.want {
				t.Errorf("formatIssues() = %q, want %q", got, test.want)
			}
		})
	}
}

// Fuzz tests.

func FuzzReadString(f *testing.F) {
	f.Add("name", "test", false, false)
	f.Add("name", "", true, false)
	f.Add("name", "test", false, true)

	f.Fuzz(func(t *testing.T, key, value string, isNull, isUnknown bool) {
		attrs := map[string]attr.Value{}

		if isNull {
			attrs[key] = types.StringNull()
		} else if isUnknown {
			attrs[key] = types.StringUnknown()
		} else {
			attrs[key] = types.StringValue(value)
		}

		got, ok := readString(attrs, key)

		if isNull || isUnknown {
			if ok {
				t.Errorf("readString() ok = true for null/unknown")
			}
		}

		if ok && got != value {
			t.Errorf("readString() = %q, want %q", got, value)
		}
	})
}

func FuzzEncodeOptionalString(f *testing.F) {
	f.Add("test")
	f.Add("")
	f.Add("production")

	f.Fuzz(func(t *testing.T, value string) {
		input := &value
		got := encodeOptionalString(input)

		if got.IsNull() {
			t.Error("encodeOptionalString() returned null for non-nil input")
		}

		strVal, ok := got.(types.String)
		if !ok {
			t.Error("encodeOptionalString() did not return a String type")
		}

		if strVal.ValueString() != value {
			t.Errorf("encodeOptionalString() = %q, want %q", strVal.ValueString(), value)
		}
	})
}
