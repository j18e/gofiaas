package converting

import (
	"fmt"
	"regexp"

	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	"github.com/j18e/gofiaas/models"
)

const LabelDeploymentID = "fiaas/deployment_id"

var reContainerImage = regexp.MustCompile(`[\w-.]+/[\w-]+/[\w-]+:[\w-.]+`)

// We use an interface so that reContainerImage will get initialized on startup.
type Converter interface {
	Convert(*fiaasv1.Application) (*models.InternalSpec, error)
}

func NewConverter() Converter {
	return &converter{}
}

type converter struct{}

func (c *converter) Convert(app *fiaasv1.Application) (*models.InternalSpec, error) {
	initMaps(app)
	if app.Name != app.Spec.Application {
		return nil, fmt.Errorf("Name does not match Spec.Name")
	}
	if app.Labels[LabelDeploymentID] == "" {
		return nil, fmt.Errorf("Labels[%s] is missing", LabelDeploymentID)
	}
	if app.UID == "" {
		return nil, fmt.Errorf("UID is missing")
	}
	if !reContainerImage.MatchString(app.Spec.Image) {
		return nil, fmt.Errorf("Spec.image does not match regex %s", reContainerImage)
	}

	spec := &models.InternalSpec{
		Labels:      *app.Spec.AdditionalLabels,
		Annotations: *app.Spec.AdditionalAnnotations,
	}

	for key, val := range app.Spec.Config {
		var ok bool
		switch key {
		case "version":
			spec.Version, ok = val.(int)
		case "replicas":
			spec.Replicas, ok = val.(models.ReplicaConfig)
		case "ingress":
			spec.Ingress, ok = val.(models.IngressConfig)
		case "healthchecks":
			spec.Healthchecks, ok = val.(models.HealthchecksConfig)
		case "resources":
			spec.Resources, ok = val.(models.ResourcesConfig)
		case "metrics":
			spec.Metrics, ok = val.(models.MetricsConfig)
		case "ports":
			spec.Ports, ok = val.(models.PortsConfig)
		case "secrets_in_environment":
			spec.SecretsInEnvironment, ok = val.(bool)
		case "admin_access":
			spec.AdminAccess, ok = val.(bool)
		case "extensions":
		// TODO
		default:
			return nil, fmt.Errorf("Spec.Config.%s: unrecognized field", key)
		}
		if !ok {
			return nil, fmt.Errorf("Spec.Config[%s]: could not convert to int", key)
		}
	}

	return spec, nil
}

func initMaps(app *fiaasv1.Application) {
	if app.Labels == nil {
		app.Labels = make(map[string]string)
	}
	if app.Annotations == nil {
		app.Annotations = make(map[string]string)
	}
	if app.Spec.AdditionalLabels == nil {
		app.Spec.AdditionalLabels = &fiaasv1.AdditionalLabelsOrAnnotations{}
	}
	if app.Spec.AdditionalAnnotations == nil {
		app.Spec.AdditionalAnnotations = &fiaasv1.AdditionalLabelsOrAnnotations{}
	}
}
