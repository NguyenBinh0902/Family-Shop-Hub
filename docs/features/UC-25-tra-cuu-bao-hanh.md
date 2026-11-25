# UC-25 Tra cứu bảo hành theo IMEI

- **Milestone:** M6 · **Context:** inventory (+ sales, customer) · **Capability:** `serial.warranty` · **Tác nhân:** Nhân viên, Quản lý, Chủ hộ
- **Quy tắc:** BR-12 · **Quyền:** `stock.view`

## Luồng chính
1. Gõ hoặc quét IMEI/số khung (chuẩn hóa như khi nhập).
2. Hệ thống tìm trong cửa hàng hiện tại (Chủ hộ: trong cả hộ).
3. Hiển thị: sản phẩm, trạng thái, ngày bán, số đơn, khách mua, hạn bảo hành, nhãn "Còn bảo hành" / "Hết bảo hành" (so với ngày nghiệp vụ hôm nay).

## Tiêu chí chấp nhận
- **Given** IMEI bán ngày 15/01/2027 bảo hành 12 tháng **When** tra ngày 10/01/2028 **Then** "Còn bảo hành đến 15/01/2028"; tra ngày 16/01/2028 **Then** "Hết bảo hành".
- **Given** IMEI chưa bán **Then** hiển thị "Chưa bán", trạng thái kho.
- **Given** IMEI bán ở tiệm khác cùng hộ **When** nhân viên tra **Then** "Không tìm thấy ở cửa hàng này"; Chủ hộ tra **Then** thấy đủ thông tin.

## API gợi ý
`GET /api/v1/stores/{storeId}/serial-units/lookup?serial=` · `GET /api/v1/households/{householdId}/serial-units/lookup?serial=`
