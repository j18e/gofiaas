package main

import (
	"context"
	"log"
	"os"

	"github.com/j18e/gofiaas/deploying"
	"github.com/j18e/gofiaas/models"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	if err != nil {
		return err
	}
	k8s, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	deployerCfg := deploying.Config{
		ServiceType: corev1.ServiceTypeNodePort,
	}
	deployer := deploying.NewDeployer(k8s, "jamietest", deployerCfg)

	app := models.Application{}
	app.Name = "something"

	if err := deployer.Deploy(context.TODO(), app); err != nil {
		return err
	}
	return nil
}
