package data

import (
	"context"
	"github.com/ZQCard/kratos-service-base/coupon/internal/biz"
	"github.com/ZQCard/kratos-service-base/coupon/internal/data/model"
	"github.com/go-kratos/kratos/v2/log"
)

type couponDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewCouponRepo(data *Data, logger log.Logger) biz.CouponRepo {
	return couponDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (c couponDataRepo) CouponList(ctx context.Context, params map[string]interface{}) ([]*biz.CouponInfo, error) {
	var list []model.Coupon
	db := c.data.db.Model(&model.Coupon{})
	if status, ok := params["status"]; ok && status.(int) != 0 {
		db = db.Where("status = ?", status)
	}
	err := db.Find(&list).Error
	var resp []*biz.CouponInfo
	for _, coupon := range list {
		resp = append(resp, &biz.CouponInfo{
			Id:           coupon.Id,
			Name:         coupon.Name,
			Type:         coupon.Type,
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
	return resp, err
}

func (c couponDataRepo) GetCoupon(ctx context.Context, params map[string]interface{}) (*biz.CouponInfo, error) {
	var couponInfo biz.CouponInfo
	db := c.data.db.Model(&model.Coupon{})
	if status, ok := params["status"]; ok && status.(int) != 0 {
		db = db.Where("status = ?", status)
	}
	if id, ok := params["id"]; ok && id.(int64) != 0 {
		db = db.Where("id = ?", id)
	}
	var record model.Coupon
	err := db.First(&record).Error
	couponInfo.Id = record.Id
	couponInfo.Name = record.Name
	couponInfo.Type = record.Type
	couponInfo.Money = record.Money
	couponInfo.LowerMoney = record.LowerMoney
	couponInfo.Status = record.Status
	couponInfo.TotalCount = record.TotalCount
	couponInfo.CollectCount = record.CollectCount
	couponInfo.StartTime = record.StartTime
	couponInfo.EndTime = record.EndTime
	couponInfo.Limit = record.Limit
	return &couponInfo, err
}
