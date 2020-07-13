package main

import (
	"flag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"os"

	"github.com/aclevername/config-map-controller/log"

	"github.com/aclevername/config-map-controller/reconciler"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.SetLevel(0)
	annotation := "x-k8s.io/curl-me-that"

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

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("failed to build kube client from: %s", *kubeconfig)
		os.Exit(1)
	}

	configMapListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "configmaps", v1.NamespaceAll, fields.Everything())

	watcher, err := configMapListWatcher.Watch(metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind: "configmaps",
		},
	})

	if err != nil {
		log.Error("failed to build watcher: %v", err)
		os.Exit(1)
	}

	configMapChannel := watcher.ResultChan()

	r := reconciler.New(clientset, annotation)

	for{
		var event watch.Event
		event = <- configMapChannel

		eventObj := event.Object.DeepCopyObject()
		res, ok := eventObj.(*v1.ConfigMap)
		if !ok {
			log.Error("resource not a configmap: %v", eventObj)
		}

		err := r.ReconcileResource(res)
		if err != nil {
			log.Error("resource failed to be reconiled: %v", err)
		}
	}
}
