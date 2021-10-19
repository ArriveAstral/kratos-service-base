package model

import "time"

const CouponStatusClaimed = 1
const CouponStatusUsed = 2

type CouponUser struct {
	Id        int64
	CouponId  int64
	UserId    int64
	OrderId   int64
	Status    int64
	StartTime *time.Time
	EndTime   *time.Time
	UseTime   *time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Coupon    Coupon `gorm:"FOREIGNKEY:coupon_id;ASSOCIATION_FOREIGNKEY:ID"`
}

func (CouponUser) TableName() string {
	return "coupon_user"
}
