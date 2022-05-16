package deploying

import (
	"context"

	"github.com/j18e/gofiaas/pkg/models"
	"github.com/j18e/gofiaas/pkg/validating"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	ServiceType corev1.ServiceType
}

type Deployer interface {
	Deploy(context.Context, models.Application) error
}

func NewDeployer(k8s kubernetes.Interface, namespace string, cfg Config) Deployer {
	return &namespacedDeployer{
		namespace: namespace,

		deployments:              newDeploymentDeployer(k8s, namespace),
		horizontalPodAutoscalers: newHorizontalPodAutoscalerDeployer(k8s, namespace),
		ingresses:                newIngressDeployer(k8s, namespace),
		serviceAccounts:          newServiceAccountDeployer(k8s, namespace),
		services:                 newServiceDeployer(k8s, namespace, cfg.ServiceType),
	}
}

type namespacedDeployer struct {
	namespace string
	validating.Validator

	deployments              *deploymentDeployer
	horizontalPodAutoscalers *horizontalPodAutoscalerDeployer
	ingresses                *ingressDeployer
	serviceAccounts          *serviceAccountDeployer
	services                 *serviceDeployer
}

func (d *namespacedDeployer) Deploy(ctx context.Context, app models.Application) error {
	return nil
}

func mergeDicts(base, overrides map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range base {
		res[k] = v
	}
	for k, v := range overrides {
		res[k] = v
	}
	return res
}
