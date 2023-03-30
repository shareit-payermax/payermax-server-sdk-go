package main

import (
	"fmt"
	"github.com/PayermaxZhanglong/payermax-server-sdk-go/payermax"
)

const (
	appId              = "6666c8b036a24579974497c2f9a33333"
	merchantNo         = "010213834784554"
	merchantPrivateKey = "MIIEpAIBAAKCAQEAynjvfZjICKDsvkLAKgWi3Vj95vLc4p62GTGieJckjgQ5iBf+bLPRpHn4vNtNGKQpKIOJfrFJPGpCzf7blkZz20XJ9Fg/HmDOQCDH2t/e+lMrvHibnau2Fu0s3ddI5bP8PNy9dwkPL/zLGIIzwESXjNb1hNh1CzddC9x3Kg1kd0DjPtMs9nmdss6n8Z0OO1TGIm2bFAKL0taVwe9SfOD+53vgBtICmrxGUV/eAUuJ/jZwmyOyyrCa7yHSt4nG8DWP238vKEvC7tRJk6aXsqSpylAejX5TXmWr01WV4q9jSOSKZx5RzHK5kxGHsSSbkuihKPwCtIbI9nAFOcFJkECkOQIDAQABAoIBAQCLzZ1wF8JnUtalOOn/Gg0u0Xffs+oZAIW0N7V7lsFC3l/rPMYMwM0MyLubi8FcNl3E/648sKVk9epS2ps81EDDMxkTgqtyil1fokLdOp94MV2Nsamh4SLGCdZlB3XqRbDxRWn9e1/lPqPttFmPdM1ADl1Q8TVAHWY9/mi5vK2WORjc915rUJFxlK6CHaPlD6WQWxXvmDYBjlzUSSFR+jYMtE+Hnk740t46MDCPYV6lVh7GuyxHKpx+qryrVPwDA/tw/WWiM7gsPr11NHAB0a2R86Jqg/nv+BhCwGpj5VBxx+eyzf5IO94EjBgUH2DaYtGNIpxSetmdFMIk/teRINIBAoGBANg3xo3kO15EtiqyBZuv+2pj44StoCh3sGkGNHDLEVmk7wiUxCJ+ZxLJVoltwKVO8L2sVbDDUfxn6ZIbuQHSsBdGHeE+fxrTc0LPEwZ48gyogxcbQcz9FHgP6w7IWJsIVUadX0DuNdAl9VYOVSq1Tgt1wWv4umBRuM7VnxHqh1DRAoGBAO+5uarkFgg0G8rlDYNpAEbEcR1sHrI67G9iBKpIF4O67BXEBy9t/WgjME3nPFZPtNQaFSiZA9tesiKRiFyGXFz4iO9q46ghOTCndwxQP8aiqlB63fqlvDw7qXO7GI0PbjIZgOAuTykZ9bipN6Qh8d1SFRWnDXGJdQcKrSHyTDbpAoGBAICfHWSHIrH/WgaoBCILXBp79XqV9rJcEPtJD6URh+616ORH6y1B2Hsafnoeaf1sqlWK0Sbn6jumbRHXoATvmoUd1uSJUv9YTjauDHlLNWJGVEVIl6oj2ytY/NG8aMlA+cmaEHIFwslh60IYIJ3ZYOX8VOWv/t8Rfki8V3ZG99whAoGAXi+1cBwfP+fhR41JCul1T1idLLcvNE2MWZLETHb4riwB1+dl/0+SsZipwOHqResZHACHcaT06/q/uG8/iULNBUYs3ww7F/K9uo0BbBgXhp6glfBASNtXIr86K5tF4R4/6HU2ul3XgkmNzpjFoLopghBe8lvpH0OndDXQojbFlQkCgYBfG99FvTY4kJm8R0uGhk+g5v1X+ZoWm4eJnpi4pDQBfaMqjB03n/E6ry1XkHz+vXGEFPL1cokK8FsGX8LTHUTJS3ojYSVrlBAQkuBws4Omd8TzXwJpphJicVAIlOXXcf2AbI7wl2yaGBMSbNCRwyxKFiIR3PNWwwRj9svuA5PBAg=="
	payermaxPublicKey  = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChMWd9o9Anc8GbSYsSgx5sJOj+l3trNSchFSeVWAX+zA7P7Q+tdSh+i58Qn+jNw2yCbNoD8ev55O9B/eHe2UfrwwtEbu6At2AKxl8Y3MJI4rxieKZI4+t/quTKKyJvuf7N9t8txxPCfNTEzbFCtRugdZj7J+Z+jM4io/QXPUkuIQIDAQAB"
)

func main() {
	client, err := payermax.CreateClient(appId, merchantNo,
		merchantPrivateKey, payermaxPublicKey, "", "", payermax.Uat)
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
