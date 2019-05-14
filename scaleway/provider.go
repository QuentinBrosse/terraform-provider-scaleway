package scaleway

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/mitchellh/go-homedir"

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
				Description: "The Scaleway access key.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					// Keep the deprecated behavior
					if accessKey := os.Getenv("SCALEWAY_ACCESS_KEY"); accessKey != "" {
						return accessKey, nil
					}
					if accessKey, exists := scwConfig.GetAccessKey(); exists {
						return accessKey, nil
					}
					return nil, nil
				}),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true, // To allow user to use deprecated `token`.
				Description: "The Scaleway secret Key.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					// No error is returned here to allow user to use deprecated `token`.
					if secretKey, exists := scwConfig.GetSecretKey(); exists {
						return secretKey, nil
					}
					return nil, nil
				}),
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true, // To allow user to use deprecated `organization`.
				Description: "The Scaleway organization ID.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					if organizationID, exist := scwConfig.GetDefaultOrganizationID(); exist {
						return organizationID, nil
					}
					// No error is returned here to allow user to use deprecated `organization`.
					return nil, nil
				}),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway default region to use for your resources.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					// Keep the deprecated behavior
					// Note: The deprecated region format conversion is handled in `config.GetDeprecatedClient`.
					if region := os.Getenv("SCALEWAY_REGION"); region != "" {
						return region, nil
					}
					if defaultRegion, exists := scwConfig.GetDefaultRegion(); exists {
						return string(defaultRegion), nil
					}
					return string(utils.RegionFrPar), nil
				}),
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway default zone to use for your resources.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					if defaultZone, exists := scwConfig.GetDefaultZone(); exists {
						return string(defaultZone), nil
					}
					return nil, nil
				}),
			},

			// Deprecated values
			"token": {
				Type:       schema.TypeString,
				Optional:   true, // To allow user to use `secret_key`.
				Deprecated: "Use `secret_key` instead.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					for _, k := range []string{"SCALEWAY_TOKEN", "SCALEWAY_ACCESS_KEY"} {
						if os.Getenv(k) != "" {
							return os.Getenv(k), nil
						}
					}
					if path, err := homedir.Expand("~/.scwrc"); err == nil {
						scwAPIKey, _, err := readDeprecatedScalewayConfig(path)
						if err != nil {
							return nil, err
						}
						return scwAPIKey, nil
					}
					// No error is returned here to allow user to use `secret_key`.
					return nil, nil
				}),
			},
			"organization": {
				Type:       schema.TypeString,
				Optional:   true, // To allow user to use `organization_id`.
				Deprecated: "Use `organization_id` instead.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					for _, k := range []string{"SCALEWAY_ORGANIZATION"} {
						if os.Getenv(k) != "" {
							return os.Getenv(k), nil
						}
					}
					if path, err := homedir.Expand("~/.scwrc"); err == nil {
						_, scwOrganization, err := readDeprecatedScalewayConfig(path)
						if err != nil {
							return nil, err
						}
						return scwOrganization, nil
					}
					// No error is returned here to allow user to use `organization_id`.
					return nil, nil
				}),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"scaleway_bucket":              resourceScalewayBucket(),
			"scaleway_user_data":           resourceScalewayUserData(),
			"scaleway_server":              resourceScalewayServer(),
			"scaleway_token":               resourceScalewayToken(),
			"scaleway_ssh_key":             resourceScalewaySSHKey(),
			"scaleway_ip":                  resourceScalewayIP(),
			"scaleway_ip_reverse_dns":      resourceScalewayIPReverseDNS(),
			"scaleway_security_group":      resourceScalewaySecurityGroup(),
			"scaleway_security_group_rule": resourceScalewaySecurityGroupRule(),
			"scaleway_volume":              resourceScalewayVolume(),
			"scaleway_volume_attachment":   resourceScalewayVolumeAttachment(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"scaleway_bootscript":     dataSourceScalewayBootscript(),
			"scaleway_image":          dataSourceScalewayImage(),
			"scaleway_security_group": dataSourceScalewaySecurityGroup(),
			"scaleway_volume":         dataSourceScalewayVolume(),
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
		config.DefaultOrganizationID = d.Get("organization").(string)
	}
	if config.SecretKey == "" && config.AccessKey != "" {
		config.SecretKey = config.AccessKey
	}

	return config.Meta()
}

type deprecatedScalewayConfig struct {
	Organization string `json:"organization"`
	Token        string `json:"token"`
	Version      string `json:"version"`
}

func readDeprecatedScalewayConfig(path string) (string, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	var data deprecatedScalewayConfig
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return "", "", err
	}
	return data.Token, data.Organization, nil
}
