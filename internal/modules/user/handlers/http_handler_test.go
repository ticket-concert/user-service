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
		Email: "irmanjuliansyah@gmail.com",
	}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.RegisterUser{
		FullName:      "irmanjuliansyah",
		Email:         "irmanjuliansyah@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "+6281281015121",
		Address:       "<string>",
		ProvinceId:    "<string>",
		CityId:        "<string>",
		DistrictId:    "<string>",
		SubdictrictId: "<string>",
		RtRw:          "<string>",
		Role:          "user",
		KKNumber:      "<string>",
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
		Email: "irmanjuliansyah@gmail.com",
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
		Email: "irmanjuliansyah@gmail.com",
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
	suite.cUC.On("RegisterUser", mock.Anything, mock.Anything).Return(&userResponse.RegisterUser{
		Email: "irmanjuliansyah@gmail.com",
	}, errors.InternalServerError("Error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.RegisterUser{
		FullName:      "irmanjuliansyah",
		Email:         "irmanjuliansyah@gmail.com",
		Password:      "Password1@",
		NIK:           "12312131131",
		MobileNumber:  "+6281281015121",
		Address:       "<string>",
		ProvinceId:    "<string>",
		CityId:        "<string>",
		DistrictId:    "<string>",
		SubdictrictId: "<string>",
		RtRw:          "<string>",
		Role:          "user",
		KKNumber:      "<string>",
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
		Email: "irmanjuliansyah@gmail.com",
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
func (suite *UserHttpHandlerTestSuite) TestVerifyUserErrorBodyParse() {
	var response *userResponse.VerifyRegister

	suite.cUC.On("VerifyRegisterUser", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := userRequest.VerifyRegisterUser{
		Email: "irmanjuliansyah@gmail.com",
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
		Email: "irmanjuliansyah@gmail.com",
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
