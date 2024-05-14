package provider

import (
	"context"
	"log"
	"os"

	"terraform-provider-modoboa/internal/client"

	"github.com/deepmap/oapi-codegen/v2/pkg/securityprovider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &modoboaProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &modoboaProvider{
			version: version,
		}
	}
}

type modoboaProvider struct{
	version string
}

func (p *modoboaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "modoboa"
	resp.Version = p.version
}

func (p *modoboaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

type modoboaProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

func (p *modoboaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Modoboa client")

	var config modoboaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Modoboa API Host",
			"The provider cannot create the Modoboa API client as there is an unknown configuration value for the Modoboa API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MODOBOA_HOST environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Modoboa API Token",
			"The provider cannot create the Modoboa API client as there is an unknown configuration value for the Modoboa API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MODOBOA_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("MODOBOA_HOST")
	token := os.Getenv("MODOBOA_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Modoboa API Host",
			"The provider cannot create the Modoboa API client as there is a missing or empty value for the Modoboa API host. "+
				"Set the host value in the configuration or use the MODOBOA_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Modoboa API Username",
			"The provider cannot create the Modoboa API client as there is a missing or empty value for the Modoboa API token. "+
				"Set the username value in the configuration or use the MODOBOA_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "modoboa_host", host)
	ctx = tflog.SetField(ctx, "modoboa_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "modoboa_token")

	tflog.Debug(ctx, "Creating Modoboa client")

	// Create a new Modoboa client using the configuration values
	apiKeyAuth, err := securityprovider.NewSecurityProviderApiKey("header", "Authorization", "Token " + token)
	if err != nil {
		log.Fatal(err)
	}

	c, err := client.NewClientWithResponses(host , client.WithRequestEditorFn(apiKeyAuth.Intercept))
	if err != nil {
    resp.Diagnostics.AddError(
      "Unable to create Modoboa Client",
      err.Error(),
    )
    return
	}

	// Make the Modoboa client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = c
	resp.ResourceData = c

	tflog.Info(ctx, "Configured Modoboa client", map[string]any{"success": true})
}

func (p *modoboaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
  return []func() datasource.DataSource {
		NewDomainsDataSource,  }
}

func (p *modoboaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDomainsResource,
	}
}
