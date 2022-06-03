package core

import (
	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Spec struct {
	UID          types.UID // used in owner references in created resources
	Name         string
	Image        string
	DeploymentID string

	Annotations fiaasv1.AdditionalLabelsOrAnnotations
	Labels      fiaasv1.AdditionalLabelsOrAnnotations

	// Fields from Application.Spec.Config
	Version              *int
	Replicas             *Replicas
	Ingress              []IngressHost
	Healthchecks         *HealthchecksConfig
	Resources            *corev1.ResourceRequirements
	Metrics              *MetricsConfig
	Ports                []Port
	SecretsInEnvironment *bool
	AdminAccess          *bool
	// TODO Extensions
}

type Replicas struct {
	Minimum                uint `json:"minimum"`
	Maximum                uint `json:"maximum"`
	CPUThresholdPercentage uint `json:"cpu_threshold_percentage"`
	Singleton              bool `json:"singleton"`
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

type Port struct {
	Protocol   string `json:"protocol"`
	Name       string `json:"name"`
	Port       uint   `json:"port"`
	TargetPort uint   `json:"target_port"`
}
