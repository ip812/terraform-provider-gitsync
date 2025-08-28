// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"context"
	"fmt"
	"os"
	"strings"
	"terraform-provider-gitsync/internal/sync"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	defaultGitBranch = "main"
)

var _ provider.Provider = &gitSyncProvider{}

type gitSyncProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// gitSyncProviderModel describes the provider data model.
type gitSyncProviderModel struct {
	Repository types.String `tfsdk:"repository"`
	Branch     types.String `tfsdk:"branch"`
	Token      types.String `tfsdk:"token"`
}

func (p *gitSyncProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gitsync"
	resp.Version = p.version
}

func (p *gitSyncProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"repository": schema.StringAttribute{
				Required: true,
			},
			"branch": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *gitSyncProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg gitSyncProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if cfg.Repository.IsUnknown() {
		resp.Diagnostics.AddAttributeError(getUnknownAttributeError("repository"))
	}

	if cfg.Branch.IsUnknown() {
		resp.Diagnostics.AddAttributeError(getUnknownAttributeError("branch"))
	}

	if cfg.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(getUnknownAttributeError("token"))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	repository := os.Getenv("GITSYNC_REPOSITORY")
	branch := os.Getenv("GITSYNC_BRANCH")
	token := os.Getenv("GITSYNC_TOKEN")

	if !cfg.Repository.IsNull() {
		repository = cfg.Repository.ValueString()
	}

	if !cfg.Branch.IsNull() {
		branch = cfg.Branch.ValueString()
	}

	if !cfg.Token.IsNull() {
		token = cfg.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if repository == "" {
		resp.Diagnostics.AddAttributeError(getMissingAttributeError("repository"))
	}

	if branch == "" {
		cfg.Branch = types.StringValue(defaultGitBranch)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(getMissingAttributeError("token"))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := sync.NewClient(&repository, &branch, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create GitSync API Client",
			fmt.Sprintf("An unexpected error was encountered trying to create the GitSync API client: %s", err.Error()),
		)
		return
	}
	resp.ResourceData = client
}

func (p *gitSyncProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewValueYamlResource,
	}
}

func (p *gitSyncProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *gitSyncProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *gitSyncProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &gitSyncProvider{
			version: version,
		}
	}
}

func getUnknownAttributeError(attr string) (path.Path, string, string) {
	return path.Root(attr),
		fmt.Sprintf("Unknown GitSync API %s", attr),
		fmt.Sprintf(
			"The provider cannot create the GitSync API client as there is an unknown configuration value for the GitSync API %s. Either target apply the source of the value first, set the value statically in the configuration, or use the GITSYNC_%s environment variable.",
			attr,
			strings.ToUpper(attr),
		)
}

func getMissingAttributeError(attr string) (path.Path, string, string) {
	return path.Root(attr),
		fmt.Sprintf("Missing GitSync API %s", attr),
		fmt.Sprintf(
			"The provider cannot create the GitSync API client as there is a missing or empty value for the GitSync API %s. Set the %s value in the configuration or use the GITSYNC_%s environment variable. If either is already set, ensure the value is not empty.",
			attr,
			attr,
			strings.ToUpper(attr),
		)
}
