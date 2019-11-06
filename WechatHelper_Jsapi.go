package wxhelper

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"github.com/mjiulee/lego"
	"io"
)

/****************************************************************************/
/* GetJsSDKSignature	-- 无accesstoken管理中心时调用
* params:微信JSAPI，计算当前网页的签名
  ---
*/
func (self *WechatHelper) GetJsSDKSignature(appid, appsec, noncestr, timestamp, url string) string {
	actoken := self.GetAccessToken(appid,appsec)
	jsapi_ticket := self.GetJsApiTicket(actoken)
	string1 := "jsapi_ticket=" + jsapi_ticket
	string1 += "&noncestr=" + noncestr
	string1 += "&timestamp=" + timestamp
	string1 += "&url=" + url
	s := sha1.New()
	io.WriteString(s, string1)
	//fmt.Println(string1)
	return fmt.Sprintf("%x", s.Sum(nil))
}

/* GetJsSDKSignatureWithTokenCenter -- 有accesstoken管理中心时调用
* params:微信JSAPI，计算当前网页的签名
  ---
*/
func (self *WechatHelper) GetJsSDKSignatureWithTokenCenter(appid, appsec, noncestr, timestamp, url,tcurl string) string {
	actoken := self.GetComentAccessToken(appid,tcurl)
	jsapi_ticket := self.GetJsApiTicket(actoken)
	string1 := "jsapi_ticket=" + jsapi_ticket
	string1 += "&noncestr=" + noncestr
	string1 += "&timestamp=" + timestamp
	string1 += "&url=" + url
	s := sha1.New()
	io.WriteString(s, string1)
	//fmt.Println(string1)
	return fmt.Sprintf("%x", s.Sum(nil))
}

/* JsApiTicket公众号网页JSAPI接口计算
* params:
  ---
*/
func (self *WechatHelper) GetJsApiTicket(actoken string) string {
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
	tiketRes, err := httpclient.Get(url, map[string]string{
		"type":         "jsapi",
		"access_token": actoken,
	})

	if err != nil {
		lego.LogError(err.Error())
		return ""
	} else {
		bodyString, _ := tiketRes.ToString()
		var result WxJsTicketRsp
		if err := json.Unmarshal([]byte(bodyString), &result); err != nil {
			return ""
		}
		if result.Errcode == 0 {
			return result.Ticket
		} else {
			lego.LogError(bodyString)
			return ""
		}
	}
}
