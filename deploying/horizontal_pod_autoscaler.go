package deploying

import (
	"context"

	"github.com/j18e/gofiaas/models"
	"k8s.io/client-go/kubernetes"
	clientautoscalingv1 "k8s.io/client-go/kubernetes/typed/autoscaling/v1"
)

type horizontalPodAutoscalerDeployer struct {
	namespace   string
	autoscalers clientautoscalingv1.HorizontalPodAutoscalerInterface
}

func newHorizontalPodAutoscalerDeployer(k8s kubernetes.Interface, namespace string) *horizontalPodAutoscalerDeployer {
	return &horizontalPodAutoscalerDeployer{
		namespace:   namespace,
		autoscalers: k8s.AutoscalingV1().HorizontalPodAutoscalers(namespace),
	}
}

func (d *horizontalPodAutoscalerDeployer) Deploy(ctx context.Context, spec models.ApplicationSpec) error {
	return nil
}
