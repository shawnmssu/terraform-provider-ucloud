package ucloud

import (
	"fmt"

	"github.com/ucloud/ucloud-sdk-go/external"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"github.com/ucloud/ucloud-sdk-go/ucloud/log"

	pudb "github.com/ucloud/ucloud-sdk-go/private/services/udb"
	"github.com/ucloud/ucloud-sdk-go/services/uaccount"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
	"github.com/ucloud/ucloud-sdk-go/services/udisk"
	"github.com/ucloud/ucloud-sdk-go/services/udpn"
	"github.com/ucloud/ucloud-sdk-go/services/uhost"
	"github.com/ucloud/ucloud-sdk-go/services/ulb"
	"github.com/ucloud/ucloud-sdk-go/services/unet"
	"github.com/ucloud/ucloud-sdk-go/services/vpc"
)

// Config is the configuration of ucloud meta data
type Config struct {
	PublicKey             string
	PrivateKey            string
	Profile               string
	SharedCredentialsFile string

	Region     string
	ProjectId  string
	Insecure   bool
	BaseURL    string
	MaxRetries int
}

// UCloudClient is the ucloud openapi client
type UCloudClient struct {
	region    string
	projectId string

	config     *ucloud.Config
	credential *auth.Credential

	uhostconn    *uhost.UHostClient
	unetconn     *unet.UNetClient
	ulbconn      *ulb.ULBClient
	vpcconn      *vpc.VPCClient
	uaccountconn *uaccount.UAccountClient
	udiskconn    *udisk.UDiskClient
	udpnconn     *udpn.UDPNClient
	udbconn      *udb.UDBClient

	// private services
	pudbconn *pudb.UDBClient
}

// Client will returns a client with connections for all product
func (c *Config) Client() (*UCloudClient, error) {
	var err error
	var client UCloudClient
	var cfg *ucloud.Config
	var cred *auth.Credential

	client.region = c.Region
	client.projectId = c.ProjectId

	cfgV := ucloud.NewConfig()
	cfg = &cfgV

	// set common attributes (region, project id, etc ...)
	cfg.Region = c.Region
	cfg.ProjectId = c.ProjectId

	// enable auto retry with http/connection error
	cfg.MaxRetries = c.MaxRetries
	cfg.LogLevel = log.DebugLevel
	cfg.UserAgent = "Terraform/1.11.x"

	// if no base url be set, get insecure http or secure https default url
	// uf base url is set, use it
	if len(c.BaseURL) == 0 && c.Insecure {
		cfg.BaseUrl = GetInsecureEndpointURL(c.Region)
	} else if len(c.BaseURL) == 0 && !c.Insecure {
		cfg.BaseUrl = GetEndpointURL(c.Region)
	} else {
		cfg.BaseUrl = c.BaseURL
	}

	if len(c.Profile) != 0 {
		// load public/private key from shared credential file
		cred, err = external.LoadUCloudCredentialFile(c.SharedCredentialsFile, c.Profile)
		if err != nil {
			return nil, fmt.Errorf("cannot load shared credential file, %s", err)
		}
	} else {
		// load public/private key from shared credential file
		credV := auth.NewCredential()
		cred = &credV
		cred.PublicKey = c.PublicKey
		cred.PrivateKey = c.PrivateKey
	}

	// initialize client connections
	client.uhostconn = uhost.NewClient(&config, &credential)
	client.unetconn = unet.NewClient(&config, &credential)
	client.ulbconn = ulb.NewClient(&config, &credential)
	client.vpcconn = vpc.NewClient(&config, &credential)
	client.uaccountconn = uaccount.NewClient(&config, &credential)
	client.udiskconn = udisk.NewClient(&config, &credential)
	client.udpnconn = udpn.NewClient(&config, &credential)
	client.udbconn = udb.NewClient(&config, &credential)

	// initialize client connections for private usage
	client.pudbconn = pudb.NewClient(&config, &credential)

	client.config = cfg
	client.credential = cred
	return &client, nil
}
