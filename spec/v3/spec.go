package v3

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/j18e/gofiaas/spec/core"
)

type ApplicationSpec struct {
	Version              uint                        `json:"version"`
	Replicas             core.ReplicasConfig         `json:"replicas"`
	Ingress              []core.IngressHost          `json:"ingress"`
	Healthchecks         core.HealthchecksConfig     `json:"healthchecks"`
	Resources            corev1.ResourceRequirements `json:"resources"`
	Metrics              core.MetricsConfig          `json:"metrics"`
	Ports                []core.PortConfig           `json:"ports"`
	SecretsInEnvironment bool                        `json:"secrets_in_environment"`
	AdminAccess          bool                        `json:"admin_access"`
}
