package deploying

import (
	"context"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	"github.com/j18e/gofiaas/pkg/models"
)

type Config struct {
	ServiceType corev1.ServiceType
}

type Deployer interface {
	Deploy(context.Context, models.InternalSpec) error
	Delete(context.Context, models.InternalSpec) error
}

func NewDeployer(k8s kubernetes.Interface, factory informers.SharedInformerFactory, namespace string, cfg Config) Deployer {
	return &namespacedDeployer{
		namespace: namespace,

		deployments:              newDeploymentDeployer(k8s, namespace),
		horizontalPodAutoscalers: newHorizontalPodAutoscalerDeployer(k8s, namespace),
		ingresses:                newIngressDeployer(k8s, namespace),
		serviceAccounts: newServiceAccountDeployer(
			k8s.CoreV1().ServiceAccounts(namespace),
			factory.Core().V1().ServiceAccounts().Lister().ServiceAccounts(namespace),
		),
		services: newServiceDeployer(k8s, namespace, cfg.ServiceType),
	}
}

type namespacedDeployer struct {
	namespace string

	deployments              *deploymentDeployer
	horizontalPodAutoscalers *horizontalPodAutoscalerDeployer
	ingresses                *ingressDeployer
	serviceAccounts          *serviceAccountDeployer
	services                 *serviceDeployer
}

func (d *namespacedDeployer) Deploy(ctx context.Context, spec models.InternalSpec) error {
	return nil
}

func (d *namespacedDeployer) Delete(ctx context.Context, spec models.InternalSpec) error {
	return nil
}

func (d *namespacedDeployer) makeLabels(spec models.InternalSpec) map[string]string {
	return map[string]string{
		"app":                 spec.Name,
		"fiaas/version":       strconv.Itoa(spec.Version),
		"fiaas/deployment_id": spec.DeploymentID,
	}
}
