package spec

import (
	"testing"

	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	"github.com/j18e/gofiaas/spec/core"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestFromFIAASV1(t *testing.T) {
	for _, tc := range []struct {
		name      string
		app       *fiaasv1.Application
		exp       *core.Spec
		shouldErr bool
	}{
		{
			name:      "goodApplication",
			app:       newApp(appParams{}),
			shouldErr: false,
			exp:       Default(),
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
		{
			name:      "noVersion",
			app:       newApp(appParams{noVersion: true}),
			shouldErr: true,
		},
		{
			name:      "badVersion",
			app:       newApp(appParams{config: msi{"version": 17}}),
			shouldErr: true,
		},
		{
			name: "customReplicas",
			app: newApp(appParams{config: msi{"replicas": core.Replicas{
				Minimum: 2, Maximum: 3,
			}}}),
			shouldErr: false,
			exp:       newSpec(specParams{replicas: &core.Replicas{Minimum: 2, Maximum: 3}}),
		},
		{
			name:      "emptyIngress",
			app:       newApp(appParams{config: msi{"ingress": []core.IngressHost{}}}),
			shouldErr: false,
			exp:       newSpec(specParams{ingress: []core.IngressHost{}}),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FromFIAASV1(tc.app)
			if tc.shouldErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.exp, got)
		})
	}
}

type msi map[string]interface{}

type specParams struct {
	replicas *core.Replicas
	ingress  []core.IngressHost
}

func newSpec(params specParams) *core.Spec {
	res := Default()
	if params.replicas != nil {
		res.Replicas = *params.replicas
	}
	if params.ingress != nil {
		res.Ingress = params.ingress
	}
	return res
}

type appParams struct {
	config             msi
	emptyUID           bool
	mismatchNames      bool
	noDeploymentID     bool
	badImage           bool
	invalidConfigField bool
	noVersion          bool
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
	for k, v := range params.config {
		app.Spec.Config[k] = v
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
	if params.noVersion {
		delete(app.Spec.Config, "version")
	}
	return app
}
