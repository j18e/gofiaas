package deploy

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/j18e/gofiaas/spec/core"
)

func ownerReferences(spec core.Spec) []metav1.OwnerReference {
	uid := spec.UID
	if uid == "" {
		return nil
	}
	return []metav1.OwnerReference{
		{
			APIVersion:         "fiaas.schibsted.io/v1",
			Kind:               "Application",
			Name:               spec.Name,
			UID:                types.UID(uid),
			Controller:         boolPtr(true),
			BlockOwnerDeletion: boolPtr(true),
		},
	}
}

func boolPtr(b bool) *bool {
	return &b
}
