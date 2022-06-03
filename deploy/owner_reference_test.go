package deploy

import (
	"testing"

	"github.com/j18e/gofiaas/spec/core"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
)

func Test_ownerReferences(t *testing.T) {
	name := "app-1"
	spec := core.Spec{Name: name}

	// test app with no UID returns empty list
	got := ownerReferences(spec)
	assert.Empty(t, got, "passing in application with no UID: expected empty list")

	uid := uuid.NewUUID()
	spec.UID = types.UID(uid)
	exp := []metav1.OwnerReference{
		{
			APIVersion:         "fiaas.schibsted.io/v1",
			Kind:               "Application",
			Name:               name,
			UID:                uid,
			Controller:         boolPtr(true),
			BlockOwnerDeletion: boolPtr(true),
		},
	}
	assert.Equal(t, exp, ownerReferences(spec))
}
