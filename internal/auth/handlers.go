package auth

import (
	"cs-market/internal/storage"
	"cs-market/internal/users"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/steam"
)

func InitAuth() {
	steamKey := os.Getenv("STEAM_API_KEY")
	callbackURL := os.Getenv("CALLBACK_URL")

	log.Printf("Initializing Steam auth with key: %s and callback: %s", steamKey, callbackURL)

	goth.UseProviders(
		steam.New(steamKey, callbackURL),
	)
}

func SteamLoginHandler(c *gin.Context) {
	provider, err := goth.GetProvider("steam")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка получения провайдера Steam"})
		return
	}

	session, err := provider.BeginAuth("")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка начала авторизации Steam"})
		return
	}

	authURL, err := session.GetAuthURL()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка получения URL авторизации"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func SteamCallbackHandler(c *gin.Context) {
	provider, err := goth.GetProvider("steam")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка получения провайдера Steam"})
		return
	}

	session, err := provider.BeginAuth("")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка начала авторизации Steam"})
		return
	}

	_, err = session.Authorize(provider, c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка авторизации Steam"})
		return
	}

	steamUser, err := provider.FetchUser(session)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка получения данных Steam"})
		return
	}

	user := users.User{
		SteamID:   steamUser.UserID,
		Username:  steamUser.NickName,
		AvatarURL: steamUser.AvatarURL,
	}

	result := storage.DB.Where(users.User{SteamID: steamUser.UserID}).FirstOrCreate(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка входа в систему"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
