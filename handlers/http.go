package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func New(code int, msg ...string) *fiber.Error {
	return fiber.NewError(code, msg...)
}

func NotFound() *fiber.Error {
	return fiber.NewError(fiber.StatusNotFound)
}

func BusinessError(msg ...string) *fiber.Error {
	return fiber.NewError(
		http.StatusInternalServerError,
		msg...,
	)
}

func BusinessErrorf(msg string, params ...any) *fiber.Error {
	return fiber.NewError(
		http.StatusInternalServerError,
		fmt.Sprintf(msg, params),
	)
}

func ParameterError(msg ...string) *fiber.Error {
	return fiber.NewError(
		http.StatusBadRequest,
		msg...,
	)
}

func ParameterErrorf(msg string, params ...any) *fiber.Error {
	return fiber.NewError(
		http.StatusBadRequest,
		fmt.Sprintf(msg, params),
	)
}

func UnknownError(msg ...string) *fiber.Error {
	return fiber.NewError(
		http.StatusForbidden,
		msg...,
	)
}

func UnknownErrorf(msg string, params ...any) *fiber.Error {
	return fiber.NewError(
		http.StatusForbidden,
		fmt.Sprintf(msg, params),
	)
}
