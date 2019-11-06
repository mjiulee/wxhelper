package wxhelper

/*
 *  公众号的微信消息
 */
type WxMsgBase struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int    `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	MsgId        string `xml:"MsgId"` //消息id，64位整型
}

/*
 *  文本消息
 */
type WxMsgText struct {
	WxMsgBase
	Content string `xml:"Content"` // 文本
}

/*
 *  图片消息
 */
type WxMsgImage struct {
	WxMsgBase
	PicUrl  string `xml:"PicUrl"`  //图片链接（由系统生成）
	MediaId string `xml:"MediaId"` //图片消息媒体id，可以调用多媒体文件下载接口拉取数据。
}

/*
 *  Link消息
 */
type WxMsgLink struct {
	WxMsgBase
	Title       string `xml:"Title"`       //消息标题
	Description string `xml:"Description"` //消息描述
	Url         string `xml:"Url"`         //消息链接
}

/*
 *  位置消息
 */
type WxMsgLocation struct {
	WxMsgBase
	Location_X string `xml:"Location_X"` //地理位置维度
	Location_Y string `xml:"Location_Y"` //地理位置经度
	Scale      string `xml:"Scale"`      //地图缩放大小
	Label      string `xml:"Label"`      //地理位置信息
}

/*
 *  音频消息
 */
type WxMsgVoice struct {
	WxMsgBase
	MediaID     string `xml:"MediaID"`     //语音消息媒体id，可以调用多媒体文件下载接口拉取该媒体
	Format      string `xml:"Format"`      //语音格式：amr
	Recognition string `xml:"Recognition"` //语音识别结果，UTF8编码
}

/*
 *  视频消息
 */
type WxMsgVideo struct {
	WxMsgBase
	MediaId      string `xml:"MediaId"`      //视频消息媒体id，可以调用多媒体文件下载接口拉取数据。
	ThumbMediaId string `xml:"ThumbMediaId"` //视频消息缩略图的媒体id，可以调用多媒体文件下载接口拉取数据
}

/*
 *  关注Event
 */
type WxMsgSubscribe struct {
	WxMsgBase
	Event string `xml:"Event"` //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
}
