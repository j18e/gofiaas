package spec

import "github.com/j18e/gofiaas/spec/core"

type Factory interface {
	Render(core.Spec) core.Spec
}

type factory struct{}

func (f *factory) Render(in core.Spec) core.Spec {
	out := Default()
	out.UID = in.UID
	out.Name = in.Name
	out.Image = in.Image
	out.Annotations = in.Annotations
	out.Labels = in.Labels

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
	return out
}
