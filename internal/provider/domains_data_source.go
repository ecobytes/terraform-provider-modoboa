package provider

import (
	"context"
	"log"
	"net/http"

	"terraform-provider-modoboa/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

  "github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
)

// domainsDataSource is the data source implementation.
type domainsDataSource struct{
  client *client.ClientWithResponses
}

// domainsDataSourceModel maps the data source schema data.
type domainsDataSourceModel struct {
  Domains []client.Domain `tfsdk:"domains"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
  _ datasource.DataSource               = &domainsDataSource{}
  _ datasource.DataSourceWithConfigure = &domainsDataSource{}
)

// NewDomainsDataSource is a helper function to simplify the provider implementation.
func NewDomainsDataSource() datasource.DataSource {
  return &domainsDataSource{}
}

// Metadata returns the data source type name.
func (d *domainsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_domains"
}

// Schema defines the schema for the data source.
func (d *domainsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
  resp.Schema = schema.Schema{
    Attributes: map[string]schema.Attribute{
      "domains": schema.ListNestedAttribute{
        Computed: true,
        NestedObject: schema.NestedAttributeObject{
          Attributes: map[string]schema.Attribute{
            "creation": schema.StringAttribute{
              MarkdownDescription: "Date and time of creation",
              CustomType: timetypes.RFC3339Type{},
              Computed: true,
              Optional: true,
            },
            "default_mailbox_quota": schema.StringAttribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "dkim_key_length": schema.Int64Attribute{
              MarkdownDescription: "",
              Computed: true,
            },
            "dkim_key_selector": schema.StringAttribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "dkim_private_key_path": schema.StringAttribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "dkim_public_key": schema.StringAttribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "dns_global_status": schema.StringAttribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "domain_admin": schema.StringAttribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "domainalias_count": schema.Int64Attribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "enable_dkim": schema.BoolAttribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "enabled": schema.BoolAttribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "last_modification": schema.StringAttribute{
              MarkdownDescription: "Date and time of last modification",
              CustomType: timetypes.RFC3339Type{},
              Computed: true,
              Optional: true,
            },
            "mailbox_count": schema.Int64Attribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "mbalias_count": schema.Int64Attribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "message_limit": schema.Int64Attribute{
              MarkdownDescription: "",
              Computed: true,
            },
            "name": schema.StringAttribute{
              MarkdownDescription: "",
              Computed: true,
            },
            "opened_alarms_count": schema.Int64Attribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "pk": schema.Int64Attribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "quota": schema.StringAttribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
            "transport": schema.StringAttribute{
              MarkdownDescription: "(do not use; serializer object for Transport model)",
              Computed: true,
              Optional: true,
            },
            "type": schema.StringAttribute{
              MarkdownDescription: "",
              Computed: true,
              Optional: true,
            },
          },
        },
      },
    },
  }
}

// Configure adds the provider configured client to the data source.
func (d *domainsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
  if req.ProviderData == nil {
    return
  }

  d.client = req.ProviderData.(*client.ClientWithResponses)
}

// Read refreshes the Terraform state with the latest data.
func (d *domainsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
  var state domainsDataSourceModel

  res, err := d.client.DomainsListWithResponse(ctx)
  if err != nil {
    resp.Diagnostics.AddError(
      "Unable to receive domains list",
      err.Error(),
    )
    return
  }

  if res.StatusCode() != http.StatusOK {
    log.Fatalf("Expected HTTP 200 but received %d", res.StatusCode())
  }

  // Map response body to model
  for _, domain := range *res.JSON200 {
    domainState := client.Domain{
      Creation:            domain.Creation,
      DefaultMailboxQuota: domain.DefaultMailboxQuota,
      DkimKeyLength:       domain.DkimKeyLength,
      DkimKeySelector:     domain.DkimKeySelector,
      DkimPrivateKeyPath:  domain.DkimPrivateKeyPath,
      DkimPublicKey:       domain.DkimPublicKey,
      DnsGlobalStatus:     domain.DnsGlobalStatus,
      DomainAdmin:         domain.DomainAdmin,
      DomainaliasCount:    domain.DomainaliasCount,
      EnableDkim:          domain.EnableDkim,
      Enabled:             domain.Enabled,
      LastModification:    domain.LastModification,
      MailboxCount:        domain.MailboxCount,
      MbaliasCount:        domain.MbaliasCount,
      MessageLimit:        domain.MessageLimit,
      Name:                domain.Name,
      OpenedAlarmsCount:   domain.OpenedAlarmsCount,
      Pk:                  domain.Pk,
      Quota:               domain.Quota,
      Type:                domain.Type,
    }
    state.Domains = append(state.Domains, domainState)
  }

  // Set state
  diags := resp.State.Set(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }
}
