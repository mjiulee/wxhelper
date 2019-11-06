package wxhelper

import (
	"crypto/tls"
	"encoding/xml"
	"github.com/ddliu/go-httpclient"
	"github.com/mjiulee/lego"
	"net/http"
	"strings"
	"github.com/mjiulee/lego/utils"
)

/* BRIEF: 微信支付相关接口
*/


type WxPayOrderRsp struct {
	// 基础返回内容
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`

	// 商户基本信息
	Appid    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`

	// 统一下单返回内容
	PrepayId  string `xml:"prepay_id"`
	TradeType string `xml:"trade_type"`

	// 退款返回内容
	TransactionId string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	OutRefundNo   string `xml:"out_refund_no"`
	RefundId      string `xml:"refund_id"`
	RefundChannel string `xml:"refund_channel"`
	RefundFee     string `xml:"refund_fee"`
}

/* 微信支付统一下单(APP 下单，公众号下单请参考快购)
* params:
  ---
*/
func (self *WechatHelper) UnifiedOrder(appid, mch_id, mch_key, openid, orderno, notifyUrl,obody,odetail string, amount int64) *WxPayOrderRsp {
	ip := utils.GetLocalIpAddress()
	nonce_str := string(utils.RandString(32, utils.KC_RAND_KIND_ALL))
	header := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	params := map[string]string{
		"appid":            appid,
		"mch_id":           mch_id,
		"trade_type":       "JSAPI",
		"body":             obody,
		"detail":           odetail,
		"out_trade_no":     orderno,
		"total_fee":        utils.Int64ToString(amount),
		"notify_url":       notifyUrl,
		"openid":           openid,
		"nonce_str":        nonce_str,
		"spbill_create_ip": ip,
	}

	sign := self.GenSign(mch_key, params)
	params["sign"] = sign

	xmlParams := utils.Map2Xml(params)
	lego.LogInfo(xmlParams)

	body := strings.NewReader(xmlParams)
	tokenRes, err := httpclient.Do("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", header, body)

	if err != nil {
		lego.LogError(err.Error())
		return nil
	} else {
		bodyString, _ := tokenRes.ToString()
		lego.LogError(bodyString)
		var wxrsp WxPayOrderRsp
		if err := xml.Unmarshal([]byte(bodyString), &wxrsp); err != nil {
			lego.LogError(err.Error())
			return nil
		}
		return &wxrsp
	}
}

/* 退款申请
* params:
  ---
*/
func (self *WechatHelper) Refund(appid, mch_id, mch_key, orderno string, total_fee, refund_fee int, tlsConfig *tls.Config) *WxPayOrderRsp {

	nonce_str := string(utils.RandString(32, utils.KC_RAND_KIND_ALL))
	params := map[string]string{
		"appid":         appid,
		"mch_id":        mch_id,
		"out_trade_no":  orderno,
		"out_refund_no": orderno,
		"total_fee":     utils.IntToString(total_fee),
		"refund_fee":    utils.IntToString(refund_fee),
		"nonce_str":     nonce_str,
		"op_user_id":    mch_id,
	}

	sign := self.GenSign(mch_key, params)
	params["sign"] = sign

	xmlContent := utils.Map2Xml(params)

	lego.LogInfo("********退款参数********")
	lego.LogInfo(xmlContent)

	body := strings.NewReader(xmlContent)

	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}

	rsp, err := client.Post("https://api.mch.weixin.qq.com/secapi/pay/refund", "text/xml", body)
	if err != nil {
		lego.LogError(err.Error())
		return nil
	} else {
		respone := Response{rsp}
		bodyString, _ := respone.ToString()

		var wxrsp WxPayOrderRsp
		if err := xml.Unmarshal([]byte(bodyString), &wxrsp); err != nil {
			lego.LogError(err.Error())
		}
		return &wxrsp
	}
}

/*****************************************************************************************************************************/
type WxPayToWxUserResponse struct {
	// 基础返回内容
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	// 商户基本信息
	Appid    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`

	// 统一下单返回内容
	PartnerTradeNo string `xml:"partner_trade_no"`
	PaymentNo      string `xml:"payment_no"`
	PaymentTime    string `xml:"payment_time"`

	// 全部返回内容
	Body string
}

/* 提现申请-企业付款-微信个人钱包
* params:
  ---
*/
func (self *WechatHelper) WithdrawToWechat(appid, mch_id, mch_key, openId, orderno string, amount int64, tlsConfig *tls.Config) *WxPayToWxUserResponse {
	nonce_str := string(utils.RandString(32, utils.KC_RAND_KIND_ALL))
	spcip := utils.GetLocalIpAddress()

	params := map[string]string{
		"mch_appid":        appid,
		"mchid":            mch_id,
		"nonce_str":        nonce_str,
		"partner_trade_no": orderno,
		"openid":           openId,
		"check_name":       "NO_CHECK",
		"amount":           utils.Int64ToString(amount),
		"desc":             "用户红包提现",
		"spbill_create_ip": spcip,
	}

	sign := self.GenSign(mch_key, params)
	params["sign"] = sign

	xmlContent := utils.Map2Xml(params)
	lego.LogInfo("********提现参数********")
	lego.LogInfo(xmlContent)

	body := strings.NewReader(xmlContent)

	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}

	rsp, err := client.Post("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", "text/xml", body)
	if err != nil {
		lego.LogError(err.Error())
		return nil
	} else {
		respone := Response{rsp}
		bodyString, _ := respone.ToString()

		var wxrsp WxPayToWxUserResponse
		if err := xml.Unmarshal([]byte(bodyString), &wxrsp); err != nil {
			lego.LogError(bodyString)
			lego.LogError(err.Error())
		}

		wxrsp.Body = bodyString
		return &wxrsp
	}
}


/****************************************************************************
/* 生成支付信息-公众号网页使用
* params:
  ---
*/
func (self *WechatHelper) GenPayInfo(appid, mchkey, prepayid string) (payinfo map[string]string) {
	// appid := config.GetIniByKey("WECHAT", "WX_APPID")
	nonce_str := string(utils.RandString(32, utils.KC_RAND_KIND_ALL))
	params := map[string]string{
		"appId":     appid,
		"timeStamp": utils.GetTimeStamp(),
		"nonceStr":  nonce_str,
		"package":   "prepay_id=" + prepayid,
		"signType":  "MD5",
	}

	sign := self.GenSign(mchkey, params)
	params["paySign"] = sign

	return params
}

/* 生成支付信息-小程序或App使用
* params:
  ---
*/
func (self *WechatHelper) GenPayInfoForApp(appid, mch_id, mchkey, prepayid string) (payinfo map[string]string) {
	nonce_str := string(utils.RandString(32, utils.KC_RAND_KIND_ALL))
	params := map[string]string{
		"appid":     appid,
		"partnerid": mch_id,
		"prepayid":  prepayid,
		"timestamp": utils.GetTimeStamp(),
		"noncestr":  nonce_str,
		"package":   "Sign=WXPay",
	}

	sign := self.GenSign(mchkey, params)
	params["sign"] = sign

	return params
}
