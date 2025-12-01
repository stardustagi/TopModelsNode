package responses

type NodeLoginResp struct {
	NodeName  string `json:"node_name"`
	Address   string `json:"address"`
	Jwt       string `json:"jwt"`
	AccessKey string `json:"access_key"`
	Once      string `json:"once"`
	Config    any    `json:"config"`
}
