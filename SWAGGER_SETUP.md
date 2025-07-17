# Swagger Documentation Setup Guide

## Tổng quan
Dự án Go-Metro đã được tích hợp Swagger để tạo documentation API tự động. Tất cả các API endpoints đã được comment với các annotation Swagger.

## Cài đặt

### 1. Cài đặt swag CLI
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. Cài đặt dependencies
```bash
go mod tidy
```

## Sử dụng

### 1. Generate Swagger documentation
```bash
swag init
```
Lệnh này sẽ tạo thư mục `docs/` chứa các file Swagger được generate tự động.

### 2. Chạy server
```bash
go run main.go
```

### 3. Truy cập Swagger UI
Mở trình duyệt và truy cập: `http://localhost:8080/swagger/index.html`

## API Documentation

### Authentication
- **POST** `/auth/register` - Đăng ký tài khoản mới
- **POST** `/auth/login` - Đăng nhập

### User Management
- **GET** `/user/profile` - Xem thông tin profile (cần authentication)
- **PUT** `/user/profile` - Cập nhật profile (cần authentication)
- **PUT** `/user/password` - Đổi mật khẩu (cần authentication)

### Card Management
- **POST** `/card` - Tạo card mới
- **GET** `/card` - Lấy danh sách tất cả cards
- **GET** `/card/{id}` - Lấy card theo ID
- **GET** `/card/cardid/{card_id}` - Lấy card theo Card ID
- **PUT** `/card/{id}` - Cập nhật card
- **DELETE** `/card/{id}` - Xóa card
- **POST** `/card/{id}/topup` - Nạp tiền vào card
- **GET** `/card/user/{user_id}` - Lấy cards theo user ID
- **GET** `/card/status/{status}` - Lấy cards theo status

### History Management
- **POST** `/history` - Tạo lịch sử giao dịch
- **GET** `/history` - Lấy tất cả lịch sử
- **GET** `/history/{id}` - Lấy lịch sử theo ID

### Admin Management (Admin only)
- **GET** `/admin/users` - Lấy tất cả users
- **GET** `/admin/users/{id}` - Lấy user theo ID
- **PUT** `/admin/users/{id}` - Cập nhật user
- **DELETE** `/admin/users/{id}` - Xóa user

### System
- **GET** `/health` - Health check

## Authentication

Hầu hết các API endpoints yêu cầu authentication thông qua JWT token. Token được trả về khi đăng nhập thành công.

### Sử dụng token:
1. Đăng nhập để lấy token
2. Thêm header: `Authorization: Bearer <your_token>`

## Response Format

Tất cả API responses đều theo format:
```json
{
  "success": true,
  "message": "Success message",
  "data": {
    // Response data
  }
}
```

## Error Handling

Các lỗi thường gặp:
- `400` - Bad Request (validation error, invalid data)
- `401` - Unauthorized (missing or invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found (resource not found)
- `500` - Internal Server Error (server error)

## Cập nhật Documentation

Khi thêm API mới hoặc thay đổi API hiện tại:

1. Thêm Swagger annotations vào handler function
2. Chạy `swag init` để regenerate documentation
3. Restart server

## Ví dụ Swagger Annotation

```go
// @Summary Create a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object true "User data"
// @Success 201 {object} utils.Response{data=models.User}
// @Failure 400 {object} utils.Response
// @Router /auth/register [post]
func Register(c *gin.Context) {
    // Handler implementation
}
```

## Troubleshooting

### Lỗi "docs not found"
- Chạy `swag init` để generate docs
- Kiểm tra import `_ "go-metro/docs"` trong main.go

### Lỗi "swag command not found"
- Cài đặt swag: `go install github.com/swaggo/swag/cmd/swag@latest`
- Thêm `$GOPATH/bin` vào PATH

### Lỗi "invalid annotation"
- Kiểm tra syntax của Swagger annotations
- Đảm bảo comment bắt đầu với `// @`
- Kiểm tra tên tag và parameter names 