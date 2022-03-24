package middleware

import (
	"log"

	"github.com/brolyssjl/clean_architecture_example/api/views"
	"github.com/brolyssjl/clean_architecture_example/pkg"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func Validate() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(viper.GetString("jwt_secret")),
		ErrorHandler: jwtError,
	})
}

func ValidateAndGetClaims(c *fiber.Ctx, role string) (map[string]interface{}, error) {
	token, err := jwt.Parse(c.Get("Authorization"), func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt_secret")), nil
	})

	if err != nil {
		log.Println("=================")
		log.Println(token)
		return nil, views.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		log.Println(">>>>>>>>>>>>>>>>>")
		log.Println(claims)
		return nil, views.ErrInvalidToken
	}

	if claims.Valid() != nil {
		log.Println("<<<<<<<<<<<<<<<<<")
		return nil, views.ErrInvalidToken
	}

	if claims["role"].(string) != role {
		log.Println(claims["role"])
		return nil, pkg.ErrUnauthorized
	}
	return claims, nil
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
