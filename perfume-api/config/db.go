// config/db.go
package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Lỗi kết nối DB: %v", err) // Thay vì log.Fatal
		// os.Exit(1) // Bạn có thể thêm lại cái này sau khi debug
		// Để debug, tạm thời bỏ os.Exit(1) để hàm có thể hoàn thành, và bạn sẽ thấy lỗi.
		// Sau đó, khi chương trình cố gắng sử dụng DB, nó sẽ panic vì DB là nil.
		// Quan trọng: Thêm một return ở đây để hàm không tiếp tục với DB là nil.
		return
	}
	DB = db
	log.Println("✅ Đã kết nối PostgreSQL thành công!")
}
