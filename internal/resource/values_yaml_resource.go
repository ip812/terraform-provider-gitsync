// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"strings"

	"terraform-provider-gitsync/internal/git"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	defaultBranch = "main"
)

var _ resource.Resource = &ValuesYamlResource{}
var _ resource.ResourceWithImportState = &ValuesYamlResource{}

func NewValueYamlResource() resource.Resource {
	return &ValuesYamlResource{}
}

type ValuesYamlResource struct {
	client git.Client
}

type ValuesYamlResourceModel struct {
	Path    types.String `tfsdk:"path"`
	Branch  types.String `tfsdk:"branch"`
	Content types.String `tfsdk:"content"`
	ID      types.String `tfsdk:"id"`
}

func (r *ValuesYamlResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_values_yaml"
}

func (r *ValuesYamlResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create or update a file in the configured Git repository.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				MarkdownDescription: "Relative path of the file to create/update in the repo.",
				Required:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Branch to commit to. Defaults to the main branch.",
				Optional:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "File content to write.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Unique ID (path + commit hash).",
			},
		},
	}
}

func (r *ValuesYamlResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(git.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data",
			fmt.Sprintf("Expected *git.Client, got %T", req.ProviderData),
		)
		return
	}
	r.client = c
}

func (r *ValuesYamlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ValuesYamlResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Create(ctx, git.ValuesYamlModel{
		Path:    data.Path.ValueString(),
		Branch:  data.Branch.ValueString(),
		Content: data.Content.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create file in GitHub",
			fmt.Sprintf(
				"An error occurred while updating %q in branch %q: %v",
				data.Path.ValueString(),
				data.Branch.ValueString(),
				err,
			),
		)
		return
	}

	data.ID = types.StringValue(fmt.Sprintf(
		"github-%s-%s-%s",
		r.client.Owner(),
		r.client.Repository(),
		strings.ReplaceAll(data.Branch.ValueString(), "/", "-"),
	))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ValuesYamlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *ValuesYamlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ValuesYamlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *ValuesYamlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
