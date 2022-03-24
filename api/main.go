package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brolyssjl/clean_architecture_example/api/handler"
	"github.com/brolyssjl/clean_architecture_example/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err.Error())
	}
	os.Setenv("secret", viper.GetString("jwt_secret"))
}

func dbConnect(host, port, user, dbname, password string) (*gorm.DB, error) {
	/*if os.Getenv("ENV") == "local" {
		return gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
	}*/

	db, err := gorm.Open(mysql.Open(fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, dbname,
	)))

	return db, err
}

func main() {
	env := viper.GetString("mode")

	// DB binding
	dbprefix := "database_" + env
	dbhost := viper.GetString(dbprefix + ".host")
	dbport := viper.GetString(dbprefix + ".port")
	dbuser := viper.GetString(dbprefix + ".user")
	dbname := viper.GetString(dbprefix + ".dbname")
	dbpassword := viper.GetString(dbprefix + ".password")

	db, err := dbConnect(dbhost, dbport, dbuser, dbname, dbpassword)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err.Error())
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	log.Println("Connected to the database")

	// migrations
	db.AutoMigrate(&user.User{})

	// initializing repos and services
	userRepo := user.NewPostgresRepo(db)
	userSvc := user.NewService(userRepo)

	// Initializing handlers
	app := fiber.New()

	//Middlewares
	app.Use(cors.New())

	handler.MakeUserHandler(app, userSvc)

	app.Get("/", func(c *fiber.Ctx) error {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": "the server is up and running :)",
		})
	})

	serverPrefix := "server_" + env
	port := viper.GetString(serverPrefix + ".port")

	log.Printf("Starting in %s mode", env)
	log.Printf("Server running on %s", port)
	log.Fatal(app.Listen(":" + port))
}
