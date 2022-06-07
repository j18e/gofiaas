package spec

import (
	"errors"
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

	// check for invalid fields
	validFields := []string{
		"version", "replicas", "ingress", "healthchecks", "resources", "metrics", "ports",
		"secrets_in_environment", "admin_access", "extensions",
	}
	for fld := range app.Spec.Config {
		if !stringInSlice(fld, validFields) {
			return nil, fmt.Errorf("Spec.Config[%s]: unrecognized field", fld)
		}
	}

	// check for required fields
	requiredFields := []string{"version"}
	for _, fld := range requiredFields {
		if _, ok := app.Spec.Config[fld]; !ok {
			return nil, fmt.Errorf("Spec.Config[%s]: missing required field", fld)
		}
	}

	// set fields in spec with value from app or default

	// the default version being set is intentionally invalid
	if err := setInt(app.Spec.Config["version"], -1, &spec.Version); err != nil {
		return nil, fmt.Errorf("Spec.Config[version]: %w", err)
	}
	if err := setReplicas(app.Spec.Config["replicas"], spec); err != nil {
		return nil, fmt.Errorf("Spec.Config[replicas]: %w", err)
	}
	if err := setIngress(app.Spec.Config["ingress"], spec); err != nil {
		return nil, fmt.Errorf("Spec.Config[ingress]: %w", err)
	}
	if err := setHealthchecks(app.Spec.Config["healthchecks"], spec); err != nil {
		return nil, fmt.Errorf("Spec.Config[healthchecks]: %w", err)
	}
	if err := setResources(app.Spec.Config["resources"], spec); err != nil {
		return nil, fmt.Errorf("Spec.Config[resources]: %w", err)
	}
	if err := setMetrics(app.Spec.Config["metrics"], spec); err != nil {
		return nil, fmt.Errorf("Spec.Config[metrics]: %w", err)
	}
	if err := setPorts(app.Spec.Config["ports"], spec); err != nil {
		return nil, fmt.Errorf("Spec.Config[ports]: %w", err)
	}
	if err := setBool(app.Spec.Config["secrets_in_environment"], defaultSecretsInEnvironment(),
		&spec.SecretsInEnvironment); err != nil {
		return nil, fmt.Errorf("Spec.Config[secrets_in_environment]: %w", err)
	}
	if err := setBool(app.Spec.Config["admin_access"], defaultAdminAccess(), &spec.AdminAccess); err != nil {
		return nil, fmt.Errorf("Spec.Config[admin_access]: %w", err)
	}

	// validate given version
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

func setInt(val interface{}, defaultVal int, target *int) error {
	if val == nil {
		*target = defaultVal
		return nil
	}
	var ok bool
	*target, ok = val.(int)
	if !ok {
		return fmt.Errorf("could not convert %v to int", val)
	}
	return nil
}

func setBool(val interface{}, defaultVal bool, target *bool) error {
	if val == nil {
		*target = defaultVal
		return nil
	}
	var ok bool
	*target, ok = val.(bool)
	if !ok {
		return fmt.Errorf("could not convert %v to bool", val)
	}
	return nil
}

func setReplicas(val interface{}, spec *core.Spec) error {
	if val == nil {
		spec.Replicas = defaultReplicas()
		return nil
	}
	var ok bool
	spec.Replicas, ok = val.(core.Replicas)
	if !ok {
		return errors.New("could not convert to core.Replicas")
	}
	return nil
}

func setIngress(val interface{}, spec *core.Spec) error {
	if val == nil {
		spec.Ingress = defaultIngress()
		return nil
	}
	var ok bool
	spec.Ingress, ok = val.([]core.IngressHost)
	if !ok {
		return errors.New("could not convert to []core.IngressHost")
	}
	return nil
}

func setHealthchecks(val interface{}, spec *core.Spec) error {
	if val == nil {
		spec.Healthchecks = defaultHealthchecks()
		return nil
	}
	var ok bool
	spec.Healthchecks, ok = val.(core.HealthchecksConfig)
	if !ok {
		return errors.New("could not convert to []core.HealthchecksConfig")
	}
	return nil
}

func setResources(val interface{}, spec *core.Spec) error {
	if val == nil {
		spec.Resources = defaultResources()
		return nil
	}
	var ok bool
	spec.Resources, ok = val.(corev1.ResourceRequirements)
	if !ok {
		return errors.New("could not convert to corev1.ResourceRequirements")
	}
	return nil
}

func setMetrics(val interface{}, spec *core.Spec) error {
	if val == nil {
		spec.Metrics = defaultMetrics()
		return nil
	}
	var ok bool
	spec.Metrics, ok = val.(core.MetricsConfig)
	if !ok {
		return errors.New("could not convert to core.MetricsConfig")
	}
	return nil
}

func setPorts(val interface{}, spec *core.Spec) error {
	if val == nil {
		spec.Ports = defaultPorts()
		return nil
	}
	var ok bool
	spec.Ports, ok = val.([]core.Port)
	if !ok {
		return errors.New("could not convert to []core.Port")
	}
	return nil
}
