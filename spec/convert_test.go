package spec

import (
	"testing"

	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestConverter_Convert(t *testing.T) {
	for _, tc := range []struct {
		name      string
		app       *fiaasv1.Application
		shouldErr bool
	}{
		{
			name:      "goodApplication",
			app:       newApp(appParams{}),
			shouldErr: false,
		},
		{
			name:      "emptyApplication",
			app:       &fiaasv1.Application{},
			shouldErr: true,
		},
		{
			name:      "emptyUID",
			app:       newApp(appParams{emptyUID: true}),
			shouldErr: true,
		},
		{
			name:      "mismatchNames",
			app:       newApp(appParams{mismatchNames: true}),
			shouldErr: true,
		},
		{
			name:      "noDeploymentID",
			app:       newApp(appParams{noDeploymentID: true}),
			shouldErr: true,
		},
		{
			name:      "badImage",
			app:       newApp(appParams{badImage: true}),
			shouldErr: true,
		},
		{
			name:      "invalidConfigField",
			app:       newApp(appParams{invalidConfigField: true}),
			shouldErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			converter := NewConverter()
			_, err := converter.Convert(tc.app)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type appParams struct {
	emptyUID           bool
	mismatchNames      bool
	noDeploymentID     bool
	badImage           bool
	invalidConfigField bool
}

func newApp(params appParams) *fiaasv1.Application {
	const (
		name      = "my-app"
		namespace = "my-namespace"
	)
	app := &fiaasv1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				LabelDeploymentID: "1234",
			},
			UID: types.UID("uid"),
		},
		Spec: fiaasv1.ApplicationSpec{
			Application: name,
			Image:       "docker.io/fiaas/fiaas-deploy-daemon:1.2",
			Config: fiaasv1.Config{
				"version": 3,
			},
		},
	}
	if params.emptyUID {
		app.UID = ""
	}
	if params.mismatchNames {
		app.Spec.Application += "-1"
	}
	if params.noDeploymentID {
		app.Labels[LabelDeploymentID] = ""
	}
	if params.badImage {
		app.Spec.Image = "iaminvalid"
	}
	if params.invalidConfigField {
		app.Spec.Config["someweirdfield"] = 123
	}
	return app
}
