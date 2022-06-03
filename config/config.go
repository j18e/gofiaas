package config

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/j18e/gofiaas/deploy"
	"github.com/j18e/gofiaas/web"
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

	Deployer deploy.Config
}

func ParseFlags() (*FlagSet, error) {
	fs := &FlagSet{}
	dCfg := deploy.NewConfig()
	wCfg := web.Config{}

	flag.IntVar(&wCfg.Port, "web.port", 5000, "Listen port for fiaas-deploy-daemon's web interface")
	flag.StringVar(&fs.secretsDir, "secrets-dir", DefaultSecretsDir, "Path to the directory containing secrets")
	flag.StringVar(&fs.environment, "environment", "", "Name of the environment being deployed to")
	serviceType := flag.String("deploy.service-type", "", "Type of Kubernetes service to create for Applications")

	flag.Func("deploy.ingress-suffix", "Suffixes to be used for hosts in created ingresses", func(s string) error {
		dCfg.Ingresses.Suffixes = append(dCfg.Ingresses.Suffixes, s)
		return nil
	})
	flag.Func("deploy.host-rewrite-rule", "<regex-pattern>=<replacement> pair for rewriting ingress hosts",
		hostRewriteRuleFn(dCfg.Ingresses.HostRewriteRules))
	flag.Func("deploy.global-env", "Extra environment variables to add to created deployments in key=val format",
		globalEnvVarsFn(dCfg.Deployments.GlobalEnvVars))

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
			return fmt.Errorf("parsing deploy.global-env %s: required format key=val", s)
		}
		if !deploy.ReEnvVar.MatchString(kvPair[0]) {
			return fmt.Errorf("parsing deploy.global-env %s: variable name must match regex %s", s, deploy.ReEnvVar)
		}
		envVars[kvPair[0]] = kvPair[1]
		return nil
	}
}

func hostRewriteRuleFn(hostMap map[*regexp.Regexp]string) func(string) error {
	return func(s string) error {
		kvPair := strings.Split(s, "=")
		if len(kvPair) != 2 {
			return fmt.Errorf("parsing deploy.host-rewrite-rule %s: required format key=val", s)
		}
		re, err := regexp.Compile(kvPair[0])
		if err != nil {
			return fmt.Errorf("compiling regexp in deploy.host-rewrite-rule %s: %w", s, err)
		}
		hostMap[re] = kvPair[1]
		return nil
	}
}
