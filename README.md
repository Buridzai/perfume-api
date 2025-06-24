# Perfume-API

## Giới thiệu dự án

Dự án **Perfume-API** là hệ thống Backend cho một ứng dụng thương mại điện tử chuyên về nước hoa. Được xây dựng bằng ngôn ngữ lập trình Go, API này cung cấp các chức năng cốt lõi cho phép người dùng tương tác với hệ thống: từ việc quản lý tài khoản cá nhân, duyệt và tìm kiếm các sản phẩm, thêm sản phẩm vào giỏ hàng, cho đến việc hoàn tất quá trình đặt hàng.

Ban đầu, dự án được phát triển theo kiến trúc monolith để dễ dàng triển khai và quản lý. Tuy nhiên, với tầm nhìn mở rộng và duy trì tính linh hoạt, dự án đang trong quá trình chuyển đổi dần sang kiến trúc microservice. Điều này giúp tách biệt các trách nhiệm, tăng khả năng mở rộng cho từng thành phần và cho phép triển khai độc lập. Để quản lý môi trường phát triển và vận hành các dịch vụ, dự án sử dụng Docker Compose, và RabbitMQ được tích hợp làm Message Broker cho giao tiếp bất đồng bộ giữa các service.

## Kiến trúc

Hiện tại, dự án bao gồm các service chính sau:

  * **Perfume-API (Order Management Service)**: Đây là phần cốt lõi ban đầu của API, hiện chịu trách nhiệm chính về quản lý đơn hàng, giỏ hàng, và điều phối tương tác với các service khác trong hệ thống. Nó tiếp nhận các yêu cầu từ frontend và điều phối luồng nghiệp vụ liên quan đến đặt hàng và giỏ hàng.
  * **Delivery Service**: Đây là một microservice độc lập. Nhiệm vụ chính của nó là lắng nghe các thông điệp về đơn hàng mới được gửi qua RabbitMQ từ Perfume-API, sau đó xử lý các nghiệp vụ liên quan đến việc giao hàng như cập nhật trạng thái vận chuyển, gửi thông báo giao hàng, v.v.

**Các Microservice được đề xuất để tách trong tương lai:**

Để tiếp tục phát triển theo kiến trúc microservice, các module sau sẽ được tách thành các service độc lập:

  * **User Management Service**: Sẽ quản lý toàn bộ các nghiệp vụ liên quan đến người dùng, bao gồm đăng ký tài khoản, đăng nhập, quản lý profile cá nhân, cũng như xử lý logic liên quan đến JWT (JSON Web Tokens) để xác thực và ủy quyền người dùng.
  * **Product Catalog Service**: Sẽ chuyên biệt hóa việc quản lý thông tin sản phẩm. Các chức năng như tạo mới (Create), đọc (Read), cập nhật (Update) và xóa (Delete) sản phẩm, cùng với các tính năng tìm kiếm và xem chi tiết sản phẩm sẽ nằm tại đây.
  * **Shopping Cart Service**: Sẽ chịu trách nhiệm quản lý tất cả các thao tác liên quan đến giỏ hàng của người dùng, bao gồm thêm sản phẩm vào giỏ, cập nhật số lượng, và xóa sản phẩm khỏi giỏ hàng.
  * **Payment Service**: Sẽ là một service chuyên biệt để xử lý các giao dịch thanh toán. Service này sẽ tích hợp với các cổng thanh toán bên thứ ba (như VNPay, Momo, Stripe, PayPal) và quản lý các quy trình thanh toán, hoàn tiền, và kiểm tra trạng thái giao dịch.

## Công nghệ sử dụng

