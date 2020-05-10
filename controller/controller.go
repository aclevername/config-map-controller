package controller

import (
	"sync"

	"github.com/aclevername/config-map-controller/log"

	apiv1 "k8s.io/api/core/v1"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type ConfigMapController struct {
	queue     workqueue.RateLimitingInterface
	informer  cache.Controller
	processor Processor
}

//go:generate counterfeiter -o fakes/fake_queue.go k8s.io/client-go/util/workqueue.RateLimitingInterface
//go:generate counterfeiter -o fakes/fake_informer.go k8s.io/client-go/tools/cache.Controller

//go:generate counterfeiter -o fakes/fake_processor.go . Processor
type Processor interface {
	ProcessResource(cm *apiv1.ConfigMap) error
}

func NewConfigMapController(queue workqueue.RateLimitingInterface, informer cache.Controller, processor Processor) *ConfigMapController {
	return &ConfigMapController{
		informer:  informer,
		queue:     queue,
		processor: processor,
	}
}

func (c *ConfigMapController) Run(stopCh chan struct{}) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		c.informer.Run(stopCh)
		wg.Done()
	}()

	defer c.queue.ShutDown()
	defer close(stopCh)

	for c.run() {
	}

	log.Debug("controller shutting down")
	wg.Wait()
}

func (c *ConfigMapController) run() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	val, ok := key.(*apiv1.ConfigMap)
	if !ok {
		return true
	}

	err := c.processor.ProcessResource(val)
	if err != nil {
		log.Error("error processing  configmap %s/%s, error: %v", key.(*apiv1.ConfigMap).Namespace, key.(*apiv1.ConfigMap).Name, err)
		return true
	}
	return true
}
