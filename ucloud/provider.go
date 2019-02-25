package ucloud

import (
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_PUBLIC_KEY", nil),
				Description: descriptions["public_key"],
			},

			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_PRIVATE_KEY", nil),
				Description: descriptions["private_key"],
			},

			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_PROFILE", nil),
				Description: descriptions["profile"],
			},

			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_SHARED_CREDENTIAL_FILE", defaultSharedCredentialsFile),
				Description: descriptions["shared_credentials_file"],
			},

			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_REGION", nil),
				Description: descriptions["region"],
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_PROJECT_ID", nil),
				Description: descriptions["project_id"],
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     defaultMaxRetries,
				Description: descriptions["max_retries"],
			},

			"insecure": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       defaultInSecure,
				Description:   descriptions["insecure"],
				ConflictsWith: []string{"base_url"},
			},

			"base_url": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       defaultBaseURL,
				Description:   descriptions["base_url"],
				ConflictsWith: []string{"insecure"},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ucloud_projects":            dataSourceUCloudProjects(),
			"ucloud_images":              dataSourceUCloudImages(),
			"ucloud_zones":               dataSourceUCloudZones(),
			"ucloud_eips":                dataSourceUCloudEips(),
			"ucloud_db_parameter_groups": dataSourceUCloudDBParameterGroups(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ucloud_instance":               resourceUCloudInstance(),
			"ucloud_eip":                    resourceUCloudEIP(),
			"ucloud_eip_association":        resourceUCloudEIPAssociation(),
			"ucloud_vpc":                    resourceUCloudVPC(),
			"ucloud_subnet":                 resourceUCloudSubnet(),
			"ucloud_vpc_peering_connection": resourceUCloudVPCPeeringConnection(),
			"ucloud_udpn_connection":        resourceUCloudUDPNConnection(),
			"ucloud_lb":                     resourceUCloudLB(),
			"ucloud_lb_listener":            resourceUCloudLBListener(),
			"ucloud_lb_attachment":          resourceUCloudLBAttachment(),
			"ucloud_lb_rule":                resourceUCloudLBRule(),
			"ucloud_disk":                   resourceUCloudDisk(),
			"ucloud_disk_attachment":        resourceUCloudDiskAttachment(),
			"ucloud_security_group":         resourceUCloudSecurityGroup(),
			"ucloud_lb_ssl":                 resourceUCloudLBSSL(),
			"ucloud_lb_ssl_attachment":      resourceUCloudLBSSLAttachment(),
			"ucloud_db_instance":            resourceUCloudDBInstance(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		PublicKey:             d.Get("public_key").(string),
		PrivateKey:            d.Get("private_key").(string),
		Region:                d.Get("region").(string),
		MaxRetries:            d.Get("max_retries").(int),
		Insecure:              d.Get("insecure").(bool),
		BaseURL:               d.Get("base_url").(string),
		Profile:               d.Get("profile").(string),
		SharedCredentialsFile: d.Get("shared_credentials_file").(string),
	}

	if projectId, ok := d.GetOk("project_id"); ok && projectId.(string) != "" {
		config.ProjectId = projectId.(string)
	}

	client, err := config.Client()
	return client, err
}

var ucloudMutexKV = mutexkv.NewMutexKV()

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"public_key":              "...",
		"private_key":             "...",
		"region":                  "...",
		"project_id":              "...",
		"max_retries":             "...",
		"insecure":                "...",
		"base_url":                "...",
		"profile":                 "...",
		"shared_credentials_file": "...",
	}
}
