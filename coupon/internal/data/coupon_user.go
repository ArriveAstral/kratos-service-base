package data

import (
	"context"
	"errors"
	userpbv1 "github.com/ZQCard/kratos-service-base/api/user/v1"
	"github.com/ZQCard/kratos-service-base/coupon/internal/biz"
	"github.com/ZQCard/kratos-service-base/coupon/internal/data/model"
	"github.com/ZQCard/kratos-service-base/coupon/internal/pkg/util/timeSugar"
	"github.com/go-kratos/kratos/v2/metadata"
	"gorm.io/gorm"
	"time"
)

func (c couponDataRepo) CouponUserList(ctx context.Context, params map[string]interface{}, page, pageSize int) ([]*biz.CouponUserInfo, int64, error) {
	var list []model.CouponUser
	db := c.data.db.Model(&model.CouponUser{})
	if userId, ok := params["user_id"]; ok && userId.(int64) != 0 {
		db = db.Where("user_id = ?", userId)
	}
	if t, ok := params["type"]; ok && t.(int64) != 0 {
		db = db.Where("type = ?", t)
	}
	if status, ok := params["status"]; ok && status.(int64) != 0 {
		db = db.Where("status = ?", status)
	}
	var count int64
	db.Count(&count)
	db.Scopes(model.Paginate(page, pageSize)).Preload("Coupon").Find(&list)
	var resp []*biz.CouponUserInfo
	for _, v := range list {
		resp = append(resp, &biz.CouponUserInfo{
			Id:           v.Id,
			Name:         v.Coupon.Name,
			Type:         v.Coupon.Type,
			Money:        v.Coupon.Money,
			LowerMoney:   v.Coupon.LowerMoney,
			Status:       v.Status,
			CouponId:     v.CouponId,
			CouponStatus: v.Coupon.Status,
			OrderId:      v.OrderId,
			UseTime:      timeSugar.ChangeTimeToYMDStr(v.UseTime),
			StartTime:    timeSugar.ChangeTimeToYMDStr(v.StartTime),
			EndTime:      timeSugar.ChangeTimeToYMDStr(v.EndTime),
			CreatedAt:    timeSugar.ChangeTimeToYMDStr(v.CreatedAt),
		})
	}
	return resp, count, nil
}

func (c couponDataRepo) RushCoupon(ctx context.Context, userId, couponId int64) error {
	// 查看用户是否有抢券资格,如果冻结,则无法抢券
	// metadata传递信息
	authorizationStr := ctx.Value("Authorization")
	ctx = metadata.AppendToClientContext(ctx, "Authorization", authorizationStr.(string))
	reply, err := c.data.userClient.CheckUserOK(ctx, &userpbv1.CheckUserOKRequest{})
	if err != nil {
		return err
	}
	if reply.IsOk == false {
		return errors.New("用户已冻结")
	}
	couponUser := model.CouponUser{}
	couponUser.UserId = userId
	couponUser.CouponId = couponId
	couponUser.Status = model.CouponStatusClaimed
	now := time.Now()
	couponUser.StartTime = &now
	end := time.Now().AddDate(0, 0, 7)
	couponUser.EndTime = &end
	db := c.data.db.Begin()
	// 保存用户优惠券信息
	err = db.Model(&model.CouponUser{}).Create(&couponUser).Error
	if err != nil {
		db.Rollback()
		return errors.New("服务器繁忙,请稍后重试")
	}
	// 更新优惠券已领张数
	affect := db.Model(&model.Coupon{}).Where("id = ? AND total_count > collect_count", couponId).UpdateColumn("collect_count", gorm.Expr("collect_count + ?", 1)).RowsAffected
	if affect == 0 {
		return errors.New("优惠券保存失败")
	}
	db.Commit()
	return nil
}

// 查询用户优惠券拥有数量
func (c couponDataRepo) GetCouponUserCount(ctx context.Context, params map[string]interface{}) (int64, error) {
	var count int64
	db := c.data.db.Model(&model.CouponUser{})
	if status, ok := params["status"]; ok && status.(int) != 0 {
		db = db.Where("status = ?", status)
	}
	if id, ok := params["id"]; ok && id.(int64) != 0 {
		db = db.Where("id = ?", id)
	}
	if userId, ok := params["user_id"]; ok && userId.(int64) != 0 {
		db = db.Where("user_id = ?", userId)
	}
	if couponId, ok := params["coupon_id"]; ok && couponId.(int64) != 0 {
		db = db.Where("coupon_id = ?", couponId)
	}
	err := db.Count(&count).Error
	return count, err
}
