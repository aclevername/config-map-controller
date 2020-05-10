package main

import (
	"flag"
	"os"

	"github.com/aclevername/config-map-controller/log"

	"github.com/aclevername/config-map-controller/reconciler"

	"github.com/aclevername/config-map-controller/controller"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"k8s.io/client-go/kubernetes"
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

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	_, informer := cache.NewIndexerInformer(configMapListWatcher, &v1.ConfigMap{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			queue.Add(obj)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			queue.Add(new)
		}}, cache.Indexers{})

	r := reconciler.New(clientset, annotation)
	configMapController := controller.NewConfigMapController(queue, informer, &r)

	stopCh := make(chan struct{})

	log.Debug("starting controller to watch for %s annotation", annotation)
	configMapController.Run(stopCh)
}
