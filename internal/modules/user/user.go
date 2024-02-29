package user

import (
	"context"
	userEntity "user-service/internal/modules/user/models/entity"
	userRequest "user-service/internal/modules/user/models/request"
	userResponse "user-service/internal/modules/user/models/response"
	wrapper "user-service/internal/pkg/helpers"
)

type UsecaseQuery interface {
	GetProfile(origCtx context.Context, payload userRequest.GetProfile) (*userResponse.GetProfile, error)
}

type UsecaseCommand interface {
	UpdateUser(origCtx context.Context, payload userRequest.UpdateUser, userId string) (string, error)
	RegisterUser(origCtx context.Context, payload userRequest.RegisterUser) (*userResponse.RegisterUser, error)
	VerifyRegisterUser(origCtx context.Context, payload userRequest.VerifyRegisterUser) (*userResponse.VerifyRegister, error)
	LoginUser(origCtx context.Context, payload userRequest.LoginUser) (*userResponse.LoginUserResp, error)
}

type MongodbRepositoryCommand interface {
	UpsertOneUserTemp(ctx context.Context, user userEntity.User) <-chan wrapper.Result
	UpsertOneUser(ctx context.Context, user userEntity.User) <-chan wrapper.Result
}

type MongodbRepositoryQuery interface {
	FindOneUserId(ctx context.Context, userId string) <-chan wrapper.Result
	FindOneByEmail(ctx context.Context, email string) <-chan wrapper.Result
	FindOneByEmailUserTemp(ctx context.Context, email string) <-chan wrapper.Result
}
