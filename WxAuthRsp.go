package wxhelper

// 公众号页面用户确定授权后，通过scrop和code请求返回的数据结构
type WxAuthRsp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}