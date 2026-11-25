# UC-26 Ghi phiếu thu chi

- **Milestone:** M5 · **Context:** cashbook · **Capability:** `core.cashbook` · **Tác nhân:** Quản lý, Chủ hộ
- **Quy tắc:** BR-23, BR-26, BR-32 · **Quyền:** `cash.entry`

## Luồng chính
1. Chọn "Chi" hoặc "Thu", danh mục (lương, mặt bằng, điện nước, chi khác / thu khác), số tiền, hình thức, ngày (mặc định hôm nay, không chọn ngày tương lai), ghi chú.
2. Lưu. Phiếu bất biến.
3. Sổ thu chi: danh sách phiếu theo khoảng ngày, lọc theo danh mục; phiếu hệ thống (trả nhà cung cấp, thu nợ, cọc) hiển thị kèm liên kết tới chứng từ gốc.
4. **Điều chỉnh:** từ một phiếu tay, tạo phiếu điều chỉnh (số tiền chênh lệch, có thể âm) với lý do bắt buộc.

## Ngoại lệ
- Số tiền ≤ 0 ở phiếu thường → lỗi.
- Ngày tương lai → lỗi.
- Điều chỉnh phiếu hệ thống → không cho phép.

## Tiêu chí chấp nhận
- **Given** ghi chi mặt bằng 5.000.000 ngày 01/12 **Then** báo cáo chi phí vận hành tháng 12 tăng 5.000.000.
- **Given** phiếu ghi nhầm 5.000.000 thay vì 4.500.000 **When** điều chỉnh −500.000 **Then** phiếu gốc không đổi, chi phí vận hành tháng còn 4.500.000, nhật ký có lý do.
- **Given** gọi API sửa hoặc xóa phiếu **Then** không tồn tại endpoint đó.
- **Given** Nhân viên **When** ghi phiếu chi **Then** 403.

## API gợi ý
`GET/POST /api/v1/stores/{storeId}/cash-entries` · `POST .../cash-entries/{id}/adjustments`
