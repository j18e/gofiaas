package deploying

import (
	"context"
	"strings"

	"github.com/j18e/gofiaas/pkg/models"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type serviceDeployer struct {
	namespace   string
	serviceType corev1.ServiceType
	services    clientcorev1.ServiceInterface
}

func newServiceDeployer(k8s kubernetes.Interface, namespace string, serviceType corev1.ServiceType) *serviceDeployer {
	return &serviceDeployer{
		services:    k8s.CoreV1().Services(namespace),
		namespace:   namespace,
		serviceType: serviceType,
	}
}

func (d *serviceDeployer) deploy(ctx context.Context, app models.Application, labels, selector map[string]string) error {
	if !d.shouldHaveService(app) {
		return d.delete(ctx, app.Name)
	}
	var exists bool
	svc, err := d.services.Get(ctx, app.Name, metav1.GetOptions{})
	switch {
	case errors.IsNotFound(err):
		exists = false // just so we're explicit
		svc = &corev1.Service{}
	case err != nil:
		return err
	default:
		exists = true
	}

	svc.Name = app.Name
	svc.Namespace = d.namespace
	svc.Labels = mergeDicts(labels, app.Spec.Config.Labels.Service)
	svc.Annotations = mergeDicts(app.Spec.Config.Annotations.Service, d.tcpPortAnnotations(app.Spec.Config.Ports))
	svc.OwnerReferences = ownerReferences(app)
	svc.Spec = corev1.ServiceSpec{
		Selector: selector,
		Type:     d.serviceType,
		Ports:    d.makePorts(app.Spec.Config.Ports),
	}
	if exists {
		return d.update(ctx, svc)
	}
	return d.create(ctx, svc)
}

func (d *serviceDeployer) create(ctx context.Context, svc *corev1.Service) error {
	_, err := d.services.Create(ctx, svc, metav1.CreateOptions{})
	return err
}

func (d *serviceDeployer) update(ctx context.Context, svc *corev1.Service) error {
	_, err := d.services.Update(ctx, svc, metav1.UpdateOptions{})
	return err
}

func (d *serviceDeployer) delete(ctx context.Context, name string) error {
	err := d.services.Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		log.Debugf("deleting service %s: service not found, doing nothing", name)
		return nil
	}
	return err
}

func (d *serviceDeployer) shouldHaveService(app models.Application) bool {
	return len(app.Spec.Config.Ports) > 0
}

func (d *serviceDeployer) tcpPortAnnotations(ports models.PortsConfig) map[string]string {
	var res []string
	for _, port := range ports {
		if port.Protocol == "tcp" {
			res = append(res, port.Name)
		}
	}
	if len(res) == 0 {
		return nil
	}
	return map[string]string{"fiaas/tcp_port_names": strings.Join(res, ",")}
}

func (d *serviceDeployer) makePorts(ports models.PortsConfig) []corev1.ServicePort {
	var res []corev1.ServicePort
	for _, p := range ports {
		res = append(res, corev1.ServicePort{
			Name:       p.Name,
			Protocol:   "TCP",
			Port:       int32(p.Port),
			TargetPort: intstr.FromInt(int(p.TargetPort)),
		})
	}
	return res
}
