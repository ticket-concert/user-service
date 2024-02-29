package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"user-service/internal/modules/address"
	addressEntity "user-service/internal/modules/address/models/entity"
	user "user-service/internal/modules/user"
	userEntity "user-service/internal/modules/user/models/entity"
	userRequest "user-service/internal/modules/user/models/request"
	userResponse "user-service/internal/modules/user/models/response"
	"user-service/internal/pkg/constants"
	"user-service/internal/pkg/errors"
	"user-service/internal/pkg/helpers"
	kafkaPkgConfluent "user-service/internal/pkg/kafka/confluent"
	"user-service/internal/pkg/log"
	"user-service/internal/pkg/redis"

	uuid "github.com/google/uuid"
	"go.elastic.co/apm"
)

type commandUsecase struct {
	userRepositoryQuery    user.MongodbRepositoryQuery
	userRepositoryCommand  user.MongodbRepositoryCommand
	logger                 log.Logger
	redis                  redis.Collections
	kafkaProducer          kafkaPkgConfluent.Producer
	jwtHelper              helpers.TokenGenerator
	addressRepositoryQuery address.MongodbRepositoryQuery
}

func NewCommandUsecase(
	umq user.MongodbRepositoryQuery, umc user.MongodbRepositoryCommand,
	log log.Logger, rc redis.Collections, kp kafkaPkgConfluent.Producer,
	jwt helpers.TokenGenerator, amq address.MongodbRepositoryQuery) user.UsecaseCommand {
	return commandUsecase{
		userRepositoryQuery:    umq,
		userRepositoryCommand:  umc,
		logger:                 log,
		redis:                  rc,
		kafkaProducer:          kp,
		jwtHelper:              jwt,
		addressRepositoryQuery: amq,
	}
}

func (c commandUsecase) UpdateUser(origCtx context.Context, payload userRequest.UpdateUser, userId string) (string, error) {
	domain := "userUsecase-UpdateUser"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	if payload.Role != userRequest.RoleUser {
		msg := "User role not found"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return "", errors.NotFound(msg)
	}

	resp := <-c.userRepositoryQuery.FindOneUserId(ctx, userId)
	if resp.Error != nil {
		return "", resp.Error
	}

	if resp.Data == nil {
		msg := "Email not found"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return "", errors.NotFound(msg)
	}

	userData, ok := resp.Data.(*userEntity.User)
	if !ok {
		return "", errors.InternalServerError("cannot parsing data")
	}

	countryId, err := strconv.Atoi(payload.CountryId)
	if err != nil {
		msg := "CountryId must integer"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return "", errors.BadRequest(msg)
	}

	country := <-c.addressRepositoryQuery.FindOneCountry(ctx, countryId)
	if country.Error != nil {
		return "", country.Error
	}
	if country.Data == nil {
		msg := "Country not found"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return "", errors.CustomError(msg, 4002, http.StatusNotFound)
	}
	countryData, ok := country.Data.(*addressEntity.Country)
	if !ok {
		msg := "cannot parsing data country"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return "", errors.InternalServerError("cannot parsing data")
	}

	var subDistrictUser userEntity.Subdistrict
	if countryData.Code == "ID" {
		subdistrict := <-c.addressRepositoryQuery.FindOneSubdistrict(ctx, payload.SubdictrictId)
		if subdistrict.Error != nil {
			return "", subdistrict.Error
		}
		if subdistrict.Data == nil {
			msg := "Subdistrict not found"
			c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			return "", errors.CustomError(msg, 4002, http.StatusNotFound)
		}
		subdistrictData, ok := subdistrict.Data.(*addressEntity.SubDistrict)
		if !ok {
			msg := "cannot parsing data subdistrict"
			c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			return "", errors.InternalServerError("cannot parsing data")
		}
		if subdistrictData != nil {
			subDistrictUser = userEntity.Subdistrict{
				Id:           subdistrictData.Id,
				Name:         subdistrictData.Name,
				DistrictId:   subdistrictData.DistrictId,
				DistrictName: subdistrictData.DistrictName,
				CityId:       subdistrictData.CityId,
				CityName:     subdistrictData.CityName,
				ProvinceId:   subdistrictData.ProvinceId,
				ProvinceName: subdistrictData.ProvinceName,
			}
		}
	}

	user := userEntity.User{
		UserId:       userData.UserId,
		FullName:     payload.FullName,
		Email:        userData.Email,
		NIK:          userData.NIK,
		MobileNumber: helpers.VerifyPhoneNumber62(payload.MobileNumber),
		Password:     userData.Password,
		Subdistrict:  subDistrictUser,
		Country: userEntity.Country{
			Id:            countryData.Id,
			Code:          countryData.Code,
			Name:          countryData.Name,
			FullName:      countryData.FullName,
			ContinentId:   countryData.ContinentCode,
			ContinentName: countryData.ContinentName,
			Latitude:      payload.Latitude,
			Longitude:     payload.Longitude,
		},
		Status:    userData.Status,
		Address:   payload.Address,
		RtRw:      payload.RtRw,
		KKNumber:  userData.KKNumber,
		Role:      payload.Role,
		LoginAt:   userData.LoginAt,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: time.Now(),
	}
	respUser := <-c.userRepositoryCommand.UpsertOneUser(ctx, user)
	if respUser.Error != nil {
		return "", respUser.Error
	}
	return "Update user success", nil
}

