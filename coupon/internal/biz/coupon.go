package biz

import (
	"github.com/ZQCard/kratos-service-base/coupon/internal/data/model"
	"github.com/ZQCard/kratos-service-base/coupon/internal/pkg/util/timeSugar"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/net/context"
)

type CouponInfo struct {
	Id           int64
	Name         string
	Type         int64
	Money        float32
	LowerMoney   float32
	Status       int64
	TotalCount   int64
	CollectCount int64
	StartTime    string
	EndTime      string
	Limit        int64
}

type CouponUserInfo struct {
	Id           int64
	Name         string
	Type         int64
	Money        float32
	LowerMoney   float32
	Status       int64
	CouponId     int64
	CouponStatus int64
	OrderId      int64
	UseTime      string
	StartTime    string
	EndTime      string
	CreatedAt    string
}

type CouponRepo interface {
	CouponList(ctx context.Context, params map[string]interface{}) ([]*CouponInfo, error)
	GetCoupon(ctx context.Context, params map[string]interface{}) (*CouponInfo, error)
	GetCouponUserCount(ctx context.Context, params map[string]interface{}) (int64, error)
	RushCoupon(ctx context.Context, userId, CouponId int64) error
	CouponUserList(ctx context.Context, params map[string]interface{}, page, pageSize int) ([]*CouponUserInfo, int64, error)
}

type CouponUseCase struct {
	repo CouponRepo
	log  *log.Helper
}

func NewCouponUseCase(repo CouponRepo, logger log.Logger) *CouponUseCase {
	return &CouponUseCase{repo: repo, log: log.NewHelper(logger)}
}

// 优惠券列表
func (c *CouponUseCase) List(ctx context.Context) ([]*CouponInfo, error) {
	// 选择可以展示的优惠券
	params := make(map[string]interface{})
	params["status"] = model.CouponStatusOk
	return c.repo.CouponList(ctx, params)
}

// 用户优惠券列表
func (c *CouponUseCase) ListUserCoupon(ctx context.Context, params map[string]interface{}, page, pageSize int64) ([]*CouponUserInfo, int64, error) {
	return c.repo.CouponUserList(ctx, params, int(page), int(pageSize))
}

// 抢券
func (c *CouponUseCase) Rush(ctx context.Context, userId, couponId int64) error {
	// 用户抢券
	params := make(map[string]interface{})
	params["id"] = couponId
	coupon, err := c.repo.GetCoupon(ctx, params)
	if err != nil {
		return err
	}

	// 是否禁用
	if coupon.Status == model.CouponStatusForbid {
		return errors.BadRequest(SystemBadRequestErrorMsg, "优惠券已禁用")
	}

	// 是否充足
	if coupon.TotalCount < coupon.CollectCount+1 {
		return errors.BadRequest(SystemBadRequestErrorMsg, "优惠券已抢完")
	}

	// 是否在可抢时间内
	if timeSugar.CurrentTimeHI() < coupon.StartTime || timeSugar.CurrentTimeHI() > coupon.EndTime {
		return errors.BadRequest(SystemBadRequestErrorMsg, "优惠券当前时间无法领取")
	}

	// 查看这个人是否抢到优惠券数量上限
	couponUserParams := make(map[string]interface{})
	couponUserParams["coupon_id"] = couponId
	couponUserParams["user_id"] = userId
	couponUserCount, err := c.repo.GetCouponUserCount(ctx, couponUserParams)
	if err != nil {
		return err
	}

	if couponUserCount >= int64(coupon.Limit) {
		return errors.BadRequest(SystemBadRequestErrorMsg, "抢券数量达到上限")
	}
	// 保存用户优惠券信息
	err = c.repo.RushCoupon(ctx, userId, couponId)
	if err != nil {
		c.log.Error("RushCoupon：" + err.Error())

		return errors.InternalServer(SystemInternalErrorMsg, "用户领券失败")
	}
	return nil
}
