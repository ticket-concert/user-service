package queries

import (
	"context"
	"fmt"
	"user-service/internal/modules/address"
	"user-service/internal/modules/address/models/entity"
	"user-service/internal/modules/address/models/request"
	"user-service/internal/pkg/databases/mongodb"
	wrapper "user-service/internal/pkg/helpers"
	"user-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) address.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindProvinces(ctx context.Context, payload request.Province) <-chan wrapper.Result {
	var province []entity.Province
	var countData int64
	output := make(chan wrapper.Result)

	go func() {
		// var filter interface{}
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &province,
			CountData:      &countData,
			CollectionName: "province",
			Filter:         bson.M{"name": primitive.Regex{Pattern: ".*" + payload.Search + ".*", Options: "i"}},
			Sort: &mongodb.Sort{
				FieldName: "name",
				By:        mongodb.SortAscending,
			},
			Page: payload.Page,
			Size: payload.Size,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindCitiesByParam(ctx context.Context, payload request.City) <-chan wrapper.Result {
	var city []entity.City
	var countData int64
	output := make(chan wrapper.Result)

	go func() {
		// var filter interface{}
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &city,
			CountData:      &countData,
			CollectionName: "city",
			Filter: bson.M{"provinceId": payload.ProvinceId,
				"name": primitive.Regex{Pattern: ".*" + payload.Search + ".*", Options: "i"}},
			Sort: &mongodb.Sort{
				FieldName: "name",
				By:        mongodb.SortAscending,
			},
			Page: payload.Page,
			Size: payload.Size,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindDistrictByParam(ctx context.Context, payload request.District) <-chan wrapper.Result {
	var district []entity.District
	var countData int64
	output := make(chan wrapper.Result)

	go func() {
		// var filter interface{}
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &district,
			CountData:      &countData,
			CollectionName: "district",
			Filter: bson.M{"cityId": payload.CityId, "provinceId": payload.ProvinceId,
				"name": primitive.Regex{Pattern: ".*" + payload.Search + ".*", Options: "i"}},
			Sort: &mongodb.Sort{
				FieldName: "name",
				By:        mongodb.SortAscending,
			},
			Page: payload.Page,
			Size: payload.Size,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindSubDistrictByParam(ctx context.Context, payload request.SubDistrict) <-chan wrapper.Result {
	var subDistrict []entity.SubDistrict
	var countData int64
	output := make(chan wrapper.Result)

	go func() {
		// var filter interface{}
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &subDistrict,
			CountData:      &countData,
			CollectionName: "subdistrict",
			Filter: bson.M{"districtId": payload.DistrictId, "cityId": payload.CityId,
				"provinceId": payload.ProvinceId, "name": bson.M{"$regex": fmt.Sprintf(".*%s.*", payload.Search), "$options": "i"}},
			Sort: &mongodb.Sort{
				FieldName: "name",
				By:        mongodb.SortAscending,
			},
			Page: payload.Page,
			Size: payload.Size,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindOneSubdistrict(ctx context.Context, id string) <-chan wrapper.Result {
	var subDistrict entity.SubDistrict
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &subDistrict,
			CollectionName: "subdistrict",
			Filter: bson.M{
				"id": id,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindCountries(ctx context.Context, payload request.Country) <-chan wrapper.Result {
	var country []entity.Country
	var countData int64
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &country,
			CountData:      &countData,
			CollectionName: "country",
			Filter:         bson.M{"name": primitive.Regex{Pattern: ".*" + payload.Search + ".*", Options: "i"}},
			Sort: &mongodb.Sort{
				FieldName: "name",
				By:        mongodb.SortAscending,
			},
			Page: payload.Page,
			Size: payload.Size,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindOneCountry(ctx context.Context, id int) <-chan wrapper.Result {
	var country entity.Country
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &country,
			CollectionName: "country",
			Filter: bson.M{
				"id": id,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindContinent(ctx context.Context) <-chan wrapper.Result {
	var continent []entity.Continent
	var countData int64
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &continent,
			CountData:      &countData,
			CollectionName: "continent",
			Filter:         bson.M{},
			Sort: &mongodb.Sort{
				FieldName: "name",
				By:        mongodb.SortAscending,
			},
			Page: 1,
			Size: 10,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