func (c commandUsecase) RegisterUser(origCtx context.Context, payload userRequest.RegisterUser) (*userResponse.RegisterUser, error) {
	domain := "userUsecase-RegisterUser"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()
	validEmail := helpers.IsEmailValid(payload.Email)
	if !validEmail {
		msg := "Incorrect email format"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.CustomError(msg, 4001, http.StatusBadRequest)
	}
	if helpers.IsBlacklistedEmail(payload.Email) {
		msg := "Email blacklist"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.CustomError(msg, 4002, http.StatusBadRequest)
	}

	if !helpers.IsValidPassword(payload.Password) {
		msg := "Password not criteria"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.CustomError(msg, 4004, http.StatusBadRequest)
	}

	if payload.Role != userRequest.RoleUser {
		msg := "User role not found"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound(msg)
	}

	resp := <-c.userRepositoryQuery.FindOneByEmail(ctx, payload.Email)
	if resp.Error != nil {
		return nil, resp.Error
	}
	if resp.Data != nil {
		msg := "Email is already registered"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.BadRequest(msg)
	}

	countryId, err := strconv.Atoi(payload.CountryId)
	if err != nil {
		msg := "CountryId must integer"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.BadRequest(msg)
	}

	country := <-c.addressRepositoryQuery.FindOneCountry(ctx, countryId)
	if country.Error != nil {
		return nil, country.Error
	}
	if country.Data == nil {
		msg := "Country not found"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.CustomError(msg, 4002, http.StatusNotFound)
	}
	countryData, ok := country.Data.(*addressEntity.Country)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data country")
	}

	var subDistrictUser userEntity.Subdistrict
	if countryData.Code == "ID" {
		subdistrict := <-c.addressRepositoryQuery.FindOneSubdistrict(ctx, payload.SubdictrictId)
		if subdistrict.Error != nil {
			return nil, subdistrict.Error
		}
		if subdistrict.Data == nil {
			msg := "Subdistrict not found"
			c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			return nil, errors.CustomError(msg, 4002, http.StatusNotFound)
		}
		subdistrictData, ok := subdistrict.Data.(*addressEntity.SubDistrict)
		if !ok {
			return nil, errors.InternalServerError("cannot parsing data subdistrict")
		}
		if subdistrictData != nil {
			subDistrictUser = userEntity.Subdistrict{
				Id:           subdistrictData.Id,
				Name:         subdistrictData.Name,
				DistrictId:   subdistrictData.DistrictId,
				DistrictName: subdistrictData.DistrictName,
				CityId:       subdistrictData.CityId,
				CityName:     subdistrictData.CityName,
				ProvinceId:   subdistrictData.ProvinceId,
				ProvinceName: subdistrictData.ProvinceName,
			}
		}
	}

	passwordHash := helpers.GeneratePassword(payload.Password)

	user := userEntity.User{
		UserId:       uuid.New().String(),
		FullName:     payload.FullName,
		Email:        payload.Email,
		NIK:          payload.NIK,
		MobileNumber: helpers.VerifyPhoneNumber62(payload.MobileNumber),
		Password:     passwordHash,
		Subdistrict:  subDistrictUser,
		Country: userEntity.Country{
			Id:            countryData.Id,
			Code:          countryData.Code,
			Name:          countryData.Name,
			FullName:      countryData.FullName,
			ContinentId:   countryData.ContinentCode,
			ContinentName: countryData.ContinentName,
			Latitude:      payload.Latitude,
			Longitude:     payload.Longitude,
		},
		Address:   payload.Address,
		RtRw:      payload.RtRw,
		KKNumber:  payload.KKNumber,
		Role:      payload.Role,
		LoginAt:   time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	respUser := <-c.userRepositoryCommand.UpsertOneUserTemp(ctx, user)
	if respUser.Error != nil {
		return nil, respUser.Error
	}

	otp := helpers.GenerateRandomOtp()
	// Send kafka data

	kafkaData := struct {
		UserId   string `json:"userId"`
		FullName string `json:"fullName"`
		Email    string `json:"email"`
		Otp      string `json:"otp"`
	}{
		UserId:   user.UserId,
		FullName: payload.FullName,
		Email:    payload.Email,
		Otp:      string(otp),
	}
	marshaledKafkaData, _ := json.Marshal(kafkaData)
	otpTopic := "concert-send-otp-user-registration"
	c.kafkaProducer.Publish(otpTopic, marshaledKafkaData, nil)
	c.logger.Info(ctx, fmt.Sprintf("Send kafka email otp : %s, topic: %s", payload.Email, otpTopic), fmt.Sprintf("%+v", payload.Email))

	c.redis.Set(ctx, fmt.Sprintf("%s:%s", constants.RedisKeyOtpRegister, payload.Email), kafkaData.Otp, 3*time.Minute)

	return &userResponse.RegisterUser{
		Email: payload.Email,
	}, nil
}

