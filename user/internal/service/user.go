package service

import (
	"context"
	pb "github.com/ZQCard/kratos-service-base/api/user/v1"
	"github.com/ZQCard/kratos-service-base/user/internal/biz"
	log "github.com/go-kratos/kratos/v2/log"
)

type UserService struct {
	pb.UnimplementedUserServer
	uc  *biz.UserUseCase
	ac  *biz.AddressUseCase
	log *log.Helper
}

func NewUserService(uc *biz.UserUseCase, ac *biz.AddressUseCase, logger log.Logger) *UserService {
	return &UserService{uc: uc, ac: ac, log: log.NewHelper(logger)}
}

// 用户注册
func (s *UserService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserReply, error) {
	response := &pb.RegisterUserReply{}

	userRegisterInfo := &biz.User{}
	userRegisterInfo.Username = req.Username
	userRegisterInfo.Password = req.Password
	userRegisterInfo.Mobile = req.Mobile
	userRegisterInfo.Nickname = req.Nickname
	userRegisterInfo.Avatar = req.Avatar
	err := s.uc.Register(ctx, userRegisterInfo)
	response.Message = "注册成功"
	return response, err
}

func (s *UserService) LoginByUsername(ctx context.Context, req *pb.LoginByUsernameRequest) (*pb.LoginUserReply, error) {
	token, err := s.uc.LoginByUsername(ctx, &biz.User{
		Username: req.Username,
		Password: req.Password,
	})
	return &pb.LoginUserReply{Token: token}, err
}

func (s *UserService) LoginByVerifyCode(ctx context.Context, req *pb.LoginByVerifyCodeRequest) (*pb.LoginUserReply, error) {
	token, err := s.uc.LoginByVerifyCode(ctx, req.Mobile, req.Code)
	return &pb.LoginUserReply{Token: token}, err
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	id := ctx.Value("x-md-global-uid")
	user, err := s.uc.GetUserById(ctx, id.(int64))
	if err != nil {
		return &pb.GetUserReply{}, err
	}

	return &pb.GetUserReply{User: &pb.UserInfo{
		Id:       user.Id,
		Username: user.Username,
		Password: "",
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}}, err
}

func (s *UserService) CheckUserOK(ctx context.Context, req *pb.CheckUserOKRequest) (*pb.CheckUserOKReply, error) {
	id := ctx.Value("x-md-global-uid")
	ok, err := s.uc.CheckUserOK(ctx, id.(int64))
	return &pb.CheckUserOKReply{
		IsOk: ok,
	}, err
}
