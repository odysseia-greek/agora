package diogenes

import (
	vault "github.com/hashicorp/vault/api"
)

func (v *Vault) SetOnetimeToken(token string) {
	v.Connection.SetToken(token)
}

func (v *Vault) LoginWithRootToken(rootToken string) error {
	v.Connection.SetToken(rootToken)
	return nil
}

func (v *Vault) GetCurrentToken() string {
	return v.Connection.Token()
}

func (v *Vault) CreateOneTimeToken(policy []string) (string, error) {
	renew := false

	tokenRequest := vault.TokenCreateRequest{
		Policies:    policy,
		TTL:         "5m",
		DisplayName: "solonCreated",
		NumUses:     1,
		Renewable:   &renew,
	}

	resp, err := v.Connection.Auth().Token().Create(&tokenRequest)
	if err != nil {
		return "", err
	}

	return resp.Auth.ClientToken, nil
}
