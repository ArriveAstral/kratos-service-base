package model

import "time"

const AddressStatusOk = int64(1)
const AddressStatusForbid = int64(2)

const AddressDefaultOK = int64(1)
const AddressDefaultNO = int64(2)

type Address struct {
	Id              int64
	UserId          int64
	Country         string
	Province        string
	City            string
	CountryDistrict string
	DetailedAddress string
	IsDefault       int64
	Status          int64
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func (Address) TableName() string {
	return "user_address"
}
