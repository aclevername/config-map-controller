package controller_test

import (
	"github.com/aclevername/config-map-controller/controller"
	"github.com/aclevername/config-map-controller/controller/fakes"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigMapController", func() {
	var (
		clientset     *fake.Clientset
		fakeProcessor *fakes.FakeProcessor
		queue         workqueue.RateLimitingInterface
		informer      cache.Controller
		configMap     *apiv1.ConfigMap
	)
	BeforeEach(func() {
		configMap = &apiv1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "configmap",
				Namespace: v1.NamespaceDefault,
			},
		}

		clientset = fake.NewSimpleClientset(configMap)

		fakeProcessor = new(fakes.FakeProcessor)

		configMapListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "configmaps", v1.NamespaceAll, fields.Everything())

		queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

		_, informer = cache.NewIndexerInformer(configMapListWatcher, &v1.ConfigMap{}, 0, cache.ResourceEventHandlerFuncs{}, cache.Indexers{})

	})

	Describe("New", func() {
		It("Builds a ConfigMapController", func() {
			configMapController := controller.NewConfigMapController(queue, informer, fakeProcessor)
			Expect(configMapController.GetQueue()).To(Equal(queue))
			Expect(configMapController.GetInformer()).To(Equal(informer))
			Expect(configMapController.GetProcessor()).To(Equal(fakeProcessor))

		})
	})

	Describe("Run", func() {
		var (
			fakeQueue    *fakes.FakeRateLimitingInterface
			fakeInformer *fakes.FakeController
			configMap    *apiv1.ConfigMap
			stopCh       chan struct{}
		)

		BeforeEach(func() {
			fakeQueue = new(fakes.FakeRateLimitingInterface)
			fakeInformer = new(fakes.FakeController)
			fakeQueue.GetReturns(nil, true)
			configMap = &apiv1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "configmap",
					Namespace: "default",
				},
			}
			stopCh = make(chan struct{})

		})

		It("processes the item until told to exit by the queue", func() {
			var callCount int
			fakeQueue.GetStub = func() (i interface{}, b bool) {
				if callCount == 0 {
					callCount++
					return configMap, false
				} else {
					return nil, true
				}
			}
			configMapController := controller.NewConfigMapController(fakeQueue, fakeInformer, fakeProcessor)
			configMapController.Run(stopCh)
			By("Starting the informer")
			Expect(fakeInformer.RunCallCount()).To(Equal(1))

			By("calling the queue")
			Expect(fakeQueue.GetCallCount()).To(Equal(2))
			Expect(fakeQueue.DoneCallCount()).To(Equal(1))
			Expect(fakeQueue.DoneArgsForCall(0)).To(Equal(configMap))

			By("processing the item")
			Expect(fakeProcessor.ProcessResourceCallCount()).To(Equal(1))
			Expect(fakeProcessor.ProcessResourceArgsForCall(0)).To(Equal(configMap))

			By("shuting down the queue")
			Expect(fakeQueue.ShutDownCallCount()).To(Equal(1))

			By("closing the channel")
			Expect(stopCh).To(BeClosed())

		})
	})
})
