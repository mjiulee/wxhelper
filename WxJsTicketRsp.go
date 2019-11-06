package wxhelper

/* 微信js接口权限ticket
 *{
	"errcode":0,
	"errmsg":"ok",
	"ticket":"bxLdikRXVbTPdHSM05e5u5sUoXNKd8-41ZO3MhKoyN5OfkWITDGgnr2fwJ0m9E8NYzWKVZvdVtaUgWvsdshFKA",
	"expires_in":7200
}
*/
type WxJsTicketRsp struct {
	Errcode   int64  `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}
