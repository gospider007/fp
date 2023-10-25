# introduce
Help the server obtain the client's Ja3 fingerprint, http2 fingerprint, Ja4 fingerprint, Ja4H fingerprint
# features
* Completely implemented by Golang without external dependencies
* Automatic certificate, automatic replacement upon expiration
* Fast integration of all frameworks related to Golang
* Both http1.1 and http2 are supported

# get started
## install
```
go get github.com/gospider007/fp
```
## quick start with gin
```go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gospider007/fp"
)

func main() {
	handle := gin.Default()
	//fp url in : https://127.0.0.1:8999/
	handle.GET("/", fp.GinHandlerFunc)
	err := fp.Server(nil, handle.Handler(), fp.Option{Addr: ":8999"})
	log.Print(err)
}
```
## browser open fp url display
```json
{
    "akamai_fp": "1:65536,2:0,4:6291456,6:262144|15663105|0|m,a,s,p",
    "http2": {
        "InitialSetting": [
            {
                "Id": 1,
                "Val": 65536
            },
            {
                "Id": 2,
                "Val": 0
            },
            {
                "Id": 4,
                "Val": 6291456
            },
            {
                "Id": 6,
                "Val": 262144
            }
        ],
        "ConnFlow": 15663105,
        "OrderHeaders": [
            ":method",
            ":authority",
            ":scheme",
            ":path",
            "sec-ch-ua",
            "sec-ch-ua-mobile",
            "sec-ch-ua-platform",
            "upgrade-insecure-requests",
            "user-agent",
            "accept",
            "sec-fetch-site",
            "sec-fetch-mode",
            "sec-fetch-user",
            "sec-fetch-dest",
            "accept-encoding",
            "accept-language"
        ],
        "Priority": {
            "StreamDep": 0,
            "Exclusive": true,
            "Weight": 255
        }
    },
    "ja3": "772,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,16-0-5-35-45-27-18-23-17513-51-65281-11-43-10-13-41,12092-29-23-24,0",
    "ja3n": "772,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,,12092-29-23-24,0",
    "ja4": "t13d1516h2_d9e802bd6bed_ea967019fa2a",
    "ja4h": "ge20nn10zhcn_0f5a7a41a252_e3b0c44298fc_e3b0c44298fc",
    "negotiatedProtocol": "h2",
    "orderHeaders": null,
    "tls": {
        "Ciphers": [
            60138,
            4865,
            4866,
            4867,
            49195,
            49199,
            49196,
            49200,
            52393,
            52392,
            49171,
            49172,
            156,
            157,
            47,
            53
        ],
        "Curves": [
            2570,
            12092,
            29,
            23,
            24
        ],
        "Extensions": [
            64250,
            16,
            0,
            5,
            35,
            45,
            27,
            18,
            23,
            17513,
            51,
            65281,
            11,
            43,
            10,
            13,
            10794,
            41
        ],
        "Points": [
            0
        ],
        "Protocols": [
            "h2",
            "http/1.1"
        ],
        "Versions": [
            2570,
            772,
            771
        ],
        "Algorithms": [
            1027,
            2052,
            1025,
            1283,
            2053,
            1281,
            2054,
            1537
        ],
        "RandomTime": "2045-11-20 13:47:01 -0500 EST",
        "RandomBytes": "94032f738ddcb26991ae036c1838f1dc8a79ed5b3fe3efefee89cec7",
        "SessionId": "9b7991ce3fc624e0c91cc15c25f359d1f20307a055df0d5d27a31f98ce4fd74d",
        "CompressionMethods": "00"
    },
    "tlsVersion": 772,
    "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.61"
}
```
# Contributing
If you have a bug report or feature request, you can [open an issue](../../issues/new)
# Contact
If you have questions, feel free to reach out to us in the following ways:
* QQ Group (Chinese): 939111384 - <a href="http://qm.qq.com/cgi-bin/qm/qr?_wv=1027&k=yI72QqgPExDqX6u_uEbzAE_XfMW6h_d3&jump_from=webapi"><img src="https://pub.idqqimg.com/wpa/images/group.png"></a>
* WeChat (Chinese): gospider007

## Sponsors
If you like and it really helps you, feel free to reward me with a cup of coffee, and don't forget to mention your github id.
<table>
    <tr>
        <td align="center">
            <img src="https://github.com/gospider007/tools/blob/master/play/wx.jpg?raw=true" height="200px" width="200px"   alt=""/>
            <br />
            <sub><b>Wechat</b></sub>
        </td>
        <td align="center">
            <img src="https://github.com/gospider007/tools/blob/master/play/qq.jpg?raw=true" height="200px" width="200px"   alt=""/>
            <br />
            <sub><b>Alipay</b></sub>
        </td>
    </tr>
</table>


