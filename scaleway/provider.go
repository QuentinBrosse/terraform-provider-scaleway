package scaleway

import (
	"fmt"
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	scwUtils "github.com/scaleway/scaleway-sdk-go/utils"
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
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					if accessKey, exist := scwConfig.GetAccessKey(); exist {
						return accessKey, nil
					}
					return nil, fmt.Errorf("no access key found")
				}),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Scaleway secret Key.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					if secretKey, exist := scwConfig.GetSecretKey(); exist {
						return secretKey, nil
					}
					return nil, fmt.Errorf("no secret key found")
				}),
			},
			"token": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use `secret_key` instead.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					if secretKey, exist := scwConfig.GetSecretKey(); exist {
						return secretKey, nil
					}
					return nil, fmt.Errorf("no secret key found")
				}),
			},
			"organization": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use `organization_id` instead.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					if organizationID, exist := scwConfig.GetDefaultOrganizationID(); exist {
						return organizationID, nil
					}
					return nil, fmt.Errorf("no organization id found")
				}),
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway organization ID.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					if organizationID, exist := scwConfig.GetDefaultOrganizationID(); exist {
						return organizationID, nil
					}
					return nil, fmt.Errorf("no organization id found")
				}),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway default region to use for your resources.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					region, exist := scwConfig.GetDefaultRegion()
					if exist {
						return region, nil
					}
					return scwUtils.RegionFrPar, nil
				}),
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway default zone to use for your resources.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					zone, exist := scwConfig.GetDefaultZone()
					if exist {
						return zone, nil
					}
					return scwUtils.ZoneFrPar1, nil
				}),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			//"scaleway_bucket":              resourceScalewayBucket(),
			//"scaleway_user_data":           resourceScalewayUserData(),
			//"scaleway_server":              resourceScalewayServer(),
			//"scaleway_token": 			  resourceScalewayToken(),
			//"scaleway_ssh_key":             resourceScalewaySSHKey(),
			//"scaleway_ip":                  resourceScalewayIP(),
			//"scaleway_ip_reverse_dns":      resourceScalewayIPReverseDNS(),
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

func providerConfigure(_ *schema.ResourceData) (interface{}, error) {
	return NewMeta()
}
