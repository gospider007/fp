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
## quick start
```go
package main

import (
	"github.com/gospider007/fp"
)

func main() {
	fp.Start("")
}
```
## browser open fp url display
```json
{
    "goSpiderSpec": "1603010758010007540303549da70a27650ac8bbd769fba4f6cde3526c8284870ca819a2d7ee59f91dd5d82016ad87de639f4d6c80f0f5e690ed09c69fb807aa501b6d36f62b7362b9c6c2d500202a2a130113021303c02bc02fc02cc030cca9cca8c013c014009c009d002f0035010006eb6a6a000044cd000500030268320005000501000000000010000e000c02683208687474702f312e3100120000003304ef04ed3a3a00010011ec04c0c71c97d646c904b105647471b33b10d804695bc2c8d3c86f47f94fdb673e26da8e2ecb3400570d48f17b78f37e734ca1602c3b54a50bc9017c9324bb779cbb03d1721b000c36e80820d21c02b8c23278bd63185110aa6fb3fb4544b83d53c04d2802bc48c59a62b61c5cf2bd7f9882c6b047e85b10efe011db940bd0d68a6ce8c63236548203857e141c39fb505471bc1aeb3beb0473fe88384c692d64a515b235b4d5d47a8af75ab0b2b47143b2b356bb314ab75e11a7938b29bb6008c3133157dba7f490bbd5822749f942ac381e57d90a06513c3f1c50ae573d9446092af116122aaa4d100c0c047f89f96e405c792bb7526ed9c95f48836ed121f9f07421da3b1c498969ab249a8b3b470808262482ae34ae8d7515f2230538907d936a98f4213e46464d1ca34c70f37bc0337fe3823e3bbc63070830fed379009dc6378288a3d760ac5c311cd564af304a38acc217955706740c4484c6ae94b374da3efa740554856a58a35ff2a243704c9481511ca93920f463b75bf63369941145f5a51cf698df301ba424186aaa6afde15a57a5bfab66495013a9c885c839543fa757b29bb34471db47295086041873ef63879752702cd1407a0904aed915e66633ee44afea607fbe9074df581abd55039b88223d002d56015df3093826d7544fda95b4d8083c7a02df40ce93683c739ca4f006148186b92b8740586a2841954b9d827610035d8b2c7472a12c0773505c9660f29b8b49b4964db015b19a171a85b4a4ec4d0ee55d786202fa9c523cbaa07aaa9b49a94c02b8af3f97c8cd40106c0bc8c3e643637180bfc78673837cca5852720a065d0125e03077c55cc0054cc2e280bbcf08a59ed037a695cab1f32ffb9233f97484a079261171bc3b56869b290155e498cc205b33b40b3ddb3f8ff1b996113ea04c70e053b9ab75b49c4102eb939c789c2029994924c1c784a7c772761b23193490f09ee06b2e213398ba1bc7e121913fc9668421c78858a6d4a5cc433b7861ab679864087ff15b09aa175e249cb1a5317b01b16137b507acbf66684a4d12cc336740ce06b04a83a10ab9cff08a044232571e35be934b78b9b42d2ab65a808c32951a98a2a55f070b24b13b9b3909bd02727e6787366a62bb651a9ac008a00c2438f680a4ac520e5961cde0aa28b853183d97831e187aeaa08f1076250d9c44ccb459d0db64a65c6548d5506b986f068804d0074040e2534e3b6365a8008940b3323b0e4b01c7573556bca6ce806287c2e719f2011db8245088b9181c5ca1c6e14b6779c51a802e2b6c6010a9b04c5aaac798a746308f663456aa49007d0059aebc9c5bf31051c87eb8e06816342206c08a7e21aa4340bdae00c96823cebb617690dccd8c77b1e3275eeee47d34722ecd861383b8304ff0bf1f9282926a1e53f4a05ddc77cafcce956766e92c99a7557506c46b67b2c579638cd4a446c7121a1c661fb0d66d10da25079b22472231b69946450308cf80026c1869c0d240d260873f54b76e28abe58c2351591ad1f69a93700be3980e463262929b9c6b5514594134eb0ba2aa4a75ea66440cac73cb71c62655929854bc08bc161160c9f998b92bd548df41056d18487a375508577cac027ecd8238b9435a137425210108f2546a31d50068b1034cafff90f2a8081db595ce962c3b4eb7a9e68d684a5a383dad2aaa8d0a49ee20a470f0c84e69c451001d0020ee82524eb01895d2d9d73b340153a2e1f2b87b3f2b6bd8a28f6a784c07842201000a000c000a3a3a11ec001d00170018000b00020100ff01000100001b0003020002002b000706eaea030403030023000000170000fe0d00da0000010001f80020ac0e6ddc19a2713d529afbab7c2b95b246acb8fa6d9a43b79cb7d0bf7cd8166100b0eb150ced81bb0865deb78c62e300306611a6137b98936036d9b86d0f8f0863c3e31a6fdd44328bc617ec2f29101f61dbfecef34ac83e867009faf129135507b1dac657bfb7d590ecff47885249ca4d4e08d8c1ebeaa1028598ddee5dbe79536cc36a963dc1aa2487561593e6d95ae9e06c9f520d8bbd33aa435320cd748238ff4521ac30ba320b2460c6afb023c3df0f983b3a856f43c8402a7829cc13105c4ce98f87600fae96d4d7b8bc99afc48734000d0012001004030804040105030805050108060601002d00020101baba00010000290094006f006916d86c1287b922b674b29381334b0a131fe978fe4ef4a3dc6b3c9553ad10fa55b6c8faa65c2ee18b2cbc1f5c04d65be3ad33e4c86c70a3635b0ead4a6692957f232ddb31cd1d00d66abe8ea49176aee52b31c7e9ae3d9e1888f011b75e450965385cf5cdf91edad3e1d86c801200212044977c15a44b7bb3ae7f1a631241a1e6d49ed1a887225148361e079824655c19@@505249202a20485454502f322e300d0a0d0a534d0d0a0d0a00001804000000000000010001000000020000000000040060000000060004000000000408000000000000ef00010001f101250000000180000000ff82418a089d5c0b8170dc79f7df878440874148b1275ad1ffb9fe749d3fd4372ed83aa4fe7efbc1fcbefff3f4a7f388e79a82a97a7b0f497f9fbef07f21659fe7e94fe6f4f61e935b4ff3f7de0fe42cb3fcff408b4148b1275ad1ad49e33505023f30408d4148b1275ad1ad5d034ca7b29f07226d61634f53224092b6b9ac1c8558d520a4b6c2ad617b5a54251f01317ad9d07f66a281b0dae053fad0321aa49d13fda992a49685340c8a6adca7e28104416e277fb521aeba0bc8b1e632586d975765c53facd8f7e8cff4a506ea5531149d4ffda97a7b0f49580b2cae05c0b814dc394761986d975765cf40884148b576d959d05f8daec2ca54927fbaec2d85aa42d94085aedb2b3a0b86aec2ca54927f53e5497ca589d34d1f43aeba0c41a4c7a98f33a69a3fdf9a68fa1d75d0620d263d4c79a68fbed00177fe8d48e62b03ee697e8d48e62b1e0b1d7f46a4731581d754df5f2c7cfdf6800bbdf43aeba0c41a4c7a9841a6a8b22c5f249c754c5fbef046cfdf6800bbbf408a4148b4a549275906497f83a8f517408a4148b4a549275a93c85f86a87dcd30d25f408a4148b4a549275ad416cf023f31408a4148b4a549275a42a13f8690e4b692d49f50929bd9abfa5242cb40d25fa523b3e94f684c9f518cf73ad7b4fd7b9fefb4005dff4086aec31ec327d785b6007d286f",
    "h2": {
        "connFlow": 15663105,
        "orderHeaders": [
            [
                "sec-ch-ua",
                "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\""
            ],
            [
                "sec-ch-ua-mobile",
                "?0"
            ],
            [
                "sec-ch-ua-platform",
                "\"macOS\""
            ],
            [
                "upgrade-insecure-requests",
                "1"
            ],
            [
                "user-agent",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
            ],
            [
                "sec-purpose",
                "prefetch;prerender"
            ],
            [
                "purpose",
                "prefetch"
            ],
            [
                "accept",
                "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
            ],
            [
                "sec-fetch-site",
                "none"
            ],
            [
                "sec-fetch-mode",
                "navigate"
            ],
            [
                "sec-fetch-user",
                "?1"
            ],
            [
                "sec-fetch-dest",
                "document"
            ],
            [
                "accept-encoding",
                "gzip, deflate, br, zstd"
            ],
            [
                "accept-language",
                "zh-CN,zh;q=0.9"
            ],
            [
                "priority",
                "u=0, i"
            ]
        ],
        "pri": "PRI * HTTP/2.0",
        "priority": {
            "exclusive": true,
            "streamDep": 0,
            "weight": 255
        },
        "settings": [
            {
                "ID": 1,
                "Val": 65536
            },
            {
                "ID": 2,
                "Val": 0
            },
            {
                "ID": 4,
                "Val": 6291456
            },
            {
                "ID": 6,
                "Val": 262144
            }
        ],
        "sm": "SM",
        "streams": [
            {
                "name": "Http2SettingsFrame",
                "settings": [
                    {
                        "id": 1,
                        "val": 65536
                    },
                    {
                        "id": 2,
                        "val": 0
                    },
                    {
                        "id": 4,
                        "val": 6291456
                    },
                    {
                        "id": 6,
                        "val": 262144
                    }
                ],
                "streamID": 0,
                "type": 4
            },
            {
                "connFlow": 15663105,
                "name": "Http2WindowUpdateFrame",
                "streamID": 0,
                "type": 8
            },
            {
                "headers": [
                    {
                        "name": "sec-ch-ua",
                        "value": "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\""
                    },
                    {
                        "name": "sec-ch-ua-mobile",
                        "value": "?0"
                    },
                    {
                        "name": "sec-ch-ua-platform",
                        "value": "\"macOS\""
                    },
                    {
                        "name": "upgrade-insecure-requests",
                        "value": "1"
                    },
                    {
                        "name": "user-agent",
                        "value": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
                    },
                    {
                        "name": "sec-purpose",
                        "value": "prefetch;prerender"
                    },
                    {
                        "name": "purpose",
                        "value": "prefetch"
                    },
                    {
                        "name": "accept",
                        "value": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
                    },
                    {
                        "name": "sec-fetch-site",
                        "value": "none"
                    },
                    {
                        "name": "sec-fetch-mode",
                        "value": "navigate"
                    },
                    {
                        "name": "sec-fetch-user",
                        "value": "?1"
                    },
                    {
                        "name": "sec-fetch-dest",
                        "value": "document"
                    },
                    {
                        "name": "accept-encoding",
                        "value": "gzip, deflate, br, zstd"
                    },
                    {
                        "name": "accept-language",
                        "value": "zh-CN,zh;q=0.9"
                    },
                    {
                        "name": "priority",
                        "value": "u=0, i"
                    }
                ],
                "name": "Http2MetaHeadersFrame",
                "priority": {
                    "exclusive": true,
                    "streamDep": 0,
                    "weight": 255
                },
                "streamID": 1,
                "type": 1
            }
        ]
    },
    "tls": {
        "algorithms": [
            1027,
            2052,
            1025,
            1283,
            2053,
            1281,
            2054,
            1537
        ],
        "cipherSuites": [
            10794,
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
        "compressionMethods": "AA==",
        "contentType": 22,
        "curves": [
            2570,
            4588,
            29,
            23,
            24
        ],
        "extensions": [
            {
                "data": "",
                "type": 27242
            },
            {
                "data": "0003026832",
                "type": 17613
            },
            {
                "data": "0100000000",
                "type": 5
            },
            {
                "data": "000c02683208687474702f312e31",
                "type": 16
            },
            {
                "data": "",
                "type": 18
            },
            {
                "data": "04ed3a3a00010011ec04c0c71c97d646c904b105647471b33b10d804695bc2c8d3c86f47f94fdb673e26da8e2ecb3400570d48f17b78f37e734ca1602c3b54a50bc9017c9324bb779cbb03d1721b000c36e80820d21c02b8c23278bd63185110aa6fb3fb4544b83d53c04d2802bc48c59a62b61c5cf2bd7f9882c6b047e85b10efe011db940bd0d68a6ce8c63236548203857e141c39fb505471bc1aeb3beb0473fe88384c692d64a515b235b4d5d47a8af75ab0b2b47143b2b356bb314ab75e11a7938b29bb6008c3133157dba7f490bbd5822749f942ac381e57d90a06513c3f1c50ae573d9446092af116122aaa4d100c0c047f89f96e405c792bb7526ed9c95f48836ed121f9f07421da3b1c498969ab249a8b3b470808262482ae34ae8d7515f2230538907d936a98f4213e46464d1ca34c70f37bc0337fe3823e3bbc63070830fed379009dc6378288a3d760ac5c311cd564af304a38acc217955706740c4484c6ae94b374da3efa740554856a58a35ff2a243704c9481511ca93920f463b75bf63369941145f5a51cf698df301ba424186aaa6afde15a57a5bfab66495013a9c885c839543fa757b29bb34471db47295086041873ef63879752702cd1407a0904aed915e66633ee44afea607fbe9074df581abd55039b88223d002d56015df3093826d7544fda95b4d8083c7a02df40ce93683c739ca4f006148186b92b8740586a2841954b9d827610035d8b2c7472a12c0773505c9660f29b8b49b4964db015b19a171a85b4a4ec4d0ee55d786202fa9c523cbaa07aaa9b49a94c02b8af3f97c8cd40106c0bc8c3e643637180bfc78673837cca5852720a065d0125e03077c55cc0054cc2e280bbcf08a59ed037a695cab1f32ffb9233f97484a079261171bc3b56869b290155e498cc205b33b40b3ddb3f8ff1b996113ea04c70e053b9ab75b49c4102eb939c789c2029994924c1c784a7c772761b23193490f09ee06b2e213398ba1bc7e121913fc9668421c78858a6d4a5cc433b7861ab679864087ff15b09aa175e249cb1a5317b01b16137b507acbf66684a4d12cc336740ce06b04a83a10ab9cff08a044232571e35be934b78b9b42d2ab65a808c32951a98a2a55f070b24b13b9b3909bd02727e6787366a62bb651a9ac008a00c2438f680a4ac520e5961cde0aa28b853183d97831e187aeaa08f1076250d9c44ccb459d0db64a65c6548d5506b986f068804d0074040e2534e3b6365a8008940b3323b0e4b01c7573556bca6ce806287c2e719f2011db8245088b9181c5ca1c6e14b6779c51a802e2b6c6010a9b04c5aaac798a746308f663456aa49007d0059aebc9c5bf31051c87eb8e06816342206c08a7e21aa4340bdae00c96823cebb617690dccd8c77b1e3275eeee47d34722ecd861383b8304ff0bf1f9282926a1e53f4a05ddc77cafcce956766e92c99a7557506c46b67b2c579638cd4a446c7121a1c661fb0d66d10da25079b22472231b69946450308cf80026c1869c0d240d260873f54b76e28abe58c2351591ad1f69a93700be3980e463262929b9c6b5514594134eb0ba2aa4a75ea66440cac73cb71c62655929854bc08bc161160c9f998b92bd548df41056d18487a375508577cac027ecd8238b9435a137425210108f2546a31d50068b1034cafff90f2a8081db595ce962c3b4eb7a9e68d684a5a383dad2aaa8d0a49ee20a470f0c84e69c451001d0020ee82524eb01895d2d9d73b340153a2e1f2b87b3f2b6bd8a28f6a784c07842201",
                "type": 51
            },
            {
                "data": "000a3a3a11ec001d00170018",
                "type": 10
            },
            {
                "data": "0100",
                "type": 11
            },
            {
                "data": "00",
                "type": 65281
            },
            {
                "data": "020002",
                "type": 27
            },
            {
                "data": "06eaea03040303",
                "type": 43
            },
            {
                "data": "",
                "type": 35
            },
            {
                "data": "",
                "type": 23
            },
            {
                "data": "0000010001f80020ac0e6ddc19a2713d529afbab7c2b95b246acb8fa6d9a43b79cb7d0bf7cd8166100b0eb150ced81bb0865deb78c62e300306611a6137b98936036d9b86d0f8f0863c3e31a6fdd44328bc617ec2f29101f61dbfecef34ac83e867009faf129135507b1dac657bfb7d590ecff47885249ca4d4e08d8c1ebeaa1028598ddee5dbe79536cc36a963dc1aa2487561593e6d95ae9e06c9f520d8bbd33aa435320cd748238ff4521ac30ba320b2460c6afb023c3df0f983b3a856f43c8402a7829cc13105c4ce98f87600fae96d4d7b8bc99afc48734",
                "type": 65037
            },
            {
                "data": "001004030804040105030805050108060601",
                "type": 13
            },
            {
                "data": "0101",
                "type": 45
            },
            {
                "data": "00",
                "type": 47802
            },
            {
                "data": "006f006916d86c1287b922b674b29381334b0a131fe978fe4ef4a3dc6b3c9553ad10fa55b6c8faa65c2ee18b2cbc1f5c04d65be3ad33e4c86c70a3635b0ead4a6692957f232ddb31cd1d00d66abe8ea49176aee52b31c7e9ae3d9e1888f011b75e450965385cf5cdf91edad3e1d86c801200212044977c15a44b7bb3ae7f1a631241a1e6d49ed1a887225148361e079824655c19",
                "type": 41
            }
        ],
        "handShakeType": 1,
        "handshakeVersion": 771,
        "messageVersion": 769,
        "points": "AA==",
        "protocols": [
            "h2",
            "http/1.1"
        ],
        "randomBytes": "J2UKyLvXafuk9s3jUmyChIcMqBmi1+5Z+R3V2A==",
        "randomTime": 1419618058,
        "sessionId": "Fq2H3mOfTWyA8PXmkO0Jxp+4B6pQG2029itzYrnGwtU=",
        "versions": [
            2570,
            772,
            771
        ]
    }
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


