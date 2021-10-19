package model

import "time"

const SmsIsExpire = int64(1)
const SmsIsExpireNO = int64(2)

const SmsTypeVerifyCode = int64(1)

type Sms struct {
	Id         int64
	Mobile     string
	Content    string
	Type       int64
	Scene      string
	IsExpire   int64
	ExpireTime *time.Time
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func (Sms) TableName() string {
	return "sms"
}
