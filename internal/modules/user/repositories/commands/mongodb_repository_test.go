package commands_test

import (
	"context"
	"testing"
	"user-service/internal/modules/user"
	userEntity "user-service/internal/modules/user/models/entity"
	mongoRC "user-service/internal/modules/user/repositories/commands"
	"user-service/internal/pkg/helpers"
	mocks "user-service/mocks/pkg/databases/mongodb"
	mocklog "user-service/mocks/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
	mockMongodb *mocks.Collections
	mockLogger  *mocklog.Logger
	repository  user.MongodbRepositoryCommand
	ctx         context.Context
}

func (suite *CommandTestSuite) SetupTest() {
	suite.mockMongodb = new(mocks.Collections)
	suite.mockLogger = &mocklog.Logger{}
	suite.repository = mongoRC.NewCommandMongodbRepository(
		suite.mockMongodb,
		suite.mockLogger,
	)
	suite.ctx = context.WithValue(context.TODO(), "key", "value")
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (suite *CommandTestSuite) TestUpsertOneUserTemp() {
	testUser := userEntity.User{
		Email: "alif@gmail.com",
	}

	// Mock UpsertOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("UpsertOne", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.UpsertOneUserTemp(suite.ctx, testUser)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert UpsertOne
	suite.mockMongodb.AssertCalled(suite.T(), "UpsertOne", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestUpsertOneUser() {
	testUser := userEntity.User{
		Email: "alif@gmail.com",
	}

	// Mock UpsertOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("UpsertOne", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.UpsertOneUser(suite.ctx, testUser)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert UpsertOne
	suite.mockMongodb.AssertCalled(suite.T(), "UpsertOne", mock.Anything, mock.Anything)
}
