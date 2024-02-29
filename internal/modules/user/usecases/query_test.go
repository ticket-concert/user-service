// query_test.go
package usecases_test

import (
	"context"
	"testing"

	"user-service/internal/modules/user"
	userEntity "user-service/internal/modules/user/models/entity"
	userRequest "user-service/internal/modules/user/models/request"
	uc "user-service/internal/modules/user/usecases"
	"user-service/internal/pkg/errors"
	"user-service/internal/pkg/helpers"
	mockcert "user-service/mocks/modules/user"
	mocklog "user-service/mocks/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type QueryUsecaseTestSuite struct {
	suite.Suite
	mockUserRepositoryQuery   *mockcert.MongodbRepositoryQuery
	mockUserRepositoryCommand *mockcert.MongodbRepositoryCommand
	mockLogger                *mocklog.Logger
	usecase                   user.UsecaseQuery
	ctx                       context.Context
}

func (suite *QueryUsecaseTestSuite) SetupTest() {
	suite.mockUserRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockUserRepositoryCommand = &mockcert.MongodbRepositoryCommand{}
	suite.mockLogger = &mocklog.Logger{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewQueryUsecase(
		suite.mockUserRepositoryQuery,
		suite.mockUserRepositoryCommand,
		suite.mockLogger,
	)
}
func TestQueryUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(QueryUsecaseTestSuite))
}

func (suite *QueryUsecaseTestSuite) TestGetProfileSuccess() {
	// Arrange
	payload := userRequest.GetProfile{
		UserId: "76142a47-40c3-44a0-a7d3-793ee09a518b",
	}

	mockUserQueryResponse := helpers.Result{
		Data: &userEntity.User{
			UserId:       "76142a47-40c3-44a0-a7d3-793ee09a518b",
			FullName:     "alif",
			Email:        "alif@gmail.com",
			NIK:          "12312131131",
			MobileNumber: "+6281281015121",
			Address:      "<string>",
			RtRw:         "<string>",
			Role:         "user",
			KKNumber:     "<string>",
		},
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, payload.UserId).Return(mockChannel(mockUserQueryResponse))

	// Act
	result, err := suite.usecase.GetProfile(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), payload.UserId, result.UserId)
	assert.Equal(suite.T(), "alif", result.FullName)
	assert.Equal(suite.T(), "alif@gmail.com", result.Email)
	assert.Equal(suite.T(), "12312131131", result.NIK)
	assert.Equal(suite.T(), "+6281281015121", result.MobileNumber)
	assert.Equal(suite.T(), "<string>", result.Address)
	assert.Equal(suite.T(), "<string>", result.RtRw)
	assert.Equal(suite.T(), "user", result.Role)
	assert.Equal(suite.T(), "<string>", result.KKNumber)
}

func (suite *QueryUsecaseTestSuite) TestGetProfileNotFound() {
	// Arrange
	payload := userRequest.GetProfile{
		UserId: "nonexistentuser",
	}

	mockUserQueryResponse := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, payload.UserId).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	result, err := suite.usecase.GetProfile(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err, "User Not Found")
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestGetProfileError() {
	// Arrange
	payload := userRequest.GetProfile{
		UserId: "76142a47-40c3-44a0-a7d3-793ee09a518b",
	}

	mockUserQueryResponse := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("Error"),
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, payload.UserId).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	result, err := suite.usecase.GetProfile(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err, "Error")
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestGetProfileParsingError() {
	// Arrange
	payload := userRequest.GetProfile{
		UserId: "76142a47-40c3-44a0-a7d3-793ee09a518b",
	}

	mockUserQueryResponse := helpers.Result{
		Data:  "invalid_data_type", // Simulate an invalid data type in the response
		Error: nil,
	}
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, payload.UserId).Return(mockChannel(mockUserQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	result, err := suite.usecase.GetProfile(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err, "cannot parsing data")
	assert.Nil(suite.T(), result)
}
