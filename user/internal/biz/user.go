package biz

import (
	"context"
	"github.com/ZQCard/kratos-service-base/user/internal/data/model"
	"github.com/ZQCard/kratos-service-base/user/internal/pkg/middleware/jwt"
	"github.com/ZQCard/kratos-service-base/user/internal/pkg/util/encryption"
	"github.com/ZQCard/kratos-service-base/user/internal/pkg/util/random"
	"github.com/ZQCard/kratos-service-base/user/internal/pkg/util/timeSugar"
	"github.com/ZQCard/kratos-service-base/user/internal/pkg/util/validator"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type User struct {
	Id          int64
	Username    string `validate:"required,min=6,max=30" label:"用户名"`
	Password    string `validate:"required,min=6,max=30" label:"密码"`
	Mobile      string `validate:"required,numeric,len=11" label:"手机号码"`
	Nickname    string `validate:"required,min=2,max=20" label:"昵称"`
	Avatar      string `validate:"required,url,max=100" label:"头像"`
	Salt        string
	LastLoginAt *time.Time
	Status      int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type UserRepo interface {
	RegisterUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, params map[string]interface{}) (*User, error)
	CheckUserLoginVerifyCode(ctx context.Context, mobile, code, scene string) (bool, error)
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, log: log.NewHelper(logger)}
}

// 用户注册业务逻辑
func (uc *UserUseCase) Register(ctx context.Context, user *User) error {
	// 参数验证
	err := validator.Validate(user)
	if err != nil {
		return errors.BadRequest(err.Error(), err.Error())
	}
	// 检查用户手机号是否已存在
	userData, err := uc.repo.GetUser(ctx, map[string]interface{}{
		"mobile": user.Mobile,
	})
	if err != nil || userData.Id != 0 {
		return errors.BadRequest(SystemBadRequestErrorMsg, "手机号已经存在")
	}
	// 检查用户手机号是否已存在
	userData, err = uc.repo.GetUser(ctx, map[string]interface{}{
		"mobile": user.Username,
	})
	if err != nil || userData.Id != 0 {
		return errors.BadRequest(SystemBadRequestErrorMsg, "用户名已经存在")
	}
	// 用户注册，生成随机盐
	user.Salt = random.RandString(10)
	// 对于用户密码进行加密
	user.Password = encryptUserPassword(user.Salt, user.Password)
	user.Status = model.UserStatusOk
	user.LastLoginAt = timeSugar.CurrentTimeYMDHISTime()
	// 用户入库
	if err := uc.repo.RegisterUser(ctx, user); err != nil {
		return err
	}

	// 用户注册成功, 向kafka发送消息，通过用户手机号码发送短信
	//partition, offset, err := myKafka.SendSyncMessage("user_register", "", user.Mobile, "", "", 0, 0, false)
	//
	//if err != nil {
	//	partition = partition
	//	offset = offset
	//	uc.log.Error("Register User Kafka Error：" + err.Error())
	//	return err
	//}
	return nil
}

func (uc *UserUseCase) LoginByUsername(ctx context.Context, user *User) (string, error) {
	userData, err := uc.repo.GetUser(ctx, map[string]interface{}{
		"username": user.Username,
	})
	if err != nil {
		// 记录日志
		uc.log.Error("LoginByUsername：" + err.Error())
		return "", errors.InternalServer(SystemInternalErrorMsg, SystemInternalErrorMsg)
	}
	if userData.Id == 0 {
		return "", errors.BadRequest(SystemBadRequestErrorMsg, "用户名或密码错误")
	}
	// 检测用户状态
	if userData.Status != model.UserStatusOk {
		return "", errors.InternalServer(SystemForbiddenErrorMsg, SystemForbiddenErrorMsg)
	}
	// 验证用户密码
	if !checkUserPassword(userData.Salt, userData.Password, user.Password) {
		return "", errors.BadRequest(SystemBadRequestErrorMsg, "密码错误")
	}
	token, err := generateUserToken(userData)
	if err != nil {
		// 记录日志
		uc.log.Error("LoginByUsername：" + err.Error())
		return "", errors.InternalServer(SystemInternalErrorMsg, SystemInternalErrorMsg)
	}
	return token, nil
}

func (uc *UserUseCase) LoginByVerifyCode(ctx context.Context, mobile, code string) (string, error) {
	userData, err := uc.repo.GetUser(ctx, map[string]interface{}{
		"mobile": mobile,
	})
	if err != nil {
		// 记录日志
		uc.log.Error("LoginByVerifyCode：" + err.Error())
		return "", errors.InternalServer(SystemInternalErrorMsg, SystemInternalErrorMsg)
	}
	if userData.Id == 0 {
		return "", errors.BadRequest(SystemBadRequestErrorMsg, "用户名或密码错误")
	}
	// 检测用户状态
	if userData.Status != model.UserStatusOk {
		return "", errors.InternalServer(SystemForbiddenErrorMsg, SystemForbiddenErrorMsg)
	}
	// 验证用户验证码
	isSuccess, err := uc.repo.CheckUserLoginVerifyCode(ctx, mobile, code, "login")
	if err != nil {
		return "", errors.InternalServer(SystemInternalErrorMsg, SystemInternalErrorMsg)
	}
	if !isSuccess {
		return "", errors.BadRequest(SystemBadRequestErrorMsg, "验证码错误")
	}
	token, err := generateUserToken(userData)
	if err != nil {
		// 记录日志
		uc.log.Error("LoginByVerifyCode：" + err.Error())
		return "", errors.InternalServer(SystemInternalErrorMsg, SystemInternalErrorMsg)
	}
	return token, nil
}

func (uc *UserUseCase) GetUserById(ctx context.Context, id int64) (*User, error) {
	user, err := uc.repo.GetUser(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		// 记录日志
		uc.log.Error("GetUserById：" + err.Error())
		return &User{}, errors.InternalServer(SystemInternalErrorMsg, SystemInternalErrorMsg)
	}
	return user, nil
}

func (uc *UserUseCase) CheckUserOK(ctx context.Context, id int64) (bool, error) {
	user, err := uc.repo.GetUser(ctx, map[string]interface{}{
		"id": id,
	})
	return user.Status == model.UserStatusOk, err
}

func generateUserToken(user *User) (string, error) {
	token, err := jwt.GenerateToken(user.Id, user.Username)
	if err != nil {
		return token, err
	}
	return token, nil
}

// 对于用户密码进行加密
func encryptUserPassword(salt, password string) string {
	return encryption.EncodeMD5(salt + password)
}

// 验证用户密码
func checkUserPassword(salt, password, originPassword string) bool {
	return encryption.EncodeMD5(salt+originPassword) == password
}
