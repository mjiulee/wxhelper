package wxhelper

// 通过accesstoken和openid获得到的微信用户信息
type WxUserInfoRsp struct {
	Openid     string   `json:"openid"`
	Unionid    string   `json:"unionid"` //": " o6_bmasdasdsad6_2sgVt7hMZOPfL"
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
}

// 关注时，用户的信息
type WxSubscribeUserInfo struct {
	Subscribe      int    `json:"subscribe"`      //": 1,
	Openid         string `json:"openid"`         //": "o6_bmjrPTlm6_2sgVt7hMZOPfL2M",
	Nickname       string `json:"nickname"`       //": "Band",
	Sex            int    `json:"sex"`            //": 1,
	Language       string `json:"language"`       //": "zh_CN",
	City           string `json:"city"`           //": "广州",
	Province       string `json:"province"`       //": "广东",
	Country        string `json:"country"`        //": "中国",
	Headimgurl     string `json:"headimgurl"`     //": "http://wx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0",
	Subscribe_time int    `json:"subscribe_time"` //": 1382694957,
	Unionid        string `json:"unionid"`        //": " o6_bmasdasdsad6_2sgVt7hMZOPfL"
	Remark         string `json:"remark"`         //": "",
	Groupid        int    `json:"groupid"`        //": 0
}
