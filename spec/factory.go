package spec

import (
	"github.com/j18e/gofiaas/spec/core"
	v3 "github.com/j18e/gofiaas/spec/v3"
)

func RenderSpec(in core.Spec) core.Spec {
	out := Default()
	out.UID = in.UID
	out.Name = in.Name
	out.Image = in.Image
	out.Annotations = in.Annotations
	out.Labels = in.Labels

	if in.Version != nil {
		out.Version = in.Version
	}
	if in.Replicas != nil {
		out.Replicas = in.Replicas
	}
	if in.Ingress != nil {
		out.Ingress = in.Ingress
	}
	if in.Healthchecks != nil {
		out.Healthchecks = in.Healthchecks
	}
	if in.Resources != nil {
		out.Resources = in.Resources
	}
	if in.Metrics != nil {
		out.Metrics = in.Metrics
	}
	if in.Ports != nil {
		out.Ports = in.Ports
	}
	if in.SecretsInEnvironment != nil {
		out.SecretsInEnvironment = in.SecretsInEnvironment
	}
	if in.AdminAccess != nil {
		out.AdminAccess = in.AdminAccess
	}
	return out
}

func ToV3(in core.Spec) v3.Spec {
	return v3.Spec{
		Version:              3,
		Replicas:             in.Replicas,
		Ingress:              in.Ingress,
		Healthchecks:         in.Healthchecks,
		Resources:            in.Resources,
		Metrics:              in.Metrics,
		Ports:                in.Ports,
		SecretsInEnvironment: *in.SecretsInEnvironment,
		AdminAccess:          *in.AdminAccess,
	}
}
