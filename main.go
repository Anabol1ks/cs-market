package main

import (
	"cs-market/internal/auth"
	"cs-market/internal/storage"
	"cs-market/internal/users"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	key := os.Getenv("TEST_ENV")
	if key == "" {
		if key == "" {
			log.Println("\nПеременной среды нет, используется .env")
			// Загружаем .env
			err := godotenv.Load()
			if err != nil {
				log.Fatal("Ошибка загрузки .env файла")
			}
		}
	}

	storage.ConnectDatabase()

	err := storage.DB.AutoMigrate(&users.User{})
	if err != nil {
		log.Fatal("Ошибка миграции: ", err)
	}

	auth.InitAuth()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Укажи адрес фронтенда React
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/auth/steam", auth.SteamLoginHandler)
	r.GET("/auth/steam/callback", auth.SteamCallbackHandler)
	r.POST("/auth/refresh", auth.RefreshTokenHandler)

	authorized := r.Group("/")
	{
		authorized.Use(auth.AuthMiddleware())
		authorized.GET("/authMud", auth.TokenProv)
	}
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
