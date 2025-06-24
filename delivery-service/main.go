// delivery-service/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

// OrderMessage định nghĩa cấu trúc của message đơn hàng mà Delivery Service sẽ nhận
type OrderMessage struct {
	OrderID uint    `json:"order_id"`
	UserID  uint    `json:"user_id"`
	Total   float64 `json:"total"`
	Address string  `json:"address"`
}

func main() {
	// Lấy RabbitMQ URL từ biến môi trường
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/" // Fallback nếu không có biến môi trường
		log.Printf("⚠️ RABBITMQ_URL không được đặt, sử dụng mặc định: %s", amqpURL)
	}

	conn, err := connectRabbitMQ(amqpURL)
	if err != nil {
		log.Fatalf("❌ Không thể kết nối RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Không thể mở kênh RabbitMQ: %v", err)
	}
	defer ch.Close()

	// Khai báo Exchange đã được sử dụng trong perfume-api/utils/rabbitmq/rabbitmq.go
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
		log.Fatalf("❌ Không thể khai báo Exchange: %v", err)
	}

	// Khai báo một queue tạm thời, tự động xóa khi kết nối đóng
	q, err := ch.QueueDeclare(
		"",    // Tên ngẫu nhiên cho queue
		false, // Durable
		false, // Delete when unused
		true,  // Exclusive (chỉ client này có thể truy cập)
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("❌ Không thể khai báo Queue: %v", err)
	}

	// Bind queue vào exchange
	err = ch.QueueBind(
		q.Name,        // Tên queue ngẫu nhiên
		"",            // Routing key (empty cho fanout exchange)
		"delivery-ex", // Tên exchange
		false,         // No-wait
		nil,           // Arguments
	)
	if err != nil {
		log.Fatalf("❌ Không thể bind Queue: %v", err)
	}

	log.Printf("✅ Delivery Service đang chờ thông điệp đơn hàng trên exchange 'delivery-ex'...")

	msgs, err := ch.Consume(
		q.Name, // Tên queue
		"",     // Consumer
		true,   // Auto-ack (tự động xác nhận sau khi nhận)
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Args
	)
	if err != nil {
		log.Fatalf("❌ Không thể đăng ký Consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var orderMsg OrderMessage
			err := json.Unmarshal(d.Body, &orderMsg)
			if err != nil {
				log.Printf("❗ Lỗi giải mã message: %s, Body: %s", err, d.Body)
				continue
			}
			log.Printf("🚚 Đã nhận đơn hàng #%d (Người dùng: %d, Tổng: %.2f, Địa chỉ: %s)",
				orderMsg.OrderID, orderMsg.UserID, orderMsg.Total, orderMsg.Address)

			// Mô phỏng quá trình xử lý giao hàng
			time.Sleep(2 * time.Second) // Giả lập thời gian xử lý
			log.Printf("✅ Đơn hàng #%d đã được xử lý giao hàng.", orderMsg.OrderID)
		}
	}()

	<-forever
}

// connectRabbitMQ cố gắng kết nối lại nhiều lần
func connectRabbitMQ(amqpURL string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	for i := 0; i < 5; i++ { // Thử kết nối 5 lần
		conn, err = amqp.Dial(amqpURL)
		if err == nil {
			return conn, nil
		}
		log.Printf("❗ Thử kết nối RabbitMQ lần %d thất bại: %v. Đang thử lại...", i+1, err)
		time.Sleep(5 * time.Second) // Đợi 5 giây trước khi thử lại
	}
	return nil, fmt.Errorf("không thể kết nối RabbitMQ sau nhiều lần thử: %w", err)
}
