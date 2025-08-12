# Chức năng
# Auth: (Có phân quyền (USER/ADMIN))
- Đăng ký ok 
- Đăng nhập ok
- Xem thông tin user (auth user) ok
- Đổi passwword (auth user) ok
- Sửa thông tin user (auth user) ok
# Quản lý người dùng
- Tạo người dùng (admin) ok
- Xem danh sách user (admin) ok
- Sửa thông tin user (admin) ok
- Xóa người dùng (admin) ok
# Quản lý thẻ RFID. ok
- Thêm ok
- Xóa ok
- Sửa ok
# Danh sách trạm ok
- Xem ok
- Xóa ok
- Sửa ok
# Cấp phát vé (Gán dữ liệu người dùng + Thẻ)
# Báo cáo thống kê
# Lịch trình chuyến


# Model
# User:
- id
- username
- password
- email
- fullname
- role (ADMIN/USER/STAFF)
- status
- createdAt
- updatedAt
# Card:
- id
- ownerid
- blance
- status (thẻ active là thẻ có owner)
- price (giá bán thẻ)
- createdAt
- updatedAt
# Station:
- id
- name
- IPAdress
- createdAt
- updatedAt
# Train:
- id
- name
- type
- company
- createdAt
- updatedAt.
# Trip
- starttime
- endtime
- direction

# SellHistory (mục tiêu xem người bán và tổng kết doanh thu)
- cardid
- sellerid
- cardPriceSold
- time
# StationHistory (theo thẻ, theo trạm, theo ngày)
- id
- action (checkin/checkout)
- time
- cardid
- stationid
- usedBalance

Tuyến (route) Chuyến (trip) Biểu đồ giờ ()