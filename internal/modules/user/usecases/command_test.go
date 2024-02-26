package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"user-service/internal/modules/user"
	userEntity "user-service/internal/modules/user/models/entity"
	userRequest "user-service/internal/modules/user/models/request"
	uc "user-service/internal/modules/user/usecases"
	"user-service/internal/pkg/errors"
	"user-service/internal/pkg/helpers"
	mockcert "user-service/mocks/modules/user"
	mockjwt "user-service/mocks/pkg/helpers"
	mockkafka "user-service/mocks/pkg/kafka"
	mocklog "user-service/mocks/pkg/log"
	mockredis "user-service/mocks/pkg/redis"
)

type CommandUsecaseTestSuite struct {
	suite.Suite
	mockUserRepositoryQuery   *mockcert.MongodbRepositoryQuery
	mockUserRepositoryCommand *mockcert.MongodbRepositoryCommand
	mockLogger                *mocklog.Logger
	mockRedis                 *mockredis.Collections
	mockKafkaProducer         *mockkafka.Producer
	mockJwt                   *mockjwt.TokenGenerator
	usecase                   user.UsecaseCommand
	ctx                       context.Context
}

func (suite *CommandUsecaseTestSuite) SetupTest() {
	suite.mockUserRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockUserRepositoryCommand = &mockcert.MongodbRepositoryCommand{}
	suite.mockLogger = &mocklog.Logger{}
	suite.mockRedis = &mockredis.Collections{}
	suite.mockKafkaProducer = &mockkafka.Producer{}
	suite.mockJwt = &mockjwt.TokenGenerator{}
	suite.ctx = context.WithValue(context.TODO(), "key", "value")
	suite.usecase = uc.NewCommandUsecase(
		suite.mockUserRepositoryQuery,
		suite.mockUserRepositoryCommand,
		suite.mockLogger,
		suite.mockRedis,
		suite.mockKafkaProducer,
		suite.mockJwt,
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
		FullName:      "Irman Juliansyah",
		Email:         "irmanjuliansyah@gmail.com",
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
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)

	assert := assert.New(suite.T())
	assert.Error(err, "Incorrect email format")

	// Test email blacklist

	payload.Email = "irmanjuliansyah@yopmail.com"
	_, err = suite.usecase.RegisterUser(suite.ctx, payload)
	assert.Error(err, "Email blacklist")

	// Test password not criteria

	payload.Email = "irmanjuliansyah@gmail.com"
	payload.Password = "password"
	_, err = suite.usecase.RegisterUser(suite.ctx, payload)
	assert.Error(err, "Password not criteria")

	// Test user role not found

	payload.Email = "irmanjuliansyah@gmail.com"
	payload.Password = "Password1@"
	payload.Role = "role bukan didalam list"
	_, err = suite.usecase.RegisterUser(suite.ctx, payload)
	assert.Error(err, "User role not found")

	// Validation find email error
	payload.Email = "irmanjuliansyah@gmail.com"
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
		Email:         "irmanjuliansyah@gmail.com",
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
		Email:         "irmanjuliansyah@gmail.com",
		Password:      "Password1@",
		FullName:      "Full Name",
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
	suite.mockUserRepositoryQuery.On("FindOneByEmail", mock.Anything, payload.Email).Return(mockChannel(mockFindOneByEmail))

	mockUpsertOneUserTemp := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("Error"),
	}
	suite.mockUserRepositoryCommand.On("UpsertOneUserTemp", mock.Anything, mock.Anything).Return(mockChannel(mockUpsertOneUserTemp))

	// Act
	_, err := suite.usecase.RegisterUser(suite.ctx, payload)
	// Assert
	assert := assert.New(suite.T())
	assert.Error(err, "Error")
}

func (suite *CommandUsecaseTestSuite) TestVerifyRegisterUserSuccess() {
	// Arrange user request verify
	payload := userRequest.VerifyRegisterUser{
		Email: "irmanjuliansyah@gmail.com",
		Otp:   "123456",
	}

	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("123456", nil))

	mockUserQueryResponse := helpers.Result{
		Data: &userEntity.User{
			Email:        "irmanjuliansyah@gmail.com",
			Address:      "<string>",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			CreatedAt:    time.Now(),
			FullName:     "irmanjuliansyah",
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
		Email: "irmanjuliansyahgmail.com",
		Otp:   "123456",
	}
	_, err := suite.usecase.VerifyRegisterUser(suite.ctx, payload)

	assert := assert.New(suite.T())
	assert.Error(err, "Incorrect email format")

	// Test email blacklist

	payload.Email = "irmanjuliansyah@yopmail.com"
	_, err = suite.usecase.VerifyRegisterUser(suite.ctx, payload)
	assert.Error(err, "Email blacklist")

	// Test redis empty
	payload.Email = "irmanjuliansyah@gmail.com"
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
		Email: "irmanjuliansyah@gmail.com",
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
		Email: "irmanjuliansyah@gmail.com",
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
		Email: "irmanjuliansyah@gmail.com",
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
		Email: "irmanjuliansyah@gmail.com",
		Otp:   "123456",
	}
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("123456", nil))
	mockUserQueryResponse := helpers.Result{
		Data: &userEntity.User{
			Email:        "irmanjuliansyah@gmail.com",
			Address:      "<string>",
			CreatedAt:    time.Now(),
			FullName:     "irmanjuliansyah",
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
