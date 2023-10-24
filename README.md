# introduce
Help the server obtain the client's Ja3 fingerprint, http2 fingerprint, and Ja4 fingerprint
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
            "cache-control",
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
    "ja3": "772,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,23-45-13-11-16-65281-27-17513-10-43-51-5-18-35-21,29-23-24,0",
    "ja3n": "772,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,,29-23-24,0",
    "ja4": "t13i1515h2_d9e802bd6bed_c72f4b657d17",
    "negotiatedProtocol": "h2",
    "tls": {
        "Ciphers": [
            31354,
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
            29,
            23,
            24
        ],
        "Extensions": [
            64250,
            23,
            45,
            13,
            11,
            16,
            65281,
            27,
            17513,
            10,
            43,
            51,
            5,
            18,
            35,
            47802,
            21
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
        "RandomTime": "2057-07-29 16:32:24 +0800 CST",
        "RandomBytes": "da910004407cbc3583d60526e1656691bfcc23cd5f9c8e9320cc6c35",
        "SessionId": "3f019637350afb7a08591b2b2bda9ca9b7a8010478af331609e275c19bd41223",
        "CompressionMethods": "00"
    },
    "tlsVersion": 772,
    "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60"
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


