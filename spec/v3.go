package spec

import (
	"github.com/j18e/gofiaas/spec/core"
	v3 "github.com/j18e/gofiaas/spec/v3"
)

func ToV3(in core.Spec) v3.Spec {
	return v3.Spec{
		Version:              3,
		Replicas:             in.Replicas,
		Ingress:              in.Ingress,
		Healthchecks:         in.Healthchecks,
		Resources:            in.Resources,
		Metrics:              in.Metrics,
		Ports:                in.Ports,
		SecretsInEnvironment: in.SecretsInEnvironment,
		AdminAccess:          in.AdminAccess,
	}
}
