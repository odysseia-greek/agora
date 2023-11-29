package diogenes

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"time"
)

// Initialize initializes the Vault server.
func (v *Vault) Initialize(shares, threshold int) (*api.InitResponse, error) {
	initRequest := &api.InitRequest{
		SecretShares:    shares,
		SecretThreshold: threshold,
	}

	return v.Connection.Sys().Init(initRequest)
}

func (v *Vault) InitializeAutoUnseal(shares, threshold int) (*api.InitResponse, error) {
	initRequest := &api.InitRequest{
		SecretShares:      0,
		SecretThreshold:   0,
		StoredShares:      0,
		RecoveryShares:    shares,
		RecoveryThreshold: threshold,
	}
	return v.Connection.Sys().Init(initRequest)
}

func (v *Vault) Leader() (*api.LeaderResponse, error) {
	return v.Connection.Sys().Leader()
}

// Status attempts to unseal the Vault using the provided keys.
func (v *Vault) Status() (*api.SealStatusResponse, error) {
	// Check if the Vault is sealed
	return v.Connection.Sys().SealStatus()
}

// Unseal attempts to unseal the Vault using the provided keys.
func (v *Vault) Unseal(keys []string) (bool, error) {
	// Check if the Vault is sealed
	status, err := v.Status()
	if err != nil {
		return false, err
	}

	if !status.Sealed {
		return true, nil
	}

	// Unseal the Vault
	for _, key := range keys {
		response, err := v.Connection.Sys().Unseal(key)
		if err != nil {
			return false, err
		}

		// Check if the Vault is unsealed after each key
		if !response.Sealed {
			return true, nil
		}

		// it will take some time for the status to change if you don't wait another key will be presented resulting in an error
		time.Sleep(1 * time.Second)

		status, err = v.Status()
		if err != nil {
			return false, err
		}

		if !status.Sealed {
			return true, nil
		}
	}

	// If the loop completes and the Vault is still sealed, return an error
	return false, fmt.Errorf("unable to unseal Vault with provided keys")
}

// AutoUnsealGCP attempts to unseal Vault using a gcp provider key
func (v *Vault) AutoUnsealGCP(keyRing, cryptoKey, location string, keys []string) (bool, error) {
	// Check if the Vault is sealed
	status, err := v.Status()
	if err != nil {
		return false, err
	}

	if !status.Sealed {
		return true, nil
	}

	for _, key := range keys {
		unsealConfig := map[string]interface{}{
			"type":               "gcpckms",
			"gcpckms_key_ring":   keyRing,
			"gcpckms_crypto_key": cryptoKey,
			"gcpckms_location":   location,
			"key":                key,
		}

		_, err = v.Connection.Logical().Write("sys/unseal", unsealConfig)
		if err != nil {
			return false, err
		}

		// it will take some time for the status to change if you don't wait another key will be presented resulting in an error
		time.Sleep(1 * time.Second)

		status, err = v.Status()
		if err != nil {
			return false, err
		}

		if !status.Sealed {
			return true, nil
		}

	}
	return false, fmt.Errorf("unable to unseal Vault with provided keys")
}

func (v *Vault) RaftJoin(leaderAddress string, cert, key, ca []byte) (*api.RaftJoinResponse, error) {
	raftCommand := &api.RaftJoinRequest{
		LeaderAPIAddr:    leaderAddress,
		LeaderClientKey:  string(key),
		LeaderClientCert: string(cert),
		LeaderCACert:     string(ca),
	}

	return v.Connection.Sys().RaftJoin(raftCommand)
}
