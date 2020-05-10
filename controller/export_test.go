package controller

import (
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func (c *ConfigMapController) GetInformer() cache.Controller {
	return c.informer
}

func (c *ConfigMapController) GetQueue() workqueue.RateLimitingInterface {
	return c.queue
}

func (c *ConfigMapController) GetProcessor() Processor {
	return c.processor
}
