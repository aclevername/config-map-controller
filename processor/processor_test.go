package processor_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aclevername/config-map-controller/processor"

	httpFakes "github.com/aclevername/config-map-controller/processor/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

//go:generate counterfeiter -o fakes/fake_read_closer.go io.ReadCloser

var _ = Describe("ProcessResource", func() {
	var (
		configMapController processor.ConfigMapProcessor
		fakeClient          *fake.Clientset
		fakeHTTPClient      *httpFakes.FakeHTTPClient
		configMap           *apiv1.ConfigMap
		namespace           = "my-namespace"
		resourceName        = "my-resource"
		annotationKey       = "my-annotation"
	)

	BeforeEach(func() {
		configMap = &apiv1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: namespace,
			},
		}

		fakeHTTPClient = new(httpFakes.FakeHTTPClient)
		fakeHTTPClient.DoReturns(&http.Response{Body: ioutil.NopCloser(strings.NewReader("hello-there")), StatusCode: http.StatusOK}, nil)

	})

	JustBeforeEach(func() {
		fakeClient = fake.NewSimpleClientset(configMap)
		configMapController = processor.New(fakeClient, annotationKey)
		configMapController.SetHTTPClient(fakeHTTPClient)
	})

	When("the annotation exists", func() {
		BeforeEach(func() {
			configMap.Annotations = map[string]string{
				annotationKey: "my-cool-value=example.com",
			}
		})

		When("the data field key has not already been set", func() {
			When("there is no existing data", func() {
				It("creates the data and adds the field with the correct value", func() {
					err := configMapController.ProcessResource(configMap)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeHTTPClient.DoCallCount()).To(Equal(1))
					Expect(fakeHTTPClient.DoArgsForCall(0).URL.String()).To(Equal("example.com"))

					updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(updatedConfigMap).To(Equal(&apiv1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name:      resourceName,
							Namespace: namespace,
							Annotations: map[string]string{
								annotationKey: "my-cool-value=example.com",
							},
						},
						Data: map[string]string{
							"my-cool-value": "hello-there",
						},
					}))
				})

			})

			When("there is existing data", func() {
				BeforeEach(func() {
					configMap.Data = map[string]string{
						"foo": "bar",
					}
				})

				It("adds the data field with the correct value to the existing data", func() {
					err := configMapController.ProcessResource(configMap)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeHTTPClient.DoCallCount()).To(Equal(1))
					Expect(fakeHTTPClient.DoArgsForCall(0).URL.String()).To(Equal("example.com"))

					updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(updatedConfigMap).To(Equal(&apiv1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name:      resourceName,
							Namespace: namespace,
							Annotations: map[string]string{
								annotationKey: "my-cool-value=example.com",
							},
						},
						Data: map[string]string{
							"my-cool-value": "hello-there",
							"foo":           "bar",
						},
					}))
				})
			})
		})

		When("the annotation value isn't a key=url format", func() {
			BeforeEach(func() {
				configMap.Annotations = map[string]string{
					annotationKey: "this looks wrong",
				}
			})

			It("does not error", func() {
				By("returning nill")
				err := configMapController.ProcessResource(configMap)
				Expect(err).NotTo(HaveOccurred())

				By("not modifying the object")
				updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedConfigMap).To(Equal(configMap))
			})
		})

		When("the http request fails", func() {
			BeforeEach(func() {
				fakeHTTPClient.DoReturns(nil, errors.New("failed"))
			})

			It("returns an error", func() {
				err := configMapController.ProcessResource(configMap)
				Expect(err).To(MatchError("failed to curl example.com, got error: failed"))

				By("not modifying the object")
				updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedConfigMap).To(Equal(configMap))
			})
		})

		When("the http request does not return 200", func() {
			BeforeEach(func() {
				fakeHTTPClient.DoReturns(&http.Response{Body: ioutil.NopCloser(strings.NewReader("hello-there")), StatusCode: http.StatusInternalServerError}, nil)
			})

			It("returns an error", func() {
				err := configMapController.ProcessResource(configMap)
				Expect(err).To(MatchError("failed to curl example.com, got status code: 500"))

				By("not modifying the object")
				updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedConfigMap).To(Equal(configMap))
			})
		})

		When("reading the response body is nil", func() {
			BeforeEach(func() {
				fakeHTTPClient.DoReturns(&http.Response{Body: nil, StatusCode: http.StatusOK}, nil)
			})

			It("returns an error", func() {
				err := configMapController.ProcessResource(configMap)
				Expect(err).To(MatchError("empty response body from example.com"))

				By("not modifying the object")
				updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedConfigMap).To(Equal(configMap))
			})
		})

		When("reading the response body fails", func() {
			BeforeEach(func() {
				fakeBody := new(httpFakes.FakeReadCloser)
				fakeBody.ReadReturns(0, errors.New("failed"))
				fakeHTTPClient.DoReturns(&http.Response{Body: fakeBody, StatusCode: http.StatusOK}, nil)
			})

			It("returns an error", func() {
				err := configMapController.ProcessResource(configMap)
				Expect(err).To(MatchError(ContainSubstring("failed to read response body: failed")))

				By("not modifying the object")
				updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedConfigMap).To(Equal(configMap))
			})
		})

		When("the data key is already filled in", func() {
			BeforeEach(func() {
				configMap.Data = map[string]string{
					"my-cool-value": "already set",
				}
			})

			It("does not error", func() {
				By("returning nil")
				err := configMapController.ProcessResource(configMap)
				Expect(err).NotTo(HaveOccurred())

				By("not modifying the object")
				updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedConfigMap).To(Equal(configMap))
			})
		})

		When("the update fails", func() {
			BeforeEach(func() {
				fakeHTTPClient.DoStub = func(arg1 *http.Request) (response *http.Response, e error) {
					//Delete the resource before the update can occur, causing the update to fail
					err := fakeClient.CoreV1().ConfigMaps(namespace).Delete(resourceName, nil)
					Expect(err).NotTo(HaveOccurred())
					return &http.Response{Body: ioutil.NopCloser(strings.NewReader("hello-there")), StatusCode: http.StatusOK}, nil
				}
			})
			It("returns an error", func() {
				err := configMapController.ProcessResource(configMap)
				Expect(err).To(MatchError(ContainSubstring("failed to update configmap: ")))
			})
		})
	})

	When("the annotation does not exist", func() {
		It("does not error", func() {
			By("returning nill")
			err := configMapController.ProcessResource(configMap)
			Expect(err).NotTo(HaveOccurred())

			By("not modifying the object")
			updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedConfigMap).To(Equal(configMap))
		})
	})

})
