package deploying

import (
	"context"

	"github.com/j18e/gofiaas/models"
	"k8s.io/client-go/kubernetes"
	clientnetworkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
)

type ingressDeployer struct {
	namespace string
	ingresses clientnetworkingv1.IngressInterface
}

func newIngressDeployer(k8s kubernetes.Interface, namespace string) *ingressDeployer {
	return &ingressDeployer{
		namespace: namespace,
		ingresses: k8s.NetworkingV1().Ingresses(namespace),
	}
}

func (d *ingressDeployer) Deploy(ctx context.Context, spec models.ApplicationSpec) error {
	return nil
}
