package inventory

import (
	"cs-market/internal/storage"
	"cs-market/internal/users"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	marketableItems := make([]interface{}, 0) // Используем interface{} для универсальности

	for _, desc := range date.Descriptions {
		if desc.Marketable == 1 && desc.Tradable == 1 {
			// Собираем только продаваемые предметы
			marketableItems = append(marketableItems, desc)
		}
	}

	// Отправляем фильтрованный список
	c.JSON(http.StatusOK, gin.H{"inventory": marketableItems})
}

type Inventory struct {
	Assets []struct {
		AssetID string `json:"assetid"`
		ClassID string `json:"classid"`
	} `json:"assets"`
	Descriptions []struct {
		ClassID    string   `json:"classid"`
		MarketName string   `json:"market_name"`
		IconURL    string   `json:"icon_url"`
		Price      *float64 `json:"price"`
		Marketable int      `json:"marketable"`
		Tradable   int      `json:"tradable"`
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

	// Собираем все market_hash_name для запроса в базу данных
	var marketNames []string
	for _, desc := range inv.Descriptions {
		marketNames = append(marketNames, desc.MarketName)
	}

	// Запрос в базу данных для получения всех скинов
	var skins []Skin
	if err := storage.DB.Where("market_hash_name IN ?", marketNames).Find(&skins).Error; err != nil {
		return nil, fmt.Errorf("ошибка при получении скинов из базы данных: %w", err)
	}

	// Преобразуем скины в мапу для быстрого поиска по MarketName
	skinMap := make(map[string]Skin)
	for _, skin := range skins {
		skinMap[skin.MarketHashName] = skin
	}

	// Присваиваем цену каждому элементу, если скин найден в базе данных
	for i := range inv.Descriptions {
		if skin, exists := skinMap[inv.Descriptions[i].MarketName]; exists {
			inv.Descriptions[i].Price = skin.MinPrice
		} else {
			// Если скин не найден в базе данных, устанавливаем цену как nil
			inv.Descriptions[i].Price = nil
		}
	}

	return &inv, nil
}

// https://api.skinport.com/v1/items?app_id=730&currency=RUB&tradable=0

func UpdatePrices(db *gorm.DB) {
	url := "https://api.skinport.com/v1/items?app_id=730&currency=RUB&tradable=0"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Accept-Encoding", "br")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Декодируем Brotli-сжатый ответ
	brReader := brotli.NewReader(resp.Body)
	body, err := io.ReadAll(brReader)
	if err != nil {
		fmt.Println("Ошибка при декодировании Brotli:", err)
		return
	}

	var skins []Skin
	if err := json.Unmarshal(body, &skins); err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		fmt.Println("Ответ сервера:", string(body)) // Логируем, если вдруг структура неверная
		return
	}

	for _, skin := range skins {
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "market_hash_name"}},
			DoUpdates: clause.AssignmentColumns([]string{"min_price", "avg_price", "max_price", "updated_at"}),
		}).Create(&skin)
	}

	fmt.Println("Цены обновлены:", time.Now())
}

func StartPriceUpdater(db *gorm.DB) {
	go func() {
		for {
			UpdatePrices(db)
			time.Sleep(10 * time.Minute)
		}
	}()
}
