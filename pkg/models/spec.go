package models

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Application is a top-level type. A client is created for it.
type Application struct {
	metav1.TypeMeta   `json:",inline"` // apiVersion, kind
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ApplicationSpec `json:"spec"`
}

type ApplicationSpec struct {
	Application           string              `json:"application"`
	Image                 string              `json:"image"`
	Config                ApplicationConfig   `json:"config"`
	AdditionalLabels      LabelsOrAnnotations `json:"additional_labels,omitempty"`
	AdditionalAnnotations LabelsOrAnnotations `json:"additional_annotations,omitempty"`
}

type ApplicationConfig struct {
	Version              uint                `json:"version"`
	Replicas             ReplicaConfig       `json:"replicas"`
	Ingress              IngressConfig       `json:"ingress"`
	Healthchecks         HealthchecksConfig  `json:"healthchecks"`
	Resources            ResourcesConfig     `json:"resources"`
	Metrics              MetricsConfig       `json:"metrics"`
	Ports                PortsConfig         `json:"ports"`
	Annotations          LabelsOrAnnotations `json:"annotations,omitempty"`
	Labels               LabelsOrAnnotations `json:"labels,omitempty"`
	SecretsInEnvironment bool                `json:"secrets_in_environment"`
	AdminAccess          bool                `json:"admin_access"`
}

type ReplicaConfig struct {
	Minimum                uint `json:"minimum"`
	Maximum                uint `json:"maximum"`
	CPUThresholdPercentage uint `json:"cpu_threshold_percentage"`
	Singleton              bool `json:"singleton"`
}

type IngressConfig []IngressHost

type IngressHost struct {
	Host        string            `json:"host"`
	Paths       IngressPaths      `json:"paths"`
	Annotations map[string]string `json:"annotations"`
}

type IngressPaths []IngressPath

type IngressPath struct {
	Path string `json:"path"`
	Port string `json:"port"`
}

type HealthchecksConfig struct {
	Liveness  Healthcheck `json:"liveness"`
	Readiness Healthcheck `json:"readiness"`
}

type Healthcheck struct {
	Execute             HealthcheckExecute `json:"execute"`
	HTTP                HealthcheckHTTP    `json:"http"`
	TCP                 HealthcheckTCP     `json:"tcp"`
	InitialDelaySeconds uint               `json:"initial_delay_seconds"`
	PeriodSeconds       uint               `json:"period_secconds"`
	SuccessThreshold    uint               `json:"success_threshold"`
	FailureThreshold    uint               `json:"failure_threshold"`
	TimeoutSeconds      uint               `json:"timeout_seconds"`
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
	Port uint `json:"port"`
}

type ResourcesConfig struct {
	Requests Resources `json:"requests"`
	Limits   Resources `json:"limits"`
}

type Resources struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type MetricsConfig struct {
	Prometheus PrometheusConfig `json:"prometheus"`
	Datadog    DatadogConfig    `json:"datadog"`
}

type PrometheusConfig struct {
	Enabled bool   `json:"enabled"`
	Port    string `json:"port"`
	Path    string `json:"path"`
}

type DatadogConfig struct {
	Enabled bool              `json:"enabled"`
	Tags    map[string]string `json:"tags"`
}

type PortsConfig []PortConfig

type PortConfig struct {
	Protocol   string `json:"protocol"`
	Name       string `json:"name"`
	Port       uint   `json:"port"`
	TargetPort uint   `json:"target_port"`
}

type LabelsOrAnnotations struct {
	Deployment              map[string]string `json:"deployment,omitempty"`
	HorizontalPodAutoscaler map[string]string `json:"horizontal_pod_autoscaler,omitempty"`
	Ingress                 map[string]string `json:"ingress,omitempty"`
	Service                 map[string]string `json:"service,omitempty"`
	ServiceAccount          map[string]string `json:"service_account,omitempty"`
	Pod                     map[string]string `json:"pod,omitempty"`
}
