package payermax

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	ErrBadResponse = errors.New("alipay: bad response")
)

type Client struct {
	mu                 sync.Mutex
	appId              string
	merchantNo         string
	spMerchantNo       string
	merchantAuthToken  string
	baseUrl            string
	httpClient         *http.Client
	merchantPrivateKey *rsa.PrivateKey // 商户私钥
	payermaxPublicKey  *rsa.PublicKey  // payermax公钥钥

}

func CreateClient(appId, merchantNo, merchantPrivateKey, payermaxPublicKey, spMerchantNo, merchantAuthToken, baseUrl string) (client *Client, err error) {
	priKey, err := DecodePrivateKey(merchantPrivateKey)
	if err != nil {
		return nil, err
	}

	pubKey, err := DecodePublicKey(payermaxPublicKey)
	if err != nil {
		return nil, err
	}

	client = &Client{}
	client.appId = appId
	client.merchantNo = merchantNo
	client.merchantPrivateKey = priKey
	client.payermaxPublicKey = pubKey
	client.baseUrl = baseUrl
	client.spMerchantNo = spMerchantNo
	client.merchantAuthToken = merchantAuthToken

	client.httpClient = &http.Client{
		Timeout: 15 * time.Second,
	}

	return client, nil
}

func (this *Client) Send(apiName, data string) (resp string, resErr error) {

	var reqBody = `{"keyVersion":"1","merchantNo":"","requestTime":"","version":"1.1","appId":"","data":{}}`

	var reqMap map[string]interface{}
	err := json.Unmarshal([]byte(reqBody), &reqMap)
	if err != nil {
		return "", err
	}

	var dataMap map[string]interface{}
	err = json.Unmarshal([]byte(data), &dataMap)
	if err != nil {
		return "", err
	}

	reqMap["merchantNo"] = this.merchantNo
	reqMap["appId"] = this.appId
	reqMap["requestTime"] = time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")
	reqMap["data"] = dataMap
	if this.spMerchantNo != "" {
		reqMap["merchantAuthToken"] = this.merchantAuthToken
	}

	resultBytes, err := json.Marshal(reqMap)
	if err != nil {
		return "", err
	}

	resultJson := string(resultBytes)
	req, err := http.NewRequest("POST", this.baseUrl+apiName, strings.NewReader(resultJson))

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	rsaSign, resErr := GetRsaSign(resultJson, this.merchantPrivateKey)
	req.Header.Set("sign", rsaSign)
	req.Header.Set("sdk-ver", "go-1.0")
	if this.spMerchantNo != "" {
		req.Header.Set("sign", rsaSign)
	}

	response, err := this.httpClient.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return "", err
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	responseSign := response.Header.Get("sign")
	respBody := string(responseBytes)

	var respMap map[string]interface{}
	err = json.Unmarshal([]byte(respBody), &respMap)
	if err != nil {
		return "", err
	}

	if respMap["code"] != "APPLY_SUCCESS" {
		return respBody, nil
	}

	if err = VerifySign(respBody, responseSign, this.payermaxPublicKey); err != nil {
		return "", err
	}

	return respBody, nil
}
