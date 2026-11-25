# Tài liệu FamilyShop Hub

Toàn bộ tài liệu viết bằng tiếng Việt; tên trong code theo `glossary.md`.

## Đọc gì khi làm một feature

1. `roadmap.md`: feature thuộc milestone nào, làm tới đâu.
2. `features/UC-xx-*.md`: luồng, ngoại lệ, tiêu chí chấp nhận (mỗi tiêu chí là ít nhất một test).
3. `rules.md`: các BR mà use case dẫn chiếu.
4. `domain/<context>.md`: thực thể, trạng thái, giao tiếp giữa context.
5. `conventions/`: cách viết code.

## Danh mục

| File | Nội dung |
| --- | --- |
| `overview.md` | Mục tiêu, người dùng, phạm vi MVP và ngoài phạm vi |
| `glossary.md` | Thuật ngữ tiếng Việt và tên trong code |
| `roles-permissions.md` | Vai trò, mã quyền, ma trận quyền |
| `rules.md` | Quy tắc nghiệp vụ BR-xx |
| `store-types.md` | Capability, loại hình, cấu hình, thuộc tính động |
| `compliance-nfr.md` | Yêu cầu thuế TAX-xx, phi chức năng NFR-xx |
| `open-questions.md` | Đã chốt, CHƯA CHỐT, mặc định BA cần xem lại |
| `roadmap.md` | Milestone M0–M8 và quy trình làm việc |
| `domain/` | 8 bounded context: identity, setup, catalog, inventory, sales, customer, cashbook, reporting |
| `features/` | 30 use case UC-01 đến UC-30 |
| `adr/` | Quyết định kiến trúc |
| `conventions/` | Quy ước backend Go và frontend |

## Khi tài liệu sai hoặc thiếu

Không tự đoán. Hỏi chủ repo, rồi sửa tài liệu trong cùng PR với code. Mục CHƯA CHỐT thì không implement.
