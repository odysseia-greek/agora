package diogenes

import (
	"fmt"
	"log"
	"time"
)

func (v *Vault) Health() (bool, error) {
	_, err := v.Connection.Logical().Read("sys/health")
	if err != nil {
		return false, fmt.Errorf("unable to connect to vault: %w", err)
	}

	return true, nil
}

func (v *Vault) CheckHealthyStatus(ticks, tick time.Duration) bool {
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
			healthy, err := v.Health()
			if err != nil {
				log.Printf("Error getting response: %s", err)
				continue
			}
			if healthy {
				return true
			}
		case <-timeout.C:
			return false
		}
	}
}
