# UC-15 Xem cảnh báo tồn thấp

- **Milestone:** M3 · **Context:** inventory · **Capability:** `inv.quantity` · **Tác nhân:** Quản lý, Chủ hộ (Nhân viên xem được danh sách)
- **Quyền:** `stock.view`

## Luồng chính
1. Màn hình tồn kho có bộ lọc "Sắp hết": sản phẩm `QUANTITY` có `min_stock` và `on_hand ≤ min_stock`.
2. Trang chủ cửa hàng hiển thị số sản phẩm sắp hết; bấm vào ra danh sách.
3. Từ danh sách, bấm "Nhập hàng" mở UC-12 với các sản phẩm đó đã điền sẵn.

## Tiêu chí chấp nhận
- **Given** trân châu `min_stock` 1.000 g, tồn 800 g **Then** xuất hiện trong "Sắp hết".
- **Given** sản phẩm không đặt `min_stock` **Then** không bao giờ xuất hiện.
- **Given** sản phẩm `DISCONTINUED` **Then** không xuất hiện.

## Ngoài phạm vi
Thông báo đẩy, gửi SMS cảnh báo.

## API gợi ý
`GET /api/v1/stores/{storeId}/stock?filter=low`
