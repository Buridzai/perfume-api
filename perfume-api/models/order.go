package models

type Order struct {
	ID      uint        `json:"id" gorm:"primaryKey"`
	UserID  uint        `json:"user_id"`
	Total   float64     `json:"total"`
	Address string      `json:"address"` // <-- Đảm bảo dòng này đã được thêm
	Items   []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}
