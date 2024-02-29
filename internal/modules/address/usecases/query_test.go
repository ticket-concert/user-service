package usecases_test

import (
	"context"
	"testing"
	"user-service/internal/modules/address"
	"user-service/internal/modules/address/models/entity"
	"user-service/internal/modules/address/models/request"
	uc "user-service/internal/modules/address/usecases"
	"user-service/internal/pkg/errors"
	"user-service/internal/pkg/helpers"
	mockcertAddress "user-service/mocks/modules/address"
	mocklog "user-service/mocks/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type QueryUsecaseTestSuite struct {
	suite.Suite
	mockAddressRepositoryQuery *mockcertAddress.MongodbRepositoryQuery
	mockLogger                 *mocklog.Logger
	usecase                    address.UsecaseQuery
	ctx                        context.Context
}

func (suite *QueryUsecaseTestSuite) SetupTest() {
	suite.mockAddressRepositoryQuery = &mockcertAddress.MongodbRepositoryQuery{}
	suite.mockLogger = &mocklog.Logger{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewQueryUsecase(
		suite.mockAddressRepositoryQuery,
		suite.mockLogger,
	)
}
func TestQueryUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(QueryUsecaseTestSuite))
}

func (suite *QueryUsecaseTestSuite) TestFindProvinceSuccess() {
	// Arrange
	payload := request.Province{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.Province{
			{
				Id:   "1",
				Name: "name",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindProvinces", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))

	// Act
	result, err := suite.usecase.FindProvinces(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindProvinceErr() {
	// Arrange
	payload := request.Province{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data:  &[]entity.Province{},
		Error: errors.InternalServerError("error"),
	}
	suite.mockAddressRepositoryQuery.On("FindProvinces", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindProvinces(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)

	mockUserQueryResponse = helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindProvinces", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err = suite.usecase.FindProvinces(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindProvinceErrParse() {
	// Arrange
	payload := request.Province{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.City{
			{
				Id:           "id",
				Name:         "name",
				ProvinceId:   "provinceId",
				ProvinceName: "provinceName",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindProvinces", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindProvinces(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)

	mockUserQueryResponse = helpers.Result{
		Data:  &entity.Province{},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindProvinces", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err = suite.usecase.FindProvinces(suite.ctx, payload)
	suite.T().Log(err)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindProvinceNil() {
	// Arrange
	payload := request.Province{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data:  &entity.Province{},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindProvinces", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindProvinces(suite.ctx, payload)

	suite.T().Log(err)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindCitiesSuccess() {
	// Arrange
	payload := request.City{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.City{
			{
				Id:           "id",
				Name:         "name",
				ProvinceId:   "provinceId",
				ProvinceName: "provinceName",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindCitiesByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))

	// Act
	result, err := suite.usecase.FindCities(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindCitiesErr() {
	// Arrange
	payload := request.City{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data:  &[]entity.City{},
		Error: errors.InternalServerError("error"),
	}
	suite.mockAddressRepositoryQuery.On("FindCitiesByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindCities(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)

	mockUserQueryResponse = helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindCitiesByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err = suite.usecase.FindCities(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindCitiesErrParse() {
	// Arrange
	payload := request.City{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.Country{
			{
				Id:   1,
				Name: "name",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindCitiesByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindCities(suite.ctx, payload)
	suite.T().Log(err)

	// Assert
	assert.Error(suite.T(), err)

}

func (suite *QueryUsecaseTestSuite) TestFindDistrictSuccess() {
	// Arrange
	payload := request.District{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.District{
			{
				Id:           "id",
				Name:         "name",
				CityId:       "id",
				CityName:     "name",
				ProvinceId:   "provinceId",
				ProvinceName: "provinceName",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindDistrictByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))

	// Act
	result, err := suite.usecase.FindDistricts(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindDistrictErr() {
	// Arrange
	payload := request.District{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data:  &[]entity.District{},
		Error: errors.InternalServerError("error"),
	}
	suite.mockAddressRepositoryQuery.On("FindDistrictByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindDistricts(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)

	mockUserQueryResponse = helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindDistrictByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err = suite.usecase.FindDistricts(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)

	mockUserQueryResponse = helpers.Result{
		Data: &entity.Province{
			Id:   "1",
			Name: "name",
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindDistrictByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err = suite.usecase.FindDistricts(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindDistrictErrParse() {
	// Arrange
	payload := request.District{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.Province{
			{
				Id:   "1",
				Name: "name",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindDistrictByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindDistricts(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindSubdistrictSuccess() {
	// Arrange
	payload := request.SubDistrict{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.SubDistrict{
			{
				Id:           "id",
				Name:         "name",
				CityId:       "id",
				CityName:     "name",
				DistrictId:   "id",
				DistrictName: "name",
				ProvinceId:   "provinceId",
				ProvinceName: "provinceName",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindSubDistrictByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))

	// Act
	result, err := suite.usecase.FindSubDistricts(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindSubdistrictErr() {
	// Arrange
	payload := request.SubDistrict{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data:  &[]entity.SubDistrict{},
		Error: errors.InternalServerError("error"),
	}
	suite.mockAddressRepositoryQuery.On("FindSubDistrictByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindSubDistricts(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)

	mockUserQueryResponse = helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindSubDistrictByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err = suite.usecase.FindSubDistricts(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindSubdistrictErrParse() {
	// Arrange
	payload := request.SubDistrict{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.Province{
			{
				Id:   "id",
				Name: "name",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindSubDistrictByParam", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindSubDistricts(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindCountriesSuccess() {
	// Arrange
	payload := request.Country{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.Country{
			{
				Id:   1,
				Code: "",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindCountries", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))

	// Act
	result, err := suite.usecase.FindCountries(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindCountriesErr() {
	// Arrange
	payload := request.Country{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data:  &[]entity.Country{},
		Error: errors.InternalServerError("error"),
	}
	suite.mockAddressRepositoryQuery.On("FindCountries", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindCountries(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)

	mockUserQueryResponse = helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindCountries", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err = suite.usecase.FindCountries(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindCountriesErrParse() {
	// Arrange
	payload := request.Country{
		Page: 1,
		Size: 10,
	}

	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.Province{
			{
				Id:   "id",
				Name: "name",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindCountries", mock.Anything, payload).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindCountries(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindContinentSuccess() {
	// Arrange
	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.Continent{
			{
				Code: "code",
				Name: "name",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindContinent", mock.Anything).Return(mockChannel(mockUserQueryResponse))

	// Act
	result, err := suite.usecase.FindContinent(suite.ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindContinentSuccessErr() {
	// Arrange
	mockUserQueryResponse := helpers.Result{
		Data:  &[]entity.Continent{},
		Error: errors.InternalServerError("error"),
	}
	suite.mockAddressRepositoryQuery.On("FindContinent", mock.Anything).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindContinent(suite.ctx)

	// Assert
	assert.Error(suite.T(), err)

	mockUserQueryResponse = helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindContinent", mock.Anything).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err = suite.usecase.FindContinent(suite.ctx)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindContinentSuccessErrParse() {
	// Arrange
	mockUserQueryResponse := helpers.Result{
		Data: &[]entity.Province{
			{
				Id:   "id",
				Name: "name",
			},
		},
		Error: nil,
	}
	suite.mockAddressRepositoryQuery.On("FindContinent", mock.Anything).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindContinent(suite.ctx)

	// Assert
	assert.Error(suite.T(), err)
}

// Helper function to create a channel
func mockChannel(result helpers.Result) <-chan helpers.Result {
	responseChan := make(chan helpers.Result)

	go func() {
		responseChan <- result
		close(responseChan)
	}()

	return responseChan
}
