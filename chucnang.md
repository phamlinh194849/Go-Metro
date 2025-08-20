- Cấp phát thẻ:
+ Tạo thẻ: POST/card
+ Lấy danh sách thẻ: GET/card
+ Xem lịch sử thẻ: GET/sell-history/card/{card_id}
+ Lấy danh sach người dung :GET/admin/users
+ Xem chi tiết người dung: GET/admin/users/{id}   // Thêm thông tin thẻ
+ Xem lịch sử người cấp phát: GET/sell-history/seller/{seller_id}

- Lịch sử cấp phát thẻ:
+ Lấy danh sách lịch sử cấp phát thẻ: GET/history  // đã có filter theo card_id và seller_id

- Quản lý kho thẻ RFID:
+ Lấy danh sách thẻ: GET/card
+ Xóa thẻ: DELETE/card
+ Cập nhật thẻ: PUT/card/{rf_id}
+ xem thẻ bang RFID: GET/card/cardid/{rf_id}
+ Xem thẻ bang trạng thái: GET/card/status/{status}
+ Xem thẻ bang ID người dung: GET/card/user/{owner_id}
+ Xem thẻ bang ID: GET/card/{id}


- Danh sách trạm: // them hình  minh họa
+ Lấy danh sách trạm:GET/station
+ Tạo trạm: POST/station
+ sửa, xóa
+ POST/station/{id}/checkin       // dùng cho phần cưng
+ POST/station/{id}/checkout	  // dùng cho phần cưng

- Lịch trình chuyến:
+ Lấy danh sách chuyến: GET/trip
+ Tao chuyến: POST/trip
+ Xem danh sách chuyến: GET/trip/active
+ Xem chuyến  bằng ID: GET/trip/{id}
+ Xem danh sách chuyến theo chiều: GET/trip/direction/{directio}
+ Xem danh sách chuyến bằng mã tàu:GET/trip/train/{train_i}
+ Sửa
+ Xóa



- Danh sách thiết bị:		// xem sau
- Báo cáo thống kê: 	//thêm API theo khoảng ngày(đếm tổng ngườidùng, tổng thẻ, doanh thu)


- Lịch sử hoạt đông:
+ Xem danh sách: GET/station-history
+ Xem chi tiêt:GET/station-history/{id}