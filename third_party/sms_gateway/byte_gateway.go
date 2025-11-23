package sms_gateway

import (
	"fmt"

	"github.com/volcengine/volc-sdk-golang/service/sms"
)

func SendByteSMS(phoneNumbers string, message string) (*sms.SmsResponse, error) {
	ak := "AKLTYTI1YmI4MDQ3MjIzNGJkYWFkYjJlOGZhNWE3ZWJlM2I"
	sk := "TWpNM00yWXlORFJpT0Rkak5HVTJObUV4TmpWallqbGtOak0wWkRVelpEUQ=="
	sms.DefaultInstance.Client.SetAccessKey(ak)
	sms.DefaultInstance.Client.SetSecretKey(sk)
	req := &sms.SmsRequest{
		SmsAccount:    "86967286",
		Sign:          "雅逐智能",
		TemplateID:    "SPT_09a29a26",
		TemplateParam: message,
		PhoneNumbers:  phoneNumbers,
	}
	result, statusCode, err := sms.DefaultInstance.Send(req)
	fmt.Println("Status Code:", statusCode)
	if err != nil {
		return nil, err
	}
	return result, nil
}
