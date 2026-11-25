# Context: Customer & Debt (`internal/customer`)

Khách hàng dùng chung trong hộ; công nợ và lịch sử theo cửa hàng.

## Thực thể

**Customer** (cấp hộ)

| Trường | Kiểu | Ghi chú |
| --- | --- | --- |
| id, household_id | UUID | |
| name | string | 1–120 ký tự |
| phone | string, nullable | `kernel.Phone`; duy nhất trong hộ nếu có (BR-20) |
| note | string | |
| created_at, created_in_store_id | | |

**Receivable** (cấp cửa hàng): `id`, `household_id`, `store_id`, `customer_id`, `order_id`, `original_amount`, `remaining_amount` (≥ 0), `status` (`OPEN` \| `PAID` \| `VOIDED`), `created_at`.

**DebtCollection** (cấp cửa hàng): `id`, `store_id`, `customer_id`, `amount`, `method`, `account_id`, `collected_at`, `created_by`, `cash_entry_id`; phân bổ `{receivable_id, amount}`.

## Quy tắc

- Tạo khách với số điện thoại đã có trong hộ → trả khách hiện có, không tạo mới (BR-20).
- Lịch sử mua của khách = đơn bán có `customer_id` (đọc từ Sales qua port); nhân viên chỉ thấy đơn của cửa hàng mình (BR-21).
- Thu nợ: phân bổ vào Receivable `OPEN` cũ nhất trước; tổng thu ≤ tổng nợ còn lại tại cửa hàng (BR-22). Mỗi lần thu tạo một CashEntry `INCOME` loại `DEBT_COLLECTION` qua port sang Cashbook (BR-31: không phải doanh thu).
- Hủy đơn bán nợ chưa thu đồng nào → Receivable `VOIDED`.
- Không xóa khách đã có giao dịch; sửa tên, ghi chú, số điện thoại được (số mới vẫn phải duy nhất).

## Cung cấp cho context khác

- Sales: `FindOrCreateByPhone`, `CreateReceivable`, `VoidReceivable`, `HasCollections(orderID)`.

## Use case

UC-22, UC-23, UC-24.
