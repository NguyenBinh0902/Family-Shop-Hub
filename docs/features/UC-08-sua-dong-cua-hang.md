# UC-08 Sửa thông tin, tạm ngưng, đóng cửa hàng

- **Milestone:** M2 · **Context:** setup · **Tác nhân:** Chủ hộ
- **Quy tắc:** BR-04, BR-05, BR-13 · **Quyền:** `store.manage`

## Luồng chính
- **Sửa:** tên, địa chỉ, số điện thoại, màu, cho phép bán âm, mức trần giảm giá của nhân viên (0–100%). Đổi chủ thể chỉ khi cửa hàng chưa có giao dịch.
- **Tạm ngưng / mở lại:** `ACTIVE ⇄ SUSPENDED`.
- **Đóng:** xác nhận bằng cách gõ lại tên cửa hàng → `CLOSED`, không mở lại được.

## Ngoại lệ
- Đổi chủ thể khi đã có giao dịch → lỗi `BR-04`.
- Cửa hàng serial bật bán âm → vẫn lưu được, nhưng sản phẩm serial luôn chặn (BR-13).
- Còn đơn `OPEN` (theo bàn) khi đóng hoặc tạm ngưng → yêu cầu xử lý xong đơn mở trước.

## Tiêu chí chấp nhận
- **Given** cửa hàng `SUSPENDED` **When** tạo đơn, nhập hàng, ghi phiếu chi **Then** đều bị từ chối; xem báo cáo vẫn được.
- **Given** cửa hàng `CLOSED` **When** gọi mở lại **Then** trả lỗi `BR-05`.
- **Given** cửa hàng đã có đơn **When** đổi `business_entity_id` **Then** trả lỗi `BR-04`.
- Đổi mức trần giảm giá có hiệu lực ngay ở đơn tiếp theo.

## API gợi ý
`PATCH /api/v1/stores/{storeId}` · `POST /api/v1/stores/{storeId}/suspend` · `.../resume` · `.../close`
