package sms_gateway

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/stardustagi/TopLib/libs/uuid"
)

// SendSMS 调用短信 API
func SendSMS(userNumber, messageContent string) (string, error) {
	apiURL := "http://47.116.77.168:8513/sms/Api/ReturnJson/Send.do"

	// 构造表单参数
	data := url.Values{}
	data.Set("SpCode", "101696")
	data.Set("LoginName", "HYCSH01")
	data.Set("Password", "8Bok3F7i44J!D")
	data.Set("UserNumber", userNumber)
	data.Set("MessageContent", fmt.Sprintf("【京东】%s", messageContent))
	data.Set("SerialNumber", uuid.GenNumberString(10))
	data.Set("ScheduleTime", "")
	//data.Set("ExtendAccessNum", "")
	//data.Set("f", "1") // 返回格式为 JSON

	// 创建请求
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
