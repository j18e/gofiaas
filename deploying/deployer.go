package deploying

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	"github.com/j18e/gofiaas/log"
	"github.com/j18e/gofiaas/models"
)

const DefaultSecretsDir = "/var/run/secrets/fiaas"

var ReEnvVar = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]+`)

type Config struct {
	ServiceType      corev1.ServiceType
	SecretsDir       string
	PreStopDelay     int
	IngressSuffixes  []string
	HostRewriteRules map[*regexp.Regexp]string
	GlobalEnvVars    map[string]string
}

type ResourceDeployer interface {
	fmt.Stringer
	Deploy(context.Context, models.InternalSpec) error
	Delete(context.Context, models.InternalSpec) error
}

func NewDeployer(k8s kubernetes.Interface, factory informers.SharedInformerFactory, namespace string, cfg Config) *Deployer {
	return &Deployer{
		deployers: []ResourceDeployer{
			newDeploymentDeployer(k8s.AppsV1().Deployments(namespace), cfg.GlobalEnvVars),
			newAutoscalerDeployer(k8s.AutoscalingV1().HorizontalPodAutoscalers(namespace)),
			newIngressDeployer(
				k8s.NetworkingV1().Ingresses(namespace),
				ingressDeployerConfig{
					suffixes:         cfg.IngressSuffixes,
					hostRewriteRules: cfg.HostRewriteRules,
				},
			),
			newServiceAccountDeployer(
				k8s.CoreV1().ServiceAccounts(namespace),
				factory.Core().V1().ServiceAccounts().Lister().ServiceAccounts(namespace),
			),
			newServiceDeployer(k8s.CoreV1().Services(namespace), cfg.ServiceType),
		},
	}
}

type Deployer struct {
	deployers []ResourceDeployer
}

func (d *Deployer) Deploy(ctx context.Context, spec models.InternalSpec) {
	for _, dep := range d.deployers {
		if err := dep.Deploy(ctx, spec); err != nil {
			log.Logger.Warnf("Error calling deploy on %s: %v", dep, err)
		}
	}
}

func (d *Deployer) Delete(ctx context.Context, spec models.InternalSpec) {
	for _, dep := range d.deployers {
		if err := dep.Delete(ctx, spec); err != nil {
			log.Logger.Warnf("Error calling deploy on %s: %v", dep, err)
		}
	}
}

func (d *Deployer) makeLabels(spec models.InternalSpec) map[string]string {
	return map[string]string{
		"app":                 spec.Name,
		"fiaas/version":       strconv.Itoa(spec.Version),
		"fiaas/deployment_id": spec.DeploymentID,
	}
}
