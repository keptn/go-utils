package utils

import (
	"fmt"
)

func DoHelmUpgrade(project string, stage string) error {
	helmChart := fmt.Sprintf("%s/helm-chart", project)
	projectStage := fmt.Sprintf("%s-%s", project, stage)
	err := ExecuteCommand("helm", []string{"init", "--client-only"})
	if err != nil {
		return err
	}
	err = ExecuteCommand("helm", []string{"dep", "update", helmChart})
	if err != nil {
		return err
	}
	return ExecuteCommand("helm", []string{"upgrade", "--install", projectStage, helmChart, "--namespace", projectStage})
}

func CheckDeploymentRolloutStatus(serviceName string, namespace string) error {

	return ExecuteCommand("kubectl", []string{"rollout", "status", "deployment/" + serviceName, "--namespace", namespace})
}
