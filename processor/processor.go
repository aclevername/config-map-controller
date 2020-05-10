package processor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aclevername/config-map-controller/log"

	apiv1 "k8s.io/api/core/v1"
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
		return fmt.Errorf("annotation value '%s' does not match expected format key=url", annotation)
	}

	key := splitAnnotation[0]
	rawUrl := splitAnnotation[1]
	u, err := url.Parse(rawUrl)
	if err != nil {
		return fmt.Errorf("invalid url provided: %s", rawUrl)
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	_, ok = configMap.Data[key]
	if ok {
		log.Debug("data field %s already set on %s/%s", key, configMap.Namespace, configMap.Name)
		return nil
	}

	value, err := curl(u.String(), c.httpClient)
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
		return fmt.Errorf("failed to update configmap: %v", err)
	}

	log.Debug("successfully updated %s/%s", configMap.Namespace, configMap.Name)

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
