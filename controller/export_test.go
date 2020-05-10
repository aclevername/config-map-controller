package controller

func (c *ConfigMapController) SetHTTPClient(client HTTPClient) {
	c.httpClient = client
}
