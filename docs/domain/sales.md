# Context: Sales (`internal/sales`)

Đơn bán, thanh toán, hủy đơn (bút toán đảo), phiếu cọc, bàn, hóa đơn. Là nơi ghi nhận doanh thu.

## Thực thể

**SalesOrder**

| Trường | Kiểu | Ghi chú |
| --- | --- | --- |
| id | UUID v7 | Do client sinh (BR-17) |
| household_id, store_id | UUID | |
| number | string | `yyMMdd-NNN`, NNN đếm lại mỗi ngày nghiệp vụ theo cửa hàng; cấp khi hoàn tất (đơn mở dùng số tạm) |
| status | `OPEN` \| `COMPLETED` \| `CANCELLED` | |
| table_id | UUID, nullable | Chỉ khi có `fnb.table` |
| customer_id | UUID, nullable | |
| subtotal | Money | Tổng dòng |
| discount_type | `NONE` \| `AMOUNT` \| `PERCENT` | Giảm giá cấp đơn |
| discount_value | int64 | Số tiền hoặc phần trăm (0–100) |
| discount_amount | Money | Đã quy đổi, làm tròn (BR-33) |
| total | Money | `subtotal − discount_amount` ≥ 0 |
| note | string | |
| version | int | Tăng mỗi lần sửa đơn `OPEN` (optimistic locking) |
| created_by, created_at | | |
| completed_by, completed_at | | Thời điểm ghi nhận doanh thu (BR-25) |
| cancelled_at, reversal_id | | |
| einvoice_ref | string, nullable | TAX-05 |

**OrderLine**: `id`, `order_id`, `product_id`, `product_name` (snapshot), `options` (jsonb snapshot `[{group, name, price_delta}]`), `quantity`, `unit_price` (giá bán + tổng price_delta), `line_total`, `unit_cost` (snapshot, BR-10), `cost_missing` (bool, BR-28), `serial_unit_id`, `serial` (snapshot).

Giá vốn của dòng món có công thức = tổng (định lượng nguyên liệu × giá vốn bình quân nguyên liệu) của công thức và tùy chọn, tính lúc hoàn tất.

Giảm giá cấp đơn không phân bổ xuống dòng trong MVP; báo cáo lợi nhuận dùng `total` của đơn và tổng `unit_cost × quantity` các dòng.

**Payment**: `id`, `order_id`, `method` (`CASH` \| `BANK_TRANSFER` \| `EWALLET` \| `DEBT` \| `DEPOSIT`), `account_id` (tài khoản/ví của chủ thể, nullable), `amount` (> 0), `tendered` (khách đưa, chỉ `CASH`, để in tiền thối).

**Reversal**: `id`, `household_id`, `store_id`, `order_id`, `reason` (bắt buộc), `created_by`, `created_at` (ngày ghi nhận âm, BR-29), dòng `{order_line_id, quantity, amount, unit_cost}` (MVP: đủ mọi dòng), `total`.

**Deposit** (tiền nhận cọc ghi bằng CashEntry `DEPOSIT_RECEIVED`, mang hình thức và tài khoản; không tạo Payment): `id`, `household_id`, `store_id`, `customer_id` (bắt buộc), `serial_unit_id`, `amount`, `status` (`ACTIVE` \| `APPLIED` \| `CANCELLED`), `applied_order_id`, `cancel_mode` (`REFUND` \| `FORFEIT`), `cancel_reason`, `cancelled_at`, `cancelled_by`, `created_by`, `created_at`.

**DiningTable**: `id`, `store_id`, `name`, `sort_order`, `status` (`ACTIVE` \| `ARCHIVED`). Một bàn có tối đa một đơn `OPEN`.

## Trạng thái đơn

```
(fnb.table) OPEN → COMPLETED       OPEN → bị xóa (chưa ghi gì, không cần bút toán)
(bán tại quầy) tạo thẳng COMPLETED
COMPLETED → CANCELLED (tạo Reversal)
```

## Hoàn tất đơn (một transaction DB)

1. Kiểm tra cửa hàng `ACTIVE`, quyền, giảm giá trong mức trần.
2. Kiểm tra thanh toán: tổng = `total` (BR-18); có `DEBT` thì phải có khách (BR-19); `DEPOSIT` phải khớp phiếu cọc của chính serial trong đơn.
3. Gọi Inventory: trừ kho, bán serial, lấy giá vốn.
4. Ghi đơn `COMPLETED`, dòng (snapshot), thanh toán; cấp số đơn.
5. Có `DEBT` → gọi Customer tạo Receivable. Có `DEPOSIT` → phiếu cọc `APPLIED`.
6. Đánh dấu `has_transactions` cho cửa hàng và sản phẩm.
7. Sau commit: gọi `InvoiceIssuer` (no-op trong MVP; lỗi không làm hỏng đơn, chỉ log).

Idempotent: `INSERT ... ON CONFLICT (id) DO NOTHING`; nếu đơn đã tồn tại với cùng `store_id` thì trả đơn đó (HTTP 200); khác cửa hàng hoặc khác nội dung → lỗi 409.

## Hủy đơn

Tạo Reversal (đủ dòng), gọi Inventory hoàn kho/serial, gọi Customer đảo Receivable chưa thu. Nếu khách đã trả một phần nợ của đơn này thì từ chối hủy trong MVP (lỗi rõ ràng; hoàn tiền nợ đã thu là ngoài phạm vi). Đơn có `DEPOSIT` khi hủy: phiếu cọc chuyển lại `ACTIVE` và serial về `RESERVED`.

## Ghi chú cho Reporting

Doanh thu = đơn `COMPLETED` (theo `completed_at`) trừ Reversal (theo `created_at`). Không lưu bảng tổng hợp riêng (BR-24).

## Use case

UC-16, UC-17, UC-18, UC-19, UC-20, UC-21.
