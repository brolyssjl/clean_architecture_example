package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Service interface {
	Register(ctx *fiber.Ctx, email, password, phoneNumber string) (*User, error)
	Login(ctx *fiber.Ctx, email, password string) (*User, error)
	ChangePassword(ctx *fiber.Ctx, email, password string) error
	BuildProfile(ctx *fiber.Ctx, user *User) (*User, error)
	GetUserProfile(ctx *fiber.Ctx, email string) (*User, error)
	IsValid(user *User) (bool, error)
	GetRepo() Repository
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) Register(ctx *fiber.Ctx, email, password, phoneNumber string) (u *User, err error) {
	exists, err := s.repo.DoesEmailExist(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("User already exists")
	}

	hasher := md5.New()
	hasher.Write([]byte(password))

	return s.repo.CreateMinimal(ctx, email, hex.EncodeToString(hasher.Sum(nil)), phoneNumber)
}

func (s *service) Login(ctx *fiber.Ctx, email, password string) (u *User, err error) {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return s.repo.FindByEmailAndPassword(ctx, email, hex.EncodeToString(hasher.Sum(nil)))
}

func (s *service) ChangePassword(ctx *fiber.Ctx, email, password string) (err error) {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return s.repo.ChangePassword(ctx, email, hex.EncodeToString(hasher.Sum(nil)))
}

func (s *service) BuildProfile(ctx *fiber.Ctx, user *User) (u *User, err error) {
	return s.repo.BuildProfile(ctx, user)
}

func (s *service) GetUserProfile(ctx *fiber.Ctx, email string) (u *User, err error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *service) IsValid(user *User) (ok bool, err error) {
	return ok, err
}

func (s *service) GetRepo() Repository {
	return s.repo
}
