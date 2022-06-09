package core

import (
	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Spec struct {
	UID          types.UID `json:"-"` // used in owner references in created resources
	Name         string    `json:"-"`
	Image        string    `json:"-"`
	DeploymentID string    `json:"-"`

	AdditionalAnnotations *fiaasv1.AdditionalLabelsOrAnnotations `json:"-"`
	AdditionalLabels      *fiaasv1.AdditionalLabelsOrAnnotations `json:"-"`

	// fiaas.yml fields
	Version              *int                         `json:"version,omitempty"`
	Replicas             *Replicas                    `json:"replicas,omitempty"`
	Ingress              []IngressHost                `json:"ingress,omitempty"`
	Healthchecks         *HealthchecksConfig          `json:"healthchecks,omitempty"`
	Resources            *corev1.ResourceRequirements `json:"resources,omitempty"`
	Metrics              *MetricsConfig               `json:"metrics,omitempty"`
	Ports                []Port                       `json:"ports,omitempty"`
	SecretsInEnvironment *bool                        `json:"secrets_in_environment,omitempty"`
	AdminAccess          *bool                        `json:"admin_access,omitempty"`
	Labels               *LabelsOrAnnotations         `json:"labels,omitempty"`
	Annotations          *LabelsOrAnnotations         `json:"annotations,omitempty"`
}

type LabelsOrAnnotations struct {
	// Global                  map[string]string `json:"global,omitempty"`
	Deployment              map[string]string `json:"deployment,omitempty"`
	HorizontalPodAutoscaler map[string]string `json:"horizontal_pod_autoscaler,omitempty"`
	Ingress                 map[string]string `json:"ingress,omitempty"`
	Service                 map[string]string `json:"service,omitempty"`
	ServiceAccount          map[string]string `json:"service_account,omitempty"`
	Pod                     map[string]string `json:"pod,omitempty"`
	// Status                  map[string]string `json:"status,omitempty"`
}

type Replicas struct {
	Minimum                int  `json:"minimum"`
	Maximum                int  `json:"maximum"`
	CPUThresholdPercentage int  `json:"cpu_threshold_percentage"`
	Singleton              bool `json:"singleton"`
}

func (r *Replicas) AutoscalingEnabled() bool {
	return r.Minimum != r.Maximum
}

type IngressHost struct {
	Host        string            `json:"host"`
	Paths       []IngressPath     `json:"paths"`
	Annotations map[string]string `json:"annotations"`
}

type IngressPath struct {
	Path string `json:"path"`
	Port string `json:"port"`
}

type MetricsConfig struct {
	Prometheus *PrometheusConfig `json:"prometheus,omitempty"`
	Datadog    *DatadogConfig    `json:"datadog,omitempty"`
}

type PrometheusConfig struct {
	Enabled bool   `json:"enabled"`
	Port    string `json:"port"`
	Path    string `json:"path"`
}

type DatadogConfig struct {
	Enabled bool              `json:"enabled"`
	Tags    map[string]string `json:"tags,omitempty"`
}

type Port struct {
	Protocol   string `json:"protocol"`
	Name       string `json:"name"`
	Port       int    `json:"port"`
	TargetPort int    `json:"target_port"`
}

type HealthchecksConfig struct {
	Liveness  *Healthcheck `json:"liveness"`
	Readiness *Healthcheck `json:"readiness"`
}

type Healthcheck struct {
	Execute             *HealthcheckExecute `json:"execute,omitempty"`
	HTTP                *HealthcheckHTTP    `json:"http,omitempty"`
	TCP                 *HealthcheckTCP     `json:"tcp,omitempty"`
	InitialDelaySeconds int                 `json:"initial_delay_seconds"`
	PeriodSeconds       int                 `json:"period_secconds"`
	SuccessThreshold    int                 `json:"success_threshold"`
	FailureThreshold    int                 `json:"failure_threshold"`
	TimeoutSeconds      int                 `json:"timeout_seconds"`
}

type HealthcheckExecute struct {
	Command string `json:"command"`
}

type HealthcheckHTTP struct {
	Path        string            `json:"path"`
	Port        string            `json:"port"`
	HTTPHeaders map[string]string `json:"http_headers"`
}

type HealthcheckTCP struct {
	Port string `json:"port"`
}
