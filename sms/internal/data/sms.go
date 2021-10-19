package data

import (
	"context"
	"github.com/ZQCard/kratos-service-base/sms/internal/biz"
	"github.com/ZQCard/kratos-service-base/sms/internal/data/model"
	"github.com/go-kratos/kratos/v2/log"
)

type smsDataRepo struct {
	data *Data
	log  *log.Helper
}

func (u smsDataRepo) SaveSmsData(ctx context.Context, sms *biz.Sms) error {
	return u.data.db.Model(&model.Sms{}).Create(&model.Sms{
		Id:         0,
		Mobile:     sms.Mobile,
		Content:    sms.Content,
		Type:       sms.Type,
		Scene:      sms.Scene,
		IsExpire:   sms.IsExpire,
		ExpireTime: sms.ExpireTime,
	}).Error
}

func (u smsDataRepo) ChangeSmsExpire(id int64) error {
	return u.data.db.Model(&model.Sms{}).Where("id = ?", id).Update("is_expire", model.SmsIsExpire).Error
}

func (u smsDataRepo) GetSmsData(ctx context.Context, params map[string]interface{}) *biz.Sms {
	var sms model.Sms
	db := u.data.db.Model(&model.Sms{})

	if mobile, ok := params["mobile"]; ok && mobile.(string) != "" {
		db = db.Where("mobile = ?", mobile)
	}

	if scene, ok := params["scene"]; ok && scene.(string) != "" {
		db = db.Where("scene = ?", scene)
	}

	if content, ok := params["content"]; ok && content.(string) != "" {
		db = db.Where("content = ?", content)
	}

	if t, ok := params["type"]; ok && t.(int64) != 0 {
		db = db.Where("type = ?", t)
	}

	if isExpire, ok := params["is_expire"]; ok && isExpire.(int64) != 0 {
		db = db.Where("is_expire = ?", isExpire)
	}

	if expireTime, ok := params["expire_time"]; ok && expireTime.(string) != "" {
		db = db.Where("expire_time >= ?", expireTime)
	}
	db.Order("id DESC").First(&sms)
	res := &biz.Sms{
		Id:         sms.Id,
		Mobile:     sms.Mobile,
		Content:    sms.Content,
		IsExpire:   sms.IsExpire,
		Type:       sms.Type,
		Scene:      sms.Scene,
		ExpireTime: sms.ExpireTime,
	}
	return res
}

func NewSmsRepo(data *Data, logger log.Logger) biz.SmsRepo {
	return smsDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
