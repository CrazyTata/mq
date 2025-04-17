package token

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"mq/common/util"
)

type Request struct {
	Version             string `json:"version"`
	MsgId               string `json:"msgid"`
	SystemTime          string `json:"systemtime"`
	StrictCheck         string `json:"strictcheck"`
	Appid               string `json:"appid"`
	ExpandParams        string `json:"expandparams"`
	Token               string `json:"token"`
	Sign                string `json:"sign"`
	EncryptionAlgorithm string `json:"encryptionalgorithm"`
}

type Response struct {
	Inresponseto string `json:"inresponseto"`
	Systemtime   string `json:"systemtime"`
	ResultCode   string `json:"resultCode"`
	Msisdn       string `json:"msisdn"`
	TaskId       string `json:"taskId"`
	Description  string `json:"desc"`
	// 如果有其他响应参数，请在此添加
}

// Validate 获取手机号码接口
func (t *Token) Validate(msgId, token string) (mobile, outRelation string, err error) {

	systemTime := time.Now().Format("20060102150405000")
	// 请求参数
	req := Request{
		Version:     t.version,
		MsgId:       msgId,
		SystemTime:  systemTime,
		StrictCheck: t.strictCheck,
		Appid:       t.appid,
		Token:       token,
		Sign:        t.sign(msgId, systemTime, token),
	}
	fmt.Println(req)
	// 响应解析变量
	var resp Response

	// 创建请求的JSON数据
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Error marshaling request: %v", err)
		return
	}
	headerData := map[string]string{
		"Content-Type": "application/json",
	}

	body, err := util.Post(t.url, jsonData, headerData)
	if err != nil {
		return
	}
	// 解析响应
	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.logger.Errorf("tokenValidateError err:%s", err.Error())
		return
	}
	//103000 返回成功
	//103101 签名错误
	//103113 token 格式错误
	//103119 appid 不存在
	//103133 sourceid 不合法（服务端需要使用调用 SDK 时使用的appid 去换取号码）
	//103211 其他错误
	//103412 无效的请求
	//103414 参数校验异常
	//103511 请求 ip 不在社区配置的服务器白名单内
	//103811 token 为空
	//104201 token 失效或不存在
	//105018 用户权限不足（使用了本机号码校验的 token 去调用本接口）105019 应用未授权（开发者社区未勾选能力）
	//105312 套餐已用完
	//105313 非法请求
	t.logger.Infof("tokenValidate request %+v", req)

	t.logger.Infof("tokenValidate response inresponseto:%s,systemtime:%s,resultCode:%s,msisdn:%s,taskId:%s,desc:%s", resp.Inresponseto, resp.Systemtime, resp.ResultCode, resp.Msisdn, resp.TaskId, resp.Description)

	if resp.ResultCode != "103000" {
		t.logger.Errorf("tokenValidateError error inresponseto:%s,systemtime:%s,resultCode:%s,msisdn:%s,taskId:%s,desc:%s", resp.Inresponseto, resp.Systemtime, resp.ResultCode, resp.Msisdn, resp.TaskId, resp.Description)
		return
	}
	mobile = resp.Msisdn
	outRelation = resp.Inresponseto
	return
}

// sign 签名
func (t *Token) sign(msgId, systemTime, token string) string {
	// 将所有参数组合在一起
	params := strings.Join([]string{t.appid, t.version, msgId, systemTime, t.strictCheck, token, t.appSecret}, "")

	hash := md5.New()
	hash.Write([]byte(params))
	md5Sum := hash.Sum(nil)

	// 将MD5哈希值转换为十六进制字符串
	return strings.ToUpper(hex.EncodeToString(md5Sum))
}
