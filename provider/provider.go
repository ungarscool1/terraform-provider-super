package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"aws_role_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The AWS role ARN to assume.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"super_api_gateway_api_keys": dataSourceAPIGatewayAPIKeys(),
		},
	}
}
