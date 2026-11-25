# UC-21 In hoặc chia sẻ hóa đơn bán

- **Milestone:** M4 · **Context:** frontend (dữ liệu từ sales) · **Tác nhân:** Nhân viên, Quản lý, Chủ hộ
- **Quyền:** `sale.create`

## Luồng chính
1. Sau khi bán, hoặc từ danh sách đơn, bấm "In".
2. Trình duyệt mở bản in dùng CSS `@media print`, khổ 58 mm hoặc 80 mm (cài đặt theo thiết bị, lưu ở trình duyệt).
3. Nội dung: tên, địa chỉ, SĐT cửa hàng; số đơn; ngày giờ (giờ Việt Nam); người bán; từng dòng (tên, tùy chọn, số lượng, đơn giá, thành tiền; serial và hạn bảo hành nếu có); tạm tính; giảm giá; tổng; các khoản thanh toán; tiền khách đưa và tiền thối (nếu tiền mặt); dòng cảm ơn.
4. "Chia sẻ": tạo ảnh hóa đơn ở client và mở Web Share API (gửi qua Zalo...). Trình duyệt không hỗ trợ thì cho tải ảnh.

## Tiêu chí chấp nhận
- **Given** đơn đã hủy **When** in lại **Then** bản in có dòng chữ "ĐÃ HỦY" nổi bật.
- **Given** khổ 58 mm **Then** không có chữ bị cắt ngang; tên sản phẩm dài tự xuống dòng.
- Bản in không chứa giá vốn.
- In lại một đơn cũ cho ra cùng nội dung (lấy từ snapshot, không lấy giá hiện tại).

## Ngoài phạm vi
Máy in nhiệt kết nối trực tiếp (Bluetooth, LAN), hóa đơn điện tử.

## API gợi ý
Dùng `GET /api/v1/stores/{storeId}/orders/{id}` (đủ dữ liệu in, kèm thông tin cửa hàng).
