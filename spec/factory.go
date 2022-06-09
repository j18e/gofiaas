package spec

import (
	"encoding/json"
	"fmt"
	"os"

	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	"github.com/j18e/gofiaas/spec/core"
	"gopkg.in/yaml.v2"
)

type Factory struct {
	defaults *core.Spec
}

func NewFactory() (*Factory, error) {
	const file = `defaults.yml`
	bs, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var spec core.Spec
	if err := yaml.Unmarshal(bs, &spec); err != nil {
		return nil, err
	}
	// we'll validate the default spec just in case
	if err := validate(&spec); err != nil {
		return nil, err
	}
	return &Factory{defaults: &spec}, nil
}

func (f *Factory) Transform(app fiaasv1.Application) (*core.Spec, error) {
	spec, err := convert(app.Spec.Config)
	if err != nil {
		return nil, err
	}
	applyDefaults(spec, f.defaults)
	if err := validate(spec); err != nil {
		return nil, err
	}
	// TODO merge additional global labels into spec
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
	if *spec.Version != Version {
		return nil, fmt.Errorf("expected version %d, got %d", Version, *spec.Version)
	}
	return &spec, nil
}
