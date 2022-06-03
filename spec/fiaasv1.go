package spec

import (
	"fmt"
	"regexp"

	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/j18e/gofiaas/spec/core"
	v3 "github.com/j18e/gofiaas/spec/v3"
)

const LabelDeploymentID = "fiaas/deployment_id"

var reContainerImage = regexp.MustCompile(`[\w-.]+/[\w-]+/[\w-]+:[\w-.]+`)

func FromFIAASV1(app *fiaasv1.Application) (*core.Spec, error) {
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

	spec := &core.Spec{
		Labels:      *app.Spec.AdditionalLabels,
		Annotations: *app.Spec.AdditionalAnnotations,
	}

	requiredFields := []string{"version"}
	for _, fld := range requiredFields {
		if _, ok := app.Spec.Config[fld]; !ok {
			return nil, fmt.Errorf("Spec.Config[%s]: missing required field", fld)
		}
	}

	for key, val := range app.Spec.Config {
		var ok bool
		switch key {
		case "version":
			spec.Version, ok = val.(*int)
		case "replicas":
			spec.Replicas, ok = val.(*core.Replicas)
		case "ingress":
			spec.Ingress, ok = val.([]core.IngressHost)
		case "healthchecks":
			spec.Healthchecks, ok = val.(*core.HealthchecksConfig)
		case "resources":
			spec.Resources, ok = val.(*corev1.ResourceRequirements)
		case "metrics":
			spec.Metrics, ok = val.(*core.MetricsConfig)
		case "ports":
			spec.Ports, ok = val.([]core.Port)
		case "secrets_in_environment":
			spec.SecretsInEnvironment, ok = val.(*bool)
		case "admin_access":
			spec.AdminAccess, ok = val.(*bool)
		case "extensions":
		// TODO
		default:
			return nil, fmt.Errorf("Spec.Config.%s: unrecognized field", key)
		}
		if !ok {
			return nil, fmt.Errorf("Spec.Config[%s]: could not convert to int", key)
		}
	}

	validVersions := []int{v3.Version}
	if !intInSlice(spec.Version, validVersions) {
		return nil, fmt.Errorf("Spec.Config[version] %d: supported values are %v", spec.Version, validVersions)
	}

	return spec, nil
}

func intInSlice(i int, ix []int) bool {
	for _, ii := range ix {
		if i == ii {
			return true
		}
	}
	return false
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
