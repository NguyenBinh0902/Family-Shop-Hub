# UC-28 Báo cáo Toàn hộ

- **Milestone:** M8 · **Context:** reporting · **Tác nhân:** Chủ hộ
- **Quy tắc:** BR-24 · **Quyền:** `report.household`

## Luồng chính
1. Vào "Toàn hộ" từ bộ chuyển (UC-07).
2. Hiển thị cùng bộ chỉ số như UC-27, gộp mọi cửa hàng của hộ (kể cả cửa hàng đã đóng, với dữ liệu trong kỳ).
3. Bảng so sánh theo cửa hàng (doanh thu, lợi nhuận gộp, lợi nhuận ròng, tỷ trọng doanh thu) và theo chủ thể kinh doanh.
4. Bấm vào một cửa hàng → mở báo cáo cửa hàng đó cùng kỳ.

## Tiêu chí chấp nhận
- **Given** cửa hàng A doanh thu 3.000.000, B 2.000.000 trong kỳ **Then** Toàn hộ 5.000.000, tỷ trọng A 60%, B 40%.
- **Given** A và B thuộc chủ thể X, C thuộc chủ thể Y **Then** bảng theo chủ thể có X = A + B, Y = C.
- **Given** Quản lý **When** gọi API Toàn hộ **Then** 403.
- Tổng Toàn hộ luôn bằng tổng các cửa hàng (kiểm bằng test), vì cùng tính từ dữ liệu gốc (BR-24).

## API gợi ý
`GET /api/v1/households/{householdId}/reports/summary?from=&to=` (kèm `byStore`, `byEntity`)
