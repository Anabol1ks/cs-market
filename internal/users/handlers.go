package users

import (
	"cs-market/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// GetUserProfileHandler godoc
// @Summary Получение профиля пользователя
// @Description Получение информации о своём профиле пользователя
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User "Информация о профиле"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Router /profile [get]
func GetUserProfileHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	var user User
	if err := storage.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}
