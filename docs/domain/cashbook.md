# Context: Cashbook (`internal/cashbook`)

Phiếu thu chi ngoài bán hàng. Là nguồn số liệu chi phí vận hành và dòng tiền.

## Thực thể

**CashEntry** (cấp cửa hàng)

| Trường | Kiểu | Ghi chú |
| --- | --- | --- |
| id, household_id, store_id | UUID | |
| direction | `INCOME` \| `EXPENSE` | |
| category | string | Mã danh mục, xem bảng dưới |
| amount | Money | > 0 |
| method | `CASH` \| `BANK_TRANSFER` \| `EWALLET` | |
| account_id | UUID, nullable | Tài khoản/ví của chủ thể |
| occurred_on | date | Ngày nghiệp vụ, ≤ hôm nay |
| note | string | |
| source | `MANUAL` \| `SYSTEM` | `SYSTEM` do context khác tạo |
| ref_type, ref_id | | Ví dụ `purchase_receipt`, `debt_collection`, `deposit` |
| adjusts_entry_id | UUID, nullable | Phiếu điều chỉnh của phiếu nào (BR-23) |
| created_by, created_at | | |

Bất biến (BR-23). Phiếu điều chỉnh có cùng `direction`, `amount` có thể âm hoặc dương, bắt buộc có ghi chú; phiếu `SYSTEM` không điều chỉnh tay được.

## Danh mục

| Mã | Hướng | Tên hiển thị | Chi phí vận hành? | Nguồn |
| --- | --- | --- | --- | --- |
| `SALARY` | EXPENSE | Lương | Có | Tay |
| `RENT` | EXPENSE | Mặt bằng | Có | Tay |
| `UTILITIES` | EXPENSE | Điện, nước, internet | Có | Tay |
| `OTHER_EXPENSE` | EXPENSE | Chi khác | Có | Tay |
| `PURCHASE_PAYMENT` | EXPENSE | Trả tiền nhập hàng | **Không** (BR-32) | Hệ thống (UC-12, UC-13) |
| `DEPOSIT_REFUND` | EXPENSE | Hoàn tiền cọc | Không | Hệ thống (UC-19) |
| `OTHER_INCOME` | INCOME | Thu khác | — | Tay |
| `DEBT_COLLECTION` | INCOME | Thu nợ | — (BR-31) | Hệ thống (UC-24) |
| `DEPOSIT_RECEIVED` | INCOME | Nhận cọc | — (BR-30) | Hệ thống (UC-19) |

Giữ tiền cọc (hủy cọc kiểu `FORFEIT`) không tạo phiếu thu chi, vì tiền đã vào lúc nhận cọc; Reporting đọc phiếu cọc để tính vào "Thu nhập khác". Danh mục tùy chỉnh theo cửa hàng để sau MVP.

## Cung cấp cho context khác

`RecordSystemEntry(ctx, scope, entry)` chạy trong transaction của bên gọi.

## Use case

UC-26.
