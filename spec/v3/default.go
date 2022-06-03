package v3

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	defaultCPULimit      = resource.MustParse("400m")
	defaultMemoryLimit   = resource.MustParse("512Mi")
	defaultCPURequest    = resource.MustParse("200m")
	defaultMemoryRequest = resource.MustParse("256Mi")
)

func DefaultSpec() ApplicationSpec {
	return ApplicationSpec{
		Replicas: ReplicasConfig{
			Minimum:                2,
			Maximum:                5,
			CPUThresholdPercentage: 50,
			Singleton:              true,
		},
		Ingress: []IngressHost{
			{
				Paths: []IngressPath{
					{
						Path: "/",
						Port: "http",
					},
				},
			},
		},
		Healthchecks: HealthchecksConfig{
			Liveness: Healthcheck{
				HTTP: HealthcheckHTTP{
					Path: "/_/health",
					Port: "http",
				},
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				SuccessThreshold:    1,
				FailureThreshold:    3,
				TimeoutSeconds:      1,
			},
			Readiness: Healthcheck{
				HTTP: HealthcheckHTTP{
					Path: "/_/health",
					Port: "http",
				},
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				SuccessThreshold:    1,
				FailureThreshold:    3,
				TimeoutSeconds:      1,
			},
		},
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    defaultCPULimit,
				corev1.ResourceMemory: defaultMemoryLimit,
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    defaultCPURequest,
				corev1.ResourceMemory: defaultMemoryRequest,
			},
		},
		Metrics: MetricsConfig{
			Prometheus: PrometheusConfig{
				Enabled: true,
				Port:    "http",
				Path:    "/_/metrics",
			},
		},
		Ports: []PortConfig{
			{
				Name:       "http",
				Protocol:   "http",
				Port:       80,
				TargetPort: 8080,
			},
		},
	}
}
