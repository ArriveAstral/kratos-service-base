package service

import (
	"context"
	"github.com/ZQCard/kratos-service-base/sms/internal/biz"
	"github.com/go-kratos/kratos/v2/log"

	pb "github.com/ZQCard/kratos-service-base/api/sms/v1"
)

type SmsService struct {
	pb.UnimplementedSmsServer
	uc  *biz.SmsUseCase
	log *log.Helper
}

func NewSmsService(uc *biz.SmsUseCase, logger log.Logger) *SmsService {
	return &SmsService{uc: uc, log: log.NewHelper(logger)}
}

func (s *SmsService) SendSmsCode(ctx context.Context, req *pb.SendSmsCodeRequest) (*pb.SendSmsCodeReply, error) {
	isSuccess := true
	err := s.uc.SendCode(ctx, req.Mobile, req.Scene)
	if err != nil {
		isSuccess = false
	}
	return &pb.SendSmsCodeReply{IsSuccess: isSuccess}, err
}

func (s *SmsService) VerifySmsCode(ctx context.Context, req *pb.VerifySmsCodeRequest) (*pb.VerifySmsCodeReply, error) {
	isSuccess, err := s.uc.VerifyCode(ctx, req.Mobile, req.Scene, req.Code)
	return &pb.VerifySmsCodeReply{IsSuccess: isSuccess}, err
}
