package futikon

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (t *TheofratosHandler) WatchConfigMapChanges() error {
	configMapClient := t.Kube.CoreV1().ConfigMaps(t.Namespace)

	// Watch for changes in the ConfigMap
	watch, err := configMapClient.Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to watch config map: %v", err)
	}

	go func() {
		for event := range watch.ResultChan() {
			if event.Type == "MODIFIED" || event.Type == "ADDED" {
				// Get the ConfigMap object
				configMap, ok := event.Object.(*v1.ConfigMap)
				if !ok {
					logging.Error("Received non-ConfigMap event, ignoring")
					continue
				}

				// Check if this is the ConfigMap we care about (e.g., by name)
				if configMap.Name == t.ConfigMapName {
					logging.Debug(fmt.Sprintf("ConfigMap %s changed: %s", configMap.Name, event.Type))

					err := t.handle(configMap.Data)
					if err != nil {
						logging.Error(fmt.Sprintf("Failed to handle config map change: %v", err))
					}
				}
			}
		}
	}()

	select {}
}
