# UC-07 Chuyển giữa các cửa hàng, góc nhìn Toàn hộ

- **Milestone:** M2 · **Context:** setup (dữ liệu), frontend (điều hướng) · **Tác nhân:** Mọi người dùng
- **Quy tắc:** BR-02

## Luồng chính
1. Người dùng mở bộ chuyển ở đầu màn hình; hệ thống liệt kê các cửa hàng người dùng có membership còn hiệu lực (Chủ hộ: mọi cửa hàng của hộ). Chủ hộ thấy thêm mục "Toàn hộ" ở đầu danh sách.
2. Chọn cửa hàng → URL đổi sang `/stores/{storeId}/...`; tên và màu cửa hàng hiển thị cố định ở đầu màn hình.
3. Chọn "Toàn hộ" → URL `/households/{householdId}/overview` (chỉ báo cáo, không bán hàng).
4. Mỗi tab trình duyệt giữ ngữ cảnh riêng theo URL. Cửa hàng chọn gần nhất được nhớ ở trình duyệt để mở lần sau.

Người thuộc nhiều hộ chuyển hộ ở cùng bộ chuyển (nhóm theo hộ).

## Ngoại lệ
- Chỉ thuộc một cửa hàng → không hiện bộ chuyển.
- Mở URL cửa hàng không có quyền → 403, đưa về danh sách cửa hàng của mình.
- Membership bị thu hồi khi đang dùng → request kế tiếp 403, đưa về danh sách cửa hàng.
- Cửa hàng `SUSPENDED`/`CLOSED` vẫn hiện trong danh sách với nhãn trạng thái; vào được để xem.

## Tiêu chí chấp nhận
- **Given** người dùng là Nhân viên ở A và Quản lý ở B **When** mở bộ chuyển **Then** thấy đúng A và B, không thấy "Toàn hộ".
- **Given** mở hai tab, tab 1 ở A, tab 2 ở B **When** tạo đơn ở tab 1 **Then** đơn thuộc A.
- **Given** Nhân viên chỉ thuộc A **When** mở URL của cửa hàng C cùng hộ **Then** trả 403.

## API gợi ý
`GET /api/v1/me/stores` (nhóm theo hộ, kèm vai trò, trạng thái, màu)
