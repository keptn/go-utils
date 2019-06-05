package utils

import (
	"fmt"
)

// DoHelmUpgrade executes a helm update and upgrade
func DoHelmUpgrade(project string, stage string) error {
	helmChart := fmt.Sprintf("%s/helm-chart", project)
	projectStage := fmt.Sprintf("%s-%s", project, stage)
	_, err := ExecuteCommand("helm", []string{"init", "--client-only"})
	if err != nil {
		return err
	}
	_, err = ExecuteCommand("helm", []string{"dep", "update", helmChart})
	if err != nil {
		return err
	}
	_, err = ExecuteCommand("helm", []string{"upgrade", "--install", projectStage, helmChart, "--namespace", projectStage})
	return err
}

// CheckDeploymentRolloutStatus checks the rollout status of the provided deployment
func CheckDeploymentRolloutStatus(serviceName string, namespace string) error {
	_, err := ExecuteCommand("kubectl", []string{"rollout", "status", "deployment/" + serviceName, "--namespace", namespace})
	return err
}
