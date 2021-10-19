package model

import "time"

const CouponStatusOk = 1
const CouponStatusForbid = 2

type Coupon struct {
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
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

func (Coupon) TableName() string {
	return "coupon"
}
