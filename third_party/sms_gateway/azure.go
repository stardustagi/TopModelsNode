package sms_gateway

import (
	"fmt"
	"os"
	"time"

	"resty.dev/v3"
)

// 获取 Azure AD 访问令牌（客户端凭证模式）
func getAzureAccessToken(tenantID, clientID, clientSecret string) (string, error) {
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	client := resty.New().
		SetTimeout(15 * time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second)

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     clientID,
			"client_secret": clientSecret,
			"scope":         "https://communication.azure.com/.default",
		}).
		SetResult(&tokenResp).
		Post(tokenURL)

	if err != nil {
		return "", fmt.Errorf("request token error: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("token request failed: %s - %s", resp.Status(), resp.String())
	}

	return tokenResp.AccessToken, nil
}

// 调用 Azure Communication Services SMS API
func sendAzureSMS(endpoint, bearer, from string, to []string, msg, tag string, report bool) error {
	type smsBody struct {
		From                 string   `json:"from"`
		To                   []string `json:"to"`
		Message              string   `json:"message"`
		EnableDeliveryReport bool     `json:"enableDeliveryReport"`
		Tag                  string   `json:"tag,omitempty"`
	}

	client := resty.New().
		SetBaseURL(endpoint).
		SetTimeout(20*time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1*time.Second).
		SetRetryMaxWaitTime(8*time.Second).
		SetHeader("Accept-Encoding", "gzip")

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthScheme("Bearer").
		SetAuthToken(bearer).
		SetBody(smsBody{
			From:                 from,
			To:                   to,
			Message:              msg,
			EnableDeliveryReport: report,
			Tag:                  tag,
		}).
		Post("/sms?api-version=2021-03-07")

	if err != nil {
		return fmt.Errorf("send sms http error: %w", err)
	}
	fmt.Println("Status:", resp.Status())
	fmt.Println("Body  :", resp.String())

	if resp.IsError() {
		return fmt.Errorf("send sms failed: %s", resp.Status())
	}
	return nil
}

func SendAzureSMS(toPhone, code string) error {
	tenantID := os.Getenv("TENANT_ID")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	resourceName := os.Getenv("ACS_RESOURCE_NAME") // 例如：my-acs-prod
	fromPhone := os.Getenv("FROM_PHONE")           // 例如：+1xxxxxxxxxx（必须是已购买/配置的发信号码）
	// toPhone := os.Getenv("TO_PHONE")               // 目标号码

	if tenantID == "" || clientID == "" || clientSecret == "" || resourceName == "" || fromPhone == "" || toPhone == "" {
		return fmt.Errorf("please set TENANT_ID, CLIENT_ID, CLIENT_SECRET, ACS_RESOURCE_NAME, FROM_PHONE, TO_PHONE")
	}

	endpoint := fmt.Sprintf("https://%s.communication.azure.com", resourceName)

	// 1) 获取 Azure AD 访问令牌
	token, err := getAzureAccessToken(tenantID, clientID, clientSecret)
	if err != nil {
		return err
	}

	sms := fmt.Sprintf("Hello from Go (resty) via Azure ACS SMS! Your verification code is: %s", code)
	// 2) 发送短信
	return sendAzureSMS(endpoint, token, fromPhone, []string{toPhone}, sms, "top maas", true)
}
