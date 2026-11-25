# UC-13 Nhập hàng serial

- **Milestone:** M6 · **Context:** inventory (+ cashbook) · **Capability:** `inv.serial` · **Tác nhân:** Nhân viên, Quản lý, Chủ hộ
- **Quy tắc:** BR-08, BR-09, BR-11, BR-32, BR-37 · **Quyền:** `stock.receive`

## Luồng chính
1. Chọn nhà cung cấp và sản phẩm `SERIAL`.
2. Quét mã vạch (camera) hoặc gõ từng IMEI/số khung; mỗi dòng kèm giá nhập riêng. Có thao tác "áp cùng giá cho các dòng đã chọn".
3. Hệ thống kiểm tra trùng ngay khi nhập (trong phiếu và trong hộ).
4. Xác nhận. Trong một transaction: tạo PurchaseReceipt, mỗi serial một `SerialUnit` `IN_STOCK` với `unit_cost` của nó, CashEntry `PURCHASE_PAYMENT`, khóa giá vốn, đánh dấu giao dịch.

Một phiếu được có cả dòng serial và dòng số lượng (ví dụ điện thoại kèm ốp lưng).

## Ngoại lệ
- Serial trùng một đơn vị `IN_STOCK`/`RESERVED` trong hộ, hoặc trùng trong cùng phiếu → báo dòng lỗi. Người dùng bỏ dòng lỗi rồi gửi lại; server không ghi phiếu khi còn dòng lỗi.
- Serial đã `SOLD` trước đây → cho nhập (máy quay lại kho là trường hợp hợp lệ), hiển thị cảnh báo.

## Tiêu chí chấp nhận
- **Given** nhập 3 iPhone với giá 18.000.000, 18.200.000, 17.900.000 **Then** có 3 SerialUnit `IN_STOCK` đúng giá từng chiếc, serial lưu chữ hoa không khoảng trắng.
- **Given** IMEI `351234567890123` đang `IN_STOCK` ở cửa hàng A **When** nhập "35 1234 5678 90123" ở cửa hàng B cùng hộ **Then** sau chuẩn hóa bị trùng, lỗi `BR-11` ở dòng đó, không có gì được ghi.
- **Given** IMEI từng `SOLD` **When** nhập lại **Then** thành công, có cảnh báo.
- Unique index bảo đảm trùng không lọt qua dù hai phiếu gửi cùng lúc.

## API gợi ý
Dùng chung `POST /api/v1/stores/{storeId}/purchase-receipts` (dòng có `serials[]`) · `POST .../serial-units/check` (kiểm tra trùng khi gõ)
