package spec

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/j18e/gofiaas/spec/core"
	v3 "github.com/j18e/gofiaas/spec/v3"
)

var (
	defaultCPULimit      = resource.MustParse("400m")
	defaultMemoryLimit   = resource.MustParse("512Mi")
	defaultCPURequest    = resource.MustParse("200m")
	defaultMemoryRequest = resource.MustParse("256Mi")
)

func Default() *core.Spec {
	return &core.Spec{
		Version:      defaultVersion(),
		Replicas:     defaultReplicas(),
		Ingress:      defaultIngress(),
		Healthchecks: defaultHealthchecks(),
		Resources:    defaultResources(),
		Metrics:      defaultMetrics(),
		Ports:        defaultPorts(),
	}
}

func defaultVersion() int {
	return v3.Version
}

func defaultReplicas() core.Replicas {
	return core.Replicas{
		Minimum:                2,
		Maximum:                5,
		CPUThresholdPercentage: 50,
		Singleton:              true,
	}
}

func defaultIngress() []core.IngressHost {
	return []core.IngressHost{
		{
			Paths: []core.IngressPath{
				{
					Path: "/",
					Port: "http",
				},
			},
		},
	}
}

func defaultHealthchecks() core.HealthchecksConfig {
	return core.HealthchecksConfig{
		Liveness: core.Healthcheck{
			HTTP: core.HealthcheckHTTP{
				Path: "/_/health",
				Port: "http",
			},
			InitialDelaySeconds: 10,
			PeriodSeconds:       10,
			SuccessThreshold:    1,
			FailureThreshold:    3,
			TimeoutSeconds:      1,
		},
		Readiness: core.Healthcheck{
			HTTP: core.HealthcheckHTTP{
				Path: "/_/health",
				Port: "http",
			},
			InitialDelaySeconds: 10,
			PeriodSeconds:       10,
			SuccessThreshold:    1,
			FailureThreshold:    3,
			TimeoutSeconds:      1,
		},
	}
}

func defaultResources() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    defaultCPULimit,
			corev1.ResourceMemory: defaultMemoryLimit,
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    defaultCPURequest,
			corev1.ResourceMemory: defaultMemoryRequest,
		},
	}
}

func defaultMetrics() core.MetricsConfig {
	return core.MetricsConfig{
		Prometheus: core.PrometheusConfig{
			Enabled: true,
			Port:    "http",
			Path:    "/_/metrics",
		},
	}
}

func defaultPorts() []core.Port {
	return []core.Port{
		{
			Name:       "http",
			Protocol:   "http",
			Port:       80,
			TargetPort: 8080,
		},
	}
}

func defaultSecretsInEnvironment() bool {
	return false
}

func defaultAdminAccess() bool {
	return false
}
