package processor

func (c *ConfigMapProcessor) SetHTTPClient(client HTTPClient) {
	c.httpClient = client
}
