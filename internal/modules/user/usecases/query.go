package usecases

import (
	"context"
	"fmt"
	"time"
	user "user-service/internal/modules/user"
	userEntity "user-service/internal/modules/user/models/entity"
	userRequest "user-service/internal/modules/user/models/request"
	userResponse "user-service/internal/modules/user/models/response"
	"user-service/internal/pkg/errors"
	"user-service/internal/pkg/log"

	"go.elastic.co/apm"
)

type queryUsecase struct {
	userRepositoryQuery   user.MongodbRepositoryQuery
	userRepositoryCommand user.MongodbRepositoryCommand
	logger                log.Logger
}

func NewQueryUsecase(
	umq user.MongodbRepositoryQuery, umc user.MongodbRepositoryCommand,
	log log.Logger) user.UsecaseQuery {
	return queryUsecase{
		userRepositoryQuery:   umq,
		userRepositoryCommand: umc,
		logger:                log,
	}
}

func (q queryUsecase) GetProfile(origCtx context.Context, payload userRequest.GetProfile) (*userResponse.GetProfile, error) {
	domain := "userUsecase-GetProfile"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()
	respUser := <-q.userRepositoryQuery.FindOneUserId(ctx, payload.UserId)
	if respUser.Error != nil {
		return nil, respUser.Error
	}
	if respUser.Data == nil {
		msg := "User Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound(msg)
	}
	userData, ok := respUser.Data.(*userEntity.User)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}
	response := userResponse.GetProfile{
		UserId:        userData.UserId,
		FullName:      userData.FullName,
		Email:         userData.Email,
		NIK:           userData.NIK,
		MobileNumber:  userData.MobileNumber,
		Address:       userData.Address,
		RtRw:          userData.RtRw,
		Role:          userData.Role,
		KKNumber:      userData.KKNumber,
		CountryCode:   userData.Country.Code,
		CountryName:   userData.Country.Name,
		ContinentName: userData.Country.ContinentName,
	}
	return &response, nil
}
