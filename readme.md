# Go-Metro

TH1: Thành công
- checkin thành công/ thất bại
+ mã thẻ
+ thời gian
+ trạng thái(thành công/ thất bại)
+ thông tin người dùng(ID người, số tiền được nạp vào thẻ)
* Lưu lại thông tin hành động check in
- Checkout thành công/ Thất bại

 
CardID, Time, Status, UserID, Blance, Action,


@Nguyen231002Hnam

## Tuỳ chọn migrate khi khởi động

Thêm biến sau vào file `.env` trong thư mục Go-Metro:

```
MIGRATE=true
```

- Nếu `MIGRATE=true` thì sẽ tự động migrate database khi khởi động.
- Nếu `MIGRATE=false` thì sẽ bỏ qua bước migrate.


@Nguyen231002Hnam

```bash
go run main.go
```