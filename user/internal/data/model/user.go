package model

import "time"

const UserStatusOk = 1
const UserStatusForbid = 2

type User struct {
	Id          int64
	Username    string
	Salt        string
	Password    string
	Mobile      string
	Nickname    string
	Avatar      string
	Status      int64
	LastLoginAt *time.Time
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (User) TableName() string {
	return "users"
}
