package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	addressEntity "user-service/internal/modules/address/models/entity"
	"user-service/internal/modules/user"
	userEntity "user-service/internal/modules/user/models/entity"
	userRequest "user-service/internal/modules/user/models/request"
	uc "user-service/internal/modules/user/usecases"
	"user-service/internal/pkg/errors"
	"user-service/internal/pkg/helpers"
	mockcertAddress "user-service/mocks/modules/address"
	mockcert "user-service/mocks/modules/user"
	mockjwt "user-service/mocks/pkg/helpers"
	mockkafka "user-service/mocks/pkg/kafka"
	mocklog "user-service/mocks/pkg/log"
	mockredis "user-service/mocks/pkg/redis"
)

type CommandUsecaseTestSuite struct {
	suite.Suite
	mockUserRepositoryQuery    *mockcert.MongodbRepositoryQuery
	mockUserRepositoryCommand  *mockcert.MongodbRepositoryCommand
	mockAddressRepositoryQuery *mockcertAddress.MongodbRepositoryQuery
	mockLogger                 *mocklog.Logger
	mockRedis                  *mockredis.Collections
	mockKafkaProducer          *mockkafka.Producer
	mockJwt                    *mockjwt.TokenGenerator
	usecase                    user.UsecaseCommand
	ctx                        context.Context
}

func (suite *CommandUsecaseTestSuite) SetupTest() {
	suite.mockUserRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockUserRepositoryCommand = &mockcert.MongodbRepositoryCommand{}
	suite.mockAddressRepositoryQuery = &mockcertAddress.MongodbRepositoryQuery{}
	suite.mockLogger = &mocklog.Logger{}
	suite.mockRedis = &mockredis.Collections{}
	suite.mockKafkaProducer = &mockkafka.Producer{}
	suite.mockJwt = &mockjwt.TokenGenerator{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewCommandUsecase(
		suite.mockUserRepositoryQuery,
		suite.mockUserRepositoryCommand,
		suite.mockLogger,
		suite.mockRedis,
		suite.mockKafkaProducer,
		suite.mockJwt,
		suite.mockAddressRepositoryQuery,
	)
	array := [][]string{{}, {"yopmail.com"}}
	helpers.CreateBlackListEmail(array)
}

func TestCommandUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CommandUsecaseTestSuite))
}

func (suite *CommandUsecaseTestSuite) TestRegisterUserSuccess() {
	// Arrange user request register
	payload := userRequest.RegisterUser{
		FullName:      "Alif Septian",
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "1",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}

	// Define a mock user repository query function
	mockFindOneByEmail := func(ctx context.Context, email string) <-chan helpers.Result {
		responseChan := make(chan helpers.Result)

		go func() {
			responseChan <- helpers.Result{
				Data:  nil,
				Error: nil,
			}
			close(responseChan)
		}()

		return responseChan
	}

	addressData := &addressEntity.Country{
		Id:            1,
		Code:          "ID",
		Name:          "Name",
		Iso3:          "Iso3",
		Number:        1,
		ContinentCode: "ContinentCode",
		ContinentName: "ContinentName",
		DisplayOrder:  1,
		FullName:      "FullName",
	}
	// Define a mock address repository query function
	mockFindOneCountry := func(ctx context.Context, id int) <-chan helpers.Result {
		responseChan := make(chan helpers.Result)

		go func() {
			responseChan <- helpers.Result{
				Data:  addressData,
				Error: nil,
			}
			close(responseChan)
		}()

		return responseChan
	}

	subdistrictData := &addressEntity.SubDistrict{
		Id:           "Id",
		Name:         "Name",
		DistrictId:   "DistricId",
		DistrictName: "DistrictName",
		CityId:       "CityId",
		CityName:     "CityName",
		ProvinceId:   "ProvinceId",
		ProvinceName: "ProvinceName",
	}
	mockFindOneSubdistrict := func(ctx context.Context, id string) <-chan helpers.Result {
		responseChan := make(chan helpers.Result)

		go func() {
			responseChan <- helpers.Result{
				Data:  subdistrictData,
				Error: nil,
			}
			close(responseChan)
		}()

		return responseChan
	}

	// Define a mock user repository command function
	mockUpsertOneUserTemp := func(ctx context.Context, user userEntity.User) <-chan helpers.Result {
		responseChan := make(chan helpers.Result)

		go func() {
			responseChan <- helpers.Result{
				Data:  user,
				Error: nil,
			}
			close(responseChan)
		}()

		return responseChan
	}
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockFindOneByEmail)
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockFindOneCountry)
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", suite.ctx, payload.SubdictrictId).Return(mockFindOneSubdistrict)
	suite.mockUserRepositoryCommand.On("UpsertOneUserTemp", mock.Anything, mock.Anything).Return(mockUpsertOneUserTemp)
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), payload.Email, payload.Email)

}

