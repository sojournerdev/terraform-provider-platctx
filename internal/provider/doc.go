// Package provider bridges Terraform types to the platctx domain model.
//
// Handles null/unknown values, validates at the boundary, and encodes results.
// Only package that imports Terraform dependencies.
package provider
