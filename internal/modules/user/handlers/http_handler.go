package handlers

import (
	user "user-service/internal/modules/user"
	userRequest "user-service/internal/modules/user/models/request"
	"user-service/internal/pkg/errors"
	"user-service/internal/pkg/helpers"
	"user-service/internal/pkg/log"
	"user-service/internal/pkg/redis"

	middlewares "user-service/configs/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHttpHandler struct {
	UserUsecaseCommand user.UsecaseCommand
	UserUsecaseQuery   user.UsecaseQuery
	Logger             log.Logger
	Validator          *validator.Validate
}

func InitUserHttpHandler(app *fiber.App, uuc user.UsecaseCommand, uuq user.UsecaseQuery, log log.Logger, redisClient redis.Collections) {
	handler := &UserHttpHandler{
		UserUsecaseCommand: uuc,
		UserUsecaseQuery:   uuq,
		Logger:             log,
		Validator:          validator.New(),
	}
	middlewares := middlewares.NewMiddlewares(redisClient)
	route := app.Group("/api/users")
	route.Post("/v1/register", middlewares.VerifyBasicAuth(), handler.RegisterUser)
	route.Post("/v1/otp/submit", middlewares.VerifyBasicAuth(), handler.VerifyRegisterUser)
	route.Post("/v1/login", middlewares.VerifyBasicAuth(), handler.Login)
	route.Put("/v1/profile", middlewares.VerifyBearer(), handler.UpdateUser)
	route.Get("/v1/profile", middlewares.VerifyBearer(), handler.GetProfile)
}

func (u UserHttpHandler) UpdateUser(c *fiber.Ctx) error {
	req := new(userRequest.UpdateUser)
	if err := c.BodyParser(req); err != nil {
		return helpers.RespError(c, u.Logger, errors.BadRequest("bad request"))
	}
	if err := u.Validator.Struct(req); err != nil {
		return helpers.RespError(c, u.Logger, errors.BadRequest(err.Error()))
	}
	userId := c.Locals("userId").(string)
	resp, err := u.UserUsecaseCommand.UpdateUser(c.Context(), *req, userId)
	if err != nil {
		return helpers.RespCustomError(c, u.Logger, err)
	}
	return helpers.RespSuccess(c, u.Logger, resp, "Update user success")
}

func (u UserHttpHandler) RegisterUser(c *fiber.Ctx) error {
	req := new(userRequest.RegisterUser)
	if err := c.BodyParser(req); err != nil {
		return helpers.RespError(c, u.Logger, errors.BadRequest("bad request"))
	}

	if err := u.Validator.Struct(req); err != nil {
		return helpers.RespError(c, u.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := u.UserUsecaseCommand.RegisterUser(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, u.Logger, err)
	}
	return helpers.RespSuccess(c, u.Logger, resp, "Register user success")
}

func (u UserHttpHandler) VerifyRegisterUser(c *fiber.Ctx) error {
	req := new(userRequest.VerifyRegisterUser)
	if err := c.BodyParser(req); err != nil {
		return helpers.RespError(c, u.Logger, errors.BadRequest("bad request"))
	}

	if err := u.Validator.Struct(req); err != nil {
		return helpers.RespError(c, u.Logger, errors.BadRequest("validation error"))
	}
	resp, err := u.UserUsecaseCommand.VerifyRegisterUser(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, u.Logger, err)
	}
	return helpers.RespSuccess(c, u.Logger, resp, "Verify register user success")
}

func (u UserHttpHandler) Login(c *fiber.Ctx) error {
	req := new(userRequest.LoginUser)
	if err := c.BodyParser(req); err != nil {
		return helpers.RespError(c, u.Logger, errors.BadRequest("bad request"))
	}
	if err := u.Validator.Struct(req); err != nil {
		return helpers.RespError(c, u.Logger, errors.BadRequest("validation error"))
	}

	resp, err := u.UserUsecaseCommand.LoginUser(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, u.Logger, err)
	}
	return helpers.RespSuccess(c, u.Logger, resp, "Login user success")
}

func (u UserHttpHandler) GetProfile(c *fiber.Ctx) error {
	req := new(userRequest.GetProfile)
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return helpers.RespError(c, u.Logger, errors.BadRequest("bad request"))
	}

	req.UserId = userId
	resp, err := u.UserUsecaseQuery.GetProfile(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, u.Logger, err)
	}
	return helpers.RespSuccess(c, u.Logger, resp, "Get Profile Success")
}
