package queries_test

import (
	"context"
	"testing"
	"user-service/internal/modules/address"
	"user-service/internal/modules/address/models/request"
	mongoRQ "user-service/internal/modules/address/repositories/queries"
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
	repository  address.MongodbRepositoryQuery
	ctx         context.Context
}

func (suite *CommandTestSuite) SetupTest() {
	suite.mockMongodb = new(mocks.Collections)
	suite.mockLogger = &mocklog.Logger{}
	suite.repository = mongoRQ.NewQueryMongodbRepository(
		suite.mockMongodb,
		suite.mockLogger,
	)
	suite.ctx = context.Background()
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (suite *CommandTestSuite) TestFindProvinces() {

	request := request.Province{}
	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindAllData", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindProvinces(suite.ctx, request)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindAllData", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestCitiesByParam() {

	request := request.City{}
	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindAllData", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindCitiesByParam(suite.ctx, request)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindAllData", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestDistrictByParam() {

	request := request.District{}
	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindAllData", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindDistrictByParam(suite.ctx, request)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindAllData", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestSubdistrictByParam() {

	request := request.SubDistrict{}
	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindAllData", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindSubDistrictByParam(suite.ctx, request)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindAllData", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestFindOneSubdistrict() {

	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindOne", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindOneSubdistrict(suite.ctx, "1")
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindOne", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestFindCountries() {

	request := request.Country{}
	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindAllData", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindCountries(suite.ctx, request)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindAllData", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestFindOneCountry() {

	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindOne", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindOneCountry(suite.ctx, 1)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindOne", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestFindContinent() {
	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindAllData", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindContinent(suite.ctx)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindAllData", mock.Anything, mock.Anything)
}
