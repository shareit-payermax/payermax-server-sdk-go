package main

import (
	"fmt"
	"github.com/PayermaxZhanglong/payermax-server-sdk-go/payermax"
)

const (
	appId              = "55c5bd4b95c5474aa9bc14a72468274a"
	merchantNo         = "P01010114284115"
	merchantPrivateKey = "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCxZ/4z+6jtGZUELTHv25GNNVDfLLj+N6iRpiPHGGc7Xg43ZzuiAhVOqUHOb/WuRd6FG60P4copoQ1JV1dkGdW3GSHa/k6tDn8F6Wva3R3lykCjFxWSx7NLVrsPWAHOh+etkLBqLyZPEtyYkQ1jAK4USYWlABakGIHbpJ54dZj2WDhrVNkydOemuF0Eh5xFbcB9jrezGn79UYhwYxQlXudhlncs93w1zdGp+SP1Fg03mhP4RR1PTylY1hi+vMalnvbsJvQBkoVTIyTKbF3ZTK3dJf0KvQAGQlobIOVvxfeTjqqp4ZU8a7bf+3Z4YTgSMKle8CkKDaDurre6xP93/GoRAgMBAAECggEAWLDgySQ7Y62ybhAaMwvWTH4nHCvDoHQ69fXVVSfCFFAFU3VsvqXD+TttBsO6U5HimTEo6BO9eepmNhBQRF8WFR+faVhSBPqGTnJP1W9ODx96CQ71Xhgwbq3Bfv5EiPgwvvji+XLw9/86AZwi0Sf533KSOdc9enQ2F5TJFPfVrPHATPOtv4unTAdhIYk4/RrNu8k2gqGTKRtXymRhrY279LxxovYNCthVzlABVMlyupW5F8IIEmeXXAruXXrO8HO0/AxzHu9hWyFRDKTx0+DGfiQQ5lQmylUdCRY034cTanC2KqBR7DelE52Xu6Y8G68ZbWOjLWWmeK9yuUai6xb0cQKBgQDkP8QK2mA0dH7DML0bFaIiay9RkOuwRhPOivIZ1Tmg0Eory5hwOuSCvbsxk6jtsMs8PmVY166BQh7I9jIFFwv/2vsrkMtPqjMl5w//XVuYxViz4wVq8VwWbEAcE3WfylmP1WQ/FGagcqddPTRaxBW13pEFo0QxKHr+mJ8VAmIf5wKBgQDG+b9KZYJQ78osBodiCirbXxvnbokF85qbAQF+sGl1rdRi40u7XpmrIdP/yELHK/sZwPm9noqL3pD5xQhMpGZazMbG8fY+jDHAlX3yluy4NxElMc3LiaIa0TNgBY270U0aXcyniXz/lhgIqCxVxghIsg6PwBqSkp5FpRdJnR/HRwKBgQDQsYT8L2MbUxC8Q4oEg6k1My+WspztFYXyqZRnDlCcuxW2KXd91jstV3EbCVnByo5ozNw5eSszQFOJh6GAiewMyhoxYTij5IjTtQspgCDJ9FcAUyiW+YwrbFrJ0PkRWyC8pG3+Rxb2yR5B4D5IZ1U6bv/GWdt13v9mXgHGYeF6LQKBgH4ZvW+q0WH19gzcvmQVyX9p0UdkuBY7gpoFkyr1JoDb/6QEJPTESWM5dx+9jQSpDOQPdrcDiQ6HqR2CD3ZzIiMdcESZB9QwCT/h/gYGPFOrIoSOAbyip0eTmZmbK24CgXjaRV9QNRXNy8GAw0hAzvMzVSMpPi3yGbhNnmqa1471AoGAYKdydHhknBYmRRCRZYrhkDvfHAhWF6O/jPmnoiKJ2hdzAvWY4wjlk5h1sAzy6YXvHAZWm6AF3O+PMLaUsaxcTWyHU9OnJcS2nh9FNcLpjgZJZ+LOiqZBNp+pDhRqdgR/NfMNGCMlD4Z5T/sWUAnzvFams5ooYfIJTh95+Sl3SlQ="
	payermaxPublicKey  = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoWJys/xvSgXRjgj/ftbZcbYifIuASKPskiayxLUtWg8ax6pj5OgY/JAUX1aqKZibwaZD3avRJ7tKjm1EJXaJhBpd/BNsvV00JU1A4/it50mTvQJFtlECkCs/RITgL+AQCc2rw4Zj3GtQzUaCgticKV4NkaVWLZrxizf7L4CX7DQNI1WkuG4b8PcGmU83eKLpT0SWpj+bhh7AMSPJYrP/en19l4q7PM5VzHcNtplEBPiqMCYu8MCf6lr1aWnGRyrffh4MYqlnnBUVcjV7EwMCR4O4XJgunrU53807sjrjJbYNEq3wkPwhMHrpeAKR2VfVhR/bkAK0yIJMezKbPBsZpwIDAQAB"
)

func main() {
	client, err := payermax.CreateClient(appId, merchantNo,
		merchantPrivateKey, payermaxPublicKey, "", "", payermax.Uat)
	if err != nil {
		fmt.Println(err)
	}

	data := `{"outTradeNo":"PAM20220109123456111617V1","subject":"hello","totalAmount":"0.99","currency":"USD","country":"HK","userId":"100000002","goodsDetails":[{"goodsId":"com.corps.gp.60","goodsName":"60鑽石","quantity":"1","price":"0.99","goodsCurrency":"USD","showUrl":"httpw://www.okgame.com"}],"language":"en","reference":"300011","frontCallbackUrl":"https://payapi.okgame.com/v2/PayerMax/result.html","notifyUrl":"https://payapi.okgame.com/v2/PayerMax/Callback.ashx"}`

	resp, err := client.Send("orderAndPay", data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