func (suite *CommandUsecaseTestSuite) TestRegisterUserErrPayloadCountry() {
	// Arrange user request register
	payload := userRequest.RegisterUser{
		FullName:      "Alif Septian",
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	mockFindOneByEmail := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("Error"),
	}

	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneByEmail))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert := assert.New(suite.T())
	assert.NotNil(err)

}

func (suite *CommandUsecaseTestSuite) TestRegisterUserErrFindCountry() {
	// Arrange user request register
	payload := userRequest.RegisterUser{
		FullName:      "Alif Septian",
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "1",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	mockFindOneByEmail := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("Error"),
	}

	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneByEmail))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert := assert.New(suite.T())
	assert.Error(err, "Error")

}

func (suite *CommandUsecaseTestSuite) TestRegisterUserErrFindCountryNil() {
	// Arrange user request register
	payload := userRequest.RegisterUser{
		FullName:      "Alif Septian",
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "1",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	mockFindOneByEmail := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneByEmail))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert := assert.New(suite.T())
	assert.Error(err, "Error")

}

func (suite *CommandUsecaseTestSuite) TestRegisterUserErrFindCountryParse() {
	// Arrange user request register
	payload := userRequest.RegisterUser{
		FullName:      "Alif Septian",
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "1",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	mockFindOneByEmail := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data:  payload,
		Error: nil,
	}

	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneByEmail))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert := assert.New(suite.T())
	assert.NotNil(err)

}

func (suite *CommandUsecaseTestSuite) TestRegisterUserErrFindSubdistrict() {
	// Arrange user request register
	payload := userRequest.RegisterUser{
		FullName:      "Alif Septian",
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "1",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	mockFindOneByEmail := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	country := &addressEntity.Country{
		Id:            1,
		Code:          "ID",
		Name:          "Name",
		Iso3:          "Iso3",
		Number:        1,
		ContinentCode: "ContinentCode",
		ContinentName: "ContinentName",
		DisplayOrder:  1,
		FullName:      "FullName",
	}

	mockFindOneCountry := helpers.Result{
		Data:  country,
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("error"),
	}

	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneByEmail))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", suite.ctx, payload.SubdictrictId).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert := assert.New(suite.T())
	assert.NotNil(err)

}

func (suite *CommandUsecaseTestSuite) TestRegisterUserErrFindSubdistrictNil() {
	// Arrange user request register
	payload := userRequest.RegisterUser{
		FullName:      "Alif Septian",
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "1",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	mockFindOneByEmail := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	country := &addressEntity.Country{
		Id:            1,
		Code:          "ID",
		Name:          "Name",
		Iso3:          "Iso3",
		Number:        1,
		ContinentCode: "ContinentCode",
		ContinentName: "ContinentName",
		DisplayOrder:  1,
		FullName:      "FullName",
	}

	mockFindOneCountry := helpers.Result{
		Data:  country,
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneByEmail))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", suite.ctx, payload.SubdictrictId).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert := assert.New(suite.T())
	assert.NotNil(err)

}

