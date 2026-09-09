package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/sojournerdev/terraform-provider-platctx/internal/platctx"
)

var _ function.Function = &CanonicalizeFunction{}

// CanonicalizeFunction adapts domain logic to Terraform Plugin Framework.
type CanonicalizeFunction struct{}

// NewCanonicalizeFunction is the factory for provider registration.
func NewCanonicalizeFunction() function.Function {
	return &CanonicalizeFunction{}
}

// Metadata registers the function name with Terraform.
func (f *CanonicalizeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "canonicalize"
}

// Definition defines the function signature for HCL consumption.
func (f *CanonicalizeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Canonicalizes platform context for service metadata",
		MarkdownDescription: "Validates and normalizes service metadata into a consistent typed structure for platform engineering.",
		Parameters: []function.Parameter{
			function.DynamicParameter{
				Name:                "context",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				MarkdownDescription: "Service metadata to canonicalize (identity, ownership, environment, governance)",
			},
		},
		Return: function.DynamicReturn{},
	}
}

// Run is the entry point: validates input, normalizes, returns typed result.
func (f *CanonicalizeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input types.Dynamic

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &input))
	if resp.Error != nil {
		return
	}

	if input.IsNull() {
		resp.Error = function.NewArgumentFuncError(0, "context must not be null. Provide a Context object with identity and ownership attributes.")
		return
	}

	if input.IsUnknown() {
		resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, unknownResult))
		return
	}

	object, ok := input.UnderlyingValue().(types.Object)
	if !ok {
		resp.Error = function.NewArgumentFuncError(0, "context must be an object with identity and ownership attributes")
		return
	}

	value, err := decode(object)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, encode(value)))
}

// decode transforms Terraform types to domain types and validates at the boundary.
func decode(object types.Object) (platctx.Context, error) {
	value, err := decodeContext(object)
	if err != nil {
		return platctx.Context{}, err
	}

	issues := platctx.Validate(value)
	if len(issues) > 0 {
		return platctx.Context{}, fmt.Errorf("%s", formatIssues(issues))
	}

	return value, nil
}

// encode normalizes and transforms domain types back to Terraform types.
func encode(value platctx.Context) *types.Dynamic {
	canonical := platctx.Canonicalize(value)
	return encodeContext(canonical)
}

func decodeContext(object types.Object) (platctx.Context, error) {
	attrs := object.Attributes()

	for key := range attrs {
		if _, ok := supportedAttributes[key]; !ok {
			return platctx.Context{}, fmt.Errorf("unsupported attribute: %s", key)
		}
	}

	identity, err := decodeIdentity(attrs)
	if err != nil {
		return platctx.Context{}, err
	}

	ownership, err := decodeOwnership(attrs)
	if err != nil {
		return platctx.Context{}, err
	}

	return platctx.Context{
		Identity:    identity,
		Ownership:   ownership,
		Environment: decodeEnvironment(attrs),
		Governance:  decodeGovernance(attrs),
	}, nil
}

func decodeIdentity(attrs map[string]attr.Value) (platctx.Identity, error) {
	identityAttrs, err := readObject(attrs, platctx.KeyIdentity)
	if err != nil {
		return platctx.Identity{}, fmt.Errorf("missing required object: %s", platctx.KeyIdentity)
	}

	name, ok := readRequiredString(identityAttrs, platctx.KeyName)
	if !ok {
		return platctx.Identity{}, fmt.Errorf("missing required value: %s.%s", platctx.KeyIdentity, platctx.KeyName)
	}

	identity := platctx.Identity{Name: platctx.SubjectName(name)}

	if ns := readOptionalString(identityAttrs, platctx.KeyNamespace); ns != nil {
		n := platctx.Namespace(*ns)
		identity.Namespace = &n
	}

	return identity, nil
}

func decodeOwnership(attrs map[string]attr.Value) (platctx.Ownership, error) {
	ownershipAttrs, err := readObject(attrs, platctx.KeyOwnership)
	if err != nil {
		return platctx.Ownership{}, fmt.Errorf("missing required object: %s", platctx.KeyOwnership)
	}

	ownedBy, ok := readRequiredString(ownershipAttrs, platctx.KeyOwnedBy)
	if !ok {
		return platctx.Ownership{}, fmt.Errorf("missing required value: %s.%s", platctx.KeyOwnership, platctx.KeyOwnedBy)
	}

	ownership := platctx.Ownership{OwnedBy: platctx.OwnerRef(ownedBy)}

	if op := readOptionalString(ownershipAttrs, platctx.KeyOperatedBy); op != nil {
		o := platctx.OperatorRef(*op)
		ownership.OperatedBy = &o
	}

	return ownership, nil
}

func decodeEnvironment(attrs map[string]attr.Value) *platctx.Environment {
	if env := readOptionalString(attrs, platctx.KeyEnvironment); env != nil {
		e := platctx.Environment(*env)
		return &e
	}
	return nil
}

