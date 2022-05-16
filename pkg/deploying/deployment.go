package deploying

import (
	"context"

	"github.com/j18e/gofiaas/pkg/models"
	"k8s.io/client-go/kubernetes"
	clientappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type deploymentDeployer struct {
	namespace   string
	deployments clientappsv1.DeploymentInterface
}

func newDeploymentDeployer(k8s kubernetes.Interface, namespace string) *deploymentDeployer {
	return &deploymentDeployer{
		namespace:   namespace,
		deployments: k8s.AppsV1().Deployments(namespace),
	}
}

func (d *deploymentDeployer) Deploy(ctx context.Context, spec models.ApplicationSpec) error {
	return nil
}