func (c commandUsecase) VerifyRegisterUser(origCtx context.Context, payload userRequest.VerifyRegisterUser) (*userResponse.VerifyRegister, error) {
	domain := "userUsecase-VerifyRegisterUser"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	validEmail := helpers.IsEmailValid(payload.Email)
	if !validEmail {
		msg := "Incorrect email format"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.CustomError(msg, 4001, http.StatusBadRequest)
	}
	if helpers.IsBlacklistedEmail(payload.Email) {
		msg := "Email blacklist"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.CustomError(msg, 4002, http.StatusBadRequest)
	}
	checkedOtp, _ := c.redis.Get(ctx, fmt.Sprintf("%s:%s", constants.RedisKeyOtpRegister, payload.Email)).Result()
	if checkedOtp == "" {
		msg := "Otp expired"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.BadRequest(msg)
	}
	if checkedOtp != payload.Otp {
		msg := "Otp not match"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.CustomError(msg, 4003, http.StatusBadRequest)
	}

	// check data user in user temp
	resp := <-c.userRepositoryQuery.FindOneByEmailUserTemp(ctx, payload.Email)
	if resp.Error != nil {
		return nil, resp.Error
	}
	if resp.Data == nil {
		msg := "Email not found"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound(msg)
	}

	userData, ok := resp.Data.(*userEntity.User)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}
	userData.Status = "active"
	respUser := <-c.userRepositoryCommand.UpsertOneUser(ctx, *userData)
	if respUser.Error != nil {
		return nil, respUser.Error
	}
	c.redis.Del(ctx, fmt.Sprintf("%s:%s", constants.RedisKeyOtpRegister, payload.Email))

	return nil, nil
}

func (c commandUsecase) LoginUser(origCtx context.Context, payload userRequest.LoginUser) (*userResponse.LoginUserResp, error) {
	domain := "userUsecase-LoginUser"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	validEmail := helpers.IsEmailValid(payload.Email)
	if !validEmail {
		msg := "Incorrect email format"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.CustomError(msg, 4001, http.StatusBadRequest)
	}
	if helpers.IsBlacklistedEmail(payload.Email) {
		msg := "Email blacklist"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.CustomError(msg, 4002, http.StatusBadRequest)
	}

	// Get attempt from redis
	attempt, _ := c.redis.Get(ctx, fmt.Sprintf("%s:%s", constants.RedisKeyLoginAttempt, payload.Email)).Result()
	attemptInt, _ := strconv.Atoi(attempt)
	if attemptInt >= 5 {
		logMessage := "You have too many attempts, please wait 10 minutes"
		c.logger.Info(ctx, logMessage, fmt.Sprintf("%+v", payload.Email))
		return nil, errors.CustomError(logMessage, 4003, http.StatusBadRequest)
	}

	resp := <-c.userRepositoryQuery.FindOneByEmail(ctx, payload.Email)
	fmt.Println("resp: ", resp)
	if resp.Error != nil {
		return nil, resp.Error
	}
	if resp.Data == nil {
		logMessage := "email / password not found"
		c.logger.Info(ctx, logMessage, fmt.Sprintf("%+v", helpers.MaskEmail(payload.Email)))
		return nil, errors.BadRequest(logMessage)
	}
	userData, ok := resp.Data.(*userEntity.User)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}
	passwordHash := helpers.GeneratePassword(payload.Password)
	if passwordHash != userData.Password {
		// Set redis attempt
		attemptInt = attemptInt + 1
		c.redis.Set(ctx, fmt.Sprintf("%s:%s", constants.RedisKeyLoginAttempt, payload.Email), attemptInt, 10*time.Minute)

		logMessage := "Username / password not match"
		c.logger.Info(ctx, logMessage, fmt.Sprintf("%+v", helpers.MaskEmail(payload.Email)))
		return nil, errors.BadRequest(logMessage)
	}

	//
	respUser := <-c.userRepositoryCommand.UpsertOneUser(ctx, *userData)
	if respUser.Error != nil {
		return nil, respUser.Error
	}
	c.redis.Del(ctx, fmt.Sprintf("%s:%s", constants.RedisKeyLoginAttempt, payload.Email))
	// Generate token
	tokenPayload := map[string]interface{}{
		"userId": userData.UserId,
		"role":   userData.Role,
	}
	jwtToken, expiredAt, err := c.jwtHelper.GenerateToken(24*time.Hour, tokenPayload)
	if err != nil {
		return nil, err
	}
	refreshToken, err := c.jwtHelper.GenerateTokenRefresh((30*24)*time.Hour, tokenPayload)
	if err != nil {
		return nil, err
	}
	return &userResponse.LoginUserResp{
		AuthToken:    jwtToken,
		RefreshToken: refreshToken,
		ExpiredAt:    expiredAt,
	}, nil

}
