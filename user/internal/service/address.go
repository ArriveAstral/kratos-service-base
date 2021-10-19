package service

import (
	"context"
	pb "github.com/ZQCard/kratos-service-base/api/user/v1"
	"github.com/ZQCard/kratos-service-base/user/internal/biz"
)

func (s *UserService) ListAddress(ctx context.Context, req *pb.ListAddressRequest) (*pb.ListAddressReply, error) {
	// 获取用户id
	userId := ctx.Value("x-md-global-uid")
	var addressList []*pb.AddressInfo
	list, count, err := s.ac.List(ctx, userId.(int64), req.Page, req.PageSize)
	for _, address := range list {
		addressList = append(addressList, &pb.AddressInfo{
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
	return &pb.ListAddressReply{Results: addressList, Count: count}, err
}
func (s *UserService) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressReply, error) {
	// 获取用户id
	userId := ctx.Value("x-md-global-uid")
	address := &biz.Address{}
	address.UserId = userId.(int64)
	address.Country = req.Country
	address.Province = req.Province
	address.City = req.City
	address.CountryDistrict = req.CountryDistrict
	address.DetailedAddress = req.DetailedAddress
	address.IsDefault = req.IsDefault
	address.Status = req.Status
	err := s.ac.Create(ctx, address)
	return &pb.CreateAddressReply{}, err
}
func (s *UserService) GetAddress(ctx context.Context, req *pb.GetAddressRequest) (*pb.GetAddressReply, error) {
	// 获取用户id
	userId := ctx.Value("x-md-global-uid")
	result, err := s.ac.Get(ctx, req.Id, userId.(int64))
	reply := &pb.GetAddressReply{
		Id:              result.Id,
		UserId:          result.UserId,
		Country:         result.Country,
		Province:        result.Province,
		City:            result.City,
		CountryDistrict: result.CountryDistrict,
		DetailedAddress: result.DetailedAddress,
		IsDefault:       result.IsDefault,
		Status:          result.Status,
	}
	return reply, err
}
func (s *UserService) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UpdateAddressReply, error) {
	// 获取用户id
	userId := ctx.Value("x-md-global-uid")
	address := &biz.Address{
		Id:              req.Id,
		UserId:          req.UserId,
		Country:         req.Country,
		Province:        req.Province,
		City:            req.City,
		CountryDistrict: req.CountryDistrict,
		DetailedAddress: req.DetailedAddress,
		IsDefault:       req.IsDefault,
		Status:          req.Status,
	}
	err := s.ac.Update(ctx, userId.(int64), address)
	return nil, err
}
func (s *UserService) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressReply, error) {
	// 获取用户id
	userId := ctx.Value("x-md-global-uid")
	err := s.ac.Delete(ctx, req.Id, userId.(int64))
	return nil, err
}
