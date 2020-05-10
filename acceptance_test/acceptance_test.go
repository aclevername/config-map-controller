package acceptance_test

import (
	"fmt"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/uuid"

	v1 "k8s.io/client-go/kubernetes/typed/core/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("Acceptance", func() {
	Context("When a config resource is annotated with 'x-k8s.io/curl-me-that'", func() {
		var (
			configMapClient v1.ConfigMapInterface
			configMapName   = "acceptance-test-with-valid-annotation"
			namespace       string
			clientset       kubernetes.Interface
		)

		BeforeEach(func() {
			u := uuid.NewUUID()
			namespace = fmt.Sprintf("acceptance-test-%s", u)

			clientset = buildK8sClient()
			_, err := clientset.CoreV1().Namespaces().Create(&apiv1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := clientset.CoreV1().Namespaces().Delete(namespace, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		When("the annotation is valid", func() {
			BeforeEach(func() {
				startHelloThereServer()
				configMapClient = clientset.CoreV1().ConfigMaps(namespace)

				configMapAnnotated := &apiv1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: configMapName,
						Annotations: map[string]string{
							"x-k8s.io/curl-me-that": "my-cool-value=http://localhost:8080",
						},
						Namespace: namespace,
					},
				}

				_, err := configMapClient.Create(configMapAnnotated)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				Expect(configMapClient.Delete(configMapName, nil)).NotTo(HaveOccurred())
			})

			It("appends data with the key and the contents of the curling the value", func() {
				Eventually(func() string {
					configMap, err := configMapClient.Get(configMapName, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					val, _ := configMap.Data["my-cool-value"]
					return val
				}, 5*time.Second, 1*time.Second).Should(Equal("hello there"))
			})
		})

		When("the annotation is invalid", func() {
			BeforeEach(func() {
				configMapClient = clientset.CoreV1().ConfigMaps(namespace)

				configMapAnnotated := &apiv1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: configMapName,
						Annotations: map[string]string{
							"x-k8s.io/curl-me-that": "not correct",
						},
						Namespace: namespace,
					},
				}

				_, err := configMapClient.Create(configMapAnnotated)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				Expect(configMapClient.Delete(configMapName, nil)).NotTo(HaveOccurred())
			})

			It("adds an event log describing why", func() {
				var eventList *apiv1.EventList
				Eventually(func() int {
					var err error
					eventList, err = clientset.CoreV1().Events(namespace).List(metav1.ListOptions{})
					Expect(err).NotTo(HaveOccurred())
					return len(eventList.Items)
				}, 5*time.Second, 1*time.Second).Should(Equal(1))

				event, err := clientset.CoreV1().Events(namespace).Get(eventList.Items[0].Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(event.Message).To(Equal("annotation value 'not correct' does not match expected format key=url"))
				Expect(event.InvolvedObject.Name).To(Equal(configMapName))
			})
		})
	})
})

func buildK8sClient() kubernetes.Interface {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())

	clientset, err := kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())

	return clientset
}

func startHelloThereServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "hello there")
		Expect(err).NotTo(HaveOccurred())
	})
	go http.ListenAndServe(":8080", nil)
}
