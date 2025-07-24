# Cấu trúc Project Go-Metro

## Tổng quan
Project được tổ chức theo mô hình MVC đơn giản, dễ bảo trì và mở rộng.

## Cấu trúc thư mục

```
go-metro/
├── main.go              # Entry point của ứng dụng
├── go.mod               # Go modules
├── go.sum               # Dependencies checksum
├── readme.md            # Documentation
├── STRUCTURE.md         # File này
├── config/
│   └── database.go      # Cấu hình database connection
├── models/
│   ├── history.go       # Model History và migration
│   └── card.go          # Model Card và migration
├── handlers/
│   ├── history.go       # Business logic cho History
│   └── card.go          # Business logic cho Card
├── routes/
│   └── routes.go        # Định nghĩa tất cả routes
└── utils/
    └── response.go      # Utility functions cho response
```

## Chi tiết từng thành phần

### 1. config/
- **database.go**: Quản lý kết nối database, đọc environment variables
- Tách biệt logic database khỏi business logic

### 2. models/
- **history.go**: Định nghĩa struct History và migration
- Chứa logic liên quan đến database schema

### 3. handlers/
- **history.go**: Xử lý HTTP requests cho History
- Chứa business logic, validation, và response
- Sử dụng utils để chuẩn hóa response format

### 4. routes/
- **routes.go**: Định nghĩa tất cả API endpoints
- Nhóm các routes theo chức năng
- Dễ dàng thêm/sửa/xóa routes

### 5. utils/
- **response.go**: Chuẩn hóa format response
- Cung cấp các helper functions cho success/error responses
- Đảm bảo consistency trong API responses

## API Endpoints

### History APIs
- `POST /history` - Tạo history mới
- `GET /history` - Lấy danh sách histories
- `GET /history/:id` - Lấy history theo ID

### Card APIs
- `POST /card` - Tạo card mới
- `GET /card` - Lấy danh sách cards
- `GET /card/:id` - Lấy card theo ID
- `GET /card/cardid/:card_id` - Lấy card theo CardID
- `PUT /card/:id` - Cập nhật card
- `DELETE /card/:id` - Xóa card
- `POST /card/:id/topup` - Nạp tiền vào card
- `GET /card/user/:user_id` - Lấy cards theo UserID
- `GET /card/status/:status` - Lấy cards theo status

### System APIs
- `GET /health` - Health check

## Lợi ích của cấu trúc mới

1. **Separation of Concerns**: Mỗi package có trách nhiệm riêng biệt
2. **Maintainability**: Dễ dàng tìm và sửa code
3. **Scalability**: Dễ dàng thêm features mới
4. **Testability**: Có thể test từng component riêng biệt
5. **Consistency**: Response format được chuẩn hóa
6. **Simplicity**: Vẫn giữ được tính đơn giản cho dự án nhỏ

## Cách thêm feature mới

1. Tạo model trong `models/`
2. Tạo handler trong `handlers/`
3. Thêm routes trong `routes/routes.go`
4. Cập nhật migration nếu cần 

## 🎯 **Model Card mới:**

```go
type Card struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    CardID    string    `gorm:"uniqueIndex;not null" json:"card_id"`
    UserID    string    `json:"user_id"`
    Balance   float64   `json:"balance" gorm:"default:0"`
    Status    string    `json:"status" gorm:"default:'active'"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## 🎯 **Các API Card mới:**

### **CRUD Operations:**
- `POST /card` - Tạo card mới
- `GET /card` - Lấy danh sách tất cả cards
- `GET /card/:id` - Lấy card theo ID
- `GET /card/cardid/:card_id` - Lấy card theo CardID
- `PUT /card/:id` - Cập nhật thông tin card
- `DELETE /card/:id` - Xóa card

### **Business Logic:**
- `POST /card/:id/topup` - Nạp tiền vào card
- `GET /card/user/:user_id` - Lấy tất cả cards của user
- `GET /card/status/:status` - Lấy cards theo trạng thái

## 🎯 **Tính năng đặc biệt:**

1. **Validation**: Kiểm tra card đã tồn tại khi tạo mới
2. **Top-up**: Có thể nạp tiền vào card với validation amount > 0
3. **Status Management**: Quản lý trạng thái card (active, inactive, blocked)
4. **User Association**: Liên kết card với user
5. **Flexible Queries**: Tìm kiếm theo nhiều tiêu chí khác nhau

## 📝 **Cách sử dụng:**

### Tạo card mới:
```json
POST /card
{
    "card_id": "CARD001",
    "user_id": "USER001",
    "balance": 0,
    "status": "active"
}
```

### Nạp tiền:
```json
POST /card/1/topup
{
    "amount": 100000
}
```

Tất cả API đều sử dụng response format chuẩn và có error handling đầy đủ! 