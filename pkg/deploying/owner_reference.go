package deploying

import (
	"github.com/j18e/gofiaas/pkg/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ownerReferences(app models.Application) []metav1.OwnerReference {
	uid := app.GetUID()
	if uid == "" {
		return nil
	}
	return []metav1.OwnerReference{
		{
			APIVersion:         "fiaas.schibsted.io/v1",
			Kind:               "Application",
			Name:               app.Name,
			UID:                uid,
			Controller:         boolPtr(true),
			BlockOwnerDeletion: boolPtr(true),
		},
	}
}

func boolPtr(b bool) *bool {
	return &b
}
