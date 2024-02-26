package queries

import (
	"context"
	user "user-service/internal/modules/user"
	userEntity "user-service/internal/modules/user/models/entity"
	"user-service/internal/pkg/databases/mongodb"
	wrapper "user-service/internal/pkg/helpers"
	"user-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) user.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindOneUserId(ctx context.Context, userId string) <-chan wrapper.Result {
	var user userEntity.User
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &user,
			CollectionName: "users",
			Filter: bson.M{
				"userId": userId,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindOneByEmail(ctx context.Context, email string) <-chan wrapper.Result {
	var user userEntity.User
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &user,
			CollectionName: "users",
			Filter: bson.M{
				"email": email,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindOneByEmailUserTemp(ctx context.Context, email string) <-chan wrapper.Result {
	var user userEntity.User
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &user,
			CollectionName: "users-temp",
			Filter: bson.M{
				"email": email,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
