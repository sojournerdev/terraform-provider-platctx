package main

import (
	"context"
	"log"
	"runtime/debug"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	intprovider "github.com/sojournerdev/terraform-provider-platctx/internal/provider"
)

// version is the provider version reported to Terraform.
//
// Populated by -X linker flags in release builds. Falls back to the Go module
// version from debug.ReadBuildInfo() for go install or local development.
var version string

func init() {
	if version != "" {
		return
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		version = info.Main.Version
	}
}

func main() {
	if err := providerserver.Serve(context.Background(), func() provider.Provider {
		return intprovider.New(version)
	}, providerserver.ServeOpts{
		Address: "registry.terraform.io/sojournerdev/platctx",
	}); err != nil {
		log.Fatal(err)
	}
}
