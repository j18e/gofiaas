package deploying

import (
	"context"

	"github.com/j18e/gofiaas/spec/core"
	corev1 "k8s.io/api/core/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type serviceDeployer struct {
	serviceType corev1.ServiceType
	services    clientcorev1.ServiceInterface
}

func newServiceDeployer(svc clientcorev1.ServiceInterface, serviceType corev1.ServiceType) *serviceDeployer {
	return &serviceDeployer{
		services:    svc,
		serviceType: serviceType,
	}
}

func (d *serviceDeployer) String() string {
	return "service-deployer"
}

func (d *serviceDeployer) Deploy(ctx context.Context, app core.Spec) error {
	return nil
}

func (d *serviceDeployer) Delete(ctx context.Context, spec core.Spec) error {
	return nil
}
