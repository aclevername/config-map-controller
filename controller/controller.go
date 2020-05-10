package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapController struct {
	clientset     kubernetes.Interface
	httpClient    HTTPClient
	annotationKey string
}

func New(clientset kubernetes.Interface, httpClient HTTPClient, annotationKey string) ConfigMapController {
	return ConfigMapController{
		clientset:     clientset,
		httpClient:    httpClient,
		annotationKey: annotationKey,
	}
}

//go:generate counterfeiter -o fakes/fake_http_client.go . HTTPClient

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func (c *ConfigMapController) ProcessItem(cm *apiv1.ConfigMap) error {
	configMap := cm.DeepCopy()
	annotation, ok := configMap.Annotations[c.annotationKey]
	if !ok {
		return nil
	}

	splitAnnotation := strings.Split(annotation, "=")
	if len(splitAnnotation) != 2 {
		return nil
	}

	key := splitAnnotation[0]
	url := splitAnnotation[1]

	_, ok = configMap.Data[key]
	if ok {
		return nil
	}

	value, err := curl(url, c.httpClient)
	if err != nil {
		return err
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
	}

	return nil
}

func curl(url string, httpClient HTTPClient) (string, error) {
	req, _ := http.NewRequest("GET", url, &bytes.Buffer{})

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to curl %s, got error: %v", url, err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to curl %s, got status code: %d", url, resp.StatusCode)
	}

	if resp.Body == nil {
		return "", fmt.Errorf("empty response body from %s", url)
	}
	respValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	return string(respValue), nil
}
