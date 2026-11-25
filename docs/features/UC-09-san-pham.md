# UC-09 Tạo, sửa, ngừng bán sản phẩm

- **Milestone:** M3 (thuộc tính serial, bảo hành bổ sung ở M6) · **Context:** catalog · **Tác nhân:** Quản lý, Chủ hộ
- **Quy tắc:** BR-06, BR-07, BR-26 · **Quyền:** `catalog.manage`

## Luồng chính
1. Nhập: tên, nhóm, đơn vị tính, giá bán, bán trực tiếp hay là nguyên liệu, kiểu theo dõi tồn, mã SKU (tùy chọn), ngưỡng tồn thấp (nếu theo số lượng), số tháng bảo hành (nếu serial và có `serial.warranty`).
2. Form hiển thị thêm các trường động của loại hình (`GET .../product-attributes`).
3. Kiểu theo dõi tồn chỉ gồm lựa chọn mà cửa hàng có capability: `NONE` luôn có; `QUANTITY` cần `inv.quantity`; `SERIAL` cần `inv.serial`.
4. Lưu. Sửa sau được, trừ kiểu theo dõi tồn.
5. **Ngừng bán:** sản phẩm chuyển `DISCONTINUED`, ẩn khỏi POS, vẫn hiện trong lịch sử và báo cáo. Bán lại được bằng thao tác "Bán lại".
6. **Xóa:** chỉ khi sản phẩm chưa có giao dịch nào.
7. Tạo nhanh nhóm sản phẩm ngay trong form.

## Ngoại lệ
- Tên trùng với sản phẩm đang bán trong cửa hàng → lỗi trường `name`.
- Trường động bắt buộc thiếu hoặc sai kiểu → lỗi theo trường.
- Đổi `stock_tracking` → lỗi `BR-07`.
- Xóa sản phẩm đã có giao dịch → lỗi `BR-06`, gợi ý ngừng bán.

## Tiêu chí chấp nhận
- **Given** cửa hàng `fnb.bakery` (không có `inv.serial`) **When** tạo sản phẩm `SERIAL` **Then** lỗi `BR-07`.
- **Given** sản phẩm đã có trong một đơn **When** xóa **Then** lỗi `BR-06`; **When** ngừng bán **Then** không còn ở danh sách POS nhưng đơn cũ vẫn hiển thị tên.
- **Given** đổi giá bán từ 25.000 lên 28.000 **Then** đơn tạo sau dùng 28.000, đơn cũ giữ 25.000, nhật ký có giá trị trước và sau.
- **Given** Nhân viên **When** gọi API tạo sản phẩm **Then** 403; **When** xem danh sách sản phẩm **Then** được, nhưng không có giá vốn.
- Tìm sản phẩm theo tên không phân biệt dấu và hoa thường ("tra sua" khớp "Trà sữa").

## API gợi ý
`GET/POST /api/v1/stores/{storeId}/products` · `GET/PATCH/DELETE .../products/{id}` · `POST .../products/{id}/discontinue` · `.../reactivate` · `GET/POST .../categories` · `GET .../product-attributes`