func (suite *CommandUsecaseTestSuite) TestRegisterUserErrFindSubdistrictParse() {
	// Arrange user request register
	payload := userRequest.RegisterUser{
		FullName:      "Alif Septian",
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "1",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	mockFindOneByEmail := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	country := &addressEntity.Country{
		Id:            1,
		Code:          "ID",
		Name:          "Name",
		Iso3:          "Iso3",
		Number:        1,
		ContinentCode: "ContinentCode",
		ContinentName: "ContinentName",
		DisplayOrder:  1,
		FullName:      "FullName",
	}

	mockFindOneCountry := helpers.Result{
		Data:  country,
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data:  country,
		Error: nil,
	}

	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneByEmail))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", suite.ctx, payload.SubdictrictId).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert := assert.New(suite.T())
	assert.NotNil(err)

}

func (suite *CommandUsecaseTestSuite) TestRegisterUserValidationFailed() {
	// Test incorrect email format

	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	payload := userRequest.RegisterUser{
		Email:         "invalidemail.com",
		Password:      "password",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		FullName:      "Full Name",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "32",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)

	assert := assert.New(suite.T())
	assert.Error(err, "Incorrect email format")

	// Test email blacklist

	payload.Email = "alif@yopmail.com"
	_, err = suite.usecase.RegisterUser(suite.ctx, payload)
	assert.Error(err, "Email blacklist")

	// Test password not criteria

	payload.Email = "alif@gmail.com"
	payload.Password = "password"
	_, err = suite.usecase.RegisterUser(suite.ctx, payload)
	assert.Error(err, "Password not criteria")

	// Test user role not found

	payload.Email = "alif@gmail.com"
	payload.Password = "Password1@"
	payload.Role = "role bukan didalam list"
	_, err = suite.usecase.RegisterUser(suite.ctx, payload)
	assert.Error(err, "User role not found")

	// Validation find email error
	payload.Email = "alif@gmail.com"
	payload.Password = "Password1@"
	payload.Role = "user"
	// Define a mock user repository query function
	mockFindOneByEmail := func(ctx context.Context, email string) <-chan helpers.Result {
		responseChan := make(chan helpers.Result)

		go func() {
			responseChan <- helpers.Result{
				Data:  nil,
				Error: errors.InternalServerError("Error"),
			}
			close(responseChan)
		}()

		return responseChan
	}
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockFindOneByEmail)
	// Act
	_, err = suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert.Error(err, "Error")
}
func (suite *CommandUsecaseTestSuite) TestRegisterUserEmailRegistered() {
	// Validation find email already register
	payload := userRequest.RegisterUser{
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		FullName:      "Full Name",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	// Define a mock user repository query function
	mockFindOneByEmail := func(ctx context.Context, email string) <-chan helpers.Result {
		responseChan := make(chan helpers.Result)

		go func() {
			responseChan <- helpers.Result{
				Data:  true,
				Error: nil,
			}
			close(responseChan)
		}()

		return responseChan
	}
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockFindOneByEmail)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert := assert.New(suite.T())
	assert.Error(err, "Email is already registered")
}

func (suite *CommandUsecaseTestSuite) TestRegisterUserUpsertOneError() {
	// Arrange
	payload := userRequest.RegisterUser{
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		FullName:      "Full Name",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "1",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}

	mockFindOneByEmail := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneByEmail))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", suite.ctx, payload.SubdictrictId).Return(mockChannel(mockFindOneSubdistrict))

	mockUpsertOneUserTemp := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("Error"),
	}
	suite.mockUserRepositoryCommand.On("UpsertOneUserTemp", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUserTemp))

	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert := assert.New(suite.T())
	suite.T().Log(err)
	assert.Error(err, "Error")
}

