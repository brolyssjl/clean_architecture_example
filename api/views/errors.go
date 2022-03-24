package views

import (
	"errors"

	"github.com/brolyssjl/clean_architecture_example/pkg"
	"github.com/gofiber/fiber/v2"
)

type ErrView struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

var (
	ErrMethodNotAllowed = errors.New("error: Method is not allowed")
	ErrInvalidToken     = errors.New("error: Invalid Authorization token")
	ErrUserExists       = errors.New("user already exists")
)

var ErrHTTPStatusMap = map[string]int{
	pkg.ErrNotFound.Error():     fiber.StatusNotFound,
	pkg.ErrInvalidSlug.Error():  fiber.StatusBadRequest,
	pkg.ErrExists.Error():       fiber.StatusConflict,
	pkg.ErrNoContent.Error():    fiber.StatusNotFound,
	pkg.ErrDatabase.Error():     fiber.StatusInternalServerError,
	pkg.ErrUnauthorized.Error(): fiber.StatusUnauthorized,
	pkg.ErrForbidden.Error():    fiber.StatusForbidden,
	ErrMethodNotAllowed.Error(): fiber.StatusMethodNotAllowed,
	ErrInvalidToken.Error():     fiber.StatusBadRequest,
	ErrUserExists.Error():       fiber.StatusConflict,
}
