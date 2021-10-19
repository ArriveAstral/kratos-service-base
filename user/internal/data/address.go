package data

import (
	"context"
	"errors"
	"github.com/ZQCard/kratos-service-base/user/internal/biz"
	"github.com/ZQCard/kratos-service-base/user/internal/data/model"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type addressDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewAddressRepo(data *Data, logger log.Logger) biz.AddressRepo {
	return &addressDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (a addressDataRepo) ListAddress(ctx context.Context, params map[string]interface{}, page, pageSize int) ([]*biz.Address, int64, error) {
	var list []model.Address
	db := a.data.db.Model(&model.Address{})

	if id, ok := params["id"]; ok && id.(int64) != 0 {
		db = db.Where("id = ?", id.(int64))
	}
	if userId, ok := params["userId"]; ok && userId.(int64) != 0 {
		db = db.Where("user_id = ?", userId.(int64))
	}
	if status, ok := params["status"]; ok && status.(int64) != 0 {
		db = db.Where("status = ?", status.(int64))
	}
	var count int64
	db.Count(&count)
	// 增加分页
	err := db.Scopes(model.Paginate(page, pageSize)).Find(&list).Error
	var resp []*biz.Address
	for _, address := range list {
		resp = append(resp, &biz.Address{
			Id:              address.Id,
			UserId:          address.UserId,
			Country:         address.Country,
			Province:        address.Province,
			City:            address.City,
			CountryDistrict: address.CountryDistrict,
			DetailedAddress: address.DetailedAddress,
			IsDefault:       address.IsDefault,
			Status:          address.Status,
		})
	}
	return resp, count, err
}

func (a addressDataRepo) GetAddress(ctx context.Context, params map[string]interface{}) (*biz.Address, error) {
	db := a.data.db.Model(&model.Address{})
	if id, ok := params["id"]; ok && id.(int64) != 0 {
		db = db.Where("id = ?", id.(int64))
	}
	if userId, ok := params["userId"]; ok && userId.(int64) != 0 {
		db = db.Where("user_id = ?", userId.(int64))
	}
	if status, ok := params["status"]; ok && status.(int64) != 0 {
		db = db.Where("status = ?", status.(int64))
	}
	var address model.Address
	if err := db.First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &biz.Address{}, nil
		}
		return nil, err
	}

	response := &biz.Address{
		Id:              address.Id,
		UserId:          address.UserId,
		Country:         address.Country,
		Province:        address.Province,
		City:            address.City,
		CountryDistrict: address.CountryDistrict,
		DetailedAddress: address.DetailedAddress,
		IsDefault:       address.IsDefault,
		Status:          address.Status,
	}
	return response, nil

}

func (a addressDataRepo) CreateAddress(ctx context.Context, address *biz.Address) error {
	// 如果该地址为默认地址，需要更新其他地址为非默认地址
	db := a.data.db.Model(&model.Address{}).Begin()
	if address.IsDefault == model.AddressDefaultOK {
		err := db.Where("user_id = ?", address.UserId).UpdateColumn("is_default", model.AddressDefaultNO).Error
		if err != nil {
			db.Rollback()
			return err
		}
	}

	err := db.Create(&model.Address{
		Id:              address.Id,
		UserId:          address.UserId,
		Country:         address.Country,
		Province:        address.Province,
		City:            address.City,
		CountryDistrict: address.CountryDistrict,
		DetailedAddress: address.DetailedAddress,
		IsDefault:       address.IsDefault,
		Status:          address.Status,
	}).Error
	if err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

func (a addressDataRepo) UpdateAddress(ctx context.Context, address *biz.Address) error {
	// 如果该地址为默认地址，需要更新其他地址为非默认地址
	db := a.data.db.Model(&model.Address{}).Begin()

	var record model.Address
	a.data.db.Model(&model.Address{}).Where("id = ?", address.Id).First(&record)
	record.Country = address.Country
	record.Province = address.Province
	record.City = address.City
	record.CountryDistrict = address.CountryDistrict
	record.DetailedAddress = address.DetailedAddress
	record.IsDefault = address.IsDefault

	if record.IsDefault == model.AddressDefaultNO && address.IsDefault == model.AddressDefaultOK {
		err := db.Where("user_id = ?", address.UserId).UpdateColumn("is_default", model.AddressDefaultNO).Error
		if err != nil {
			db.Rollback()
			return err
		}
	}
	err := db.Where("id = ?", address.Id).Save(&record).Error
	if err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

func (a addressDataRepo) DeleteAddress(ctx context.Context, id int64) error {
	return a.data.db.Model(&model.Address{}).Where("id = ?", id).Update("status", model.AddressStatusForbid).Error
}
