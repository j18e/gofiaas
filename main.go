package main

import (
	"fmt"
	"log"

	"github.com/j18e/gofiaas/config"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flags, err := config.ParseFlags()
	if err != nil {
		return err
	}
	fmt.Println(flags)
	return nil
}
