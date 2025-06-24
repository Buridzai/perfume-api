# Perfume-API

Dự án **Perfume-API** là hệ thống Backend cho một ứng dụng thương mại điện tử chuyên về nước hoa, được xây dựng bằng ngôn ngữ lập trình Go. API này cung cấp các chức năng cốt lõi cho phép người dùng quản lý tài khoản, duyệt sản phẩm, thêm vào giỏ hàng và đặt hàng. Dự án được phát triển với kiến trúc microservice, sử dụng Docker Compose để quản lý các thành phần và RabbitMQ cho giao tiếp bất đồng bộ giữa các service.

## Kiến trúc

Ban đầu được phát triển theo kiến trúc monolith, dự án đang trong quá trình chuyển đổi dần sang kiến trúc microservice để tăng tính module hóa, khả năng mở rộng và dễ dàng triển khai.

**Các Microservice chính (hoạt động hoặc đang được đề xuất):**

  * **Perfume-API (Order Management Service)**: Core API ban đầu, hiện chịu trách nhiệm chính về quản lý đơn hàng, giỏ hàng, và tương tác với các service khác.
  * **Delivery Service**: Microservice độc lập chịu trách nhiệm lắng nghe và xử lý các thông điệp đơn hàng để thực hiện giao hàng.
  * **(Đề xuất) User Management Service**: Sẽ quản lý toàn bộ các nghiệp vụ liên quan đến người dùng (đăng ký, đăng nhập, quản lý profile, JWT).
  * **(Đề xuất) Product Catalog Service**: Sẽ quản lý các chức năng liên quan đến sản phẩm (CRUD sản phẩm, tìm kiếm).
  * **(Đề xuất) Shopping Cart Service**: Sẽ quản lý các thao tác giỏ hàng.
  * **(Đề xuất) Payment Service**: Sẽ chuyên biệt hóa việc xử lý các giao dịch thanh toán.

