package main

import (
	"fmt"
	"github.com/shareit-payermax/payermax-server-sdk-go/payermax"
	"github.com/sony/gobreaker"
	"time"
)

const (
	appId              = "merchant appId"
	merchantNo         = "merchant no"
	merchantPrivateKey = "merchant private key"
	payermaxPublicKey  = "payermax public key"
)

func main() {

	sendToServer()
	//sendToServerAutoSwitchUrl()

}

func sendToServer() {
	settings := payermax.ClientSettings{
		BaseUrl:       payermax.Uat,
		ClientTimeout: 5 * time.Second,
	}
	client, err := payermax.CreateClient(appId, merchantNo,
		merchantPrivateKey, payermaxPublicKey, "", "", settings)
	if err != nil {
		fmt.Println(err)
	}

	data := `{"outTradeNo":"PAM2022010912345611217V2","subject":"hello","totalAmount":"0.99","currency":"USD","country":"HK","userId":"100000002","goodsDetails":[{"goodsId":"com.corps.gp.60","goodsName":"60鑽石","quantity":"1","price":"0.99","goodsCurrency":"USD","showUrl":"httpw://www.okgame.com"}],"language":"en","reference":"300011","frontCallbackUrl":"https://payapi.okgame.com/v2/PayerMax/result.html","notifyUrl":"https://payapi.okgame.com/v2/PayerMax/Callback.ashx"}`

	resp, err := client.Send("orderAndPay", data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func sendToServerAutoSwitchUrl() {
	var cb gobreaker.Settings
	//The circuit breaker must have the 'name' attribute set to take effect.
	cb.Name = "payermax"

	settings := payermax.ClientSettings{
		CbSettings:    cb,
		BaseUrl:       payermax.Uat,
		ClientTimeout: 5 * time.Second,
	}

	client, err := payermax.CreateClient(appId, merchantNo,
		merchantPrivateKey, payermaxPublicKey, "", "", settings)
	if err != nil {
		fmt.Println(err)
	}

	data := `{"outTradeNo":"PAM2022010912345611217V2","subject":"hello","totalAmount":"0.99","currency":"USD","country":"HK","userId":"100000002","goodsDetails":[{"goodsId":"com.corps.gp.60","goodsName":"60鑽石","quantity":"1","price":"0.99","goodsCurrency":"USD","showUrl":"httpw://www.okgame.com"}],"language":"en","reference":"300011","frontCallbackUrl":"https://payapi.okgame.com/v2/PayerMax/result.html","notifyUrl":"https://payapi.okgame.com/v2/PayerMax/Callback.ashx"}`
	resp, err := client.SendWithAutoSwitchUrl("orderAndPay", data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)

}
