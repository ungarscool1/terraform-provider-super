package provider

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAPIGatewayAPIKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAPIGatewayAPIKeysRead,
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_updated_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stage_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAPIGatewayAPIKeysRead(d *schema.ResourceData, meta interface{}) error {
	ctx := context.TODO()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// Create API Gateway client
	client := apigateway.NewFromConfig(cfg)

	var apiKeys []map[string]interface{}
	includeValues := false
	if v, ok := d.GetOk("include_values"); ok {
		includeValues = v.(bool)
	}

	paginator := apigateway.NewGetApiKeysPaginator(client, &apigateway.GetApiKeysInput{
		IncludeValues: &includeValues,
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, item := range output.Items {
			apiKey := map[string]interface{}{
				"id":                *item.Id,
				"name":              *item.Name,
				"enabled":           item.Enabled,
				"created_date":      item.CreatedDate.String(),
				"last_updated_date": item.LastUpdatedDate.String(),
				"stage_keys":        item.StageKeys,
				"tags":              item.Tags,
				"value":             aws.ToString(item.Value),
			}
			apiKeys = append(apiKeys, apiKey)
		}
	}

	if err := d.Set("items", apiKeys); err != nil {
		return err
	}
	// Set the data source ID to ensure Terraform sees it as stable
	d.SetId("super_api_gateway_api_keys")

	log.Printf("[INFO] Retrieved %d API Gateway API Keys", len(apiKeys))
	return nil
}
