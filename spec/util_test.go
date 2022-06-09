package spec

import (
	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/j18e/gofiaas/spec/core"
)

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
