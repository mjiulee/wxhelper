package wxhelper

import (
	"bytes"
	"fmt"
	"encoding/json"
	"github.com/ddliu/go-httpclient"
	"time"

	"github.com/mjiulee/lego"
)

/****************************************************************************/
/* 获取接口调用时，API的accesstoken参数，单应用的情况下调用
* params:
  ---
*/
func (self *WechatHelper) GetAccessToken(appid, appsecret string) (rsp string) {
	actoken ,ok := _actokenMap[appid]
	if ok {
		if false == self.ifAccessTokenExpire(actoken) {
			return actoken.AccessToken
		}
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, appsecret)
	//1. 获取 accesstoken
	tokenRes, err := httpclient.Get(url, map[string]string{})

	if err != nil {
		// logger.Info(err.Error())
		fmt.Println(err)
		return ""
	} else {
		bodyString, _ := tokenRes.ToString()
		var result WxAccessTokenRsp
		if err := json.Unmarshal([]byte(bodyString), &result); err != nil {
			return ""
		}
		if len(result.AccessToken) > 0 { // 大于0的情况下才会是成功
			result.RequestTime = time.Now().Unix()
			_actokenMap[appid] = &result
			return result.AccessToken
		} else {
			lego.LogError(bodyString)
			return ""
		}

	}
}

/* 获取接口调用时，API的accesstoken参数，多应用的情况下调用中控管理的token
* 存在一个公众号，在不同的应用上做模块开发的情况下，accesstoken要从同一的token管理中心获取，否则会错乱
* params:
 ---
*/
func (self *WechatHelper) GetComentAccessToken(appid ,tokenCenterurl string) (token string) {
	//1. 获取 accesstoken
	tokenRes, err := httpclient.Get(tokenCenterurl, map[string]string{
		"appid": appid,
	})

	if err != nil {
		lego.LogError(err.Error())
		return ""
	} else {
		bodyString, _ := tokenRes.ToString()
		//fmt.Println(bodyString)
		type gatapirsp struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data string `json:"data"`
		}

		var rsp gatapirsp
		bodybyte := bytes.TrimPrefix([]byte(bodyString), []byte("\xef\xbb\xbf"))
		if err := json.Unmarshal(bodybyte, &rsp); err != nil {
			fmt.Println(err)
			return ""
		}

		if rsp.Code != 0 {
			return ""
		} else {
			return rsp.Data
		}
	}
}

/* 判断token是否失效
* params:
  ---
*/
func (self *WechatHelper) ifAccessTokenExpire(ret *WxAccessTokenRsp) (expire bool) {
	if ret == nil {
		return true
	}

	nowTimeStamp := time.Now().Unix()
	if nowTimeStamp >= ret.RequestTime+ret.ExpiresIn-500 {
		return true
	}

	return false
}
