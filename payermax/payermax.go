package payermax

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/sony/gobreaker"
	"io"
	"net/http"
	"strings"
	"time"
)

var cb *gobreaker.CircuitBreaker

func init() {
	var st gobreaker.Settings
	st.Name = "payermax"
	//半开状态连续请求成功数量大于这个值则把熔断器关闭
	st.MaxRequests = 5
	//熔断的条件
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		//请求数量大于3且失败率超过50就进行熔断
		return counts.Requests >= 3 && failureRatio >= 0.5
	}
	//统计请求数量和比例的周期，这里表示统计1分钟内的请求数量和比例
	st.Interval = time.Minute
	//熔断后多长时间开始把开关设置成半开状态，好检测主域名是否正常
	st.Timeout = 30 * time.Second
	//st.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
	//}
	//判断调用是否成功，可以精细定义各种异常信息
	st.IsSuccessful = func(err error) bool {
		return err == nil
	}

	cb = gobreaker.NewCircuitBreaker(st)
}

type Client struct {
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

func (this *Client) SendWithUrl(apiName, data string, baseUrl string) (resp string, resErr error) {
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
		reqMap["spMerchantNo"] = this.spMerchantNo
	}

	resultBytes, err := json.Marshal(reqMap)
	if err != nil {
		return "", err
	}

	resultJson := string(resultBytes)
	req, err := http.NewRequest("POST", baseUrl+apiName, strings.NewReader(resultJson))

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	rsaSign, resErr := GetRsaSign(resultJson, this.merchantPrivateKey)
	req.Header.Set("sign", rsaSign)
	req.Header.Set("sdk-ver", "go-1.0")
	if this.spMerchantNo != "" {
		req.Header.Set("merchantAuthToken", this.merchantAuthToken)
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

func (this *Client) Send(apiName, data string) (resp string, resErr error) {
	return this.SendWithUrl(apiName, data, this.baseUrl)
}

func (this *Client) SendWithAutoSwitchUrl(apiName, data string) (resp string, resErr error) {
	if this.baseUrl == Uat {
		return this.Send(apiName, data)
	}

	respBody, err := cb.Execute(func() (interface{}, error) {
		body, err := this.Send(apiName, data)
		if err != nil {
			return nil, err
		}

		return body, err
	})

	//熔断异常
	if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
		//降级到备用域名
		return this.SendWithUrl(apiName, data, ProdBackUp)
	}

	if err != nil {
		return "", err
	}

	return respBody.(string), nil
}
