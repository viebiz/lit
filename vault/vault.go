package vault

import (
	"net/http"

	pkgerrors "github.com/pkg/errors"

	"github.com/hashicorp/vault-client-go"

	"github.com/viebiz/lit/monitoring"
)

type Client struct {
	vaultClient *vault.Client
	info        monitoring.ExternalServiceInfo
}

// NewClient creates a new Vault Client to access the secret vault
func NewClient(address string, httpClient *http.Client) (*Client, error) {
	vc, err := vault.New(
		vault.WithAddress(address),
		vault.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	return &Client{
		vaultClient: vc,
		info:        monitoring.NewExternalServiceInfo(address),
	}, nil
}
