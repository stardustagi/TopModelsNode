package responses

type NodeLoginResp struct {
	NodeId    int64  `json:"node_id"`
	NodeName  string `json:"node_name"`
	Address   string `json:"address"`
	Jwt       string `json:"jwt"`
	AccessKey string `json:"access_key"`
	Once      string `json:"once"`
	Config    any    `json:"config"`
}