func (suite *CommandUsecaseTestSuite) TestVerifyRegisterUserSuccess() {
	// Arrange user request verify
	payload := userRequest.VerifyRegisterUser{
		Email: "alif@gmail.com",
		Otp:   "123456",
	}

	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("123456", nil))

	mockUserQueryResponse := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "/P8zIA/HX7pQew0m4SLWcBl/ivugz8wTFyKmFPQiAaA=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneByEmailUserTemp", mock.Anything, payload.Email).Return(mockChannel(mockUserQueryResponse))

	mockUserCommandResponse := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUserCommandResponse))
	suite.mockRedis.On("Del", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewIntResult(1, nil))
	// suite.mockJwt.On("GenerateToken", mock.Anything, mock.Anything).Return("mockedToken", "mockedExpiredAt", nil)
	// suite.mockJwt.On("GenerateTokenRefresh", mock.Anything, mock.Anything).Return("mockedToken", nil)

	// Act
	_, err := suite.usecase.VerifyRegisterUser(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestVerifyRegisterUserValidationFailed() {
	// Test incorrect email format

	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Arrange user request verify
	payload := userRequest.VerifyRegisterUser{
		Email: "alifgmail.com",
		Otp:   "123456",
	}
	_, err := suite.usecase.VerifyRegisterUser(suite.ctx, payload)

	assert := assert.New(suite.T())
	assert.Error(err, "Incorrect email format")

	// Test email blacklist

	payload.Email = "alif@yopmail.com"
	_, err = suite.usecase.VerifyRegisterUser(suite.ctx, payload)
	assert.Error(err, "Email blacklist")

	// Test redis empty
	payload.Email = "alif@gmail.com"
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("", nil))
	_, err = suite.usecase.VerifyRegisterUser(suite.ctx, payload)
	assert.Error(err, "Otp expired")

	// Test Otp not match
	suite.mockRedis.ExpectedCalls = nil // Reset redis mock
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("123564s", nil))
	_, err = suite.usecase.VerifyRegisterUser(suite.ctx, payload)
	assert.Error(err, "Otp not match")
}

func (suite *CommandUsecaseTestSuite) TestVerifyRegisterUserErrorFindEmailUserTemp() {
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Arrange user request verify
	payload := userRequest.VerifyRegisterUser{
		Email: "alif@gmail.com",
		Otp:   "123456",
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("123456", nil))
	mockUserQueryResponse := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("Error"),
	}
	suite.mockUserRepositoryQuery.On("FindOneByEmailUserTemp", mock.Anything, payload.Email).Return(mockChannel(mockUserQueryResponse))

	_, err := suite.usecase.VerifyRegisterUser(suite.ctx, payload)
	assert.NotNil(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestVerifyRegisterUserNotFoundFindEmailUserTemp() {
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Arrange user request verify
	payload := userRequest.VerifyRegisterUser{
		Email: "alif@gmail.com",
		Otp:   "123456",
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("123456", nil))
	mockUserQueryResponse := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneByEmailUserTemp", mock.Anything, payload.Email).Return(mockChannel(mockUserQueryResponse))

	_, err := suite.usecase.VerifyRegisterUser(suite.ctx, payload)
	assert := assert.New(suite.T())
	assert.Error(err, "Email not found")
}

func (suite *CommandUsecaseTestSuite) TestVerifyRegisterUserCantParsingData() {
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Arrange user request verify
	payload := userRequest.VerifyRegisterUser{
		Email: "alif@gmail.com",
		Otp:   "123456",
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("123456", nil))
	mockUserQueryResponse := helpers.Result{
		Data:  "",
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneByEmailUserTemp", mock.Anything, payload.Email).Return(mockChannel(mockUserQueryResponse))

	_, err := suite.usecase.VerifyRegisterUser(suite.ctx, payload)
	assert := assert.New(suite.T())
	assert.Error(err, "cannot parsing data")
}

func (suite *CommandUsecaseTestSuite) TestVerifyRegisterUserErrUpsertOne() {
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Arrange user request verify
	payload := userRequest.VerifyRegisterUser{
		Email: "alif@gmail.com",
		Otp:   "123456",
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("123456", nil))
	mockUserQueryResponse := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "/P8zIA/HX7pQew0m4SLWcBl/ivugz8wTFyKmFPQiAaA=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneByEmailUserTemp", mock.Anything, payload.Email).Return(mockChannel(mockUserQueryResponse))

	mockUserCommandResponse := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("Error"),
	}

	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUserCommandResponse))

	_, err := suite.usecase.VerifyRegisterUser(suite.ctx, payload)
	assert.NotNil(suite.T(), err)
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

