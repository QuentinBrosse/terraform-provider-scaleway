package scaleway

import (
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scaleway/scaleway-sdk-go/utils"
)

var mu = sync.Mutex{}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway access Key.",
				// This is breaking change
				Default: "",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway secret Key.",
				// This is breaking change
				Default: "",
			},
			"token": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use `secret_key` instead.",
				// This is breaking change
				Default: "",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway organization ID.",
				// This is breaking change
				Default: "",
			},
			"organization": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use `organization_id` instead.",
				// This is breaking change
				Default: "",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway default region to use for your resources.",
				// This is breaking change
				Default: "",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway default zone to use for your resources.",
				// This is breaking change
				Default: "",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			//"scaleway_bucket":              resourceScalewayBucket(),
			//"scaleway_user_data":           resourceScalewayUserData(),
			//"scaleway_server":              resourceScalewayServer(),
			"scaleway_token": resourceScalewayToken(),
			//"scaleway_ssh_key":             resourceScalewaySSHKey(),
			"scaleway_ip":             resourceScalewayIP(),
			"scaleway_ip_reverse_dns": resourceScalewayIPReverseDNS(),
			//"scaleway_security_group":      resourceScalewaySecurityGroup(),
			//"scaleway_security_group_rule": resourceScalewaySecurityGroupRule(),
			//"scaleway_volume":              resourceScalewayVolume(),
			//"scaleway_volume_attachment":   resourceScalewayVolumeAttachment(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			//"scaleway_bootscript":     dataSourceScalewayBootscript(),
			//"scaleway_image":          dataSourceScalewayImage(),
			//"scaleway_security_group": dataSourceScalewaySecurityGroup(),
			//"scaleway_volume":         dataSourceScalewayVolume(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		AccessKey:             d.Get("access_key").(string),
		SecretKey:             d.Get("secret_key").(string),
		DefaultOrganizationID: d.Get("organization_id").(string),
		DefaultRegion:         utils.Region(d.Get("region").(string)),
		DefaultZone:           utils.Zone(d.Get("zone").(string)),
	}

	// Handle deprecated values
	if config.SecretKey == "" {
		config.SecretKey = d.Get("token").(string)
	}
	if config.DefaultOrganizationID == "" {
		config.DefaultOrganizationID = d.Get("organization_id").(string)
	}

	return config.Meta()
}