func decodeGovernance(attrs map[string]attr.Value) *platctx.Governance {
	governanceAttrs, err := readObject(attrs, platctx.KeyGovernance)
	if err != nil {
		return nil
	}

	criticality := readOptionalString(governanceAttrs, platctx.KeyCriticality)
	costCenter := readOptionalString(governanceAttrs, platctx.KeyCostCenter)
	classification := readOptionalString(governanceAttrs, platctx.KeyDataClassification)

	if criticality == nil && costCenter == nil && classification == nil {
		return nil
	}

	governance := platctx.Governance{}

	if criticality != nil {
		c := platctx.Criticality(*criticality)
		governance.Criticality = &c
	}

	if costCenter != nil {
		cc := platctx.CostCenter(*costCenter)
		governance.CostCenter = &cc
	}

	if classification != nil {
		dc := platctx.DataClassification(*classification)
		governance.DataClassification = &dc
	}

	return &governance
}

func readObject(attrs map[string]attr.Value, key string) (map[string]attr.Value, error) {
	val, ok := attrs[key]
	if !ok {
		return nil, fmt.Errorf("missing: %s", key)
	}

	object, ok := val.(types.Object)
	if !ok {
		return nil, fmt.Errorf("invalid type: %s", key)
	}

	return object.Attributes(), nil
}

// readString extracts a string attribute, returning false if missing, wrong type, or null/unknown.
func readString(attrs map[string]attr.Value, key string) (string, bool) {
	val, ok := attrs[key]
	if !ok {
		return "", false
	}

	str, ok := val.(types.String)
	if !ok {
		return "", false
	}

	if str.IsNull() || str.IsUnknown() {
		return "", false
	}

	return str.ValueString(), true
}

func readRequiredString(attrs map[string]attr.Value, key string) (string, bool) {
	return readString(attrs, key)
}

func readOptionalString(attrs map[string]attr.Value, key string) *string {
	value, ok := readString(attrs, key)
	if !ok {
		return nil
	}
	return &value
}

func encodeContext(value platctx.Context) *types.Dynamic {
	identity := map[string]attr.Value{
		platctx.KeyName:      basetypes.NewStringValue(string(value.Identity.Name)),
		platctx.KeyNamespace: encodeOptionalString(value.Identity.Namespace),
	}

	ownership := map[string]attr.Value{
		platctx.KeyOwnedBy:    basetypes.NewStringValue(string(value.Ownership.OwnedBy)),
		platctx.KeyOperatedBy: encodeOptionalString(value.Ownership.OperatedBy),
	}

	object := types.ObjectValueMust(
		contextResultType.AttrTypes,
		map[string]attr.Value{
			platctx.KeyIdentity:    types.ObjectValueMust(identityTypes, identity),
			platctx.KeyOwnership:   types.ObjectValueMust(ownershipTypes, ownership),
			platctx.KeyEnvironment: encodeOptionalString(value.Environment),
			platctx.KeyGovernance:  encodeGovernance(value.Governance),
		},
	)

	result := basetypes.NewDynamicValue(object)
	return &result
}

func encodeGovernance(governance *platctx.Governance) attr.Value {
	if governance == nil {
		return types.ObjectNull(governanceTypes)
	}

	return types.ObjectValueMust(
		governanceTypes,
		map[string]attr.Value{
			platctx.KeyCriticality:        encodeOptionalString(governance.Criticality),
			platctx.KeyCostCenter:         encodeOptionalString(governance.CostCenter),
			platctx.KeyDataClassification: encodeOptionalString(governance.DataClassification),
		},
	)
}

func encodeOptionalString[T ~string](value *T) attr.Value {
	if value == nil {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(string(*value))
}

func formatIssues(issues []platctx.ValidationIssue) string {
	var b strings.Builder

	for i, issue := range issues {
		if i > 0 {
			b.WriteString("; ")
		}

		b.WriteString(issue.Path)
		b.WriteString(": ")
		b.WriteString(issue.Message)
	}

	return b.String()
}

var supportedAttributes = map[string]struct{}{
	platctx.KeyIdentity:    {},
	platctx.KeyOwnership:   {},
	platctx.KeyEnvironment: {},
	platctx.KeyGovernance:  {},
}

var unknownResult = types.ObjectValueMust(
	contextResultType.AttrTypes,
	map[string]attr.Value{
		platctx.KeyIdentity:    types.ObjectUnknown(identityTypes),
		platctx.KeyOwnership:   types.ObjectUnknown(ownershipTypes),
		platctx.KeyEnvironment: basetypes.NewStringUnknown(),
		platctx.KeyGovernance:  types.ObjectUnknown(governanceTypes),
	},
)

var (
	identityTypes = map[string]attr.Type{
		platctx.KeyName:      types.StringType,
		platctx.KeyNamespace: types.StringType,
	}

	ownershipTypes = map[string]attr.Type{
		platctx.KeyOwnedBy:    types.StringType,
		platctx.KeyOperatedBy: types.StringType,
	}

	governanceTypes = map[string]attr.Type{
		platctx.KeyCriticality:        types.StringType,
		platctx.KeyCostCenter:         types.StringType,
		platctx.KeyDataClassification: types.StringType,
	}

	contextResultType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			platctx.KeyIdentity:    types.ObjectType{AttrTypes: identityTypes},
			platctx.KeyOwnership:   types.ObjectType{AttrTypes: ownershipTypes},
			platctx.KeyEnvironment: types.StringType,
			platctx.KeyGovernance:  types.ObjectType{AttrTypes: governanceTypes},
		},
	}
)
