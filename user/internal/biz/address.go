package biz

import (
	"context"
	"github.com/ZQCard/kratos-service-base/user/internal/data/model"
	"github.com/ZQCard/kratos-service-base/user/internal/pkg/util/validator"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type Address struct {
	Id              int64
	UserId          int64
	Country         string `validate:"required,min=1,max=50" label:"国家"`
	Province        string `validate:"required,min=1,max=50" label:"省份"`
	City            string `validate:"required,min=1,max=50" label:"城市"`
	CountryDistrict string `validate:"required,min=1,max=50" label:"县/区"`
	DetailedAddress string `validate:"required,min=1,max=50" label:"详细地址"`
	IsDefault       int64  `validate:"oneof=1 2" label:"是否默认"`
	Status          int64  `validate:"oneof=1 2" label:"状态"`
}

type AddressRepo interface {
	ListAddress(ctx context.Context, params map[string]interface{}, page, pageSize int) ([]*Address, int64, error)
	GetAddress(ctx context.Context, params map[string]interface{}) (*Address, error)
	CreateAddress(ctx context.Context, addr *Address) error
	UpdateAddress(ctx context.Context, addr *Address) error
	DeleteAddress(ctx context.Context, id int64) error
}

type AddressUseCase struct {
	repo AddressRepo
	log  *log.Helper
}

func NewAddressUseCase(repo AddressRepo, logger log.Logger) *AddressUseCase {
	return &AddressUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *AddressUseCase) List(ctx context.Context, userId int64, page, pageSize int64) ([]*Address, int64, error) {
	params := make(map[string]interface{})

	params["page"] = page
	params["pageSize"] = pageSize
	params["userId"] = userId
	params["status"] = model.AddressStatusOk
	return uc.repo.ListAddress(ctx, params, int(page), int(pageSize))
}

func (uc *AddressUseCase) Get(ctx context.Context, id, userId int64) (*Address, error) {
	params := make(map[string]interface{})
	params["id"] = id
	params["userId"] = userId
	params["status"] = model.AddressStatusOk
	address, err := uc.repo.GetAddress(ctx, params)
	if err != nil {
		// 记录日志
		uc.log.Error("GetAddress：" + err.Error())
		return nil, err
	}
	if address.Id == 0 {
		err = errors.NotFound(SystemNotfoundErrorMsg, SystemNotfoundErrorMsg)
	}
	return address, err
}

func (uc *AddressUseCase) Create(ctx context.Context, addr *Address) error {
	// 参数验证
	err := validator.Validate(addr)
	if err != nil {
		return errors.BadRequest(err.Error(), err.Error())
	}
	err = uc.repo.CreateAddress(ctx, addr)
	if err != nil {
		// 记录日志
		uc.log.Error("CreateAddress：" + err.Error())
		return errors.InternalServer(SystemInternalErrorMsg, SystemInternalErrorMsg)
	}
	return nil
}

func (uc *AddressUseCase) Update(ctx context.Context, userId int64, addr *Address) error {
	// 参数验证
	err := validator.Validate(addr)
	if err != nil {
		return errors.BadRequest(err.Error(), err.Error())
	}
	// 查询用户地址，查看是否有效
	if userId != addr.UserId {
		return errors.InternalServer(SystemBadRequestErrorMsg, SystemBadRequestErrorMsg)
	}
	address, err := uc.Get(ctx, addr.Id, userId)
	if err != nil {
		// 记录日志
		uc.log.Error("GetAddress：" + err.Error())
		return err
	}
	// 修改必要信息
	address.Country = addr.Country
	address.Province = addr.Province
	address.City = addr.City
	address.CountryDistrict = addr.CountryDistrict
	address.DetailedAddress = addr.DetailedAddress
	address.IsDefault = addr.IsDefault
	return uc.repo.UpdateAddress(ctx, address)
}

func (uc *AddressUseCase) Delete(ctx context.Context, id, userId int64) error {
	_, err := uc.Get(ctx, id, userId)
	if err != nil {
		// 记录日志
		uc.log.Error("GetAddress：" + err.Error())
		return err
	}
	return uc.repo.DeleteAddress(ctx, id)
}
