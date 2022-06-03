package config

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/j18e/gofiaas/deploying"
	corev1 "k8s.io/api/core/v1"
)

const (
	DefaultConfigFile = "/var/run/config/fiaas/cluster_config.yaml"
	DefaultSecretsDir = "/var/run/secrets/fiaas"
)

type FlagSet struct {
	secretsDir  string
	environment string
	webPort     int

	Deployer deploying.Config
}

func ParseFlags() (*FlagSet, error) {
	fs := &FlagSet{}

	flag.IntVar(&fs.webPort, "web.port", 5000, "Listen port for fiaas-deploy-daemon's web interface")
	flag.StringVar(&fs.secretsDir, "secrets-dir", DefaultSecretsDir, "Path to the directory containing secrets")
	flag.StringVar(&fs.environment, "environment", "", "Name of the environment being deployed to")
	serviceType := flag.String("app.service-type", "", "Type of Kubernetes service to create for Applications")

	flag.Func("app.ingress-suffix", "Suffixes to be used for hosts in created ingresses", func(s string) error {
		fs.Deployer.IngressSuffixes = append(fs.Deployer.IngressSuffixes, s)
		return nil
	})
	fs.Deployer.HostRewriteRules = make(map[*regexp.Regexp]string)
	flag.Func("app.host-rewrite-rule", "<regex-pattern>=<replacement> pair for rewriting ingress hosts",
		hostRewriteRuleFn(fs.Deployer.HostRewriteRules))
	fs.Deployer.GlobalEnvVars = make(map[string]string)
	flag.Func("app.global-env", "Extra environment variables to add to created deployments in key=val format",
		globalEnvVarsFn(fs.Deployer.GlobalEnvVars))

	flag.Parse()

	if fs.environment == "" {
		return nil, errors.New("flag environment required but not specified")
	}
	if *serviceType == "" {
		return nil, errors.New("flag service-type required but not specified")
	}
	switch corev1.ServiceType(*serviceType) {
	case corev1.ServiceTypeClusterIP,
		corev1.ServiceTypeNodePort,
		corev1.ServiceTypeLoadBalancer:
	default:
		return nil, fmt.Errorf("unrecognized service-type %s", *serviceType)
	}
	fs.Deployer.ServiceType = corev1.ServiceType(*serviceType)
	return fs, nil
}

func globalEnvVarsFn(envVars map[string]string) func(string) error {
	return func(s string) error {
		kvPair := strings.Split(s, "=")
		if len(kvPair) != 2 {
			return fmt.Errorf("parsing app.global-env %s: required format key=val", s)
		}
		if !deploying.ReEnvVar.MatchString(kvPair[0]) {
			return fmt.Errorf("parsing app.global-env %s: variable name must match regex %s", s, deploying.ReEnvVar)
		}
		envVars[kvPair[0]] = kvPair[1]
		return nil
	}
}

func hostRewriteRuleFn(hostMap map[*regexp.Regexp]string) func(string) error {
	return func(s string) error {
		kvPair := strings.Split(s, "=")
		if len(kvPair) != 2 {
			return fmt.Errorf("parsing app.host-rewrite-rule %s: required format key=val", s)
		}
		re, err := regexp.Compile(kvPair[0])
		if err != nil {
			return fmt.Errorf("compiling regexp in app.host-rewrite-rule %s: %w", s, err)
		}
		if hostMap == nil {
			hostMap = make(map[*regexp.Regexp]string)
		}
		hostMap[re] = kvPair[1]
		return nil
	}
}
