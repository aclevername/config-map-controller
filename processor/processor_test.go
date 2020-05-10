package processor_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"k8s.io/apimachinery/pkg/types"

	"k8s.io/client-go/kubernetes"

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
				UID:       "config-map-id",
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
				annotationKey: "my-cool-value=https://example.com",
			}
		})

		When("the data field key has not already been set", func() {
			When("there is no existing data", func() {
				It("creates the data and adds the field with the correct value", func() {
					err := configMapController.ProcessResource(configMap)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeHTTPClient.DoCallCount()).To(Equal(1))
					Expect(fakeHTTPClient.DoArgsForCall(0).URL.String()).To(Equal("https://example.com"))

					updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(updatedConfigMap).To(Equal(&apiv1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name:      resourceName,
							Namespace: namespace,
							Annotations: map[string]string{
								annotationKey: "my-cool-value=https://example.com",
							},
							UID: "config-map-id",
						},
						Data: map[string]string{
							"my-cool-value": "hello-there",
						},
					}))
				})

				When("there is no schema in the URL", func() {
					BeforeEach(func() {
						configMap.Annotations = map[string]string{
							annotationKey: "my-cool-value=example.com",
						}
					})

					It("defaults to https, creates the data and adds the field with the correct value", func() {
						err := configMapController.ProcessResource(configMap)
						Expect(err).NotTo(HaveOccurred())

						Expect(fakeHTTPClient.DoCallCount()).To(Equal(1))
						Expect(fakeHTTPClient.DoArgsForCall(0).URL.String()).To(Equal("https://example.com"))

						updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
						Expect(err).NotTo(HaveOccurred())
						Expect(updatedConfigMap).To(Equal(&apiv1.ConfigMap{
							ObjectMeta: metav1.ObjectMeta{
								Name:      resourceName,
								Namespace: namespace,
								Annotations: map[string]string{
									annotationKey: "my-cool-value=example.com",
								},
								UID: "config-map-id",
							},
							Data: map[string]string{
								"my-cool-value": "hello-there",
							},
						}))
					})
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
					Expect(fakeHTTPClient.DoArgsForCall(0).URL.String()).To(Equal("https://example.com"))

					updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(updatedConfigMap).To(Equal(&apiv1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name:      resourceName,
							Namespace: namespace,
							Annotations: map[string]string{
								annotationKey: "my-cool-value=https://example.com",
							},
							UID: "config-map-id",
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

			It("returns an error", func() {
				By("returning an error")
				err := configMapController.ProcessResource(configMap)
				Expect(err).To(MatchError("annotation value 'this looks wrong' does not match expected format key=url"))

				By("not modifying the object")
				updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedConfigMap).To(Equal(configMap))

				By("adding an event describing what happened")
				event := getEvent(fakeClient, namespace)
				Expect(event.Message).To(Equal("annotation value 'this looks wrong' does not match expected format key=url"))
				assertStandardEventFieldsSet(event, resourceName, namespace)
			})
		})

		When("the URL is invalid", func() {
			When("because it contain invalid characters", func() {
				BeforeEach(func() {
					configMap.Annotations = map[string]string{
						annotationKey: "my-cool-value=!@£%",
					}
				})

				It("returns an error", func() {
					err := configMapController.ProcessResource(configMap)
					Expect(err).To(MatchError("invalid url provided: !@£%"))

					By("not modifying the object")
					updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(updatedConfigMap).To(Equal(configMap))

					By("adding an event describing what happened")
					event := getEvent(fakeClient, namespace)
					Expect(event.Message).To(Equal("invalid url provided: !@£%"))
					assertStandardEventFieldsSet(event, resourceName, namespace)
				})
			})

			When("because it contain spaces", func() {
				BeforeEach(func() {
					configMap.Annotations = map[string]string{
						annotationKey: "my-cool-value=hello world",
					}
				})

				It("returns an error", func() {
					err := configMapController.ProcessResource(configMap)
					Expect(err).To(MatchError(ContainSubstring("failed to create http request, err: ")))

					By("not modifying the object")
					updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(updatedConfigMap).To(Equal(configMap))

					By("adding an event describing what happened")
					event := getEvent(fakeClient, namespace)
					Expect(event.Message).To(ContainSubstring("failed to create http request, err: "))
					assertStandardEventFieldsSet(event, resourceName, namespace)
				})
			})

		})

		When("the http request fails", func() {
			BeforeEach(func() {
				fakeHTTPClient.DoReturns(nil, errors.New("failed"))
			})

			It("returns an error", func() {
				err := configMapController.ProcessResource(configMap)
				Expect(err).To(MatchError("failed to curl https://example.com, got error: failed"))

				By("not modifying the object")
				updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedConfigMap).To(Equal(configMap))

				By("adding an event describing what happened")
				event := getEvent(fakeClient, namespace)
				Expect(event.Message).To(Equal("failed to curl https://example.com, got error: failed"))
				assertStandardEventFieldsSet(event, resourceName, namespace)
			})
		})

		When("the http request does not return 200", func() {
			BeforeEach(func() {
				fakeHTTPClient.DoReturns(&http.Response{Body: ioutil.NopCloser(strings.NewReader("hello-there")), StatusCode: http.StatusInternalServerError}, nil)
			})

			It("returns an error", func() {
				err := configMapController.ProcessResource(configMap)
				Expect(err).To(MatchError("failed to curl https://example.com, got status code: 500"))

				By("not modifying the object")
				updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedConfigMap).To(Equal(configMap))

				By("adding an event describing what happened")
				event := getEvent(fakeClient, namespace)
				Expect(event.Message).To(Equal("failed to curl https://example.com, got status code: 500"))
				assertStandardEventFieldsSet(event, resourceName, namespace)
			})
		})

		When("reading the response body is nil", func() {
			BeforeEach(func() {
				fakeHTTPClient.DoReturns(&http.Response{Body: nil, StatusCode: http.StatusOK}, nil)
			})

			It("returns an error", func() {
				err := configMapController.ProcessResource(configMap)
				Expect(err).To(MatchError("empty response body from https://example.com"))

				By("not modifying the object")
				updatedConfigMap, err := fakeClient.CoreV1().ConfigMaps(namespace).Get(resourceName, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedConfigMap).To(Equal(configMap))

				By("adding an event describing what happened")
				event := getEvent(fakeClient, namespace)
				Expect(event.Message).To(Equal("empty response body from https://example.com"))
				assertStandardEventFieldsSet(event, resourceName, namespace)
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

				By("adding an event describing what happened")
				event := getEvent(fakeClient, namespace)
				Expect(event.Message).To(Equal("failed to read response body: failed"))
				assertStandardEventFieldsSet(event, resourceName, namespace)
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

				By("adding an event describing what happened")
				event := getEvent(fakeClient, namespace)
				Expect(event.Message).To(ContainSubstring("failed to update configmap: "))
				assertStandardEventFieldsSet(event, resourceName, namespace)
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

func getEvent(fakeClient kubernetes.Interface, namespace string) *apiv1.Event {
	eventList, err := fakeClient.CoreV1().Events(namespace).List(metav1.ListOptions{})
	Expect(err).NotTo(HaveOccurred())
	Expect(eventList.Items).To(HaveLen(1))
	event, err := fakeClient.CoreV1().Events(namespace).Get(eventList.Items[0].Name, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	return event
}

func assertStandardEventFieldsSet(event *apiv1.Event, resourceName, namespace string) {
	Expect(event.Name).To(MatchRegexp("config-map-controller-*"))
	Expect(event.Reason).To(Equal("-"))
	Expect(event.Type).To(Equal("error"))
	Expect(event.FirstTimestamp.String()).ToNot(BeEmpty())
	Expect(event.Source.Component).To(Equal("config-map-controller"))
	Expect(event.InvolvedObject.Kind).To(Equal("ConfigMap"))
	Expect(event.InvolvedObject.Namespace).To(Equal(namespace))
	Expect(event.InvolvedObject.Name).To(Equal(resourceName))
	Expect(event.InvolvedObject.UID).To(Equal(types.UID("config-map-id")))

}
