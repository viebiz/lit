package vault

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	"github.com/hashicorp/vault-client-go"
)

const (
	secretMountPath = "secret"
)

func (client *Client) GetSecrets(ctx context.Context, path string) (map[string]interface{}, error) {
	resp, err := client.vaultClient.Secrets.KvV2Read(ctx, path, vault.WithMountPath(secretMountPath))
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	return resp.Data.Data, nil
}
