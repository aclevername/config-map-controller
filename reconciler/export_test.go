package reconciler

func (c *ConfigMapReconciler) SetHTTPClient(client HTTPClient) {
	c.httpClient = client
}
