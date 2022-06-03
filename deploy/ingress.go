package deploy

import (
	"context"
	"regexp"

	clientnetworkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"

	"github.com/j18e/gofiaas/spec/core"
)

type IngressDeployerConfig struct {
	Suffixes         []string
	HostRewriteRules map[*regexp.Regexp]string
}

type ingressDeployer struct {
	ingresses clientnetworkingv1.IngressInterface
	IngressDeployerConfig
}

func newIngressDeployer(ingresses clientnetworkingv1.IngressInterface, cfg IngressDeployerConfig) *ingressDeployer {
	return &ingressDeployer{
		ingresses:             ingresses,
		IngressDeployerConfig: cfg,
	}
}

func (d *ingressDeployer) String() string {
	return "ingress-deployer"
}

func (d *ingressDeployer) Deploy(ctx context.Context, spec core.Spec) error {
	return nil
}

func (d *ingressDeployer) Delete(ctx context.Context, spec core.Spec) error {
	return nil
}
