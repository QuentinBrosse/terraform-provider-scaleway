package scaleway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/scwconfig"
)

var scwConfig scwconfig.Config

func init() {
	config, err := scwconfig.Load()
	if err != nil {
		log.Fatalf("error: cannot load configuration: %s", err)
	}
	scwConfig = config
}

type Meta struct {
	client *scw.Client
}

func NewMeta() (*Meta, error) {
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

	// TODO: Use retryablehttp client here
	client, err := scw.NewClient(scw.WithConfig(scwConfig), scw.WithHTTPClient(cl.HTTPClient))
	if err != nil {
		return nil, fmt.Errorf("error: cannot create SDK client: %s", err)
	}

	config := &Meta{
		client: client,
	}

	// TODO: Uncomment me when server resource will be implemented.
	/*
		err = fetchServerAvailabilities(client)
		if err != nil {
			log.Printf("error: cannot fetch server availabilities: %s", err)
		}
	*/

	return config, nil
}
