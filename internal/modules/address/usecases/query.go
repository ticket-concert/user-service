package usecases

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/modules/address"
	"user-service/internal/modules/address/models/entity"
	"user-service/internal/modules/address/models/request"
	"user-service/internal/modules/address/models/response"
	"user-service/internal/pkg/errors"
	"user-service/internal/pkg/helpers"
	"user-service/internal/pkg/log"

	"go.elastic.co/apm"
)

type queryUsecase struct {
	addressRepositoryQuery address.MongodbRepositoryQuery
	logger                 log.Logger
}

func NewQueryUsecase(amq address.MongodbRepositoryQuery, log log.Logger) address.UsecaseQuery {
	return queryUsecase{
		addressRepositoryQuery: amq,
		logger:                 log,
	}
}

func (q queryUsecase) FindProvinces(origCtx context.Context, payload request.Province) (*response.ProvinceResp, error) {
	domain := "addressUsecase-FindProvinces"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	resp := <-q.addressRepositoryQuery.FindProvinces(ctx, payload)
	if resp.Error != nil {
		msg := "Error query province"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
		return nil, resp.Error
	}

	if resp.Data == nil {
		msg := "Province Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound("province not found")
	}

	province, ok := resp.Data.(*[]entity.Province)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	var collectionData = make([]response.Province, 0)
	for _, value := range *province {
		collectionData = append(collectionData, response.Province{
			Id:   value.Id,
			Name: value.Name,
		})
	}

	return &response.ProvinceResp{
		CollectionData: collectionData,
		MetaData:       helpers.GenerateMetaData(resp.Count, int64(len(*province)), payload.Page, payload.Size),
	}, nil

}

func (q queryUsecase) FindCities(origCtx context.Context, payload request.City) (*response.CityResp, error) {
	domain := "addressUsecase-FindCities"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	resp := <-q.addressRepositoryQuery.FindCitiesByParam(ctx, payload)
	if resp.Error != nil {
		msg := "Error query city"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
		return nil, resp.Error
	}

	if resp.Data == nil {
		msg := "City Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound("city not found")
	}

	city, ok := resp.Data.(*[]entity.City)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	var collectionData = make([]response.City, 0)
	for _, value := range *city {
		collectionData = append(collectionData, response.City{
			Id:           value.Id,
			Name:         value.Name,
			ProvinceId:   value.ProvinceId,
			ProvinceName: value.ProvinceName,
		})
	}

	return &response.CityResp{
		CollectionData: collectionData,
		MetaData:       helpers.GenerateMetaData(resp.Count, int64(len(*city)), payload.Page, payload.Size),
	}, nil

}

func (q queryUsecase) FindDistricts(origCtx context.Context, payload request.District) (*response.DistrictResp, error) {
	domain := "addressUsecase-FindDistricts"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	resp := <-q.addressRepositoryQuery.FindDistrictByParam(ctx, payload)
	if resp.Error != nil {
		msg := "Error query district"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
		return nil, resp.Error
	}

	if resp.Data == nil {
		msg := "District Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound("district not found")
	}

	district, ok := resp.Data.(*[]entity.District)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	var collectionData = make([]response.District, 0)
	for _, value := range *district {
		collectionData = append(collectionData, response.District{
			Id:           value.Id,
			Name:         value.Name,
			CityId:       value.CityId,
			CityName:     value.CityName,
			ProvinceId:   value.ProvinceId,
			ProvinceName: value.ProvinceName,
		})
	}

	return &response.DistrictResp{
		CollectionData: collectionData,
		MetaData:       helpers.GenerateMetaData(resp.Count, int64(len(*district)), payload.Page, payload.Size),
	}, nil
}

func (q queryUsecase) FindSubDistricts(origCtx context.Context, payload request.SubDistrict) (*response.SubDistrictResp, error) {
	domain := "addressUsecase-FindSubDistricts"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	resp := <-q.addressRepositoryQuery.FindSubDistrictByParam(ctx, payload)
	if resp.Error != nil {
		msg := "Error query subdistrict"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
		return nil, resp.Error
	}

	if resp.Data == nil {
		msg := "Subistrict Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound("subdistrict not found")
	}

	subdistrict, ok := resp.Data.(*[]entity.SubDistrict)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	var collectionData = make([]response.SubDistrict, 0)
	for _, value := range *subdistrict {
		collectionData = append(collectionData, response.SubDistrict{
			Id:           value.Id,
			Name:         value.Name,
			DistrictId:   value.DistrictId,
			DistrictName: value.DistrictName,
			CityId:       value.CityId,
			CityName:     value.CityName,
			ProvinceId:   value.ProvinceId,
			ProvinceName: value.ProvinceName,
		})
	}

	return &response.SubDistrictResp{
		CollectionData: collectionData,
		MetaData:       helpers.GenerateMetaData(resp.Count, int64(len(*subdistrict)), payload.Page, payload.Size),
	}, nil

}

func (q queryUsecase) FindCountries(origCtx context.Context, payload request.Country) (*response.CountryResp, error) {
	domain := "addressUsecase-FindCountries"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	resp := <-q.addressRepositoryQuery.FindCountries(ctx, payload)
	if resp.Error != nil {
		msg := "Error query country"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
		return nil, resp.Error
	}

	if resp.Data == nil {
		msg := "Country Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound("country not found")
	}

	country, ok := resp.Data.(*[]entity.Country)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	var collectionData = make([]response.Country, 0)
	for _, value := range *country {
		collectionData = append(collectionData, response.Country{
			Id:            value.Id,
			Code:          value.Code,
			Name:          value.Name,
			ContinentCode: value.ContinentCode,
			ContinentName: value.ContinentName,
			FullName:      value.FullName,
		})
	}

	return &response.CountryResp{
		CollectionData: collectionData,
		MetaData:       helpers.GenerateMetaData(resp.Count, int64(len(*country)), payload.Page, payload.Size),
	}, nil

}

func (q queryUsecase) FindContinent(origCtx context.Context) (*response.ContinentResp, error) {
	domain := "addressUsecase-FindContinent"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	resp := <-q.addressRepositoryQuery.FindContinent(ctx)
	if resp.Error != nil {
		msg := "Error query continent"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
		return nil, resp.Error
	}

	if resp.Data == nil {
		msg := "Continent Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp))
		return nil, errors.NotFound("continent not found")
	}

	continent, ok := resp.Data.(*[]entity.Continent)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	var collectionData = make([]response.Continent, 0)
	for _, value := range *continent {
		collectionData = append(collectionData, response.Continent{
			Code: value.Code,
			Name: value.Name,
		})
	}

	return &response.ContinentResp{
		CollectionData: collectionData,
		MetaData:       helpers.GenerateMetaData(resp.Count, int64(len(*continent)), 1, int64(len(*continent))),
	}, nil

}
