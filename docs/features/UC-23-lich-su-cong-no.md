# UC-23 Xem lịch sử mua và công nợ của khách

- **Milestone:** M5 · **Context:** customer (+ sales) · **Capability:** `core.customer` · **Tác nhân:** Nhân viên, Quản lý, Chủ hộ
- **Quy tắc:** BR-21 · **Quyền:** `customer.history.store`, `customer.history.household`

## Luồng chính
1. Mở chi tiết khách: tên, số điện thoại, ghi chú.
2. Phần "Tại cửa hàng này": danh sách đơn (ngày, số đơn, tổng, trạng thái), tổng đã mua, nợ còn lại, danh sách khoản nợ đang mở.
3. Chủ hộ thấy thêm phần "Toàn hộ": tổng đã mua và tổng nợ theo từng cửa hàng; bấm vào cửa hàng để xem chi tiết.
4. (M6) Danh sách máy đã mua kèm hạn bảo hành.

## Tiêu chí chấp nhận
- **Given** khách nợ 200.000 ở tiệm điện thoại và 50.000 ở quán trà sữa **When** Chủ hộ xem **Then** thấy tổng nợ 250.000 tách theo hai cửa hàng.
- **Given** cùng khách **When** nhân viên quán trà sữa xem **Then** chỉ thấy nợ 50.000 và đơn của quán.
- Đơn đã hủy hiển thị trạng thái "Đã hủy" và không tính vào tổng đã mua.

## API gợi ý
`GET /api/v1/stores/{storeId}/customers/{id}` · `GET /api/v1/stores/{storeId}/customers/{id}/orders` · `GET /api/v1/households/{householdId}/customers/{id}/summary`
