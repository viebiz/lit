package vault

func (client *Client) SetToken(token string) error {
	if err := client.vaultClient.SetToken(token); err != nil {
		return err
	}

	return nil
}
