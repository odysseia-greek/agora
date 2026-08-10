package futikon

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/logging"
)

func (t *TheofratosHandler) createILMPolicy(ctx context.Context, policy, phase, rollOver string) error {
	logging.Debug(fmt.Sprintf("creating policy %s with phase %s and rollover: %s", policy, phase, rollOver))

	var policyCreated *models.IndexCreateResult
	var err error

	if rollOver != "" {
		policyCreated, err = t.Elastic.Policy().CreatePolicyWithRollOver(ctx, policy, rollOver, phase)
	} else {
		policyCreated, err = t.Elastic.Policy().CreatePolicy(ctx, policy, phase)
	}

	if err != nil {
		return err
	}

	if policyCreated != nil {
		logging.Info(fmt.Sprintf("created policy: %s %v", policy, policyCreated.Acknowledged))
	}

	return nil
}
