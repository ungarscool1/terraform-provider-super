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
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

	// Fetch API Keys
	var apiKeys []string
	paginator := apigateway.NewGetApiKeysPaginator(client, &apigateway.GetApiKeysInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, key := range output.Items {
			apiKeys = append(apiKeys, aws.ToString(key.Id))
		}
	}

	// Set data source attributes
	if err := d.Set("ids", apiKeys); err != nil {
		return err
	}

	// Set the data source ID to ensure Terraform sees it as stable
	d.SetId("aws_api_gateway_api_keys")

	log.Printf("[INFO] Retrieved %d API Gateway API Keys", len(apiKeys))
	return nil
}

