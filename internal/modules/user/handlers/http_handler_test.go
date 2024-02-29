// user_http_handler_test.go

package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"user-service/internal/modules/user/handlers"
	userRequest "user-service/internal/modules/user/models/request"
	userResponse "user-service/internal/modules/user/models/response"
	"user-service/internal/pkg/errors"
	mockcert "user-service/mocks/modules/user"
	mocklog "user-service/mocks/pkg/log"
	mockredis "user-service/mocks/pkg/redis"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type UserHttpHandlerTestSuite struct {
	suite.Suite

	cUC       *mockcert.UsecaseCommand
	cUQ       *mockcert.UsecaseQuery
	cLog      *mocklog.Logger
	validator *validator.Validate
	cRedis    *mockredis.Collections
	handler   *handlers.UserHttpHandler
	app       *fiber.App
}

func (suite *UserHttpHandlerTestSuite) SetupTest() {
	suite.cUC = new(mockcert.UsecaseCommand)
	suite.cUQ = new(mockcert.UsecaseQuery)
	suite.cLog = new(mocklog.Logger)
	suite.validator = validator.New()
	suite.cRedis = new(mockredis.Collections)
	suite.handler = &handlers.UserHttpHandler{
		UserUsecaseCommand: suite.cUC,
		UserUsecaseQuery:   suite.cUQ,
		Logger:             suite.cLog,
		Validator:          suite.validator,
	}
	suite.app = fiber.New()
	handlers.InitUserHttpHandler(suite.app, suite.cUC, suite.cUQ, suite.cLog, suite.cRedis)
}

func TestUserHttpHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHttpHandlerTestSuite))
}

