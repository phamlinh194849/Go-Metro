# Go Metro API Documentation

## Tổng quan
Hệ thống quản lý thẻ metro với các API đầy đủ cho việc quản lý người dùng, thẻ, lịch sử giao dịch, và các hoạt động liên quan đến metro.

## Các API đã được bổ sung

### 1. Sell History APIs (Read-only)
Quản lý lịch sử bán thẻ - **Tự động tạo khi tạo card mới**

#### Endpoints:
- `GET /sell-history` - Lấy danh sách lịch sử bán thẻ (có phân trang và filter)
- `GET /sell-history/:id` - Lấy lịch sử bán thẻ theo ID
- `GET /sell-history/card/:card_id` - Lấy lịch sử bán thẻ theo card ID
- `GET /sell-history/seller/:seller_id` - Lấy lịch sử bán thẻ theo seller ID

#### Query Parameters cho GET /sell-history:
- `page` (int): Số trang (mặc định: 1)
- `limit` (int): Số bản ghi mỗi trang (mặc định: 10)
- `card_id` (string): Lọc theo card ID
- `seller_id` (string): Lọc theo seller ID

### 2. Station History APIs (Read-only)
Quản lý lịch sử check-in/check-out tại các trạm - **Tự động tạo khi check-in/check-out**

#### Endpoints:
- `GET /station-history` - Lấy danh sách lịch sử trạm (có phân trang và filter)
- `GET /station-history/:id` - Lấy lịch sử trạm theo ID
- `GET /station-history/card/:card_id` - Lấy lịch sử trạm theo card ID
- `GET /station-history/station/:station_id` - Lấy lịch sử trạm theo station ID
- `GET /station-history/action/:action` - Lấy lịch sử trạm theo action (checkin/checkout)

#### Query Parameters cho GET /station-history:
- `page` (int): Số trang (mặc định: 1)
- `limit` (int): Số bản ghi mỗi trang (mặc định: 10)
- `card_id` (string): Lọc theo card ID
- `station_id` (int): Lọc theo station ID
- `action` (string): Lọc theo action (checkin/checkout)

### 3. History APIs (Read-only)
Quản lý lịch sử giao dịch chung - **Tự động tạo khi có giao dịch**

#### Endpoints:
- `GET /history` - Lấy danh sách lịch sử giao dịch
- `GET /history/:id` - Lấy lịch sử giao dịch theo ID

### 4. Station Check-in/Check-out APIs
Quản lý check-in/check-out tại các trạm

#### Endpoints:
- `POST /station/:id/checkin` - Check-in tại trạm
- `POST /station/:id/checkout` - Check-out tại trạm

#### Check-in Request:
```json
{
  "card_id": "GM1234567890"
}
```

#### Check-out Request:
```json
{
  "card_id": "GM1234567890"
}
```

### 5. Trip APIs
Quản lý các chuyến tàu

#### Endpoints:
- `POST /trip` - Tạo chuyến tàu mới
- `GET /trip` - Lấy danh sách chuyến tàu (có phân trang và filter)
- `GET /trip/:id` - Lấy chuyến tàu theo ID
- `PUT /trip/:id` - Cập nhật chuyến tàu
- `DELETE /trip/:id` - Xóa chuyến tàu
- `GET /trip/train/:train_id` - Lấy chuyến tàu theo train ID
- `GET /trip/direction/:direction` - Lấy chuyến tàu theo hướng
- `GET /trip/active` - Lấy các chuyến tàu đang hoạt động

#### Query Parameters cho GET /trip:
- `page` (int): Số trang (mặc định: 1)
- `limit` (int): Số bản ghi mỗi trang (mặc định: 10)
- `train_id` (int): Lọc theo train ID
- `direction` (string): Lọc theo hướng

### 6. Train APIs
Quản lý thông tin tàu

#### Endpoints:
- `POST /train` - Tạo tàu mới
- `GET /train` - Lấy danh sách tàu (có phân trang và filter)
- `GET /train/:id` - Lấy tàu theo ID
- `PUT /train/:id` - Cập nhật tàu
- `DELETE /train/:id` - Xóa tàu
- `GET /train/type/:type` - Lấy tàu theo loại
- `GET /train/company/:company` - Lấy tàu theo công ty

#### Query Parameters cho GET /train:
- `page` (int): Số trang (mặc định: 1)
- `limit` (int): Số bản ghi mỗi trang (mặc định: 10)
- `type` (string): Lọc theo loại tàu
- `company` (string): Lọc theo công ty

## Tích hợp History tự động

### 1. SellHistory tự động tạo khi:
- **Tạo card mới** (`POST /card`) - Tự động tạo SellHistory với seller là owner của card

### 2. StationHistory tự động tạo khi:
- **Check-in tại trạm** (`POST /station/:id/checkin`) - Tạo StationHistory với action="checkin"
- **Check-out tại trạm** (`POST /station/:id/checkout`) - Tạo StationHistory với action="checkout" và trừ tiền

### 3. History tự động tạo khi:
- **Nạp tiền thẻ** (`POST /card/:rf_id/topup`) - Tạo History với CardAction=topup
- **Check-out tại trạm** (`POST /station/:id/checkout`) - Tạo History với CardAction=pay

## Cấu trúc Database

### Các bảng chính:
1. **users** - Thông tin người dùng
2. **cards** - Thông tin thẻ metro
3. **stations** - Thông tin các trạm
4. **trains** - Thông tin tàu
5. **trips** - Thông tin chuyến tàu
6. **histories** - Lịch sử giao dịch chung
7. **sell_histories** - Lịch sử bán thẻ
8. **station_histories** - Lịch sử check-in/check-out tại trạm

### Quan hệ khóa ngoại:
- `cards` → `users` (Username)
- `sell_histories` → `cards` (CardID), `users` (SellerID)
- `station_histories` → `cards` (CardID), `stations` (StationID)
- `trips` → `trains` (TrainID)

## Cách sử dụng

### 1. Khởi động ứng dụng:
```bash
# Set biến môi trường để chạy migration
export MIGRATE=true

# Chạy ứng dụng
go run main.go
```

### 2. Truy cập Swagger UI:
```
http://localhost:8080/swagger/index.html
```

### 3. Test API:
Tất cả các API đều có Swagger documentation đầy đủ với các ví dụ request/response.

## Tính năng đặc biệt

### 1. Phân trang:
Hầu hết các API GET đều hỗ trợ phân trang với parameters `page` và `limit`.

### 2. Filtering:
Các API GET hỗ trợ nhiều loại filter khác nhau tùy theo context.

### 3. Preloading:
Các API trả về dữ liệu liên quan (foreign key) đều được preload để tối ưu performance.

### 4. Error Handling:
Tất cả API đều có error handling chuẩn với HTTP status codes phù hợp.

### 5. Authentication:
Một số API yêu cầu authentication và authorization (JWT token).

### 6. Transaction Safety:
Các API tạo history đều sử dụng database transaction để đảm bảo tính nhất quán dữ liệu.

## Lưu ý quan trọng

1. **Migration**: Đảm bảo set `MIGRATE=true` khi chạy lần đầu để tạo các bảng database.
2. **Foreign Keys**: Tất cả các quan hệ đều có constraint `ON UPDATE CASCADE, ON DELETE SET NULL`.
3. **Timestamps**: Tất cả các bảng đều có `created_at` và `updated_at` tự động.
4. **Validation**: Các API đều có validation cơ bản cho input data.
5. **History Logging**: Các history được tạo tự động, không cần gọi API riêng.
6. **Transaction**: Tất cả các thao tác tạo history đều được wrap trong transaction. 