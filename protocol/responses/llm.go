package responses

type NodeRegisterResp struct {
	NodeId    string `json:"node_id"`
	Jwt       string `json:"jwt"`
	AccessKey string `json:"access_key"`
	Once      string `json:"once"`
}
