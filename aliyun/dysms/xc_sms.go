package dysms

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/GiterLab/urllib"
	"github.com/tobyzxj/uuid"
)

// setACLClient 配置默认的服务权限信息
func setACLClient(accessid, accesskey string) (cli *Client) {
	cli = new(Client)
	cli.SetAccessID(accessid)
	cli.SetAccessKey(accesskey)

	cli.SetVersion("2017-05-25")
	cli.SetRegion("cn-hangzhou")
	cli.SetEndPoint("http://dysmsapi.aliyuncs.com/")

	if urllib.GetDefaultSetting().Transport == nil {
		// set default setting for urllib
		trans := &http.Transport{
			MaxIdleConnsPerHost: 500,
			Dial: (&net.Dialer{
				Timeout: time.Duration(15) * time.Second,
			}).Dial,
		}

		urlSetting := urllib.HttpSettings{
			ShowDebug:        false,            // ShowDebug
			UserAgent:        "XCTHINGS",       // UserAgent
			ConnectTimeout:   15 * time.Second, // ConnectTimeout
			ReadWriteTimeout: 30 * time.Second, // ReadWriteTimeout
			TlsClientConfig:  nil,              // TlsClientConfig
			Proxy:            nil,              // Proxy
			Transport:        trans,            // Transport
			EnableCookie:     false,            // EnableCookie
			Gzip:             true,             // Gzip
			DumpBody:         true,             // DumpBody
		}
		if cli.SocketTimeout != 0 {
			urlSetting.ConnectTimeout = time.Duration(cli.SocketTimeout) * time.Second
			urlSetting.ReadWriteTimeout = time.Duration(cli.SocketTimeout) * time.Second
		}
		if HTTPDebugEnable {
			urlSetting.ShowDebug = true
		} else {
			urlSetting.ShowDebug = false
		}
		urllib.SetDefaultSetting(urlSetting)
	}

	return
}

// sendSMS 发送短信消息
// SendSms 发送短信接口
// businessID 设置业务请求流水号，必填。
// phoneNumbers 短信发送的号码列表，必填。 多手机号使用,分割
// signName 短信签名
// templateCode 申请的短信模板编码,必填
// templateParam 短信模板变量参数
func (c *Client) sendSMS(mobile, signName, tplID, tplParam string) (r *SendSmsResponse, businessID string, err error) {
	req := &Request{cli: c, Param: make(map[string]string)}
	// 1. 系统参数
	req.Put("SignatureMethod", "HMAC-SHA1")
	req.Put("SignatureNonce", uuid.New())
	req.Put("AccessKeyId", c.AccessID)
	req.Put("SignatureVersion", "1.0")
	req.Put("Timestamp", time.Now().UTC().Format(time.RFC3339))
	req.Put("Format", "JSON")

	// 2. 业务API参数
	req.Put("Version", c.Version)
	req.Put("RegionId", c.Region)
	// req.Put("PhoneNumbers", "your_phonenumbers")
	// req.Put("SignName", "your_signname")
	// req.Put("TemplateParam", "your_ParamString")
	// req.Put("TemplateCode", "your_templatecode")
	// req.Put("OutId", "your_outid")

	req.Put("Version", "2017-05-25")
	req.Put("Action", "SendSms")

	businessID = uuid.New()
	sr := &SendSmsRequest{Request: req}
	sr.SetOutID(businessID)    // 设置业务请求流水号，必填
	sr.SetSignName(signName)   // 短信签名
	sr.SetPhoneNumbers(mobile) // 短信发送的号码列表，必填。
	sr.SetTemplateCode(tplID)  // 短信模板
	if tplParam != "" {
		sr.SetTemplateParam(tplParam) // 短信模板参数
	}
	r, err = sr.DoActionWithException()

	return
}

// QuerySendDetails 短信发送记录查询接口
// bizID 可选 - 流水号
// phoneNumber 查询的手机号码
// pageSize 必填 - 页大小
// currentPage 必填 - 当前页码从1开始计数
// sendDate 必填 - 发送日期 支持30天内记录查询，格式yyyyMMdd
func (c *Client) QuerySendDetails(bizID, phoneNumber, pageSize, currentPage, sendDate string) *QuerySendDetailsRequest {
	req := newRequset(c)
	req.Put("Version", "2017-05-25")
	req.Put("Action", "QuerySendDetails")

	r := &QuerySendDetailsRequest{Request: req}
	r.SetPhoneNumber(phoneNumber) // 查询的手机号码
	if bizID != "" {
		r.SetBizID(bizID) // 可选 - 流水号
	}
	r.SetSendDate(sendDate)       // 必填 - 发送日期 支持30天内记录查询，格式yyyyMMdd
	r.SetCurrentPage(currentPage) // 必填 - 当前页码从1开始计数
	r.SetPageSize(pageSize)       // 必填 - 页大小
	return r
}

// AliSMS aliSMS
func AliSMS(aid, ak, mobile, sign, tplid, code string) (r *SendSmsResponse, businessID string, err error) {
	cli := setACLClient(aid, ak)
	rbody := fmt.Sprintf("{\"code\":\"%s\"}", code)
	r, businessID, err = cli.sendSMS(mobile, sign, tplid, rbody)
	return
}
