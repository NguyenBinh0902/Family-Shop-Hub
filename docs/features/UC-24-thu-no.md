# UC-24 Thu nợ

- **Milestone:** M5 · **Context:** customer (+ cashbook) · **Capability:** `core.customer` · **Tác nhân:** Nhân viên, Quản lý, Chủ hộ
- **Quy tắc:** BR-22, BR-31 · **Quyền:** `debt.collect`

## Luồng chính
1. Tìm khách, hệ thống hiển thị nợ còn lại tại cửa hàng hiện tại và các khoản nợ (cũ nhất trước).
2. Nhập số tiền thu, hình thức (tiền mặt, chuyển khoản, ví).
3. Trong một transaction: tạo DebtCollection, phân bổ vào Receivable cũ nhất trước (khoản trả đủ → `PAID`), tạo CashEntry `INCOME` loại `DEBT_COLLECTION`.
4. In hoặc chia sẻ biên nhận thu nợ (tùy chọn).

## Ngoại lệ
- Số tiền > tổng nợ còn lại tại cửa hàng → lỗi `BR-22`.
- Số tiền ≤ 0 → lỗi.
- Thu nợ của cửa hàng khác → không thể (thao tác luôn trong phạm vi cửa hàng hiện tại).

## Tiêu chí chấp nhận
- **Given** nợ đơn A 100.000 (01/11) và đơn B 80.000 (05/11) **When** thu 150.000 **Then** A `PAID`, B còn 30.000, có CashEntry `DEBT_COLLECTION` 150.000, doanh thu không đổi (BR-31).
- **Given** nợ còn 30.000 **When** thu 50.000 **Then** lỗi `BR-22`.
- **Given** đơn A đã được thu một phần **When** hủy đơn A **Then** bị từ chối (xem UC-20).

## API gợi ý
`GET /api/v1/stores/{storeId}/customers/{id}/receivables` · `POST /api/v1/stores/{storeId}/debt-collections`
