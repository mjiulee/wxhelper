package wxhelper


// 接口调用的时候的access_token: https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140183
type WxAccessTokenRsp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	RequestTime int64
}