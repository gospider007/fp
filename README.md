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
	err := fp.Serve(handle.Handler(), fp.Option{Addr: ":8999"})
	log.Print(err)
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
            <img src="./play/wx.jpg" height="200px" width="200px"   alt=""/>
            <br />
            <sub><b>Wechat</b></sub>
        </td>
        <td align="center">
            <img src="./play/qq.jpg" height="200px" width="200px"   alt=""/>
            <br />
            <sub><b>Alipay</b></sub>
        </td>
    </tr>
</table>


