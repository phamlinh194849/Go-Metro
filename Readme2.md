# Chức năng
# Auth: (Có phân quyền (USER/ADMIN))
- Đăng ký
- Đăng nhập
- Xem thông tin user (auth user)
- Đổi passwword (auth user)
- Sửa thông tin user (auth user)
# Quản lý người dùng
- Tạo người dùng (admin)
- Xem danh sách user (admin)
- Xóa người dùng
# Quản lý thẻ RFID.
- Thêm 
- Xóa
- Sửa
# Danh sách trạm
- Xem
- Xóa 
- Sửa
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
- status
- price
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