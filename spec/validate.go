package spec

import (
	"fmt"

	"github.com/j18e/gofiaas/spec/core"
)

func validate(spec *core.Spec) error {
	portNames := make(map[string]bool)
	if err := validateHealthcheck(spec.Healthchecks.Liveness, portNames); err != nil {
		return fmt.Errorf("liveness check: %w", err)
	}
	if err := validateHealthcheck(spec.Healthchecks.Readiness, portNames); err != nil {
		return fmt.Errorf("readiness check: %w", err)
	}
	return nil
}

func validateIngressHosts(hosts []core.IngressHost, ports map[string]bool) error {
	for _, h := range hosts {
		if h.Host == "" {
			return fmt.Errorf("IngressHost.Name must not be empty")
		}
		for _, p := range h.Paths {
			if p.Path == "" {
				return fmt.Errorf("IngressPath.Path must not be empty")
			}
			if p.Port == "" {
				return fmt.Errorf("IngressPath.Port must not be empty")
			}
			if !ports[p.Port] {
				return fmt.Errorf("port %s not found in application ports", p.Port)
			}
		}
	}
	return nil
}

func validateHealthcheck(hc *core.Healthcheck, ports map[string]bool) error {
	if hc == nil {
		return nil
	}
	if hc.HTTP != nil {
		if hc.HTTP.Path == "" {
			return fmt.Errorf("HTTP healthcheck: path cannot be empty")
		}
		if err := validateHealthcheckPort(hc.HTTP.Port, ports); err != nil {
			return fmt.Errorf("HTTP healthcheck: %w", err)
		}
	}
	if hc.TCP != nil {
		if err := validateHealthcheckPort(hc.TCP.Port, ports); err != nil {
			return fmt.Errorf("TCP healthcheck: %w", err)
		}
	}
	if hc.Execute != nil {
		if hc.Execute.Command == "" {
			return fmt.Errorf("execute healthcheck: command cannot be empty")
		}
	}
	return nil
}

func validateHealthcheckPort(port string, ports map[string]bool) error {
	if port == "" {
		return fmt.Errorf("port name cannot be empty")
	}
	if !ports[port] {
		return fmt.Errorf("port %s not found in application ports", port)
	}
	return nil
}
