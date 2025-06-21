package transport

type Meta struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}
