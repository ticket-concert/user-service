package queries_test

import (
	"context"
	"testing"
	"user-service/internal/modules/user"
	mongoRQ "user-service/internal/modules/user/repositories/queries"
	mocks "user-service/mocks/pkg/databases/mongodb"
	mocklog "user-service/mocks/pkg/log"

	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
	mockMongodb *mocks.Collections
	mockLogger  *mocklog.Logger
	repository  user.MongodbRepositoryQuery
	ctx         context.Context
}

func (suite *CommandTestSuite) SetupTest() {
	suite.mockMongodb = new(mocks.Collections)
	suite.mockLogger = &mocklog.Logger{}
	suite.repository = mongoRQ.NewQueryMongodbRepository(
		suite.mockMongodb,
		suite.mockLogger,
	)
	suite.ctx = context.WithValue(context.TODO(), "key", "value")
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}
