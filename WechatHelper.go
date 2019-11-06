package wxhelper

import (
	//"crypto/aes"
	//"crypto/cipher"
	//"encoding/base64"
	//"encoding/json"
	//"fmt"

	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"

	// "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"errors"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ddliu/go-httpclient"

	"github.com/mjiulee/lego"
	"github.com/mjiulee/lego/utils"
)

const (
	// 微信授权：snsapi_base
	kWxAuthUrl = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	// 微信授权：snsapi_userinfo
	kWxUserInfoUrl = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
)

var _actokenMap map[string]*WxAccessTokenRsp // 全局
var _tlsConfig *tls.Config                   // 退款tsl配置

func init() {
	_actokenMap = make(map[string]*WxAccessTokenRsp)
}

/*****************************************************************************************************************************/
/* 微信接口交互帮助类
 */
type WechatHelper struct{}

/****************************************************************************/
/* 生成小程序2维码
* params:
  ---
*/
func (self *WechatHelper) GenMiniQrCode(urlpath string, orderNo string, token string) (rsp string) {

	url := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=" + token
	header := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}

	line_color := map[string]interface{}{
		"r": "0",
		"g": "0",
		"b": "0",
	}

	params := map[string]interface{}{
		"scene":      orderNo,
		"path":       urlpath,
		"width":      430,
		"auto_color": true,
		"line_color": line_color,
	}

	body, err := json.Marshal(params)
	lego.LogInfo(string(body))

	//1. 获取 accesstoken
	tokenRes, err := httpclient.Do("POST", url, header, strings.NewReader(string(body)))

	if err != nil {
		lego.LogError(err.Error())
		return ""
	} else {
		bodyString, _ := tokenRes.ToString()
		// logger.Println(bodyString)

		// 保存图片
		now := time.Now().Format("20060102")
		saveDir := filepath.Join("./upload", path.Clean("/"+now))
		if exist, _ := utils.PathExists(saveDir); !exist {
			mkdirerr := os.Mkdir(saveDir, os.ModePerm)
			if mkdirerr != nil {
				lego.LogError("2维码文件保存目录创建失败" + err.Error())
				return ""
			}
		}

		timestamp := utils.Int64ToString(time.Now().Unix())

		savePath := path.Clean("/" + now + "/" + timestamp + ".jpg")
		lego.LogInfo("savePath=" + savePath)

		// 打开保存文件句柄
		fileName := filepath.Join("./upload", savePath)
		lego.LogInfo("fileName=" + fileName)

		fp, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			lego.LogError("2维码文件创建失败" + err.Error())
			return ""
		}

		if _, err = io.Copy(fp, strings.NewReader(string(bodyString))); err != nil {
			lego.LogError("2维码文件保存失败" + err.Error())
			return ""
		}

		defer fp.Close()
		return "/static" + savePath
	}
}

/****************************************************************************/
/* 根据页面上传文件时的serviceid，下载微信上传图片到本地服务器
* params:
  ---
*/
func (self *WechatHelper) DownLoadWechatMedia(actoken string, mediaid string) (bool, string) {
	url := "http://file.api.weixin.qq.com/cgi-bin/media/get"
	tiketRes, err := httpclient.Get(url, map[string]string{
		"media_id":     mediaid,
		"access_token": actoken,
	})

	if err != nil {
		lego.LogError(err.Error())
		return false, ""
	} else {
		cnttype := tiketRes.Header.Get("Content-Type")
		// 做下图片判断
		if strings.Contains(cnttype, "image") {
			bodyString, _ := tiketRes.ToString()
			savePath := "upload/" + mediaid + ".jpg"

			err2 := ioutil.WriteFile(savePath, []byte(bodyString), 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
			if err2 != nil {
				fmt.Println(err2)
				return false, ""
			}

			return true, savePath
		} else {
			bodystr, _ := tiketRes.ToString()
			lego.LogError("微信图片下载失败：" + bodystr)
			return false, ""
		}
	}
}

/****************************************************************************/
/* 需要支付证书的接口，调用时的tsl配置
* params:
  ---
*/
func (self *WechatHelper) getTLSConfig(certPath, keyPath, wechatCAPath string) (*tls.Config, error) {
	if _tlsConfig != nil {
		return _tlsConfig, nil
	}

	// load cert
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		lego.LogError("load wechat keys fail:" + err.Error())
		return nil, err
	}

	// load root ca
	caData, err := ioutil.ReadFile(wechatCAPath)
	if err != nil {
		lego.LogError("read wechat ca fail:" + err.Error())
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	_tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	return _tlsConfig, nil
}

/* GenSign：
* 微信-请求时需要的sign签名计算
* params:
  ---
*/
func (self *WechatHelper) GenSign(wxkey string, params map[string]string) (rsp string) {
	sorted_keys := make([]string, 0)
	for k, _ := range params {
		sorted_keys = append(sorted_keys, k)
	}

	// sort 'string' key in increasing order
	sort.Strings(sorted_keys)

	ptext := ""
	for _, k := range sorted_keys {
		v := params[k]
		if len(v) > 0 {
			ptext = ptext + k + "=" + params[k] + "&"
		}
	}

	if len(ptext) > 0 {
		lidx := len(ptext) - 1
		ptext = ptext[0:lidx]
	}
	ptext = ptext + "&key=" + wxkey
	lego.LogInfo("sign ptext=" + ptext)
	sign := strings.ToUpper(utils.Md5(ptext))
	return sign
}

/* 发送模板消息
* params:
  ---
*/
func (self *WechatHelper) SendTemplateMessage(touser, template_id, page, form_id, accesstoken string, data map[string]interface{}) error {
	params := map[string]interface{}{
		"touser":           touser,
		"template_id":      template_id,
		"page":             page,
		"form_id":          form_id,
		"data":             data,
		"emphasis_keyword": "",
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=%s", accesstoken)
	header := map[string]string{
		"Content-Type": "application/json",
	}
	databyte, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body := strings.NewReader(string(databyte))
	rsp, err := httpclient.Do("POST", url, header, body)
	if err != nil {
		lego.LogError(err.Error())
		return err
	} else {
		body, err := rsp.ToString()
		if err != nil {
			lego.LogError(err.Error())
			return err
		}

		var wxrsp map[string]interface{}
		if err := json.Unmarshal([]byte(body), &wxrsp); err != nil {
			lego.LogError(err.Error())
			return err
		}
		if int64(wxrsp["errcode"].(float64)) == 0 {
			return nil
		} else {
			return errors.New(wxrsp["errmsg"].(string))
		}
	}
}


// 小程序，授权获取电话号码时，对返回的数据进行解密，获取电话号码
// 解密
func (self *WechatHelper) AesCBCDncrypt(encryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		//panic("ciphertext too short")
		return nil,errors.New("ciphertext too short")
	}
	if len(encryptData)%blockSize != 0 {
		//panic("ciphertext is not a multiple of the block size")
		return nil,errors.New("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	// 解填充
	encryptData = self.PKCS7UnPadding(encryptData)
	return encryptData, nil
}

//解密
/**
* rawData 原始加密数据
* key  密钥
* iv  向量
 */
func (self *WechatHelper) Dncrypt(rawData, key, iv string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	key_b, err_1 := base64.StdEncoding.DecodeString(key)
	iv_b, _ := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", err
	}
	if err_1 != nil {
		return "", err_1
	}
	dnData, err := self.AesCBCDncrypt(data, key_b, iv_b)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

//去除填充
func (self *WechatHelper) PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}