package scaleway

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform/helper/logging"
	api "github.com/nicolai86/scaleway-sdk"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/scwconfig"
	"github.com/scaleway/scaleway-sdk-go/utils"
)

var scwConfig scwconfig.Config

func init() {
	config, err := scwconfig.Load()
	if err != nil {
		log.Fatalf("error: cannot load configuration: %s", err)
	}
	scwConfig = config
}

// Config is a configuration for a client.
type Config struct {
	AccessKey             string
	SecretKey             string
	DefaultOrganizationID string
	DefaultRegion         utils.Region
	DefaultZone           utils.Zone
}

// Meta contains the SDK client and a deprecated version of the SDK temporary.
//
// This meta value returned by this function is passed into all resources.
type Meta struct {
	client           *scw.Client
	deprecatedClient *api.API
}

// Meta creates a meta instance from a client configuration.
func (c *Config) Meta() (*Meta, error) {
	meta := &Meta{}

	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}
	meta.client = client

	deprecatedClient, err := c.GetDeprecatedClient()
	if err != nil {
		return nil, fmt.Errorf("error: cannot create deprecated client: %s", err)
	}
	meta.deprecatedClient = deprecatedClient

	err = fetchServerAvailabilities(client)
	if err != nil {
		log.Printf("error: cannot fetch server availabilities: %s", err)
	}

	return meta, nil
}

func (c *Config) GetClient() (*scw.Client, error) {
	cl := createRetryableHTTPClient()

	options := []scw.ClientOption{
		scw.WithHTTPClient(cl.HTTPClient),
	}

	// The access key is not used for API authentications.
	if c.SecretKey != "" {
		options = append(options, scw.WithAuth(c.AccessKey, c.SecretKey))
	}

	if c.DefaultOrganizationID != "" {
		options = append(options, scw.WithDefaultOrganizationID(c.DefaultOrganizationID))
	}

	if c.DefaultRegion != "" {
		options = append(options, scw.WithDefaultRegion(c.DefaultRegion))
	}

	if c.DefaultZone != "" {
		options = append(options, scw.WithDefaultZone(c.DefaultZone))
	}

	// TODO: Use retryablehttp client here
	client, err := scw.NewClient(options...)
	if err != nil {
		return nil, fmt.Errorf("error: cannot create SDK client: %s", err)
	}

	return client, err
}

// createRetryableHTTPClient create a retryablehttp.Client.
func createRetryableHTTPClient() *retryablehttp.Client {
	cl := retryablehttp.NewClient()

	cl.HTTPClient.Transport = logging.NewTransport("Scaleway", cl.HTTPClient.Transport)
	cl.RetryMax = 3
	cl.RetryWaitMax = 2 * time.Minute
	cl.Logger = log.New(os.Stderr, "", 0)
	cl.RetryWaitMin = time.Minute
	cl.CheckRetry = func(_ context.Context, resp *http.Response, err error) (bool, error) {
		if resp == nil {
			return true, err
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			return true, err
		}
		return retryablehttp.DefaultRetryPolicy(context.TODO(), resp, err)
	}

	return cl
}

//
// DEPRECATED ZONE
//

// client is a bridge between sdk.HTTPClient interface and retryablehttp.Client
type client struct {
	*retryablehttp.Client
}

func (c *client) Do(r *http.Request) (*http.Response, error) {
	var body io.ReadSeeker
	if r.Body != nil {
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bs)
	}
	req, err := retryablehttp.NewRequest(r.Method, r.URL.String(), body)
	for key, val := range r.Header {
		req.Header.Set(key, val[0])
	}
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Config) GetDeprecatedClient() (*api.API, error) {
	options := func(sdkApi *api.API) {
		sdkApi.Client = &client{retryablehttp.NewClient()}
	}

	// TODO: Replace by a parsing with error handling.
	region := ""
	if c.DefaultRegion == utils.RegionFrPar || c.DefaultZone == utils.ZoneFrPar1 {
		region = "par1"
	}
	if c.DefaultRegion == utils.RegionNlAms || c.DefaultZone == utils.ZoneNlAms1 {
		region = "ams1"
	}

	return api.New(
		c.DefaultOrganizationID,
		c.SecretKey,
		region,
		options,
	)
}
