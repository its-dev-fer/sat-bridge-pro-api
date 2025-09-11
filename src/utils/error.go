package utils

import (
	"app/src/response"
	"app/src/validation"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if errorsMap := validation.CustomErrorMessages(err); len(errorsMap) > 0 {
		return response.Error(c, fiber.StatusBadRequest, "Bad Request", errorsMap)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return response.Error(c, fiberErr.Code, fiberErr.Message, nil)
	}

	return response.Error(c, fiber.StatusInternalServerError, "Internal Server Error", nil)
}

func NotFoundHandler(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusNotFound, "Endpoint Not Found", nil)
}

type CommonErrorResponse struct {
    Code    int    `json:"code" example:"400"`
    Status  string `json:"status" example:"error"`
    Message string `json:"message" example:"Invalid request"`
}

type DetailedErrorResponse struct {
    Code    int         `json:"code" example:"422"`
    Status  string      `json:"status" example:"error"`
    Message string      `json:"message" example:"Validation failed"`
	Errors  interface{} `json:"errors"`

}


