# UC-30 Báo cáo hàng bán chạy và tồn kho

- **Milestone:** M8 · **Context:** reporting · **Tác nhân:** Quản lý, Chủ hộ
- **Quy tắc:** BR-24 · **Quyền:** `report.store`, `cost.view` (giá trị tồn)

## Luồng chính
1. **Hàng bán chạy:** chọn kỳ; bảng sản phẩm theo số lượng bán và doanh thu (đã trừ đơn hủy trong kỳ), sắp xếp được; mặc định top 20. Món có tùy chọn gộp theo sản phẩm.
2. **Tồn kho:** danh sách sản phẩm `QUANTITY` với tồn hiện tại, giá vốn bình quân, giá trị tồn; sản phẩm serial với số chiếc `IN_STOCK`, `RESERVED` và tổng giá vốn. Tổng giá trị tồn kho của cửa hàng.

## Tiêu chí chấp nhận
- **Given** trong kỳ bán 30 trà sữa, 12 cà phê, hủy 2 trà sữa **Then** trà sữa 28, cà phê 12, trà sữa đứng đầu.
- **Given** sữa tồn 15.000 ml giá vốn 32 đ/ml và 2 điện thoại giá vốn 18.000.000 và 18.200.000 **Then** giá trị tồn = 480.000 + 36.200.000 = 36.680.000.
- **Given** sản phẩm tồn âm **Then** hiển thị số âm, không tính vào giá trị tồn.

## API gợi ý
`GET /api/v1/stores/{storeId}/reports/top-products?from=&to=&limit=` · `GET /api/v1/stores/{storeId}/reports/inventory-value`
