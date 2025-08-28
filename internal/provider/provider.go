// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"terraform-provider-gitsync/internal/git"

	gsresource "terraform-provider-gitsync/internal/resource"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &gitSyncProvider{}

type gitSyncProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &gitSyncProvider{
			version: version,
		}
	}
}

// gitSyncProviderModel describes the provider data model.
type gitSyncProviderModel struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

func (p *gitSyncProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gitsync"
	resp.Version = p.version
}

func (p *gitSyncProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *gitSyncProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	url := os.Getenv("GITSYNC_URL")
	token := os.Getenv("GITSYNC_TOKEN")

	var data gitSyncProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.URL.ValueString() != "" {
		url = data.URL.ValueString()
	}

	if data.Token.ValueString() != "" {
		token = data.Token.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddError(getMissingAttributeError("url"))
	}

	if token == "" {
		resp.Diagnostics.AddError(getMissingAttributeError("token"))
	}

	client, err := git.NewClient(ctx, url, token)
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
		gsresource.NewValueYamlResource,
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

func getMissingAttributeError(attr string) (string, string) {
	return fmt.Sprintf("Missing GitSync API %s", attr),
		fmt.Sprintf(
			"The provider cannot create the GitSync API client as there is a missing or empty value for the GitSync API %s. Set the %s value in the configuration or use the GITSYNC_%s environment variable. If either is already set, ensure the value is not empty.",
			attr,
			attr,
			strings.ToUpper(attr),
		)
}