func (suite *UserHttpHandlerTestSuite) TestRegisterUser() {
	suite.cUC.On("RegisterUser", mock.Anything, mock.Anything).Return(&userResponse.RegisterUser{
		Email: "alif@gmail.com",
	}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.RegisterUser{
		FullName:      "Alif Septian",
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "1",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/register")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.RegisterUser(ctx)
	assert.Nil(suite.T(), err)
}
func (suite *UserHttpHandlerTestSuite) TestRegisterUserErrorBodyParse() {
	suite.cUC.On("RegisterUser", mock.Anything, mock.Anything).Return(&userResponse.RegisterUser{
		Email: "alif@gmail.com",
	}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.RegisterUser{}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/register")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")

	// block this code for error body parse
	// ctx.Request().SetBody(requestBody)

	err := suite.handler.RegisterUser(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestRegisterUserErrorValidation() {
	suite.cUC.On("RegisterUser", mock.Anything, mock.Anything).Return(&userResponse.RegisterUser{
		Email: "alif@gmail.com",
	}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.RegisterUser{}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/register")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.RegisterUser(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestRegisterUserError() {
	suite.cUC.On("RegisterUser", mock.Anything, mock.Anything).Return(nil, errors.InternalServerError("Error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.RegisterUser{
		FullName:      "Alif Septian",
		Email:         "alif@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "081281015121",
		ProvinceId:    "123",
		CityId:        "123",
		DistrictId:    "123",
		SubdictrictId: "123",
		CountryId:     "1",
		Address:       "Jalan jalan",
		RtRw:          "12/12",
		Role:          "user",
		KKNumber:      "1212121212",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/register")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.RegisterUser(ctx)
	assert.Nil(suite.T(), err)
}

// Verify register user

func (suite *UserHttpHandlerTestSuite) TestVerifyRegisterUser() {
	var response *userResponse.VerifyRegister

	suite.cUC.On("VerifyRegisterUser", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.VerifyRegisterUser{
		Email: "alif@gmail.com",
		Otp:   "123456",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/otp/submit", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/otp/submit")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.VerifyRegisterUser(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestUpdateUser() {
	var response string

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")

	suite.cUC.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPut, "/v1/profile", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx.Request().SetRequestURI("/v1/profile")
	ctx.Request().Header.SetMethod(fiber.MethodPut)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.UpdateUser(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestUpdateUserErrBodyParser() {
	var response string

	suite.cUC.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPut, "/v1/profile", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/profile")
	ctx.Request().Header.SetMethod(fiber.MethodPut)
	ctx.Request().Header.SetContentType("application/json")
	// ctx.Request().SetBody(requestBody)

	err := suite.handler.UpdateUser(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestUpdateUserErrValidation() {
	var response string

	suite.cUC.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.UpdateUser{}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPut, "/v1/profile", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/profile")
	ctx.Request().Header.SetMethod(fiber.MethodPut)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.UpdateUser(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestLoginUser() {
	var response *userResponse.LoginUserResp

	suite.cUC.On("LoginUser", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.LoginUser{
		Email:    "email",
		Password: "password",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/login")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.Login(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestUpdateUserError() {
	suite.cUC.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return("", errors.InternalServerError("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.UpdateUser{
		FullName:      "FullName",
		MobileNumber:  "+6281281015121",
		Address:       "Address",
		ProvinceId:    "ProvinceId",
		CityId:        "CityId",
		DistrictId:    "DistrictId",
		SubdictrictId: "SubdictrictId",
		CountryId:     "1",
		RtRw:          "RtRw",
		Role:          "user",
		Latitude:      "Latitude",
		Longitude:     "Longitude",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPut, "/v1/profile", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/profile")
	ctx.Request().Header.SetMethod(fiber.MethodPut)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.UpdateUser(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestLoginUserError() {
	suite.cUC.On("LoginUser", mock.Anything, mock.Anything).Return(nil, errors.InternalServerError("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.LoginUser{
		Email:    "email",
		Password: "password",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/login")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.Login(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestLoginUserErrValidation() {
	var response *userResponse.LoginUserResp

	suite.cUC.On("LoginUser", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.LoginUser{}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/login")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.Login(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestLoginUserErrBodyParser() {
	var response *userResponse.LoginUserResp

	suite.cUC.On("LoginUser", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.LoginUser{
		Email:    "email",
		Password: "password",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/login")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	// ctx.Request().SetBody(requestBody)

	err := suite.handler.Login(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestVerifyUserErrorBodyParse() {
	var response *userResponse.VerifyRegister

	suite.cUC.On("VerifyRegisterUser", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.VerifyRegisterUser{
		Email: "alif@gmail.com",
		Otp:   "123456",
	}

	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/otp/submit", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/otp/submit")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")

	// block this code for error body parse
	// ctx.Request().SetBody(requestBody)

	err := suite.handler.VerifyRegisterUser(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestVerifyUserErrorValidation() {
	var response *userResponse.VerifyRegister

	suite.cUC.On("VerifyRegisterUser", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.VerifyRegisterUser{}

	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/otp/submit", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/otp/submit")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.VerifyRegisterUser(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestVerifyUserError() {
	suite.cUC.On("VerifyRegisterUser", mock.Anything, mock.Anything).Return(nil, errors.InternalServerError("Error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.VerifyRegisterUser{
		Email: "alif@gmail.com",
		Otp:   "123456",
	}

	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/otp/submit", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/otp/submit")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.VerifyRegisterUser(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetProfile() {
	var response *userResponse.GetProfile

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")

	suite.cUQ.On("GetProfile", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/get-profile", nil)
	req.Header.Set("Content-Type", "application/json")

	err := suite.handler.GetProfile(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetProfileError() {
	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")

	suite.cUQ.On("GetProfile", mock.Anything, mock.Anything).Return(nil, errors.InternalServerError("Error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/get-profile", nil)
	req.Header.Set("Content-Type", "application/json")

	err := suite.handler.GetProfile(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetProfileErrorParse() {
	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", 12345)

	suite.cUQ.On("GetProfile", mock.Anything, mock.Anything).Return(nil, errors.InternalServerError("Error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/get-profile", nil)
	req.Header.Set("Content-Type", "application/json")

	err := suite.handler.GetProfile(ctx)
	assert.Nil(suite.T(), err)
}
