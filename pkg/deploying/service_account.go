package deploying

import (
	"context"

	"github.com/j18e/gofiaas/pkg/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type serviceAccountDeployer struct {
	namespace       string
	serviceAccounts clientcorev1.ServiceAccountInterface
}

func newServiceAccountDeployer(k8s kubernetes.Interface, namespace string) *serviceAccountDeployer {
	return &serviceAccountDeployer{
		namespace:       namespace,
		serviceAccounts: k8s.CoreV1().ServiceAccounts(namespace),
	}
}

func (d *serviceAccountDeployer) Deploy(ctx context.Context, spec models.ApplicationSpec) error {
	return nil
}

func fiaasOwned(meta metav1.ObjectMeta) bool {
	for _, ref := range meta.OwnerReferences {
		if ref.APIVersion == "fiaas.schibsted.io/v1" && ref.Kind == "Application" {
			return true
		}
	}
	return false
}
