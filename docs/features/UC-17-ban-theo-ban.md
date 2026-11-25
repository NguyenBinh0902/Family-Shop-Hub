# UC-17 Bán theo bàn

- **Milestone:** M7 · **Context:** sales · **Capability:** `fnb.table` · **Tác nhân:** Nhân viên, Quản lý, Chủ hộ
- **Quy tắc:** BR-35, cùng các quy tắc của UC-16 khi hoàn tất · **Quyền:** `sale.create`

## Luồng chính
1. Quản lý khai báo danh sách bàn (tên, thứ tự).
2. Màn hình bàn: mỗi bàn hiện trống hoặc đang có đơn (tổng tiền, thời gian mở).
3. Chọn bàn trống → tạo đơn `OPEN`. Thêm, sửa, bớt món tự do (BR-35); chưa trừ kho.
4. **Chuyển bàn:** đổi `table_id` sang bàn trống.
5. **Gộp bàn:** dồn dòng của đơn mở bàn B vào đơn mở bàn A, xóa đơn B.
6. **Tách đơn:** chọn các dòng (hoặc một phần số lượng) chuyển sang đơn mở mới, để thanh toán riêng.
7. **Thanh toán:** như bước 4–8 của UC-16; đơn chuyển `COMPLETED`, bàn trống.
8. **Hủy đơn mở:** xóa đơn `OPEN` (chưa ghi doanh thu, không cần bút toán); ghi nhật ký.

## Ngoại lệ
- Bàn đã có đơn mở khi tạo mới → mở đơn hiện có.
- Hai máy cùng sửa một đơn mở → dùng `version` (optimistic locking); máy đến sau nhận 409 và tải lại.

## Tiêu chí chấp nhận
- **Given** đơn mở bàn 3 có 2 phở **When** thêm 1 trà đá **Then** đơn vẫn `OPEN`, tồn kho không đổi.
- **Given** đơn mở bàn 3 **When** thanh toán **Then** đơn `COMPLETED`, kho bị trừ, bàn 3 trống.
- **Given** bàn 1 và bàn 2 đều có đơn mở **When** gộp 2 vào 1 **Then** chỉ còn đơn bàn 1 chứa đủ dòng, bàn 2 trống.
- **Given** tách 1 trong 3 phở sang đơn mới **Then** đơn gốc còn 2, đơn mới 1, cả hai `OPEN`.
- **Given** hai cập nhật cùng `version` **Then** cái thứ hai trả 409.

## API gợi ý
`GET/POST /api/v1/stores/{storeId}/tables` · `POST .../orders` với `status=OPEN` · `PATCH .../orders/{id}` (có `version`) · `POST .../orders/{id}/complete` · `POST .../orders/{id}/move` · `.../merge` · `.../split` · `DELETE .../orders/{id}` (chỉ `OPEN`)
