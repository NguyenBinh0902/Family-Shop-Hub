# Context: Identity & Access (`internal/identity`)

Quản lý tài khoản, hộ, thành viên, phiên đăng nhập, lời mời. Cung cấp cho các context khác: xác thực request và `tenant.Scope`.

## Thực thể

**User** (toàn hệ thống, không có `household_id`)

| Trường | Kiểu | Ghi chú |
| --- | --- | --- |
| id | UUID v7 | |
| phone | string | Đã chuẩn hóa, duy nhất toàn hệ thống |
| display_name | string | 1–80 ký tự |
| password_hash | string | argon2id |
| status | `ACTIVE` \| `DISABLED` | |
| created_at | timestamptz | |

**Household**

| Trường | Kiểu | Ghi chú |
| --- | --- | --- |
| id | UUID v7 | |
| name | string | 1–120 ký tự |
| plan_code | string | MVP: `FREE` |
| max_stores, max_members | int | Từ gói; mặc định 5 và 20 (BR-27) |
| status | `ACTIVE` \| `SUSPENDED` | |

**Membership**

| Trường | Kiểu | Ghi chú |
| --- | --- | --- |
| id | UUID v7 | |
| household_id | UUID | |
| user_id | UUID | |
| store_id | UUID, nullable | Rỗng khi `role = OWNER`; bắt buộc khi `MANAGER`/`STAFF` |
| role | `OWNER` \| `MANAGER` \| `STAFF` | |
| status | `ACTIVE` \| `REVOKED` | |
| created_at, revoked_at, revoked_by | | |

Ràng buộc: duy nhất `(household_id, user_id, store_id)` trong các membership `ACTIVE` (store_id rỗng coi là một giá trị). Một người có thể vừa là OWNER vừa có membership cửa hàng; OWNER đã đủ quyền nên UI không cho gán thêm.

**Session**: `token_hash` (SHA-256 của token ngẫu nhiên 32 byte), `user_id`, `created_at`, `last_seen_at`, `expires_at` (trượt, 30 ngày kể từ lần dùng cuối), `revoked_at`. Token gốc chỉ nằm trong cookie.

**Invitation**: `household_id`, `role`, `store_id` (nullable như Membership), `code_hash`, `created_by`, `expires_at` (7 ngày), `status` (`PENDING` \| `ACCEPTED` \| `REVOKED` \| `EXPIRED`), `accepted_by`.

**PasswordResetOTP**: `phone`, `code_hash`, `expires_at` (5 phút), `attempts` (tối đa 5), `used_at`.

## Chuẩn hóa số điện thoại

Bỏ khoảng trắng, dấu chấm, gạch; `+84` hoặc `84` ở đầu đổi thành `0`; kết quả phải là 10 chữ số bắt đầu bằng `0`. Đặt trong `kernel.Phone`, dùng chung với Customer.

## Bất biến

- Hộ luôn có ít nhất một OWNER `ACTIVE` (BR-36).
- Số người dùng khác nhau có membership `ACTIVE` trong hộ ≤ `max_members` (BR-27).
- Đổi mật khẩu bằng OTP thu hồi mọi phiên của người đó.

## Cung cấp cho context khác

- Middleware xác thực: cookie → Session → User.
- Middleware tenant: `storeId` hoặc `householdId` trên URL → kiểm tra membership → `tenant.Scope{HouseholdID, StoreID, UserID, Role}`.
- `permission.Can(scope, perm)` theo `roles-permissions.md`; kiểm tra capability gọi sang `setup`.

## Use case

UC-01, UC-02, UC-03, UC-04.
