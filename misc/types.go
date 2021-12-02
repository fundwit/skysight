package misc

type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PagedBody struct {
	Items interface{} `json:"items"`
	Total uint64      `json:"total"`
}
