package acceptance_test

import (
	"fmt"
	"net/http"
	"time"

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
		)

		BeforeEach(func() {
			startHelloThereServer()
			configMapClient = buildK8sConfigMapClient()

			configMapAnnotated := &apiv1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: configMapName,
					Annotations: map[string]string{
						"x-k8s.io/curl-me-that": "my-cool-value=http://localhost:8080",
					},
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
})

func buildK8sConfigMapClient() v1.ConfigMapInterface {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())

	clientset, err := kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())

	return clientset.CoreV1().ConfigMaps(apiv1.NamespaceDefault)
}

func startHelloThereServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "hello there")
		Expect(err).NotTo(HaveOccurred())
	})
	go http.ListenAndServe(":8080", nil)
}
