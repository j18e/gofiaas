package deploying

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/j18e/gofiaas/pkg/models"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
)

func newTestService(name, ns string, labels, selector map[string]string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   ns,
			Labels:      labels,
			Annotations: map[string]string{},
		},
		Spec: corev1.ServiceSpec{
			Selector: selector,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
}

func Test_serviceDeployer_deploy(t *testing.T) {
	var (
		name = "app1"
		ns   = "default"
		ctx  = context.Background()
		ts   = strconv.Itoa(int(time.Now().Unix()))
	)
	app := newTestApp(name, ns)
	labels := map[string]string{
		"app":                   app.Name,
		"fiaas/app_deployed_at": ts,
		"fiaas/deployment_id":   string(app.UID),
	}
	selector := map[string]string{"app": app.Name}

	exp := newTestService(name, ns, labels, selector)
	exp.Spec.Ports = append(exp.Spec.Ports, corev1.ServicePort{
		Name:     "http",
		Protocol: "TCP",
		Port:     80,
	})
	exp.OwnerReferences = ownerReferences(app)

	k8s := fake.NewSimpleClientset()
	sd := newServiceDeployer(k8s, ns, corev1.ServiceTypeClusterIP)
	err := sd.deploy(ctx, app, labels, selector)
	assert.NoError(t, err)

	got, err := k8s.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
	assert.NoError(t, nil)
	assert.Equal(t, exp, got)
}

func Test_serviceDeployer_shouldHaveService(t *testing.T) {
	sd := serviceDeployer{}
	app := models.Application{}
	got := sd.shouldHaveService(app)
	assert.Equal(t, false, got, "application without any ports should not have service")

	app.Spec.Config.Ports = append(app.Spec.Config.Ports, models.PortConfig{})
	got = sd.shouldHaveService(app)
	assert.Equal(t, true, got, "application without 1+ ports should have service")
}

func Test_serviceDeployer_tcpPortAnnotations(t *testing.T) {
	exp := map[string]string{"fiaas/tcp_port_names": "one,two"}
	ports := models.PortsConfig{
		{
			Name:     "http",
			Protocol: "http",
			Port:     80,
		},
		{
			Name:     "one",
			Protocol: "tcp",
			Port:     21,
		},
		{
			Name:     "two",
			Protocol: "tcp",
			Port:     22,
		},
	}
	sd := &serviceDeployer{}
	got := sd.tcpPortAnnotations(ports)
	assert.Equal(t, exp, got)
}

func Test_serviceDeployer_makePorts(t *testing.T) {
	ports := models.PortsConfig{{
		Name:       "http",
		Protocol:   "http",
		Port:       80,
		TargetPort: 8080,
	}}
	exp := []corev1.ServicePort{{
		Name:       "http",
		Protocol:   "TCP",
		Port:       80,
		TargetPort: intstr.FromInt(8080),
	}}
	sd := &serviceDeployer{}
	assert.Equal(t, exp, sd.makePorts(ports))
}
