// user_http_handler_test.go

package handlers_test

import (
	"net/http/httptest"
	"testing"
	"user-service/internal/modules/address/handlers"
	"user-service/internal/modules/address/models/response"
	"user-service/internal/pkg/constants"
	"user-service/internal/pkg/errors"
	mockcert "user-service/mocks/modules/address"
	mocklog "user-service/mocks/pkg/log"
	mockredis "user-service/mocks/pkg/redis"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type AddressHttpHandlerTestSuite struct {
	suite.Suite

	cUQ       *mockcert.UsecaseQuery
	cLog      *mocklog.Logger
	validator *validator.Validate
	handler   *handlers.AddressHttpHandler
	cRedis    *mockredis.Collections
	app       *fiber.App
}

func (suite *AddressHttpHandlerTestSuite) SetupTest() {
	suite.cUQ = new(mockcert.UsecaseQuery)
	suite.cLog = new(mocklog.Logger)
	suite.validator = validator.New()
	suite.cRedis = new(mockredis.Collections)
	suite.handler = &handlers.AddressHttpHandler{
		AddressUsecaseQuery: suite.cUQ,
		Logger:              suite.cLog,
		Validator:           suite.validator,
	}
	suite.app = fiber.New()
	handlers.InitAddressHttpHandler(suite.app, suite.cUQ, suite.cLog, suite.cRedis)
}

func TestUserHttpHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AddressHttpHandlerTestSuite))
}

func (suite *AddressHttpHandlerTestSuite) TestGetProvinces() {

	response := &response.ProvinceResp{
		CollectionData: []response.Province{
			{
				Id:   "id",
				Name: "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindProvinces", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/provinces?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/provinces?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetProvinces(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetProvincesErrQueryParser() {

	response := &response.ProvinceResp{
		CollectionData: []response.Province{
			{
				Id:   "id",
				Name: "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindProvinces", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/provinces?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/provinces?page=aa&size=bb")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetProvinces(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetProvincesErrValidator() {

	response := &response.ProvinceResp{
		CollectionData: []response.Province{
			{
				Id:   "id",
				Name: "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindProvinces", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/provinces?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/provinces?page=&size=")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetProvinces(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetProvincesErr() {
	suite.cUQ.On("FindProvinces", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/provinces?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/provinces?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetProvinces(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetCities() {

	response := &response.CityResp{
		CollectionData: []response.City{
			{
				Id:           "id",
				Name:         "name",
				ProvinceId:   "ProvinceId",
				ProvinceName: "ProvinceName",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindCities", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/cities?provinceId=1&page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/cities?provinceId=1&page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetCities(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetCitiesErrQueryParser() {

	response := &response.CityResp{
		CollectionData: []response.City{
			{
				Id:   "id",
				Name: "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindCities", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/cities?page=aa&size=bb", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/cities?page=aa&size=bb")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetCities(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetCitiesErrValidator() {

	response := &response.CityResp{
		CollectionData: []response.City{
			{
				Id:   "id",
				Name: "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindCities", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/cities?provinceId=1&page=&size=", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/cities?provinceId=&page=&size=")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetCities(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetCitiesErr() {
	suite.cUQ.On("FindCities", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/cities?provinceId=1&page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/cities?provinceId=1&page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetCities(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetDistricts() {

	response := &response.DistrictResp{
		CollectionData: []response.District{
			{
				Id:           "id",
				Name:         "name",
				ProvinceId:   "ProvinceId",
				ProvinceName: "ProvinceName",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindDistricts", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/subdistricts?provinceId=1&cityId=1&page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/subdistricts?provinceId=1&cityId=1&page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetDistricts(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetDistrictsErrQueryParser() {

	response := &response.DistrictResp{
		CollectionData: []response.District{
			{
				Id:           "id",
				Name:         "name",
				ProvinceId:   "ProvinceId",
				ProvinceName: "ProvinceName",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindDistricts", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/subdistricts?page=aa&size=bb", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/subdistricts?page=aa&size=bb")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetDistricts(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetDistrictsErrValidator() {

	response := &response.DistrictResp{
		CollectionData: []response.District{
			{
				Id:           "id",
				Name:         "name",
				ProvinceId:   "ProvinceId",
				ProvinceName: "ProvinceName",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindDistricts", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/subdistricts?provinceId=&cityId=&page=&size=", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/subdistricts?provinceId=&cityId=&page=&size=")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetDistricts(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetDistrictsErr() {
	suite.cUQ.On("FindDistricts", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/subdistricts?provinceId=1&cityId=1&page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/subdistricts?provinceId=1&cityId=1&page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetDistricts(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetSubdistricts() {

	response := &response.SubDistrictResp{
		CollectionData: []response.SubDistrict{
			{
				Id:           "id",
				Name:         "name",
				ProvinceId:   "ProvinceId",
				ProvinceName: "ProvinceName",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindSubDistricts", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/subdistricts?provinceId=1&cityId=1&districtId=1&page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/subdistricts?provinceId=1&cityId=1&districtId=1&page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetSubDistricts(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetSubdistrictsErrQueryParser() {

	response := &response.SubDistrictResp{
		CollectionData: []response.SubDistrict{
			{
				Id:           "id",
				Name:         "name",
				ProvinceId:   "ProvinceId",
				ProvinceName: "ProvinceName",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindSubDistricts", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/subdistricts?page=aa&size=bb", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/subdistricts?page=aa&size=bb")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetSubDistricts(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetSubdistrictsErrValidator() {

	response := &response.SubDistrictResp{
		CollectionData: []response.SubDistrict{
			{
				Id:           "id",
				Name:         "name",
				ProvinceId:   "ProvinceId",
				ProvinceName: "ProvinceName",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindSubDistricts", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/subdistricts?provinceId=&cityId=&districtId=&page=&size=", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/subdistricts?provinceId=&cityId=&districtId=&page=&size=")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetSubDistricts(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetSubdistrictsErr() {
	suite.cUQ.On("FindSubDistricts", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/subdistricts?provinceId=1&cityId=1&districtId=1&page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/subdistricts?provinceId=1&cityId=1&districtId=1&page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetSubDistricts(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetCountries() {

	response := &response.CountryResp{
		CollectionData: []response.Country{
			{
				Id:   1,
				Name: "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindCountries", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/countries?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/countries?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetCountries(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetCountriesErrQueryParser() {

	response := &response.CountryResp{
		CollectionData: []response.Country{
			{
				Id:   1,
				Name: "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindCountries", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/countries?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/countries?page=aa&size=bb")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetCountries(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetCountriesErrValidator() {

	response := &response.CountryResp{
		CollectionData: []response.Country{
			{
				Id:   1,
				Name: "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindCountries", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/countries?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/countries?page=&size=")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetCountries(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetCountriesErr() {
	suite.cUQ.On("FindCountries", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/countries?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/countries?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetCountries(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetContinent() {

	response := &response.ContinentResp{
		CollectionData: []response.Continent{
			{
				Code: "1",
				Name: "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindContinent", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/continent", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/continent")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetContinent(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *AddressHttpHandlerTestSuite) TestGetContinentErr() {
	suite.cUQ.On("FindContinent", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/continent", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/continent")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetContinent(ctx)
	assert.Nil(suite.T(), err)
}
