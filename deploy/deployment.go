package deploy

import (
	"context"
	"sort"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"github.com/j18e/gofiaas/spec/core"
)

type deploymentDeployer struct {
	deployments clientappsv1.DeploymentInterface
	DeploymentDeployerConfig
}

type DeploymentDeployerConfig struct {
	GlobalEnvVars map[string]string
	PreStopDelay  int
}

func newDeploymentDeployer(deployments clientappsv1.DeploymentInterface, cfg DeploymentDeployerConfig) *deploymentDeployer {
	return &deploymentDeployer{
		deployments:              deployments,
		DeploymentDeployerConfig: cfg,
	}
}

func (d *deploymentDeployer) String() string {
	return "deployment-deployer"
}

func (d *deploymentDeployer) Deploy(ctx context.Context, spec core.Spec) error {
	return nil
}

func (d *deploymentDeployer) Delete(ctx context.Context, spec core.Spec) error {
	err := d.deployments.Delete(ctx, spec.Name, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *deploymentDeployer) applyPrometheusAnnotations(dep *appsv1.Deployment, pc core.PrometheusConfig) {
	if !pc.Enabled {
		return
	}
	dep.Spec.Template.Annotations["prometheus.io/scrape"] = "true"
	dep.Spec.Template.Annotations["prometheus.io/port"] = pc.Port
	dep.Spec.Template.Annotations["prometheus.io/path"] = pc.Path
}

func (d *deploymentDeployer) constructEnvVars(vars map[string]string) []corev1.EnvVar {
	var res []corev1.EnvVar
	for name, val := range vars {
		res = append(res, corev1.EnvVar{
			Name:  name,
			Value: val,
		})
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}
