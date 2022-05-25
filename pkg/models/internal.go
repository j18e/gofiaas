package models

import fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"

type InternalSpec struct {
	UID          string // used in owner references in created resources
	Name         string
	Image        string
	DeploymentID string

	Annotations fiaasv1.AdditionalLabelsOrAnnotations
	Labels      fiaasv1.AdditionalLabelsOrAnnotations

	// Fields from Application.Spec.Config
	Version              int
	Replicas             ReplicaConfig
	Ingress              IngressConfig
	Healthchecks         HealthchecksConfig
	Resources            ResourcesConfig
	Metrics              MetricsConfig
	Ports                PortsConfig
	SecretsInEnvironment bool
	AdminAccess          bool
	// TODO Extensions
}
