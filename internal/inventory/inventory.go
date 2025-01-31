package inventory

import (
	"cs-market/internal/storage"
	"cs-market/internal/users"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// GetMyInventoryHandler godoc
// @Summary Получение инвентаря пользователя
// @Description Получение инвентаря пользователя
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse "Информация об инвентаре"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Router /profile/inventory [get]
func GetMyInventoryHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	var user users.User
	if err := storage.DB.Where("steam_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Получение инвентаря пользователя
	url := fmt.Sprintf("https://steamcommunity.com/inventory/%s/730/2?l=english&count=5000", user.SteamID)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения инвентаря"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения ответа"})
		return
	}

	date, err := ParseInventory(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка парсинга инвентаря"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"inventory": date})
}

type Inventory struct {
	Assets []struct {
		AssetID string `json:"assetid"`
		ClassID string `json:"classid"`
	} `json:"assets"`
	Descriptions []struct {
		ClassID    string `json:"classid"`
		MarketName string `json:"market_name"`
		IconURL    string `json:"icon_url"`
	} `json:"descriptions"`
}

func ParseInventory(data []byte) (*Inventory, error) {
	var inv Inventory
	if err := json.Unmarshal(data, &inv); err != nil {
		return nil, err
	}

	// Добавляем префикс к IconURL
	const iconPrefix = "https://community.cloudflare.steamstatic.com/economy/image/"
	for i := range inv.Descriptions {
		inv.Descriptions[i].IconURL = iconPrefix + inv.Descriptions[i].IconURL
	}

	return &inv, nil
}
