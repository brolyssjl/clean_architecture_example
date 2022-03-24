package user

import (
	"github.com/brolyssjl/clean_architecture_example/pkg"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{
		DB: db,
	}
}

func (r *repo) FindByID(c *fiber.Ctx, id uint) (user *User, err error) {
	return user, err
}

func (r *repo) BuildProfile(ctx *fiber.Ctx, user *User) (u *User, err error) {
	result := r.DB.Table("users").Where("email = ?", user.Email).Updates(map[string]interface{}{
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"phone_number": user.PhoneNumber,
		"address":      user.Address,
		"profile_pic":  user.ProfilePic,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	switch result.Error {
	case nil:
		return user, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) CreateMinimal(c *fiber.Ctx, email, password, phoneNumber string) (u *User, err error) {
	u = &User{
		Email:       email,
		Password:    password,
		PhoneNumber: phoneNumber,
	}
	result := r.DB.Create(u)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func (r *repo) FindByEmailAndPassword(c *fiber.Ctx, email, password string) (u *User, err error) {
	u = &User{}
	result := r.DB.Where("email = ? AND password = ?", email, password).First(u)

	switch result.Error {
	case nil:
		return u, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) DoesEmailExist(c *fiber.Ctx, email string) (bool, error) {
	u := &User{}
	query := r.DB.Where("email = ?", email).First(u)
	if query.Error == gorm.ErrRecordNotFound {
		return false, nil
	}

	return true, nil
}

func (r *repo) FindByEmail(c *fiber.Ctx, email string) (u *User, err error) {
	u = &User{}
	projection := "id, email, created_at, updated_at, deleted_at, phone_number, first_name, last_name, address, profile_pic"
	result := r.DB.Select(projection).Where("email = ?", email).First(u)

	switch result.Error {
	case nil:
		return u, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) ChangePassword(c *fiber.Ctx, email, password string) error {
	result := r.DB.Table("users").Where("email = ?", email).Update("password", password)

	switch result.Error {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	default:
		return pkg.ErrDatabase
	}
}
