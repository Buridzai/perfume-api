// utils/rabbitmq/rabbitmq.go
package rabbitmq

import (
	"encoding/json"
	"log"
	"os"   // Thêm import os để dùng biến môi trường
	"time" // Thêm import time để dùng time.Sleep

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var ch *amqp.Channel

func InitRabbitMQ() {
	var err error
	amqpURL := os.Getenv("RABBITMQ_URL") // Lấy từ biến môi trường
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/" // Fallback
	}

	// Thêm logic retry cho kết nối RabbitMQ
	for i := 0; i < 10; i++ { // Thử kết nối 10 lần
		log.Printf("ℹ️ Đang cố gắng kết nối RabbitMQ (lần %d/%d)...", i+1, 10)
		conn, err = amqp.Dial(amqpURL)
		if err == nil {
			log.Println("✅ Đã kết nối RabbitMQ thành công!")
			break // Kết nối thành công, thoát vòng lặp
		}
		log.Printf("❌ Không thể kết nối RabbitMQ: %v. Đang chờ 5 giây để thử lại...", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		// Nếu sau nhiều lần thử vẫn không kết nối được, in ra lỗi và thoát
		log.Fatalf("❌ FATAL: Không thể kết nối RabbitMQ sau nhiều lần thử: %v", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("❌ FATAL: Không thể tạo channel RabbitMQ: %v", err)
	}

	err = ch.ExchangeDeclare(
		"delivery-ex", // Tên Exchange
		"fanout",      // Loại Exchange
		true,          // Durable
		false,         // Auto-deleted
		false,         // Internal
		false,         // No-wait
		nil,           // Arguments
	)
	if err != nil {
		log.Fatalf("❌ FATAL: Không thể khai báo Exchange: %v", err)
	}
}

func Publish(queue string, payload interface{}) {
	// Kiểm tra nếu channel bị nil (do lỗi khởi tạo)
	if ch == nil {
		log.Println("❌ Không thể gửi message: RabbitMQ channel chưa được khởi tạo.")
		return
	}
	body, _ := json.Marshal(payload)
	err := ch.Publish("delivery-ex", "", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		log.Println("❌ Gửi message lỗi:", err)
		// Không panic ở đây để ứng dụng không sập nếu RabbitMQ bị lỗi tạm thời
	}
}
