// controllers/order_controller.go
package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/perfume-api/config"
	"github.com/yourusername/perfume-api/models"
	"github.com/yourusername/perfume-api/utils/rabbitmq"
)

type OrderItemInput struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type OrderRequest struct {
	Items   []OrderItemInput `json:"items" binding:"required"`
	Address string           `json:"address" binding:"required"`
}

func CreateOrder(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input OrderRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}
	if len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Danh sách sản phẩm không được rỗng"})
		return
	}

	var total float64
	var orderItems []models.OrderItem

	for _, item := range input.Items {
		var product models.Product
		if err := config.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Sản phẩm không tồn tại: ProductID " + fmt.Sprintf("%d", item.ProductID)})
			return
		}

		total += product.Price * float64(item.Quantity)

		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	order := models.Order{
		UserID:  userID,
		Total:   total,
		Address: input.Address, // <-- GÁN ĐỊA CHỈ TỪ INPUT VÀO ORDER MODEL
		Items:   orderItems,
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo đơn hàng: " + err.Error()})
		return
	}

	// Xoá giỏ hàng sau khi đặt (nếu có)
	config.DB.Where("user_id = ?", userID).Delete(&models.Cart{})

	// Gửi thông tin đơn hàng qua RabbitMQ
	payload := map[string]interface{}{
		"order_id": order.ID,
		"user_id":  order.UserID,
		"total":    order.Total,
		"address":  input.Address,
	}
	rabbitmq.Publish("delivery-ex", payload)
	fmt.Println("📤 Đã gửi thông tin đơn hàng vào RabbitMQ:", payload)

	c.JSON(http.StatusOK, gin.H{"message": "Đặt hàng thành công", "order": order})
}

func GetOrders(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var orders []models.Order
	err := config.DB.
		Preload("Items").
		Where("user_id = ?", userID).
		Find(&orders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy lịch sử đơn hàng"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
