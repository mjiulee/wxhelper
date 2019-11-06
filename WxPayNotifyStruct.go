package wxhelper

// 支付回调
type WxPayNotifyStruct struct {
	// 基础返回内容
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`

	// 商户基本信息
	Appid    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`

	// 统一下单返回内容
	OpenId    string `xml:"openid"`
	TradeType string `xml:"trade_type"`

	// 退款返回内容
	BankType      string `xml:"bank_type"`
	TotalFee      string `xml:"total_fee"`
	CashFee       string `xml:"cash_fee"`
	TransactionId string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	TimeEnd       string `xml:"time_end"`
}
