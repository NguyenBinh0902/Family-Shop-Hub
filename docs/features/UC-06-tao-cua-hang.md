# UC-06 Tạo cửa hàng theo loại hình

- **Milestone:** M2 · **Context:** setup · **Tác nhân:** Chủ hộ
- **Quy tắc:** BR-03, BR-08, BR-13, BR-27 · **Quyền:** `store.manage`

## Luồng chính
1. Chủ hộ chọn "Thêm cửa hàng".
2. Hệ thống hiển thị các loại hình `available`, mỗi loại có tên, mô tả, danh sách chức năng bằng lời thường.
3. Chủ hộ chọn loại hình; nhập tên, địa chỉ, số điện thoại; chọn màu; chọn chủ thể kinh doanh hoặc tạo mới ngay trong bước này (UC-05 rút gọn).
4. Hệ thống hiển thị lựa chọn không đổi được về sau: phương pháp giá vốn (MVP chỉ "Bình quân", hiển thị để người dùng biết) và loại hình, kèm cảnh báo.
5. Xác nhận. Trong một transaction: tạo Store (ghim version loại hình, `allow_negative_stock` theo mặc định loại hình), nạp dữ liệu mẫu `seed` (nhóm sản phẩm; nhóm tùy chọn nếu có `fnb.options`).
6. Chuyển vào cửa hàng mới.

## Ngoại lệ
- Hộ đã có 5 cửa hàng chưa đóng → lỗi `BR-27`, không tạo.
- Loại hình hoặc version không tồn tại → lỗi 422.

## Tiêu chí chấp nhận
- **Given** hộ có 1 chủ thể **When** tạo cửa hàng loại `fnb.milk_tea` **Then** Store có `store_type_code=fnb.milk_tea`, `store_type_version=1`, `allow_negative_stock=true`, và có các nhóm sản phẩm mẫu "Trà sữa", "Cà phê", "Topping", "Nguyên liệu".
- **Given** tạo cửa hàng `serial.phone` **Then** `allow_negative_stock=false`.
- **Given** hộ có 4 cửa hàng đang hoạt động và 1 cửa hàng đã đóng **When** tạo cửa hàng mới **Then** thành công; tạo thêm cái thứ 6 chưa đóng **Then** lỗi `BR-27`.
- **Given** cửa hàng vừa tạo **When** gọi `GET /api/v1/stores/{id}` **Then** trả kèm danh sách capability của loại hình.
- Seed lỗi giữa chừng thì không có cửa hàng nào được tạo (rollback).

## API gợi ý
`GET /api/v1/store-types` · `POST /api/v1/households/{householdId}/stores` · `GET /api/v1/stores/{storeId}`

## Ghi chú UI
Chọn loại hình bằng thẻ lớn có biểu tượng, không dùng dropdown. Tối đa 3 bước.
