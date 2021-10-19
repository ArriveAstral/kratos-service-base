package service

import (
	"context"
	"github.com/ZQCard/kratos-service-base/coupon/internal/biz"
	"github.com/go-kratos/kratos/v2/log"

	pb "github.com/ZQCard/kratos-service-base/api/coupon/v1"
)

type CouponService struct {
	pb.UnimplementedCouponServer
	cc  *biz.CouponUseCase
	log *log.Helper
}

func NewCouponService(cc *biz.CouponUseCase, logger log.Logger) *CouponService {
	return &CouponService{cc: cc, log: log.NewHelper(logger)}
}

func (s *CouponService) ListCoupon(ctx context.Context, req *pb.ListCouponRequest) (*pb.ListCouponReply, error) {
	var results []*pb.CouponInfo
	list, err := s.cc.List(ctx)
	for _, coupon := range list {
		results = append(results, &pb.CouponInfo{
			Id:           coupon.Id,
			Name:         coupon.Name,
			Type:         int64(coupon.Type),
			Money:        coupon.Money,
			LowerMoney:   coupon.LowerMoney,
			Status:       coupon.Status,
			TotalCount:   coupon.TotalCount,
			CollectCount: coupon.CollectCount,
			StartTime:    coupon.StartTime,
			EndTime:      coupon.EndTime,
			Limit:        coupon.Limit,
		})
	}
	return &pb.ListCouponReply{Results: results}, err
}
func (s *CouponService) RushCollectCoupon(ctx context.Context, req *pb.RushCollectCouponRequest) (*pb.RushCollectCouponReply, error) {
	// 获取用户id
	id := ctx.Value("x-md-global-uid")
	err := s.cc.Rush(ctx, id.(int64), req.CouponId)
	message := ""
	if err == nil {
		message = "领券成功"
	}
	return &pb.RushCollectCouponReply{Message: message}, err
}

func (s *CouponService) ListUserCoupon(ctx context.Context, req *pb.ListUserCouponRequest) (*pb.ListUserCouponReply, error) {
	// 获取用户id
	userId := ctx.Value("x-md-global-uid")
	var result []*pb.ListUserCouponReply_ListUserInfo
	params := make(map[string]interface{})
	params["status"] = req.Status
	params["type"] = req.Type
	params["userId"] = userId.(int64)
	list, count, err := s.cc.ListUserCoupon(ctx, params, req.Page, req.PageSize)
	for _, v := range list {
		result = append(result, &pb.ListUserCouponReply_ListUserInfo{
			Id:           v.Id,
			Name:         v.Name,
			Type:         v.Type,
			Money:        v.Money,
			LowerMoney:   v.LowerMoney,
			Status:       v.Status,
			CouponId:     v.CouponId,
			CouponStatus: v.CouponStatus,
			OrderId:      v.OrderId,
			UseTime:      v.UseTime,
			StartTime:    v.StartTime,
			EndTime:      v.EndTime,
			CreatedAt:    v.CreatedAt,
		})
	}
	return &pb.ListUserCouponReply{Results: result, Count: count}, err
}
