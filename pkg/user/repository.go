package user

import "github.com/gofiber/fiber/v2"

type Repository interface {
	FindByID(ctx *fiber.Ctx, id uint) (*User, error)
	BuildProfile(ctx *fiber.Ctx, user *User) (*User, error)
	CreateMinimal(ctx *fiber.Ctx, email, password, phoneNumber string) (*User, error)
	FindByEmailAndPassword(ctx *fiber.Ctx, email, password string) (*User, error)
	FindByEmail(ctx *fiber.Ctx, email string) (*User, error)
	DoesEmailExist(ctx *fiber.Ctx, email string) (bool, error)
	ChangePassword(ctx *fiber.Ctx, email, password string) error
}
