package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &Provider{}

// New takes a version because Terraform reads it from provider Metadata.
func New(version string) provider.Provider {
	return &Provider{version: version}
}

type Provider struct {
	version string
}

func (p *Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "platctx"
	resp.Version = p.version
}

func (p *Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
}

func (p *Provider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *Provider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewCanonicalizeFunction,
	}
}
