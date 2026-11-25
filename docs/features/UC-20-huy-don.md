# UC-20 Hủy đơn đã hoàn tất

- **Milestone:** M4 · **Context:** sales (+ inventory, customer) · **Tác nhân:** Quản lý, Chủ hộ
- **Quy tắc:** BR-15, BR-26, BR-29, BR-37 · **Quyền:** `sale.cancel`

## Luồng chính
1. Mở đơn `COMPLETED`, chọn "Hủy đơn", nhập lý do (bắt buộc).
2. Trong một transaction:
   - tạo Reversal với đủ mọi dòng (số lượng, thành tiền, giá vốn snapshot), `created_at` = bây giờ (BR-29);
   - hoàn kho: dòng `QUANTITY` và nguyên liệu → `StockMovement` `SALE_REVERSAL` với `unit_cost` = giá vốn snapshot; serial → `IN_STOCK`;
   - (M5) Receivable của đơn → `VOIDED`;
   - (M6) phiếu cọc đã áp → `ACTIVE`, serial → `RESERVED`;
   - đơn → `CANCELLED`, gắn `reversal_id`; nội dung đơn không đổi.
3. Ghi nhật ký.

## Ngoại lệ
- Đơn không `COMPLETED` → lỗi.
- Đơn bán nợ đã thu một phần → từ chối (ngoài phạm vi MVP).
- Đơn đã có `einvoice_ref` (sau MVP) → phải xử lý hóa đơn điện tử trước.

## Tiêu chí chấp nhận
- **Given** đơn 30.000 hoàn tất ngày 01/12, 2 bánh mì **When** hủy ngày 02/12 **Then** đơn `CANCELLED`, tồn tăng lại 2, báo cáo ngày 01/12 vẫn 30.000, ngày 02/12 ghi −30.000.
- **Given** đơn đã hủy **When** hủy lần nữa **Then** lỗi.
- **Given** Nhân viên **When** gọi API hủy **Then** 403.
- **Given** gọi hủy **Then** không có câu lệnh UPDATE nào thay đổi `total`, dòng đơn hay thanh toán (kiểm bằng test so sánh trước và sau).

## API gợi ý
`POST /api/v1/stores/{storeId}/orders/{id}/cancel` (body: `reason`) · `GET .../orders?status=` (danh sách đơn theo ngày)
