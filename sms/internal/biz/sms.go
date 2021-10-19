package biz

import (
	"context"
	"fmt"
	"github.com/ZQCard/kratos-service-base/sms/internal/data/model"
	aliSms "github.com/ZQCard/kratos-service-base/sms/internal/pkg/sms"
	"github.com/ZQCard/kratos-service-base/sms/internal/pkg/util/random"
	"github.com/ZQCard/kratos-service-base/sms/internal/pkg/util/timeSugar"
	"github.com/ZQCard/kratos-service-base/sms/internal/pkg/util/validator"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type Sms struct {
	Id         int64
	Mobile     string
	Content    string
	IsExpire   int64
	Type       int64
	Scene      string
	ExpireTime *time.Time
}

type SmsRepo interface {
	SaveSmsData(ctx context.Context, sms *Sms) error
	GetSmsData(ctx context.Context, params map[string]interface{}) *Sms
	ChangeSmsExpire(id int64) error
}

type SmsUseCase struct {
	repo SmsRepo
	log  *log.Helper
}

func NewSmsUseCase(repo SmsRepo, logger log.Logger) *SmsUseCase {
	return &SmsUseCase{repo: repo, log: log.NewHelper(logger)}
}

// 用户注册业务逻辑
type sendCode struct {
	Mobile string `validate:"numeric,len=11" label:"手机号码"`
	Scene  string `validate:"required" label:"应用场景"`
}

func (uc *SmsUseCase) SendCode(ctx context.Context, mobile, scene string) error {
	var send sendCode
	send.Mobile = mobile
	send.Scene = scene
	// 参数验证
	err := validator.Validate(send)
	if err != nil {
		return errors.BadRequest(err.Error(), err.Error())
	}
	// 生成验证码
	code := random.GenerateNumber(6)
	var sms Sms
	sms.Scene = scene
	sms.Mobile = mobile
	sms.Content = code
	sms.IsExpire = model.SmsIsExpireNO
	sms.Type = model.SmsTypeVerifyCode
	// 验证码时长5分钟
	fiveMinuteAfter := time.Now().Add(time.Duration(5) * time.Minute)
	sms.ExpireTime = &fiveMinuteAfter
	if err := uc.repo.SaveSmsData(ctx, &sms); err != nil {
		// 记录日志
		uc.log.Error("SendCode：" + err.Error())
		return errors.InternalServer(err.Error(), err.Error())
	}
	fmt.Println("发送短信验证码")

	// 发送短信验证码
	if err := aliSms.SendAliSmsVerifyCode(mobile, code); err != nil {

		// 记录日志
		uc.log.Error("SendCode：" + err.Error())
		return errors.InternalServer(err.Error(), err.Error())
	}
	fmt.Println("发送短信验证码成功")

	return nil
}

func (uc *SmsUseCase) VerifyCode(ctx context.Context, mobile, scene, code string) (bool, error) {
	var send sendCode
	send.Mobile = mobile
	send.Scene = scene
	// 参数验证
	err := validator.Validate(send)
	if err != nil {
		return false, errors.BadRequest(err.Error(), err.Error())
	}
	params := make(map[string]interface{})
	params["mobile"] = mobile
	params["scene"] = scene
	params["content"] = code
	params["type"] = model.SmsTypeVerifyCode
	params["expire_time"] = timeSugar.CurrentTimeYMDHIS()

	smsInfo := uc.repo.GetSmsData(ctx, params)
	if smsInfo.Id != 0 && smsInfo.IsExpire == model.SmsIsExpireNO {
		// 验证码成功,更改状态
		if err := uc.repo.ChangeSmsExpire(smsInfo.Id); err != nil {
			// 记录日志
			uc.log.Error("VerifyCode：" + err.Error())
		}
		return true, nil
	}
	return false, nil
}
