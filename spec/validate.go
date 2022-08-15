package spec

import (
	"fmt"
	"path"

	"github.com/j18e/gofiaas/spec/core"
)

func validate(spec *core.Spec) error {
	if err := validatePorts(spec.Ports); err != nil {
		return fmt.Errorf("ports: %w", err)
	}
	if err := validateIngressHosts(spec.Ingress, spec.Ports); err != nil {
		return fmt.Errorf("ingress: %w", err)
	}
	if err := validateHealthcheck(spec.Healthchecks.Liveness, spec.Ports); err != nil {
		return fmt.Errorf("liveness check: %w", err)
	}
	if err := validateHealthcheck(spec.Healthchecks.Readiness, spec.Ports); err != nil {
		return fmt.Errorf("readiness check: %w", err)
	}
	return nil
}

func validatePorts(ports []core.Port) error {
	names := make(map[string]bool)
	nums := make(map[int]bool)
	for _, p := range ports {
		switch {
		case p.Protocol != "tcp" && p.Protocol != "http":
			return fmt.Errorf("unknown protocol %s", p.Protocol)
		case p.Name == "":
			return fmt.Errorf("port name cannot be empty")
		case p.Port == 0:
			return fmt.Errorf("port number cannot be empty or 0")
		case names[p.Name]:
			return fmt.Errorf("found duplicate port name %s", p.Name)
		case nums[p.Port]:
			return fmt.Errorf("found duplicate port %d", p.Port)
		}
		names[p.Name] = true
		nums[p.Port] = true
	}
	return nil
}

func validateIngressHosts(hosts []core.IngressHost, ports portList) error {
	for _, h := range hosts {
		if h.Host == "" {
			return fmt.Errorf("host name must not be empty")
		}
		for _, p := range h.Paths {
			if p.Path == "" {
				return fmt.Errorf("host %s: path must not be empty", h.Host)
			}
			if p.Path[0] != '/' {
				return fmt.Errorf("host %s: path must start with /", h.Host)
			}
			hostPath := path.Join(h.Host, p.Path)
			switch {
			case p.Port.Name != "" && !ports.nameFound("http", p.Port.Name):
				return fmt.Errorf("%s: http port %s not found in application ports", hostPath, p.Port.Name)
			case p.Port.Number != 0 && !ports.numFound("http", p.Port.Number):
				return fmt.Errorf("%s: http port %d not found in application ports", hostPath, p.Port.Number)
			case p.Port.Name == "" && p.Port.Number == 0:
				return fmt.Errorf("%s: port name or number is required", hostPath)
			}
		}
	}
	return nil
}

func validateHealthcheck(hc *core.Healthcheck, ports portList) error {
	if hc == nil {
		return nil
	}
	if hc.HTTP != nil {
		if hc.HTTP.Path == "" {
			return fmt.Errorf("HTTP healthcheck: path cannot be empty")
		}
		if err := validateHealthcheckPort("http", hc.HTTP.Port, ports); err != nil {
			return fmt.Errorf("HTTP healthcheck: %w", err)
		}
	}
	if hc.TCP != nil {
		if err := validateHealthcheckPort("tcp", hc.TCP.Port, ports); err != nil {
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

func validateHealthcheckPort(protocol, port string, ports portList) error {
	if port == "" {
		return fmt.Errorf("port name cannot be empty")
	}
	if !ports.nameFound(protocol, port) {
		return fmt.Errorf("port %s not found in application ports", port)
	}
	return nil
}

type portList []core.Port

func (l *portList) nameFound(protocol, name string) bool {
	for _, p := range *l {
		if p.Protocol == protocol && p.Name == name {
			return true
		}
	}
	return false
}

func (l *portList) numFound(protocol string, num int) bool {
	for _, p := range *l {
		if p.Protocol == protocol && p.Port == num {
			return true
		}
	}
	return false
}
