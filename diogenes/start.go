package diogenes

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	"time"
)

// Initialize initializes the Vault server.
func (v *Vault) Initialize(ctx context.Context, shares, threshold int) (*api.InitResponse, error) {
	initRequest := &api.InitRequest{
		SecretShares:    shares,
		SecretThreshold: threshold,
	}

	return v.Connection.Sys().InitWithContext(ctx, initRequest)
}

func (v *Vault) InitializeAutoUnseal(ctx context.Context, shares, threshold int) (*api.InitResponse, error) {
	initRequest := &api.InitRequest{
		SecretShares:      0,
		SecretThreshold:   0,
		StoredShares:      0,
		RecoveryShares:    shares,
		RecoveryThreshold: threshold,
	}
	return v.Connection.Sys().InitWithContext(ctx, initRequest)
}

func (v *Vault) Leader(ctx context.Context) (*api.LeaderResponse, error) {
	return v.Connection.Sys().LeaderWithContext(ctx)
}

// Status attempts to unseal the Vault using the provided keys.
func (v *Vault) Status(ctx context.Context) (*api.SealStatusResponse, error) {
	// Check if the Vault is sealed
	return v.Connection.Sys().SealStatusWithContext(ctx)
}

// Unseal attempts to unseal the Vault using the provided keys.
func (v *Vault) Unseal(ctx context.Context, keys []string) (bool, error) {
	// Check if the Vault is sealed
	status, err := v.Status(ctx)
	if err != nil {
		return false, err
	}

	if !status.Sealed {
		return true, nil
	}

	// Unseal the Vault
	for _, key := range keys {
		response, err := v.Connection.Sys().UnsealWithContext(ctx, key)
		if err != nil {
			return false, err
		}

		// Check if the Vault is unsealed after each key
		if !response.Sealed {
			return true, nil
		}

		// it will take some time for the status to change if you don't wait another key will be presented resulting in an error
		if err := waitForContext(ctx, time.Second); err != nil {
			return false, err
		}

		status, err = v.Status(ctx)
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
func (v *Vault) AutoUnsealGCP(ctx context.Context, keyRing, cryptoKey, location string, keys []string) (bool, error) {
	// Check if the Vault is sealed
	status, err := v.Status(ctx)
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

		_, err = v.Connection.Logical().WriteWithContext(ctx, "sys/unseal", unsealConfig)
		if err != nil {
			return false, err
		}

		// it will take some time for the status to change if you don't wait another key will be presented resulting in an error
		if err := waitForContext(ctx, time.Second); err != nil {
			return false, err
		}

		status, err = v.Status(ctx)
		if err != nil {
			return false, err
		}

		if !status.Sealed {
			return true, nil
		}

	}
	return false, fmt.Errorf("unable to unseal Vault with provided keys")
}

func (v *Vault) RaftJoin(ctx context.Context, leaderAddress string, cert, key, ca []byte) (*api.RaftJoinResponse, error) {
	raftCommand := &api.RaftJoinRequest{
		LeaderAPIAddr:    leaderAddress,
		LeaderClientKey:  string(key),
		LeaderClientCert: string(cert),
		LeaderCACert:     string(ca),
	}

	return v.Connection.Sys().RaftJoinWithContext(ctx, raftCommand)
}

func waitForContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
