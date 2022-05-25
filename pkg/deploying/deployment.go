package deploying

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"github.com/j18e/gofiaas/pkg/models"
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

func (d *deploymentDeployer) Deploy(ctx context.Context, spec models.InternalSpec) error {
	return nil
}

func (d *deploymentDeployer) Delete(ctx context.Context, spec models.InternalSpec) error {
	err := d.deployments.Delete(ctx, spec.Name, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *deploymentDeployer) applyPrometheusAnnotations(dep *appsv1.Deployment, pc models.PrometheusConfig) {
	if !pc.Enabled {
		return
	}
	dep.Spec.Template.Annotations["prometheus.io/scrape"] = "true"
	dep.Spec.Template.Annotations["prometheus.io/port"] = pc.Port
	dep.Spec.Template.Annotations["prometheus.io/path"] = pc.Path
}
