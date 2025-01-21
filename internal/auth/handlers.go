package auth

import (
	"cs-market/internal/storage"
	"cs-market/internal/users"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

// SteamLoginHandler godoc
// @Summary Авторизация через Steam
// @Description Авторизация через Steam и получение токенов доступа (для теста требуется подключение в steam хоста с https)
// @Tags auth
// @Accept json
// @Produce json
// @Success 303 {string} string "Redirect URL"
// @Failure 400 {object} response.ErrorResponse "Ошибка начала авторизации Steam"
// @Router /auth/steam [get]
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

// SteamCallbackHandler godoc
// @Summary Обработчик коллбэка после авторизации через Steam
// @Description Обработчик коллбэка после авторизации через Steam и получение токенов доступа
// @Tags auth
// @Accept json
// @Produce json
// @Success 303 {string} string "Ссылка с указанием токенов доступа"
// @Failure 400 {object} response.ErrorResponse "Ошибка авторизации Steam"
// @Failure 500 {object} response.ErrorResponse "Ошибка генерации токенов" "Ошибка входа в систему"
// @Router /auth/steam/callback [get]
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

	accessToken, refreshToken, err := GenerateTokensJWT(user.SteamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	frontUrl := os.Getenv("FRONT_URL")
	redirectURL := frontUrl + "/auth?" +
		"access_token=" + accessToken + "&refresh_token=" + refreshToken

	c.Redirect(http.StatusSeeOther, redirectURL)
}

var jwtSecret = []byte(os.Getenv("JWT_KEY"))

var jwtSecretRefresh = []byte(os.Getenv("JWT_KEY_REFRESH"))

func GenerateTokensJWT(stramID string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": stramID,
		// "exp":     time.Now().Add(15 * time.Minute).Unix(),
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": stramID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(jwtSecretRefresh)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshTokenHandler godoc
// @Summary Обновление токена доступа
// @Description Обновление токена доступа с помощью refresh_token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.TokenResponse "Успешное обновление токена доступа"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения refresh_token из куки"
// @Failure 401 {object} response.ErrorResponse "Неверный refresh_token"
// @Failure 500 {object} response.ErrorResponse "Ошибка генерации токенов"
// @Router /auth/refresh [post]
func RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка получения refresh_token из куки"})
		return
	}

	claims, err := ValidToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный refresh_token"})
		return
	}

	accsessToken, refreshToken, err := GenerateTokensJWT(claims["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токенов"})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"access_token": accsessToken})
}

func TokenProv(c *gin.Context) {
	userID := c.GetString("user_id")

	c.JSON(http.StatusOK, gin.H{
		"message": "Доступ разрешён",
		"user_id": userID,
	})
}

// @Security BearerAuth
// VerifyTokenHandler godoc
// @Summary Проверка токена доступа
// @Description Проверка токена доступа
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.TokenResponse "Токен доступа действителен"
// @Failure 401 {object} response.ErrorResponse "Невалидный токен"
// @Router /auth/verify [get]
func VerifyTokenHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"valid": true, "user_id": userID})
}
