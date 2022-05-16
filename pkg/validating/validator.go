package validating

import (
	"errors"
	"fmt"
	"regexp"

	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
)

var reImage = regexp.MustCompile(`[\w-.]+/[\w-]/[\w-]:\w+`)

type Validator interface {
	Validate(fiaasv1.Application) error
}

func NewValidator() Validator {
	return &validator{}
}

type validator struct{}

func (v *validator) Validate(app fiaasv1.Application) error {
	if app.Name != app.Spec.Application {
		return errors.New("application's metadata.name does not equal spec.application")
	}
	if !reImage.MatchString(app.Spec.Image) {
		return fmt.Errorf("image %s does not match regex %s", app.Spec.Image, reImage)
	}
	return nil
}

func (v *validator) validateConfig(cfg map[string]interface{}) error {
	validFields := []string{
		"version",
		"replicas",
		"ingress",
		"healthchecks",
		"resources",
		"metrics",
		"ports",
		"annotations",
		"labels",
		"secrets_in_environment",
		"admin_access",
		"extensions",
	}
	for k := range cfg {
		if !stringInList(k, validFields) {
			return fmt.Errorf("unrecognized field %s", k)
		}
	}
	return nil
}

func stringInList(s string, sx []string) bool {
	for _, ss := range sx {
		if s == ss {
			return true
		}
	}
	return false
}

// apiVersion: fiaas.schibsted.io/v1
// kind: Application
// metadata:
//   labels:
//     app: fiaas-canary
//     fiaas/app_deployed_at: "1648046198"
//     fiaas/deployment_id: 40dda0259f04e1e695189085ad08a85847b988aa8e8859358d6eed22
//     owner: infrastructure
//     tags.fiaas/fiaas: "true"
//     teams.fiaas/infrastruktur: "true"
//   name: fiaas-canary
// spec:
//   additional_annotations:
//     status:
//       pipeline.finn.no/CallbackURL: https://pipeline.finntech.no/deployment/artifact/974787
//   additional_labels:
//     global:
//       owner: infrastructure
//       tags.fiaas/fiaas: "true"
//       teams.fiaas/infrastruktur: "true"
//   application: fiaas-canary
//   config:
//     admin_access: true
//     annotations:
//       ingress:
//         ingress.kubernetes.io/whitelist-source-range: 172.22.0.0/16
//     healthchecks:
//       liveness:
//         http:
//           path: /ping
//       readiness:
//         http:
//           path: /ping
//     metrics:
//       prometheus:
//         path: /internal-backstage/prometheus
//     replicas:
//       maximum: 1
//       minimum: 1
//     resources:
//       limits:
//         cpu: 1800m
//         memory: 3500Mi
//       requests:
//         cpu: 1600m
//         memory: 3Gi
//     version: 3
//   image: containers.schibsted.io/finntech/fiaas-canary:4e606ec5b9a62722f33dd67c081434dab6b388b4
