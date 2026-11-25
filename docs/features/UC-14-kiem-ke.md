# UC-14 Kiểm kê và điều chỉnh tồn

- **Milestone:** M3 · **Context:** inventory · **Capability:** `inv.quantity` · **Tác nhân:** Quản lý, Chủ hộ
- **Quy tắc:** BR-26, BR-34, BR-37 · **Quyền:** `stock.adjust`

## Luồng chính
1. Tạo phiếu kiểm kê, chọn sản phẩm cần đếm (một nhóm hoặc tất cả sản phẩm `QUANTITY`).
2. Hệ thống hiển thị tồn trên hệ thống tại thời điểm mở phiếu; người dùng nhập số đếm thực tế.
3. Nhập lý do (bắt buộc), xác nhận.
4. Mỗi dòng có chênh lệch: `StockMovement` `ADJUSTMENT` với `unit_cost` = giá vốn bình quân hiện tại; cập nhật `on_hand`. Giá vốn bình quân không đổi. Ghi nhật ký.

## Ngoại lệ
- Có giao dịch kho xảy ra sau khi mở phiếu → chênh lệch tính theo tồn **lúc xác nhận**, và hiển thị cho người dùng thấy tồn đã thay đổi trước khi xác nhận.
- Số đếm âm → lỗi.

## Tiêu chí chấp nhận
- **Given** sữa trên hệ thống 2.300 ml **When** đếm được 2.000 ml **Then** movement `ADJUSTMENT` −300 ml, tồn 2.000 ml, nhật ký có lý do.
- **Given** sản phẩm số đếm bằng tồn **Then** không tạo movement.
- **Given** Nhân viên **When** gọi API kiểm kê **Then** 403.

## Ngoài phạm vi
Kiểm kê hàng serial.

## API gợi ý
`POST /api/v1/stores/{storeId}/stock-counts` · `GET .../stock-counts` · `GET /api/v1/stores/{storeId}/stock` (danh sách tồn)
