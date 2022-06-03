package deploying

import (
	"context"

	"github.com/j18e/gofiaas/models"
	clientautoscalingv1 "k8s.io/client-go/kubernetes/typed/autoscaling/v1"
)

type autoscalerDeployer struct {
	autoscalers clientautoscalingv1.HorizontalPodAutoscalerInterface
}

func newAutoscalerDeployer(autoscalers clientautoscalingv1.HorizontalPodAutoscalerInterface) *autoscalerDeployer {
	return &autoscalerDeployer{
		autoscalers: autoscalers,
	}
}

func (d *autoscalerDeployer) String() string {
	return "autoscaler-deployer"
}

func (d *autoscalerDeployer) Deploy(ctx context.Context, spec models.InternalSpec) error {
	return nil
}

func (d *autoscalerDeployer) Delete(ctx context.Context, spec models.InternalSpec) error {
	return nil
}
