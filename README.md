# PayerMAX Server sdk

### install
```go
go get github.com/payermax/go-sdk
```

### use demo
```go

import (
	"fmt"
    "github.com/shareit-payermax/payermax-server-sdk-go/payermax"
)

const (
	appId = "your appId"
	merchantNo = "your merchantNo"
	merchantPrivateKey = "your privateKey"
	payermaxPublicKey = "payermax public key"
) 

func main() {
	client, err := payermax.CreateClient(appId, merchantNo, 
		merchantPrivateKey, payermaxPublicKey, payermax.Uat)
	if err != nil {
		fmt.Println(err)
	}

	data := `{"outTradeNo":"PAM20220109123456111617V1","subject":"hello","totalAmount":"0.99","currency":"USD","country":"HK","userId":"100000002","goodsDetails":[{"goodsId":"com.corps.gp.60","goodsName":"60鑽石","quantity":"1","price":"0.99","goodsCurrency":"USD","showUrl":"httpw://www.okgame.com"}],"language":"en","reference":"300011","frontCallbackUrl":"https://payapi.okgame.com/v2/PayerMax/result.html","notifyUrl":"https://payapi.okgame.com/v2/PayerMax/Callback.ashx"}`
	
	resp, err := client.Send("orderAndPay", data)
	//auto switch backup url
	//resp, err := client.SendWithAutoSwitchUrl("orderAndPay", data)
	
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
```
### verify notify sign

```go
ref VerifySign function
```
