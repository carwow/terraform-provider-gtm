package gtm

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"golang.org/x/oauth2/google"
)

const testFakeCredentialsPath = "./test-fixtures/fake_account.json"
const testOauthScope = "https://www.googleapis.com/auth/tagmanager.publish"

func TestAccConfigLoadValidate_credentials(t *testing.T) {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Network access not allowed; use %s=1 to enable", resource.TestEnvVar))
	}
	testAccPreCheck(t)

	creds := getTestCredsFromEnv()

	config := Config{
		Credentials: creds,
	}

	err := config.loadAndValidate()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestAccConfigLoadValidate_accessToken(t *testing.T) {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Network access not allowed; use %s=1 to enable", resource.TestEnvVar))
	}
	testAccPreCheck(t)

	creds := getTestCredsFromEnv()
	fmt.Println("!!!!!")
	fmt.Println(creds)

	c, err := google.CredentialsFromJSON(context.Background(), []byte(creds), testOauthScope)
	if err != nil {
		t.Fatalf("invalid test credentials: %s", err)
	}

	token, err := c.TokenSource.Token()
	if err != nil {
		t.Fatalf("Unable to generate test access token: %s", err)
	}

	config := Config{
		AccessToken: token.AccessToken,
	}

	err = config.loadAndValidate()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
