// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"terraform-provider-gitsync/internal/git/factory"
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
				Required:            true,
				MarkdownDescription: "The URL of your Git repository. Currently only GitHub and GitLab are supported. If the url is from github.com, the GitHub API will be used; otherwise, the GitLab API will be used. Notice that self-hosted GitLab instances with custom domains are also supported and that is the reason why there is no hard rule host to containt gitlab.com in order to use the GitLab API.",
			},
			"token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The personal access token used to authenticate with the Git provider API. The token must have sufficient permissions to create, update, and delete files in the target repository.",
			},
		},
	}
}

func (p *gitSyncProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data gitSyncProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	url := data.URL.ValueString()
	token := data.Token.ValueString()
	if data.Token.IsNull() || data.Token.IsUnknown() || token == "" {
		token = os.Getenv("GITSYNC_TOKEN")
	}
	resp.Diagnostics.AddWarning(fmt.Sprintf("URL: %s", url), fmt.Sprintf("Token: %s", token))

	if url == "" {
		resp.Diagnostics.AddError(getMissingAttributeError("url"))
	}
	if token == "" {
		resp.Diagnostics.AddError(getMissingAttributeError("token"))
	}

	f := factory.NewFactory()
	client, err := f.CreateClient(ctx, url, token)
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
		gsresource.NewValueJsonResource,
		gsresource.NewValueFileResource,
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