## Công nghệ sử dụng

  * **Backend**: Go (Golang)
      * **Web Framework**: [Gin Gonic](https://gin-gonic.com/en/docs/)
      * **ORM**: [GORM](https://gorm.io/docs/index.html) (cho PostgreSQL)
      * **Hashing**: `golang.org/x/crypto/bcrypt`
      * **JWT**: Xác thực và quản lý session
  * **Database**: [PostgreSQL](https://www.postgresql.org/docs/current/index.html)
  * **Cache/Session**: [Redis](https://redis.io/docs/latest/)
  * **Message Broker**: [RabbitMQ](https://www.rabbitmq.com/documentation.html)
  * **Containerization**: [Docker & Docker Compose](https://docs.docker.com/)
  * **Hot-Reload (Dev)**: [Air](https://github.com/cosmtrek/air)

## Cấu trúc thư mục

```
.
├── Backend/
│   ├── delivery-service/     # Microservice xử lý giao hàng
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── main.go           # Logic của Delivery Service
│   │   └── Dockerfile
│   ├── perfume-api/          # API chính (Order Management Service)
│   │   ├── config/           # Cấu hình DB, Redis, RabbitMQ
│   │   │   ├── db.go
│   │   │   └── redis.go
│   │   ├── controllers/      # Logic xử lý các request HTTP
│   │   ├── docs/             # Tài liệu Swagger/OpenAPI
│   │   ├── middlewares/      # Các middleware (JWT, Auth, Admin)
│   │   ├── models/           # Định nghĩa cấu trúc dữ liệu (structs)
│   │   ├── routes/           # Định nghĩa các tuyến đường API
│   │   ├── tmp/              # Thư mục tạm thời của Air
│   │   └── utils/            # Các tiện ích chung (Hash, JWT)
│   │       └── rabbitmq/     # Tiện ích liên quan đến RabbitMQ
│   │           └── rabbitmq.go
│   │   ├── .air.toml         # Cấu hình Air cho hot-reload
│   │   └── Dockerfile
│   └── docker-compose.yml    # Cấu hình Docker Compose cho toàn bộ hệ thống
└── frontend/                 # (Chưa cung cấp code, nhưng dự kiến sẽ dùng React.js, Vite, CSS)
```

## Yêu cầu cài đặt

  * [Docker Desktop](https://www.docker.com/products/docker-desktop) (bao gồm Docker Engine và Docker Compose)

## Hướng dẫn cài đặt và chạy dự án

1.  **Clone repository**:

    ```bash
    git clone https://github.com/yourusername/perfume-shop.git
    cd perfume-shop/Backend
    ```

2.  **Chạy các dịch vụ bằng Docker Compose**:
    Đảm bảo bạn đang ở thư mục `Backend` (nơi chứa file `docker-compose.yml`).

    ```bash
    docker-compose up --build
    ```

      * Lệnh này sẽ build các Docker image cho `perfume-api` và `delivery-service`, sau đó khởi động tất cả các service (PostgreSQL, Redis, RabbitMQ, perfume-api, delivery-service).
      * Quá trình khởi động có thể mất vài phút. Hãy chờ cho đến khi bạn thấy các log báo hiệu rằng tất cả các service đã "Listening" hoặc "Ready".
      * Nếu gặp lỗi `port is already allocated`, hãy giải phóng cổng 6379 trên hệ thống của bạn (xem mục [Xử lý lỗi](https://www.google.com/search?q=%23x%E1%BB%AD-l%C3%BD-l%E1%BB%97i) bên dưới).
      * Nếu gặp lỗi kết nối RabbitMQ, service `api` sẽ tự động retry nhờ vào logic đã được thêm vào `rabbitmq.go`.

3.  **Kiểm tra trạng thái các service**:
    Mở một terminal mới và chạy:

    ```bash
    docker-compose ps
    ```

    Đảm bảo tất cả các service đều ở trạng thái `Up`.

## API Endpoints (Perfume-API - Cổng 8080)

Bạn có thể sử dụng [Postman](https://www.postman.com/downloads/) hoặc các công cụ HTTP client khác để kiểm thử các API này.

### Authentication

  * **Đăng ký người dùng**:
      * `POST /api/auth/register`
      * **Body (JSON)**: `{"name": "...", "email": "...", "password": "..."}`
  * **Đăng nhập**:
      * `POST /api/auth/login`
      * **Body (JSON)**: `{"email": "...", "password": "..."}`
      * **Lưu ý**: Phản hồi sẽ bao gồm một JWT `token`. Hãy sử dụng token này trong header `Authorization: Bearer <TOKEN>` cho các request cần xác thực.

### Product Management

  * **Xem tất cả sản phẩm**:
      * `GET /api/products` (yêu cầu JWT token)
  * **Xem chi tiết sản phẩm theo ID**:
      * `GET /api/products/:id` (yêu cầu JWT token)
  * **Tạo sản phẩm mới (Admin Only)**:
      * `POST /api/products` (yêu cầu JWT token của Admin)
      * **Body (JSON)**: `{"name": "...", "description": "...", "price": ..., "image": "..."}`
  * **Cập nhật sản phẩm (Admin Only)**:
      * `PUT /api/products/:id` (yêu cầu JWT token của Admin)
      * **Body (JSON)**: `{"name": "...", "price": ...}`
  * **Xóa sản phẩm (Admin Only)**:
      * `DELETE /api/products/:id` (yêu cầu JWT token của Admin)

### Cart Management

  * **Xem giỏ hàng**:
      * `GET /api/cart` (yêu cầu JWT token)
  * **Thêm sản phẩm vào giỏ**:
      * `POST /api/cart` (yêu cầu JWT token)
      * **Body (JSON)**: `{"product_id": ..., "quantity": ...}`
  * **Cập nhật số lượng mặt hàng trong giỏ**:
      * `PUT /api/cart/:id` (yêu cầu JWT token)
      * **Body (JSON)**: `{"quantity": ...}`
  * **Xóa mặt hàng khỏi giỏ**:
      * `DELETE /api/cart/:id` (yêu cầu JWT token)

### Order Management

  * **Tạo đơn hàng**:
      * `POST /api/orders` (yêu cầu JWT token)
      * **Body (JSON)**: `{"items": [{"product_id": ..., "quantity": ...}], "address": "..."}`
      * **Kiểm thử Delivery Service**: Sau khi gọi API này, hãy quan sát log của container `delivery-service` trong terminal mà bạn đã chạy `docker-compose up`. Bạn sẽ thấy log báo hiệu `delivery-service` đã nhận và xử lý thông tin đơn hàng.
  * **Xem lịch sử đơn hàng**:
      * `GET /api/orders` (yêu cầu JWT token)

### Admin Specific (Qua API chính)

  * **Xem tất cả người dùng (Admin Only)**:
      * `GET /admin/users` (yêu cầu JWT token của Admin)

## Xử lý lỗi

  * **`Error: connect ECONNREFUSED 127.0.0.1:8080` (hoặc cổng khác)**: Nghĩa là không có dịch vụ nào đang chạy trên cổng đó hoặc kết nối bị từ chối.
      * **Giải pháp**: Đảm bảo service tương ứng đang chạy và không có lỗi khởi động. Xem logs của service bằng `docker-compose logs <tên_service>` để chẩn đoán nguyên nhân.
  * **`Bind for 0.0.0.0:6379 failed: port is already allocated`**: Cổng `6379` trên máy host của bạn đang bị chiếm dụng.
      * **Giải pháp**: Tìm và dừng tiến trình đang chiếm cổng `6379`.
          * Trên Windows (PowerShell): `netstat -ano | findstr :6379` để tìm PID, sau đó `Stop-Process -Id <PID> -Force`.
          * Trên Linux/macOS: `sudo lsof -i :6379` để tìm PID, sau đó `kill <PID>`.
  * **`<service-name> exited with code 1` mà không có log chi tiết**: Ứng dụng Go bên trong container thoát ngay lập tức.
      * **Giải pháp**: Tạm thời loại bỏ `command: air -c .air.toml` trong `docker-compose.yml` của service đó để cho phép `CMD ["./main"]` trong Dockerfile chạy. Nếu vẫn không có log, hãy sửa tạm thời `log.Fatal` hoặc `panic` trong code Go thành `log.Printf` để thấy lỗi cụ thể hơn.
  * **`Cannot connect to RabbitMQ: connect: connection refused`**: RabbitMQ chưa sẵn sàng chấp nhận kết nối.
      * **Giải pháp**: Đảm bảo logic retry (vòng lặp thử lại kết nối) đã được triển khai trong hàm `InitRabbitMQ` của ứng dụng Go.

