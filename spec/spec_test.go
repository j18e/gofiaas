package spec

import (
	"os"
	"path/filepath"
	"testing"

	fiaasv1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/j18e/gofiaas/spec/core"
)

func Test_convert(t *testing.T) {
	for _, tc := range []struct {
		name string
		exp  core.Spec
	}{
		{
			name: "autoscaling_disabled",
			exp: core.Spec{
				Version: intPtr(Version),
				Replicas: &core.Replicas{
					Minimum: 3,
					Maximum: 3,
				},
				Ingress: nil,
			},
		},
		{
			name: "default_tcp_healthcheck",
			exp: core.Spec{
				Version: intPtr(Version),
				Healthchecks: &core.HealthchecksConfig{
					Liveness: core.Healthcheck{
						TCP: &core.HealthcheckTCP{Port: "liveness-port"},
					},
				},
				Ports: []core.Port{
					{
						Protocol:   "tcp",
						Name:       "liveness-port",
						Port:       8889,
						TargetPort: 8882,
					},
				},
				Ingress: nil,
			},
		},
		{
			name: "exec_check",
			exp: core.Spec{
				Version: intPtr(Version),
				Healthchecks: &core.HealthchecksConfig{
					Liveness: core.Healthcheck{
						Execute: &core.HealthcheckExecute{Command: "/bin/alive"},
					},
					Readiness: core.Healthcheck{
						Execute: &core.HealthcheckExecute{Command: "/bin/ready"},
					},
				},
				Ingress: nil,
			},
		},
		{
			name: "ingress_empty",
			exp: core.Spec{
				Version: intPtr(Version),
				Ingress: []core.IngressHost{},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join("testdata", tc.name+".yml")
			bs, err := os.ReadFile(path)
			require.NoError(t, err)
			var config fiaasv1.Config
			err = yaml.Unmarshal(bs, &config)
			require.NoError(t, err)

			got, err := convert(config)
			require.NoError(t, err)
			assert.Equal(t, tc.exp, *got)
		})
	}
}

func intPtr(i int) *int {
	return &i
}
