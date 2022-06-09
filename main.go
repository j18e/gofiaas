package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	fiaasclientset "github.com/fiaas/fiaas-go-client/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/j18e/gofiaas/config"
	"github.com/j18e/gofiaas/log"
	v3 "github.com/j18e/gofiaas/spec/v3"
)

func main() {
	if err := run(); err != nil {
		log.Logger.Fatal(err)
	}
}

func run() error {
	flags, err := config.ParseFlags()
	if err != nil {
		return err
	}
	fmt.Println(flags)

	fc, err := fiaasClientset(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err != nil {
		return err
	}
	apps := fc.FiaasV1().Applications("default")
	app, err := apps.Get(context.TODO(), "consent-braze-integration", metav1.GetOptions{})
	if err != nil {
		return err
	}

	bs, err := json.Marshal(app.Spec.Config)
	if err != nil {
		return err
	}
	fmt.Println(string(bs))

	factory, err := v3.NewFactory()
	if err != nil {
		return err
	}
	spec, err := factory.Transform(*app)
	if err != nil {
		return err
	}
	fmt.Println(spec)
	return nil
}

func fiaasClientset(kubeconfigPath string) (fiaasclientset.Interface, error) {
	if kubeconfigPath == "" {
		log.Logger.Info("using in-cluster configuration")
		k8s, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		return fiaasclientset.NewForConfig(k8s)
	}
	log.Logger.Infof("using configuration from '%s'", kubeconfigPath)
	k8s, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return fiaasclientset.NewForConfig(k8s)
}
