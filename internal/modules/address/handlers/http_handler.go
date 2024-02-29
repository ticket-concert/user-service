package handlers

import (
	"user-service/internal/modules/address"
	"user-service/internal/modules/address/models/request"
	"user-service/internal/pkg/errors"
	"user-service/internal/pkg/helpers"
	"user-service/internal/pkg/log"
	"user-service/internal/pkg/redis"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AddressHttpHandler struct {
	AddressUsecaseQuery address.UsecaseQuery
	Logger              log.Logger
	Validator           *validator.Validate
}

func InitAddressHttpHandler(app *fiber.App, auq address.UsecaseQuery, log log.Logger, redisClient redis.Collections) {
	handler := &AddressHttpHandler{
		AddressUsecaseQuery: auq,
		Logger:              log,
		Validator:           validator.New(),
	}
	// middlewares := middlewares.NewMiddlewares(redisClient)
	route := app.Group("/api/users/address")

	route.Get("/v1/provinces", handler.GetProvinces)
	route.Get("/v1/cities", handler.GetCities)
	route.Get("/v1/districts", handler.GetDistricts)
	route.Get("/v1/subdistricts", handler.GetSubDistricts)
	route.Get("/v1/countries", handler.GetCountries)
	route.Get("/v1/continent", handler.GetContinent)
}

func (a AddressHttpHandler) GetProvinces(c *fiber.Ctx) error {
	req := new(request.Province)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, a.Logger, errors.BadRequest("bad request"))
	}

	if err := a.Validator.Struct(req); err != nil {
		return helpers.RespError(c, a.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := a.AddressUsecaseQuery.FindProvinces(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, a.Logger, err)
	}
	return helpers.RespPagination(c, a.Logger, resp.CollectionData, resp.MetaData, "Get province success")
}

func (a AddressHttpHandler) GetCities(c *fiber.Ctx) error {
	req := new(request.City)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, a.Logger, errors.BadRequest("bad request"))
	}

	if err := a.Validator.Struct(req); err != nil {
		return helpers.RespError(c, a.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := a.AddressUsecaseQuery.FindCities(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, a.Logger, err)
	}
	return helpers.RespPagination(c, a.Logger, resp.CollectionData, resp.MetaData, "Get city success")
}

func (a AddressHttpHandler) GetDistricts(c *fiber.Ctx) error {
	req := new(request.District)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, a.Logger, errors.BadRequest("bad request"))
	}

	if err := a.Validator.Struct(req); err != nil {
		return helpers.RespError(c, a.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := a.AddressUsecaseQuery.FindDistricts(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, a.Logger, err)
	}
	return helpers.RespPagination(c, a.Logger, resp.CollectionData, resp.MetaData, "Get district success")
}

func (a AddressHttpHandler) GetSubDistricts(c *fiber.Ctx) error {
	req := new(request.SubDistrict)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, a.Logger, errors.BadRequest("bad request"))
	}

	if err := a.Validator.Struct(req); err != nil {
		return helpers.RespError(c, a.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := a.AddressUsecaseQuery.FindSubDistricts(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, a.Logger, err)
	}
	return helpers.RespPagination(c, a.Logger, resp.CollectionData, resp.MetaData, "Get subdistrict success")
}

func (a AddressHttpHandler) GetCountries(c *fiber.Ctx) error {
	req := new(request.Country)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, a.Logger, errors.BadRequest("bad request"))
	}

	if err := a.Validator.Struct(req); err != nil {
		return helpers.RespError(c, a.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := a.AddressUsecaseQuery.FindCountries(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, a.Logger, err)
	}
	return helpers.RespPagination(c, a.Logger, resp.CollectionData, resp.MetaData, "Get country success")
}

func (a AddressHttpHandler) GetContinent(c *fiber.Ctx) error {
	resp, err := a.AddressUsecaseQuery.FindContinent(c.Context())
	if err != nil {
		return helpers.RespCustomError(c, a.Logger, err)
	}
	return helpers.RespPagination(c, a.Logger, resp.CollectionData, resp.MetaData, "Get continent success")
}
