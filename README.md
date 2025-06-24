Dự án Perfume-API là hệ thống Backend cho một ứng dụng thương mại điện tử chuyên về nước hoa, được xây dựng bằng ngôn ngữ lập trình Go. API này cung cấp các chức năng cốt lõi cho phép người dùng quản lý tài khoản, duyệt sản phẩm, thêm vào giỏ hàng và đặt hàng. Dự án được phát triển với kiến trúc microservice, sử dụng Docker Compose để quản lý các thành phần và RabbitMQ cho giao tiếp bất đồng bộ giữa các service.

Kiến trúc
Ban đầu được phát triển theo kiến trúc monolith, dự án đang trong quá trình chuyển đổi dần sang kiến trúc microservice để tăng tính module hóa, khả năng mở rộng và dễ dàng triển khai.

Các Microservice chính (hoạt động hoặc đang được tách):

Perfume-API (Order Management Service): Core API ban đầu, hiện chịu trách nhiệm chính về quản lý đơn hàng, giỏ hàng, và tương tác với các service khác.
Delivery Service: Microservice độc lập chịu trách nhiệm lắng nghe và xử lý các thông điệp đơn hàng để thực hiện giao hàng.
(Đề xuất) User Management Service: Sẽ quản lý toàn bộ các nghiệp vụ liên quan đến người dùng (đăng ký, đăng nhập, quản lý profile, JWT).
(Đề xuất) Product Catalog Service: Sẽ quản lý các chức năng liên quan đến sản phẩm (CRUD sản phẩm, tìm kiếm).
(Đề xuất) Shopping Cart Service: Sẽ quản lý các thao tác giỏ hàng.
(Đề xuất) Payment Service: Sẽ chuyên biệt hóa việc xử lý các giao dịch thanh toán.
Công nghệ sử dụng
Backend: Go (Golang)
Web Framework: Gin Gonic
ORM: GORM (cho PostgreSQL)
Hashing: golang.org/x/crypto/bcrypt
JWT: Xác thực và quản lý session
Database: PostgreSQL
Cache/Session: Redis
Message Broker: RabbitMQ
Containerization: Docker & Docker Compose
Hot-Reload (Dev): Air
