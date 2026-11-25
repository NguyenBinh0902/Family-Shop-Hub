# UC-11 Khai báo công thức nguyên liệu

- **Milestone:** M7 · **Context:** catalog · **Capability:** `fnb.recipe` · **Tác nhân:** Quản lý, Chủ hộ
- **Quy tắc:** BR-10, BR-14, BR-34 · **Quyền:** `catalog.manage`

## Luồng chính
1. Mở sản phẩm, chọn "Công thức".
2. Thêm dòng: nguyên liệu (sản phẩm `QUANTITY` cùng cửa hàng) và định lượng theo đơn vị của nguyên liệu (ví dụ 200 ml sữa, 15 g trà).
3. Hệ thống hiển thị giá vốn ước tính của một phần theo giá vốn bình quân hiện tại (chỉ người có `cost.view`).
4. Lưu. Áp dụng cho đơn hoàn tất sau thời điểm lưu.

## Ngoại lệ
- Nguyên liệu chính là sản phẩm đó, hoặc nguyên liệu có công thức → lỗi (không lồng công thức).
- Định lượng ≤ 0 → lỗi.

## Tiêu chí chấp nhận
- **Given** "Trà sữa" có công thức 200 ml sữa + 15 g trà **When** hoàn tất đơn 3 ly **Then** kho sữa giảm 600 ml, trà giảm 45 g, mỗi nguyên liệu có một `StockMovement` loại `RECIPE_CONSUMPTION`.
- **Given** sữa giá vốn 30 đ/ml, trà 400 đ/g **Then** giá vốn snapshot của dòng là 200 × 30 + 15 × 400 = 12.000 đ/ly.
- **Given** định lượng 0,5 g **Then** lưu chính xác là 500 phần nghìn (BR-34).
- **Given** sửa công thức sau khi đã bán **Then** đơn cũ không đổi giá vốn.

## API gợi ý
`GET/PUT /api/v1/stores/{storeId}/products/{id}/recipe`
