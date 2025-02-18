package vault

import (
	"github.com/hashicorp/vault-client-go"

	"github.com/viebiz/lit/monitoring"
)

type Client struct {
	vaultClient *vault.Client
	info        monitoring.ExternalServiceInfo
}
