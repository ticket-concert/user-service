package address

import (
	"context"
	"user-service/internal/modules/address/models/request"
	"user-service/internal/modules/address/models/response"
	wrapper "user-service/internal/pkg/helpers"
)

type UsecaseQuery interface {
	FindProvinces(origCtx context.Context, payload request.Province) (*response.ProvinceResp, error)
	FindCities(origCtx context.Context, payload request.City) (*response.CityResp, error)
	FindDistricts(origCtx context.Context, payload request.District) (*response.DistrictResp, error)
	FindSubDistricts(origCtx context.Context, payload request.SubDistrict) (*response.SubDistrictResp, error)
	FindCountries(origCtx context.Context, payload request.Country) (*response.CountryResp, error)
	FindContinent(origCtx context.Context) (*response.ContinentResp, error)
}

type MongodbRepositoryQuery interface {
	FindProvinces(ctx context.Context, payload request.Province) <-chan wrapper.Result
	FindCitiesByParam(ctx context.Context, payload request.City) <-chan wrapper.Result
	FindDistrictByParam(ctx context.Context, payload request.District) <-chan wrapper.Result
	FindSubDistrictByParam(ctx context.Context, payload request.SubDistrict) <-chan wrapper.Result
	FindOneSubdistrict(ctx context.Context, id string) <-chan wrapper.Result
	FindCountries(ctx context.Context, payload request.Country) <-chan wrapper.Result
	FindOneCountry(ctx context.Context, id int) <-chan wrapper.Result
	FindContinent(ctx context.Context) <-chan wrapper.Result
}
