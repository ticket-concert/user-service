package commands

import (
	"context"
	user "user-service/internal/modules/user"
	userEntity "user-service/internal/modules/user/models/entity"
	"user-service/internal/pkg/databases/mongodb"
	wrapper "user-service/internal/pkg/helpers"
	"user-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type commandMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewCommandMongodbRepository(mongodb mongodb.Collections, log log.Logger) user.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (c commandMongodbRepository) UpsertOneUserTemp(ctx context.Context, user userEntity.User) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		resp := <-c.mongoDb.UpsertOne(mongodb.UpdateOne{
			CollectionName: "users-temp",
			Document:       user,
			Filter: bson.M{
				"email": user.Email,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (c commandMongodbRepository) UpsertOneUser(ctx context.Context, user userEntity.User) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		resp := <-c.mongoDb.UpsertOne(mongodb.UpdateOne{
			CollectionName: "users",
			Document:       user,
			Filter: bson.M{
				"email": user.Email,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
