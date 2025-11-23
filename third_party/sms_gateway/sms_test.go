package sms_gateway

import "testing"

func TestSMSTest(t *testing.T) {
	response, err := SendSMS("18964479925", "您的验证码是123456。如非本人操作，请忽略本短信。")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Logf("Response: %s", response)
}

func TestByteSMSTest(t *testing.T) {
	response, err := SendByteSMS("18964479925", "{\"code\": \"202510\"}")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Logf("Response: %v", response)
}