func (suite *CommandUsecaseTestSuite) TestLoginUserSuccess() {
	// Arrange user request register
	payload := userRequest.LoginUser{
		Email:    "alif@gmail.com",
		Password: "Password1@",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("3", nil))
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneUser))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mockRedis.On("Del", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mockJwt.On("GenerateToken", mock.Anything, mock.Anything).Return("mockedToken", "mockedExpiredAt", nil)
	suite.mockJwt.On("GenerateTokenRefresh", mock.Anything, mock.Anything).Return("mockedToken", nil)
	// Act
	_, err := suite.usecase.LoginUser(suite.ctx, payload)
	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), payload.Email, payload.Email)

}

func (suite *CommandUsecaseTestSuite) TestLoginUserValidationFailed() {
	// Test incorrect email format

	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Arrange user request verify
	payload := userRequest.LoginUser{
		Email:    "alifgmail.com",
		Password: "123456",
	}
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	_, err := suite.usecase.LoginUser(suite.ctx, payload)

	assert := assert.New(suite.T())
	assert.Error(err, "Incorrect email format")

	// Test email blacklist

	payload.Email = "alif@yopmail.com"
	_, err = suite.usecase.LoginUser(suite.ctx, payload)
	assert.Error(err, "Email blacklist")

	// Test redis over attempt
	payload.Email = "alif@gmail.com"
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("7", nil))
	_, err = suite.usecase.LoginUser(suite.ctx, payload)
	assert.Error(err, "You have too many attempts, please wait 10 minutes")

}

func (suite *CommandUsecaseTestSuite) TestLoginUserErrFindEmail() {
	// Arrange user request register
	payload := userRequest.LoginUser{
		Email:    "alif@gmail.com",
		Password: "Password1@",
	}

	// Test Error resp.Error == nil
	mockFindOneUser := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("error"),
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("3", nil))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneUser))
	// Act
	_, err := suite.usecase.LoginUser(suite.ctx, payload)
	// Assert
	assert.Error(suite.T(), err)

	// Test Error resp.Data == nil
	mockFindOneUser = helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("3", nil))
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneUser))
	// Act
	_, err = suite.usecase.LoginUser(suite.ctx, payload)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestLoginUserErrFindEmailParse() {
	// Arrange user request register
	payload := userRequest.LoginUser{
		Email:    "alif@gmail.com",
		Password: "Password1@",
	}

	// Test Error unmarshal data user
	mockFindOneUser := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("3", nil))
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneUser))
	// Act
	_, err := suite.usecase.LoginUser(suite.ctx, payload)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestLoginUserErrPassword() {
	// Arrange user request register
	payload := userRequest.LoginUser{
		Email:    "alif@gmail.com",
		Password: "Password1@",
	}

	// Test Error password data user
	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("3", nil))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneUser))
	// Act
	_, err := suite.usecase.LoginUser(suite.ctx, payload)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestLoginUserErrUpsert() {
	// Arrange user request register
	payload := userRequest.LoginUser{
		Email:    "alif@gmail.com",
		Password: "Password1@",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("error"),
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("3", nil))
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneUser))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mockRedis.On("Del", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mockJwt.On("GenerateToken", mock.Anything, mock.Anything).Return("mockedToken", "mockedExpiredAt", nil)
	suite.mockJwt.On("GenerateTokenRefresh", mock.Anything, mock.Anything).Return("mockedToken", nil)
	// Act
	_, err := suite.usecase.LoginUser(suite.ctx, payload)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestLoginUserErrGenToken() {
	// Arrange user request register
	payload := userRequest.LoginUser{
		Email:    "alif@gmail.com",
		Password: "Password1@",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("3", nil))
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneUser))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mockRedis.On("Del", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mockJwt.On("GenerateToken", mock.Anything, mock.Anything).Return("mockedToken", "mockedExpiredAt", errors.BadRequest("error"))
	// suite.mockJwt.On("GenerateTokenRefresh", mock.Anything, mock.Anything).Return("mockedToken", nil)
	// Act
	_, err := suite.usecase.LoginUser(suite.ctx, payload)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestLoginUserErrGenTokenRefresh() {
	// Arrange user request register
	payload := userRequest.LoginUser{
		Email:    "alif@gmail.com",
		Password: "Password1@",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("3", nil))
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneUser))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mockRedis.On("Del", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mockJwt.On("GenerateToken", mock.Anything, mock.Anything).Return("mockedToken", "mockedExpiredAt", nil)
	suite.mockJwt.On("GenerateTokenRefresh", mock.Anything, mock.Anything).Return("mockedToken", errors.BadRequest("error"))
	// Act
	_, err := suite.usecase.LoginUser(suite.ctx, payload)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestUpdateUser() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.NoError(suite.T(), err)

	// Test error Role
	payload.Role = "role"
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	// Act
	_, err = suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)

	// Test error email not found
	payload.Role = "user"
	mockFindOneUser.Error = errors.BadRequest("error")
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	// Act
	_, err = suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)

}

