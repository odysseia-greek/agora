package diogenes

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (v *Vault) Health(ctx context.Context) (bool, error) {
	_, err := v.Connection.Sys().HealthWithContext(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to connect to vault: %w", err)
	}

	return true, nil
}

func (v *Vault) CheckHealthyStatus(ctx context.Context, ticks, tick time.Duration) bool {
	if ticks <= 0 || tick <= 0 {
		return false
	}

	ticker := time.NewTicker(tick)
	defer ticker.Stop()
	timeout := time.NewTimer(ticks)
	defer timeout.Stop()

	for {
		select {
		case t := <-ticker.C:
			log.Printf("tick: %s", t)
			healthy, err := v.Health(ctx)
			if err != nil {
				log.Printf("Error getting response: %s", err)
				continue
			}
			if healthy {
				return true
			}
		case <-timeout.C:
			return false
		case <-ctx.Done():
			return false
		}
	}
}
