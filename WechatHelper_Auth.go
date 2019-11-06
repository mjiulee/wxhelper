package wxhelper

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ddliu/go-httpclient"
	"github.com/mjiulee/lego"
)

/*微信网页授权相关接口*/

/* 授权页面地址
 * wxproxy -- 是否需要桥接代理，当出现域名绑定20个受限的时候使用桥接域名做页面授权
*/
func (self *WechatHelper) GetWxAuthPath(WxAppId, domain, prefix, state, scope ,wxproxy string) string {
	wxrequrl := "https://open.weixin.qq.com/connect/oauth2/authorize?"
	domain = url.QueryEscape(domain)
	state = url.QueryEscape(state)

	finalurl := ""
	if scope == "snsapi_base" {
		notifyPath := ""
		if len(prefix) > 0 {
			notifyPath = fmt.Sprintf("/%s/api/wx/authbase", prefix)
		} else {
			notifyPath = "/api/wx/authbase"
		}
		notifyPath = url.QueryEscape(notifyPath)
		if len(wxproxy) > 0 {
			finalurl = wxproxy + "?appid=" + WxAppId + "&redirect_uri=" + domain + notifyPath
		} else {
			finalurl = wxrequrl + "?appid=" + WxAppId + "&redirect_uri=" + domain + notifyPath
		}
	} else {
		notifyPath := ""
		if len(prefix) > 0 {
			notifyPath = fmt.Sprintf("/%s/api/wx/authinfo", prefix)
		} else {
			notifyPath = "/api/wx/authinfo"
		}

		notifyPath = url.QueryEscape(notifyPath)
		if len(wxproxy) > 0 {
			finalurl = wxproxy + "?appid=" + WxAppId + "&redirect_uri=" + domain + notifyPath
		} else {
			finalurl = wxrequrl + "?appid=" + WxAppId + "&redirect_uri=" + domain + notifyPath
		}
	}
	finalurl = finalurl + "&response_type=code"
	finalurl = finalurl + "&scope=" + scope
	finalurl = finalurl + "&state=" + state + "#wechat_redirect"
	//fmt.Println("********")
	//fmt.Println("wxrequrl=" + finalurl)
	return finalurl
}

/* 获取用户页面授权时的Accesstoken
* params:
  ---
*/
func (self *WechatHelper) GetAuthOpenId(WxAppId, WxAppSecret, xcode string) (rsp *WxAuthRsp) {
	//1.  调用微信授权接口获取access_token和openid
	authUrl := fmt.Sprintf(kWxAuthUrl, WxAppId, WxAppSecret, xcode)
	tokenRes, err := httpclient.Get(authUrl, map[string]string{})
	if err != nil {
		lego.LogError(err)
		return nil
	} else {
		bodyString, _ := tokenRes.ToString()

		var wxrsp WxAuthRsp
		if err := json.Unmarshal([]byte(bodyString), &wxrsp); err != nil {
			lego.LogInfo(bodyString)
			lego.LogError(err.Error())
		}
		return &wxrsp
	}
}

/*获取微信用户信息*/
func (self *WechatHelper) GetAuthUserInfo(accesstoken, openid string) (rsp *WxUserInfoRsp) {
	//1.  调用微信授权接口获取access_token和openid
	authUrl := fmt.Sprintf(kWxUserInfoUrl, accesstoken, openid)
	tokenRes, err := httpclient.Get(authUrl, map[string]string{})
	if err != nil {
		lego.LogError(err)
		return nil
	} else {
		bodyString, _ := tokenRes.ToString()
		lego.LogInfo(bodyString)

		var wxrsp WxUserInfoRsp
		if err := json.Unmarshal([]byte(bodyString), &wxrsp); err != nil {
			lego.LogError(err.Error())
			return nil
		}
		return &wxrsp
	}
}

/* 获取通过openid 获取用户信息，在用户关注的时候使用
* params:
  ---
*/
func (self *WechatHelper) GetUserInfoOnSubscribe(actoken, openid string) (rsp *WxSubscribeUserInfo) {
	url := "https://api.weixin.qq.com/cgi-bin/user/info"
	//1. 获取 accesstoken
	infoRes, err := httpclient.Get(url, map[string]string{
		"access_token": actoken,
		"openid":       openid,
		"lang":         "zh_CN",
	})

	if err != nil {
		lego.LogError(err.Error())
		return nil
	} else {
		bodyString, _ := infoRes.ToString()
		//fmt.Println(bodyString)
		var result WxSubscribeUserInfo
		if err := json.Unmarshal([]byte(bodyString), &result); err != nil {
			fmt.Println(err)
			return nil
		}
		return &result
	}
}

/* 获取sessionKey（小程序）
 * params:
   code --- 小程序传过来的code
*/
func (self *WechatHelper) GetSessionKey(appid, appsecret, code string) (rsp string) {
	//1. 获取 accesstoken
	tokenRes, err := httpclient.Get("https://api.weixin.qq.com/sns/jscode2session", map[string]string{
		"appid":      appid,
		"secret":     appsecret,
		"js_code":    code,
		"grant_type": "authorization_code",
	})

	if err != nil {
		lego.LogError(err.Error())
		return ""
	} else {
		bodyString, _ := tokenRes.ToString()
		return bodyString
	}
}
