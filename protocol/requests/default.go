package requests

type DefaultRequest struct {
}

type DefaultWsRequest struct {
}

// swagger:request PageReq
type PageReq struct {
	Skip  int    `json:"skip"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
}

type BasePageRequest struct {
	Data interface{} `json:"data,omitempty"`
	Page PageReq     `json:"page"`
}
