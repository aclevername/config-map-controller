package processor

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"github.com/aclevername/config-map-controller/log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type ConfigMapProcessor struct {
	clientset     kubernetes.Interface
	httpClient    HTTPClient
	annotationKey string
}

func New(clientset kubernetes.Interface, annotationKey string) ConfigMapProcessor {
	return ConfigMapProcessor{
		clientset:     clientset,
		httpClient:    &http.Client{},
		annotationKey: annotationKey,
	}
}

//go:generate counterfeiter -o fakes/fake_http_client.go . HTTPClient

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func (c *ConfigMapProcessor) ProcessResource(cm *apiv1.ConfigMap) error {
	configMap := cm.DeepCopy()
	annotation, ok := configMap.Annotations[c.annotationKey]
	if !ok {
		log.Debug("no annotation found on %s/%s", configMap.Namespace, configMap.Name)
		return nil
	}

	splitAnnotation := strings.Split(annotation, "=")
	if len(splitAnnotation) != 2 {
		return c.addEventLogAndError(
			fmt.Sprintf("annotation value '%s' does not match expected format key=url", annotation),
			configMap,
		)
	}

	key := splitAnnotation[0]
	rawUrl := splitAnnotation[1]
	u, err := url.Parse(rawUrl)
	if err != nil {
		return c.addEventLogAndError(
			fmt.Sprintf("invalid url provided: %s", rawUrl),
			configMap,
		)
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	_, ok = configMap.Data[key]
	if ok {
		log.Debug("data field %s already set on %s/%s", key, configMap.Namespace, configMap.Name)
		return nil
	}

	value, errMsg := curl(u.String(), c.httpClient)
	if errMsg != "" {
		return c.addEventLogAndError(
			errMsg,
			configMap,
		)
	}

	if configMap.Data == nil {
		configMap.Data = map[string]string{
			key: value,
		}
	} else {
		configMap.Data[key] = value
	}

	_, err = c.clientset.CoreV1().ConfigMaps(configMap.ObjectMeta.Namespace).Update(configMap)
	if err != nil {
		return c.addEventLogAndError(
			fmt.Sprintf("failed to update configmap: %v", err),
			configMap,
		)
	}

	log.Debug("successfully updated %s/%s", configMap.Namespace, configMap.Name)

	return nil
}

func (c *ConfigMapProcessor) addEventLogAndError(errMsg string, configMap *apiv1.ConfigMap) error {
	uniqueID := uuid.New()

	var event apiv1.Event
	event.Source = apiv1.EventSource{Component: "config-map-controller"}
	event.Name = "config-map-controller" + uniqueID.String()
	event.Message = errMsg
	event.Reason = "-"
	event.Type = "error"
	event.FirstTimestamp = metav1.Now()
	event.InvolvedObject = apiv1.ObjectReference{
		Kind:      "ConfigMap",
		Namespace: configMap.Namespace,
		Name:      configMap.Name,
		UID:       configMap.UID,
	}

	_, err := c.clientset.CoreV1().Events(configMap.ObjectMeta.Namespace).Create(&event)
	if err != nil {
		log.Error("error creating event: %s", err.Error())
	}
	return errors.New(errMsg)
}

func curl(url string, httpClient HTTPClient) (string, string) {
	req, _ := http.NewRequest("GET", url, &bytes.Buffer{})

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Sprintf("failed to curl %s, got error: %v", url, err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Sprintf("failed to curl %s, got status code: %d", url, resp.StatusCode)
	}

	if resp.Body == nil {
		return "", fmt.Sprintf("empty response body from %s", url)
	}
	respValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Sprintf("failed to read response body: %v", err)
	}
	return string(respValue), ""
}
