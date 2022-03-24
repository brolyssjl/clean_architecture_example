package handler

import (
	"github.com/brolyssjl/clean_architecture_example/api/middleware"
	"github.com/brolyssjl/clean_architecture_example/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func register(svc user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user user.User

		if err := c.BodyParser(&user); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "invalid request",
				"error":   err.Error(),
			})
		}

		userRegistered, err := svc.Register(c, user.Email, user.Password, user.PhoneNumber)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "we cannot register this user",
				"error":   err.Error(),
			})
		}

		c.Status(fiber.StatusCreated)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": userRegistered.Email,
			"id":    userRegistered.ID,
			"role":  "user",
		})
		tokenString, err := token.SignedString([]byte(viper.GetString("jwt_secret")))
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.JSON(fiber.Map{
				"message": "forbidden",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"token": tokenString,
			"user":  userRegistered,
		})
	}
}

func login(svc user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user user.User

		if err := c.BodyParser(&user); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "invalid request",
				"error":   err.Error(),
			})
		}

		userLoggedIn, err := svc.Login(c, user.Email, user.Password)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.JSON(fiber.Map{
				"message": "forbidden",
				"error":   err.Error(),
			})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": userLoggedIn.Email,
			"id":    userLoggedIn.ID,
			"role":  "user",
		})
		tokenString, err := token.SignedString([]byte(viper.GetString("jwt_secret")))
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.JSON(fiber.Map{
				"message": "forbidden",
				"error":   err.Error(),
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"token": tokenString,
			"user":  userLoggedIn,
		})
	}
}

func profile(svc user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodPost {
			var user user.User

			if err := c.BodyParser(&user); err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"message": "invalid request",
					"error":   err.Error(),
				})
			}

			claims, err := middleware.ValidateAndGetClaims(c, "user")
			if err != nil {
				c.Status(fiber.StatusUnauthorized)
				return c.JSON(fiber.Map{
					"message": "unathorized",
					"error":   err.Error(),
				})
			}

			user.Email = claims["email"].(string)
			userProfile, err := svc.BuildProfile(c, &user)
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"message": "bad request",
					"error":   err.Error(),
				})
			}

			c.Status(fiber.StatusOK)
			return c.JSON(userProfile)
		} else if c.Method() == fiber.MethodGet {
			claims, err := middleware.ValidateAndGetClaims(c, "user")
			if err != nil {
				c.Status(fiber.StatusUnauthorized)
				return c.JSON(fiber.Map{
					"message": "unathorized",
					"error":   err.Error(),
				})
			}
			userProfile, err := svc.GetUserProfile(c, claims["email"].(string))
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"message": "bad request",
					"error":   err.Error(),
				})
			}

			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"message": "user profile",
				"data":    userProfile,
			})
		}

		c.Status(fiber.StatusMethodNotAllowed)
		return c.JSON(fiber.Map{
			"error": "method not allowed",
		})
	}
}

func changePassword(svc user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user user.User

		if err := c.BodyParser(&user); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "invalid request",
				"error":   err.Error(),
			})
		}

		claims, err := middleware.ValidateAndGetClaims(c, "user")
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"message": "unathorized",
				"error":   err.Error(),
			})
		}
		if err := svc.ChangePassword(c, claims["email"].(string), user.Password); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "we couldn't change password",
				"error":   err.Error(),
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": "password changed",
		})
	}
}

func MakeUserHandler(app *fiber.App, svc user.Service) {
	v1 := app.Group("/v1")
	app.Use(logger.New())

	v1.Get("/user/ping", func(c *fiber.Ctx) error {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})
	v1.Post("/user/register", register(svc))
	v1.Post("/user/login", login(svc))
	v1.Get("/user/profile", profile(svc), middleware.Validate())
	v1.Post("/user/profile", profile(svc), middleware.Validate())
	v1.Patch("/user/pwd", changePassword(svc), middleware.Validate())
}
