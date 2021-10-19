package aliSms

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

var RegionId string
var AccessKey string
var AccessSecret string
var SignName string
var VerifyCodeTemplate string


// 阿里云短信发送
func SendAliSmsVerifyCode(mobile string, code string) (err error) {
	data, _ := json.Marshal(map[string]string{"code": code})

	client, err := dysmsapi.NewClientWithAccessKey(RegionId, AccessKey, AccessSecret)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = mobile
	request.SignName = SignName
	request.TemplateCode = VerifyCodeTemplate
	request.TemplateParam = string(data)
	response, err := client.SendSms(request)
	if err != nil {
		return errors.New("阿里云短信发送失败" + err.Error())
	}

	if response.Code != "OK" {
		if response.Code == "isv.BUSINESS_LIMIT_CONTROL" {
			return errors.New("获取短信过于频繁，请稍后再试")
		}
		return errors.New("阿里云短信发送失败" + response.Message)
	}
	return nil
}
