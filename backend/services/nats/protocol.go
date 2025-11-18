package message

type BaseMsg struct {
	MsgId    string      `json:"msg_id"`
	MainCode int         `json:"main_code"`
	SubCode  int         `json:"sub_code"`
	Payload  interface{} `json:"payload"`
}