Dự án này sử dụng một bộ các công nghệ hiện đại để xây dựng cả Backend và quản lý môi trường phát triển/triển khai:

  * **Backend**:
      * **Ngôn ngữ lập trình**: Go (Golang)
      * **Web Framework**: [Gin Gonic](https://gin-gonic.com/en/docs/) - Một framework HTTP nhanh và hiệu quả cho Go, được sử dụng để xây dựng các API RESTful.
      * **ORM (Object-Relational Mapping)**: [GORM](https://gorm.io/docs/index.html) - Thư viện ORM thân thiện với nhà phát triển cho Go, giúp tương tác với cơ sở dữ liệu PostgreSQL một cách dễ dàng và linh hoạt.
      * **Mã hóa mật khẩu (Hashing)**: Thư viện `golang.org/x/crypto/bcrypt` được sử dụng để băm mật khẩu người dùng, đảm bảo an toàn thông tin.
      * **JSON Web Tokens (JWT)**: Được dùng để xác thực người dùng sau khi đăng nhập và quản lý các phiên làm việc, đảm bảo tính bảo mật và trạng thái của người dùng trên hệ thống.
  * **Database**: [PostgreSQL](https://www.postgresql.org/docs/current/index.html) - Hệ quản trị cơ sở dữ liệu quan hệ mã nguồn mở mạnh mẽ và đáng tin cậy.
  * **Cache/Session Store**: [Redis](https://redis.io/docs/latest/) - Một kho dữ liệu trong bộ nhớ (in-memory data store) được sử dụng để lưu trữ các phiên làm việc (session) và dữ liệu cache, cải thiện tốc độ phản hồi của ứng dụng.
  * **Message Broker**: [RabbitMQ](https://www.rabbitmq.com/documentation.html) - Một hệ thống hàng đợi tin nhắn (message queue) mã nguồn mở, được sử dụng để xây dựng các luồng giao tiếp bất đồng bộ giữa các microservice, đặc biệt là giữa Order Management Service và Delivery Service.
  * **Containerization**: [Docker](https://docs.docker.com/) và [Docker Compose](https://docs.docker.com/compose/) - Công nghệ đóng gói ứng dụng và môi trường chạy vào các container độc lập, giúp đảm bảo tính nhất quán và đơn giản hóa quá trình triển khai trên mọi môi trường.
  * **Hot-Reload (trong phát triển)**: [Air](https://github.com/cosmtrek/air) - Công cụ giúp tự động reload ứng dụng Go khi có thay đổi trong mã nguồn trong quá trình phát triển, tăng tốc độ phát triển.

## Cấu trúc thư mục

Dự án được tổ chức thành các thư mục rõ ràng để phân tách các module và service:

```
.
├── Backend/
│   ├── delivery-service/     # Thư mục chứa mã nguồn cho Microservice xử lý giao hàng
│   │   ├── go.mod            # File module Go cho Delivery Service
│   │   ├── go.sum            # Checksum dependencies
│   │   ├── main.go           # Logic chính của Delivery Service, bao gồm kết nối RabbitMQ và xử lý message
│   │   └── Dockerfile        # Dockerfile để build image cho Delivery Service
│   ├── perfume-api/          # Thư mục chứa mã nguồn cho API chính (Order Management Service)
│   │   ├── config/           # Các file cấu hình chung của ứng dụng
│   │   │   ├── db.go         # Cấu hình kết nối cơ sở dữ liệu PostgreSQL
│   │   │   └── redis.go      # Cấu hình kết nối Redis
│   │   ├── controllers/      # Chứa các hàm xử lý logic nghiệp vụ cho từng API endpoint
│   │   ├── docs/             # Chứa các file tài liệu API (Swagger/OpenAPI) được tự động tạo
│   │   ├── middlewares/      # Các middleware xử lý xác thực (JWT), ủy quyền (Admin, Permission)
│   │   ├── models/           # Định nghĩa cấu trúc dữ liệu (structs) cho các đối tượng kinh doanh (User, Product, Order, Cart)
│   │   ├── routes/           # Định nghĩa các tuyến đường (routes) API và liên kết với các controllers
│   │   ├── tmp/              # Thư mục tạm thời được sử dụng bởi công cụ Air
│   │   └── utils/            # Các hàm tiện ích chung của ứng dụng
│   │       └── rabbitmq/     # Chứa logic kết nối và publish message lên RabbitMQ
│   │           └── rabbitmq.go
│   │   ├── .air.toml         # File cấu hình cho Air, cho phép hot-reload trong quá trình phát triển
│   │   └── Dockerfile        # Dockerfile để build image cho Perfume-API
│   └── docker-compose.yml    # File cấu hình chính cho Docker Compose, định nghĩa và liên kết tất cả các service của hệ thống
└── frontend/                 # Thư mục dự kiến chứa mã nguồn cho ứng dụng Frontend (React.js, Vite, CSS)
```

## Yêu cầu cài đặt

Để chạy dự án này, bạn cần cài đặt:

  * [Docker Desktop](https://www.docker.com/products/docker-desktop) (bao gồm Docker Engine và Docker Compose)

## Hướng dẫn cài đặt và chạy dự án

Thực hiện các bước sau để thiết lập và chạy dự án Perfume-API trong môi trường Docker:

1.  **Clone repository**:
    Bắt đầu bằng cách clone mã nguồn dự án từ GitHub về máy tính của bạn và điều hướng vào thư mục `Backend`:

    ```bash
    git clone https://github.com/yourusername/perfume-shop.git
    cd perfume-shop/Backend
    ```

2.  **Chạy các dịch vụ bằng Docker Compose**:
    Đảm bảo bạn đang ở thư mục `Backend` (nơi chứa file `docker-compose.yml`). Chạy lệnh sau để build các Docker image và khởi động tất cả các service đã được định nghĩa:

    ```bash
    docker-compose up --build
    ```

      * Lệnh này sẽ tự động build các Docker image cho `perfume-api` và `delivery-service` dựa trên Dockerfile của mỗi service.
      * Sau đó, nó sẽ khởi động tất cả các service đã khai báo: PostgreSQL (`db`), Redis (`redis`), RabbitMQ (`rabbitmq`), API chính (`perfume-api`), và Delivery Service (`delivery-service`).
      * Quá trình khởi động có thể mất vài phút tùy thuộc vào tốc độ mạng và cấu hình máy tính của bạn. Hãy chờ cho đến khi bạn thấy các log báo hiệu rằng tất cả các service đã "Listening" hoặc "Ready".
      * **Xử lý lỗi khởi động phổ biến**:
          * Nếu gặp lỗi `Bind for 0.0.0.0:6379 failed: port is already allocated`: Điều này có nghĩa là cổng 6379 trên máy host của bạn đang bị một tiến trình khác chiếm dụng. Hãy tìm và dừng tiến trình đó (ví dụ: trên Windows PowerShell sử dụng `netstat -ano | findstr :6379` để tìm PID, sau đó `Stop-Process -Id <PID> -Force`).
          * Nếu gặp lỗi `Cannot connect to RabbitMQ: connect: connection refused` (đặc biệt là ở các lần thử đầu tiên): Service `api` của bạn được cấu hình với cơ chế retry trong `rabbitmq.go`, cho phép nó thử kết nối lại nhiều lần. Hãy kiên nhẫn, nó thường sẽ kết nối thành công sau vài lần thử khi RabbitMQ sẵn sàng hoàn toàn.

3.  **Kiểm tra trạng thái các service**:
    Mở một terminal mới (hoặc sử dụng terminal hiện có sau khi lệnh `up` đã chạy) và chạy lệnh sau để kiểm tra trạng thái của tất cả các container:

    ```bash
    docker-compose ps
    ```

    Đảm bảo rằng tất cả các service (`db`, `redis`, `rabbitmq`, `api`, `delivery-service`) đều ở trạng thái `Up`.

## API Endpoints (Perfume-API - Cổng 8080)

Sau khi các dịch vụ đã chạy, API chính của bạn có thể truy cập được tại `http://localhost:8080`. Bạn có thể sử dụng các công cụ như [Postman](https://www.postman.com/downloads/) hoặc [Insomnia](https://insomnia.rest/) để kiểm thử các API này.

### Authentication

  * **Đăng ký người dùng mới**:
      * **URL**: `http://localhost:8080/api/auth/register`
      * **Method**: `POST`
      * **Body (JSON)**:
        ```json
        {
            "name": "Ten Nguoi Dung",
            "email": "email@example.com",
            "password": "matkhaucuaban"
        }
        ```
  * **Đăng nhập**:
      * **URL**: `http://localhost:8080/api/auth/login`
      * **Method**: `POST`
      * **Body (JSON)**:
        ```json
        {
            "email": "email@example.com",
            "password": "matkhaucuaban"
        }
        ```
      * **Lưu ý**: Phản hồi thành công sẽ bao gồm một JWT `token`. Hãy lưu lại token này để sử dụng trong header `Authorization: Bearer <TOKEN>` cho tất cả các request cần xác thực.

### Product Management

  * **Xem tất cả sản phẩm**:
      * **URL**: `http://localhost:8080/api/products`
      * **Method**: `GET`
      * **Headers**: `Authorization: Bearer <TOKEN>` (token của người dùng hoặc admin)
  * **Xem chi tiết sản phẩm theo ID**:
      * **URL**: `http://localhost:8080/api/products/:id` (thay `:id` bằng ID sản phẩm)
      * **Method**: `GET`
      * **Headers**: `Authorization: Bearer <TOKEN>`
  * **Tạo sản phẩm mới (Admin Only)**:
      * **URL**: `http://localhost:8080/api/products`
      * **Method**: `POST`
      * **Headers**: `Authorization: Bearer <ADMIN_TOKEN>`, `Content-Type: application/json`
      * **Body (JSON)**:
        ```json
        {
            "name": "Ten San Pham",
            "description": "Mo ta chi tiet san pham.",
            "price": 1500000,
            "image": "http://example.com/image.jpg"
        }
        ```
  * **Cập nhật sản phẩm (Admin Only)**:
      * **URL**: `http://localhost:8080/api/products/:id` (thay `:id` bằng ID sản phẩm)
      * **Method**: `PUT`
      * **Headers**: `Authorization: Bearer <ADMIN_TOKEN>`, `Content-Type: application/json`
      * **Body (JSON)**:
        ```json
        {
            "name": "Ten San Pham Cap Nhat",
            "price": 1600000
        }
        ```
  * **Xóa sản phẩm (Admin Only)**:
      * **URL**: `http://localhost:8080/api/products/:id` (thay `:id` bằng ID sản phẩm)
      * **Method**: `DELETE`
      * **Headers**: `Authorization: Bearer <ADMIN_TOKEN>`

### Cart Management

  * **Xem giỏ hàng**:
      * **URL**: `http://localhost:8080/api/cart`
      * **Method**: `GET`
      * **Headers**: `Authorization: Bearer <TOKEN>`
  * **Thêm sản phẩm vào giỏ**:
      * **URL**: `http://localhost:8080/api/cart`
      * **Method**: `POST`
      * **Headers**: `Authorization: Bearer <TOKEN>`, `Content-Type: application/json`
      * **Body (JSON)**:
        ```json
        {
            "product_id": 1,
            "quantity": 2
        }
        ```
  * **Cập nhật số lượng mặt hàng trong giỏ**:
      * **URL**: `http://localhost:8080/api/cart/:id` (thay `:id` bằng ID của mặt hàng trong giỏ)
      * **Method**: `PUT`
      * **Headers**: `Authorization: Bearer <TOKEN>`, `Content-Type: application/json`
      * **Body (JSON)**:
        ```json
        {
            "quantity": 3
        }
        ```
  * **Xóa mặt hàng khỏi giỏ**:
      * **URL**: `http://localhost:8080/api/cart/:id` (thay `:id` bằng ID của mặt hàng trong giỏ)
      * **Method**: `DELETE`
      * **Headers**: `Authorization: Bearer <TOKEN>`

### Order Management

  * **Tạo đơn hàng**:
      * **URL**: `http://localhost:8080/api/orders`
      * **Method**: `POST`
      * **Headers**: `Authorization: Bearer <TOKEN>`, `Content-Type: application/json`
      * **Body (JSON)**:
        ```json
        {
            "items": [
                {"product_id": 1, "quantity": 1},
                {"product_id": 2, "quantity": 2}
            ],
            "address": "123 Đường ABC, Phường XYZ, Quận 1, TP. Hồ Chí Minh"
        }
        ```
      * **Kiểm thử Delivery Service**: Sau khi gọi API này, hãy quan sát log của container `delivery-service` trong terminal mà bạn đã chạy `docker-compose up`. Bạn sẽ thấy log báo hiệu `delivery-service` đã nhận và xử lý thông tin đơn hàng.
  * **Xem lịch sử đơn hàng**:
      * **URL**: `http://localhost:8080/api/orders`
      * **Method**: `GET`
      * **Headers**: `Authorization: Bearer <TOKEN>`

### Admin Specific (Qua API chính)

  * **Xem tất cả người dùng (Admin Only)**:
      * **URL**: `http://localhost:8080/admin/users`
      * **Method**: `GET`
      * **Headers**: `Authorization: Bearer <ADMIN_TOKEN>`
