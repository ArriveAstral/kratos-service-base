package data

import (
	"context"
	"errors"
	smsv1 "github.com/ZQCard/kratos-service-base/api/sms/v1"
	"github.com/ZQCard/kratos-service-base/user/internal/data/model"
	"gorm.io/gorm"

	"github.com/ZQCard/kratos-service-base/user/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type userDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return userDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (u userDataRepo) RegisterUser(ctx context.Context, user *biz.User) error {
	return u.data.db.Model(&model.User{}).Create(&model.User{
		Id:          0,
		Username:    user.Username,
		Password:    user.Password,
		Mobile:      user.Mobile,
		Salt:        user.Salt,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
	}).Error
}

func (u userDataRepo) GetUser(ctx context.Context, params map[string]interface{}) (*biz.User, error) {
	db := u.data.db.Model(&model.User{})
	if id, ok := params["id"]; ok && id.(int64) != 0 {
		db = db.Where("id = ?", id.(int64))
	}
	if username, ok := params["username"]; ok && username.(string) != "" {
		db = db.Where("username = ?", username.(string))
	}
	if mobile, ok := params["mobile"]; ok && mobile.(string) != "" {
		db = db.Where("mobile = ?", mobile.(string))
	}
	var user model.User
	if err := db.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &biz.User{}, nil
		}
		return nil, err
	}

	response := &biz.User{}
	response.Id = user.Id
	response.Username = user.Username
	response.Password = user.Password
	response.Salt = user.Salt
	response.Mobile = user.Mobile
	response.Nickname = user.Nickname
	response.Avatar = user.Avatar
	response.Status = user.Status
	response.LastLoginAt = user.LastLoginAt
	response.CreatedAt = user.CreatedAt
	response.UpdatedAt = user.UpdatedAt
	return response, nil
}

// 调用短信验证码,查看是否正确
func (u userDataRepo) CheckUserLoginVerifyCode(ctx context.Context, mobile, code, scene string) (bool, error) {
	resp, err := u.data.smsClient.VerifySmsCode(ctx, &smsv1.VerifySmsCodeRequest{
		Mobile: mobile,
		Scene:  scene,
		Code:   code,
	})

	if err != nil {
		return false, err
	}
	return resp.IsSuccess, err
}
