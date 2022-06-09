package spec

import (
	"fmt"
	"regexp"

	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"

	"github.com/j18e/gofiaas/spec/core"
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

	// set fields in spec with value from app or default
	spec, err := renderV3Spec(app.Spec.Config)
	if err != nil {
		return nil, err
	}

	// merge the labels and annotations with those of the fiaasv1.Application.Spec
	mergeLabelsOrAnnotations(spec.Labels, app.Spec.AdditionalLabels.Global)
	mergeLabelsOrAnnotations(spec.Annotations, app.Spec.AdditionalAnnotations.Global)

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

func stringInSlice(s string, sx []string) bool {
	for _, ss := range sx {
		if s == ss {
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

func mergeLabelsOrAnnotations(base fiaasv1.AdditionalLabelsOrAnnotations, overrides map[string]string) {
	for key, val := range overrides {
		base.Deployment[key] = val
		base.HorizontalPodAutoscaler[key] = val
		base.Ingress[key] = val
		base.Service[key] = val
		base.ServiceAccount[key] = val
		base.Pod[key] = val
	}
}
