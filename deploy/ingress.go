package deploy

import (
	"context"
	"regexp"

	clientnetworkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"

	"github.com/j18e/gofiaas/spec/core"
)

type ingressDeployerConfig struct {
	suffixes         []string
	hostRewriteRules map[*regexp.Regexp]string
}

type ingressDeployer struct {
	ingresses clientnetworkingv1.IngressInterface
	ingressDeployerConfig
}

func newIngressDeployer(ingresses clientnetworkingv1.IngressInterface, cfg ingressDeployerConfig) *ingressDeployer {
	return &ingressDeployer{
		ingresses:             ingresses,
		ingressDeployerConfig: cfg,
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
