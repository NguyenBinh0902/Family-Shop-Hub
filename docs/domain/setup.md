# Context: Business Setup (`internal/setup`)

Quản lý chủ thể kinh doanh, cửa hàng, loại hình và capability. Là nơi duy nhất trả lời câu hỏi "cửa hàng này có capability X không".

## Thực thể

**BusinessEntity** (cấp hộ)

| Trường | Kiểu | Ghi chú |
| --- | --- | --- |
| id, household_id | UUID | |
| legal_name | string | Tên hộ kinh doanh trên giấy đăng ký |
| owner_name | string | Người đứng tên |
| tax_code | string, nullable | Chỉ chữ số, 10–13 ký tự; có thể bổ sung sau |
| address | string | |
| bank_accounts | danh sách | `{id, bank_name, account_number, account_holder}` |
| wallets | danh sách | `{id, provider, account}` (MoMo, ZaloPay...) |
| created_at, updated_at | | |

Tài khoản và ví có `id` riêng để `Payment` tham chiếu (TAX-02). Xóa một tài khoản đã được tham chiếu thì chỉ đánh dấu ngưng dùng.

**Store** (cấp cửa hàng, `household_id` + `id`)

| Trường | Kiểu | Ghi chú |
| --- | --- | --- |
| id, household_id | UUID | |
| business_entity_id | UUID | BR-03, BR-04 |
| name, address, phone | string | |
| color | string | Một trong bảng màu cố định (8 màu) để phân biệt cửa hàng |
| store_type_code, store_type_version | string, int | BR-03 |
| costing_method | `WEIGHTED_AVERAGE` | Enum để mở rộng sau |
| costing_locked_at | timestamptz, nullable | Đặt khi có giao dịch kho đầu tiên (BR-08) |
| allow_negative_stock | bool | Mặc định theo loại hình (BR-13) |
| staff_discount_cap_percent | int 0–100 | Mặc định 0 |
| has_transactions | bool | Đặt `true` ở giao dịch đầu tiên bất kỳ; khóa đổi loại hình, chủ thể |
| status | `ACTIVE` \| `SUSPENDED` \| `CLOSED` | |
| created_at, closed_at | | |

**StoreType**: đọc từ file cấu hình nhúng (xem `store-types.md`), không lưu DB.

## Trạng thái cửa hàng

```
ACTIVE ⇄ SUSPENDED
ACTIVE → CLOSED, SUSPENDED → CLOSED (không quay lại)
```

`SUSPENDED` và `CLOSED`: chỉ đọc. `CLOSED` không tính vào giới hạn số cửa hàng.

## Cung cấp cho context khác

```go
// Port mà các context khác khai báo trong app của mình, adapter gọi sang setup.
type StoreInfo interface {
    Store(ctx context.Context, scope tenant.Scope) (StoreSnapshot, error)
}
// StoreSnapshot gồm: capabilities, allow_negative_stock, staff_discount_cap_percent,
// costing_method, status, business_entity_id.
```

Inventory gọi `MarkCostingLocked`, mọi context ghi giao dịch gọi `MarkHasTransactions` trong cùng transaction.

## Use case

UC-05, UC-06, UC-07, UC-08.
