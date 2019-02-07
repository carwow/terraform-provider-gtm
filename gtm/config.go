package gtm

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/carwow/terraform-provider-gtm/version"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/httpclient"

	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
	"google.golang.org/api/tagmanager/v2"
)

// Config is the configuration structure used to instantiate the Google
// Tag Manager API.
type Config struct {
	Credentials string
	AccessToken string
	Scopes      []string

	client    *http.Client
	userAgent string

	tokenSource oauth2.TokenSource

	clientTagManager *tagmanager.Service
}

var defaultClientScopes = []string{
	"https://www.googleapis.com/auth/tagmanager.publish",
}

func (c *Config) loadAndValidate() error {
	tokenSource, err := c.getTokenSource(c.Scopes)
	if err != nil {
		return err
	}
	c.tokenSource = tokenSource

	client := oauth2.NewClient(context.Background(), tokenSource)
	client.Transport = logging.NewTransport("GTM", client.Transport)

	terraformVersion := httpclient.UserAgentString()
	providerVersion := fmt.Sprintf("terraform-provider-gtm/%s", version.ProviderVersion)
	terraformWebsite := "(+https://www.terraform.io)"
	userAgent := fmt.Sprintf("%s %s %s", terraformVersion, terraformWebsite, providerVersion)

	c.client = client
	c.userAgent = userAgent

	log.Printf("[INFO] Instantiating GTM client...")
	c.clientTagManager, err = tagmanager.New(client)
	if err != nil {
		return err
	}
	c.clientTagManager.UserAgent = userAgent

	return nil
}

func (c *Config) getTokenSource(clientScopes []string) (oauth2.TokenSource, error) {
	if c.AccessToken != "" {
		contents, _, err := pathorcontents.Read(c.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("Error loading access token: %s", err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'access_token'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		token := &oauth2.Token{AccessToken: contents}
		return oauth2.StaticTokenSource(token), nil
	}

	if c.Credentials != "" {
		contents, _, err := pathorcontents.Read(c.Credentials)
		if err != nil {
			return nil, fmt.Errorf("Error loading credentials: %s", err)
		}

		creds, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(contents), clientScopes...)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse credentials from '%s': %s", contents, err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'credentials'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return creds.TokenSource, nil
	}

	log.Printf("[INFO] Authenticating using DefaultClient...")
	log.Printf("[INFO]   -- Scopes: %s", clientScopes)
	return googleoauth.DefaultTokenSource(context.Background(), clientScopes...)
}
