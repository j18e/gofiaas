package deploying

import (
	"testing"

	"github.com/j18e/gofiaas/pkg/models"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

func Test_ownerReferences(t *testing.T) {
	name := "app-1"
	var app models.Application
	app.Name = name

	// test app with no UID returns empty list
	got := ownerReferences(app)
	assert.Empty(t, got, "passing in application with no UID: expected empty list")

	uid := uuid.NewUUID()
	app.UID = uid
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
	assert.Equal(t, exp, ownerReferences(app))
}
