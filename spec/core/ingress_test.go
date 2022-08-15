package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIngressPath_MarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		name      string
		shouldErr bool
		inp       IngressPath
		exp       string
	}{
		{
			name: "portName",
			inp: IngressPath{
				Path: "/foo",
				Port: IngressPort{Name: "foo"},
			},
			exp: `{"path":"/foo","port":"foo"}`,
		},
		{
			name: "portNum",
			inp: IngressPath{
				Path: "/foo",
				Port: IngressPort{Number: 80},
			},
			exp: `{"path":"/foo","port":80}`,
		},
		{
			name:      "missingPath",
			shouldErr: true,
			inp: IngressPath{
				Port: IngressPort{Name: "foo"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(&tc.inp)
			if tc.shouldErr {
				assert.NoError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.exp, string(got))
		})
	}
}

func TestIngressPath_UnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		name      string
		shouldErr bool
		inp       string
		exp       IngressPath
	}{
		{
			name: "portName",
			inp:  `{"path":"/foo","port":"foo"}`,
			exp: IngressPath{
				Path: "/foo",
				Port: IngressPort{Name: "foo"},
			},
		},
		{
			name: "portNum",
			inp:  `{"path":"/foo","port":80}`,
			exp: IngressPath{
				Path: "/foo",
				Port: IngressPort{Number: 80},
			},
		},
		{
			name:      "noPort",
			shouldErr: true,
			inp:       `{"path":"/foo"}`,
		},
		{
			name:      "emptyPortStr",
			shouldErr: true,
			inp:       `{"path":"/foo","port":""}`,
		},
		{
			name:      "emptyPortNum",
			shouldErr: true,
			inp:       `{"path":"/foo","port":0}`,
		},
		{
			name:      "badPortType",
			shouldErr: true,
			inp:       `{"path":"/foo","port":true}`,
		},
		{
			name:      "nullPort",
			shouldErr: true,
			inp:       `{"path":"/foo","port":null}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got IngressPath
			err := json.Unmarshal([]byte(tc.inp), &got)
			if tc.shouldErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.exp, got)
		})
	}
}
