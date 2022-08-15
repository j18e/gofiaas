package spec

import (
	"encoding/json"
	"fmt"
	"regexp"

	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	"github.com/j18e/gofiaas/spec/core"
	v3 "github.com/j18e/gofiaas/spec/v3"
	"gopkg.in/yaml.v2"
)

const LabelDeploymentID = "fiaas/deployment_id"

var reContainerImage = regexp.MustCompile(`[\w-.]+/[\w-]+/[\w-]+:[\w-.]+`)

type Factory interface {
	Transform(fiaasv1.Application) (*core.Spec, error)
}

func NewFactory(version int) (Factory, error) {
	var bs []byte
	switch version {
	case v3.Version:
		bs = v3.Defaults
	default:
		return nil, fmt.Errorf("version %d unrecognized", version)
	}
	var spec core.Spec
	if err := yaml.Unmarshal(bs, &spec); err != nil {
		return nil, err
	}
	// we'll validate the default spec just in case
	if err := validate(&spec); err != nil {
		return nil, err
	}
	return &specFactory{defaults: &spec}, nil
}

type specFactory struct {
	defaults *core.Spec
}

func (f *specFactory) Transform(app fiaasv1.Application) (*core.Spec, error) {
	if app.Name != app.Spec.Application {
		return nil, fmt.Errorf("Name does not match Spec.Name")
	}
	if app.Labels == nil || app.Labels[LabelDeploymentID] == "" {
		return nil, fmt.Errorf("Labels[%s] is missing", LabelDeploymentID)
	}
	if app.UID == "" {
		return nil, fmt.Errorf("UID is missing")
	}
	if !reContainerImage.MatchString(app.Spec.Image) {
		return nil, fmt.Errorf("Spec.image does not match regex %s", reContainerImage)
	}
	if app.Spec.Config == nil {
		return nil, fmt.Errorf("missing Spec.Config")
	}

	spec, err := convert(app.Spec.Config)
	if err != nil {
		return nil, err
	}
	applyDefaults(spec, f.defaults)
	if err := validate(spec); err != nil {
		return nil, err
	}

	spec.Name = app.Name
	spec.DeploymentID = app.Labels[LabelDeploymentID]
	spec.UID = app.UID
	spec.Image = app.Spec.Image

	for k, v := range app.Spec.AdditionalLabels.Global {
		spec.AddGlobalLabel(k, v)
	}
	// TODO what to do with app.Spec.AdditionalAnnotations.Status
	return spec, nil
}

func applyDefaults(spec, defaults *core.Spec) {
	if spec.Version == nil {
		spec.Version = defaults.Version
	}
	if spec.Replicas == nil {
		spec.Replicas = defaults.Replicas
	}
	if spec.Ingress == nil {
		spec.Ingress = defaults.Ingress
	}
	if spec.Healthchecks == nil {
		spec.Healthchecks = defaults.Healthchecks
	} else {
		if spec.Healthchecks.Liveness != nil && spec.Healthchecks.Readiness == nil {
			spec.Healthchecks.Readiness = spec.Healthchecks.Liveness
		}
	}
	if spec.Resources == nil {
		spec.Resources = defaults.Resources
	}
	if spec.Metrics == nil {
		spec.Metrics = defaults.Metrics
	}
	if spec.Ports == nil {
		spec.Ports = defaults.Ports
	}
	if spec.SecretsInEnvironment == nil {
		spec.SecretsInEnvironment = defaults.SecretsInEnvironment
	}
	if spec.AdminAccess == nil {
		spec.AdminAccess = defaults.AdminAccess
	}
	if spec.Labels == nil {
		spec.Labels = defaults.Labels
	}
	if spec.Annotations == nil {
		spec.Annotations = defaults.Annotations
	}
}

// convert does the work of converting between fiaas-go-client's Config type to
// our v3.Spec. Since fiaasv1.Config is a map[string]interface{} it's easiest
// to marshal it into json before unmarshaling it into a v3.Spec.
func convert(config fiaasv1.Config) (*core.Spec, error) {
	bs, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("getting json representation of spec: %w", err)
	}
	var spec core.Spec
	if err := json.Unmarshal(bs, &spec); err != nil {
		return nil, fmt.Errorf("unmarshaling json representation of spec: %w", err)
	}

	if spec.Version == nil {
		return nil, fmt.Errorf("version field required but missing")
	}
	switch *spec.Version {
	case v3.Version:
	default:
		return nil, fmt.Errorf("invalid version %d", *spec.Version)
	}
	return &spec, nil
}
