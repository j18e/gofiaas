package deploy

import (
	"context"

	"github.com/j18e/gofiaas/spec/core"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
)

type serviceAccountDeployer struct {
	serviceAccounts clientcorev1.ServiceAccountInterface
	lister          listerscorev1.ServiceAccountNamespaceLister
}

func newServiceAccountDeployer(serviceAccounts clientcorev1.ServiceAccountInterface, lister listerscorev1.ServiceAccountNamespaceLister) *serviceAccountDeployer {
	return &serviceAccountDeployer{
		serviceAccounts: serviceAccounts,
		lister:          lister,
	}
}

func (d *serviceAccountDeployer) String() string {
	return "service-account-deployer"
}

func (d *serviceAccountDeployer) Deploy(ctx context.Context, spec core.Spec) error {
	return nil
}

func (d *serviceAccountDeployer) Delete(ctx context.Context, spec core.Spec) error {
	err := d.serviceAccounts.Delete(ctx, spec.Name, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *serviceAccountDeployer) fiaasOwned(meta metav1.ObjectMeta) bool {
	for _, ref := range meta.OwnerReferences {
		if ref.APIVersion == "fiaas.schibsted.io/v1" && ref.Kind == "Application" {
			return true
		}
	}
	return false
}