func (suite *CommandUsecaseTestSuite) TestUpdateUserErr() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: errors.InternalServerError("error"),
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestUpdateUserErrParse() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestUpdateUserErrCountry() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "country",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestUpdateUserErrCountryDb() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: errors.BadRequest("error"),
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestUpdateUserErrCountryNil() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestUpdateUserErrCountryParse() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestUpdateUserErrSubdistrict() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: errors.BadRequest("error"),
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestUpdateUserErrSubdistrictNil() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestUpdateUserErrSubdistrictParse() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestUpdateUserErrUpsert() {
	// Arrange user request register
	payload := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}

	mockFindOneUser := helpers.Result{
		Data: &userEntity.User{
			Email:        "alif@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "alif",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			KKNumber:     "<string>",
			LoginAt:      time.Now(),
			Password:     "PzWCUGI/iepF6Xyz1dKIgfQYRwkVTN5AdXTTl9Yz+W8=",
			Role:         "user",
			RtRw:         "<string>",
			Subdistrict: userEntity.Subdistrict{
				Name:         "Desa kkn",
				DistrictId:   "1",
				DistrictName: "Kelurahan kkn",
				CityId:       "1",
				CityName:     "Kota kkn",
				ProvinceId:   "1",
				ProvinceName: "Provinsi kkn",
				Id:           "1",
			},
			UpdatedAt: time.Now(),
			UserId:    "a1d7e6c6-a4b4-48b0-b436-c882a9cb7980",
		},
		Error: nil,
	}

	mockFindOneCountry := helpers.Result{
		Data: &addressEntity.Country{
			Id:            1,
			Code:          "ID",
			Name:          "Name",
			Iso3:          "Iso3",
			Number:        1,
			ContinentCode: "ContinentCode",
			ContinentName: "ContinentName",
			DisplayOrder:  1,
			FullName:      "FullName",
		},
		Error: nil,
	}

	mockFindOneSubdistrict := helpers.Result{
		Data: &addressEntity.SubDistrict{
			Id:           "Id",
			Name:         "Name",
			DistrictId:   "DistricId",
			DistrictName: "DistrictName",
			CityId:       "CityId",
			CityName:     "CityName",
			ProvinceId:   "ProvinceId",
			ProvinceName: "ProvinceName",
		},
		Error: nil,
	}

	// Define a mock user repository command function
	mockUpsertOneUser := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneUser))
	suite.mockAddressRepositoryQuery.On("FindOneCountry", mock.Anything, 1).Return(mockChannel(mockFindOneCountry))
	suite.mockAddressRepositoryQuery.On("FindOneSubdistrict", mock.Anything, mock.Anything).Return(mockChannel(mockFindOneSubdistrict))
	suite.mockUserRepositoryCommand.On("UpsertOneUser", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUser))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	// Act
	_, err := suite.usecase.UpdateUser(suite.ctx, payload, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}
