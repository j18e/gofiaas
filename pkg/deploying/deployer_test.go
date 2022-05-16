package deploying

import (
	"github.com/j18e/gofiaas/pkg/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

func newTestApp(name, namespace string) models.Application {
	return models.Application{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, UID: uuid.NewUUID()},
		Spec: models.ApplicationSpec{
			Application: name,
			Config: models.ApplicationConfig{
				Ports: models.PortsConfig{
					{Name: "http", Protocol: "TCP", Port: 80},
				},
			},
		},
	}
}
