package main

import (
	"flag"
	"os"

	"github.com/aclevername/config-map-controller/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.SetLevel(0)

	kubeconfig := flag.String("kubeconfig", "", "path to kubeconfig")
	flag.Parse()

	if *kubeconfig == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Error("failed to build client config from: %s", *kubeconfig)
		os.Exit(1)
	}

	_, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("failed to build kube client from: %s", *kubeconfig)
		os.Exit(1)
	}
}
