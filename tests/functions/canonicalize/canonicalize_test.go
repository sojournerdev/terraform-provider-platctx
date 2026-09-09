package canonicalize_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/sojournerdev/terraform-provider-platctx/internal/provider"
)

func newProviderFactory() func() (tfprotov6.ProviderServer, error) {
	return func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6(provider.New("test"))(), nil
	}
}

func TestCanonicalizeValid(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"platctx": newProviderFactory(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					locals {
						context = provider::platctx::canonicalize({
							identity = {
								name      = "payments"
								namespace = "platform"
							}
							ownership = {
								owned_by    = "team:platform"
								operated_by = "team:sre"
							}
							environment = "production"
						})
					}

					output "context" {
						value = local.context
					}
				`,
			},
		},
	})
}

func TestCanonicalizeOptionalFieldsOnly(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"platctx": newProviderFactory(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					locals {
						context = provider::platctx::canonicalize({
							identity = {
								name = "payments"
							}
							ownership = {
								owned_by = "team:platform"
							}
						})
					}

					output "context" {
						value = local.context
					}
				`,
			},
		},
	})
}

func TestCanonicalizeGovernanceNormalization(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"platctx": newProviderFactory(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					locals {
						context = provider::platctx::canonicalize({
							identity = {
								name = "payments"
							}
							ownership = {
								owned_by = "team:platform"
							}
							governance = {}
						})
					}

					output "context" {
						value = local.context
					}
				`,
			},
		},
	})
}

func TestCanonicalizeUnknownPropagation(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"platctx": newProviderFactory(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					locals {
						context = provider::platctx::canonicalize({
							identity = {
								name = "payments"
							}
							ownership = {
								owned_by = "team:platform"
							}
						})
					}

					output "context" {
						value = local.context
					}
				`,
			},
		},
	})
}

func TestCanonicalizeInvalidEmptyName(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"platctx": newProviderFactory(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					locals {
						context = provider::platctx::canonicalize({
							identity = {
								name = ""
							}
							ownership = {
								owned_by = "team:platform"
							}
						})
					}

					output "context" {
						value = local.context
					}
				`,
				ExpectError: regexp.MustCompile(`identity.name.*must be non-empty`),
			},
		},
	})
}

func TestCanonicalizeInvalidWhitespaceName(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"platctx": newProviderFactory(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					locals {
						context = provider::platctx::canonicalize({
							identity = {
								name = " payments"
							}
							ownership = {
								owned_by = "team:platform"
							}
						})
					}

					output "context" {
						value = local.context
					}
				`,
				ExpectError: regexp.MustCompile(`identity.name.*must be free of leading`),
			},
		},
	})
}

func TestCanonicalizeInvalidCriticality(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"platctx": newProviderFactory(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					locals {
						context = provider::platctx::canonicalize({
							identity = {
								name = "payments"
							}
							ownership = {
								owned_by = "team:platform"
							}
							governance = {
								criticality = "severe"
							}
						})
					}

					output "context" {
						value = local.context
					}
				`,
				ExpectError: regexp.MustCompile(`governance.criticality.*must be one`),
			},
		},
	})
}

func TestCanonicalizeMultipleErrors(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"platctx": newProviderFactory(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					locals {
						context = provider::platctx::canonicalize({
							identity = {
								name = ""
							}
							ownership = {
								owned_by = ""
							}
						})
					}

					output "context" {
						value = local.context
					}
				`,
				ExpectError: regexp.MustCompile(`identity.name.*must be non-empty`),
			},
		},
	})
}

func TestCanonicalizeUnsupportedAttribute(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"platctx": newProviderFactory(),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					locals {
						context = provider::platctx::canonicalize({
							identity = {
								name = "payments"
							}
							ownership = {
								owned_by = "team:platform"
							}
							unknown_field = "test"
						})
					}

					output "context" {
						value = local.context
					}
				`,
				ExpectError: regexp.MustCompile(`unsupported attribute`),
			},
		},
	})
}
