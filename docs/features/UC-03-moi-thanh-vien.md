# UC-03 Mời thành viên, gán vai trò theo cửa hàng

- **Milestone:** M1 (gán theo cửa hàng hoàn thiện khi có cửa hàng ở M2) · **Context:** identity · **Tác nhân:** Chủ hộ
- **Quy tắc:** BR-02, BR-26, BR-27 · **Quyền:** `member.manage`

## Luồng chính: mời người mới
1. Chủ hộ chọn vai trò: Chủ hộ, hoặc Quản lý/Nhân viên của một cửa hàng.
2. Hệ thống tạo Invitation (hạn 7 ngày, dùng một lần) và trả đường link có mã.
3. Chủ hộ tự gửi link (Zalo, tin nhắn). Ứng dụng có nút "Sao chép link" và "Chia sẻ" (Web Share API).
4. Người nhận mở link:
   - đã đăng nhập → xác nhận tham gia;
   - chưa có tài khoản → form đăng ký (như UC-01 nhưng không tạo hộ) rồi tham gia;
   - có tài khoản nhưng chưa đăng nhập → đăng nhập rồi tham gia.
5. Hệ thống tạo Membership, Invitation chuyển `ACCEPTED`.

## Luồng phụ: gán thêm cửa hàng cho thành viên hiện có
Chủ hộ chọn thành viên, thêm vai trò ở một cửa hàng khác. Tạo Membership trực tiếp, không cần lời mời.

## Ngoại lệ
- Hộ đã đủ 20 thành viên (BR-27) → không tạo lời mời, và kiểm tra lại khi chấp nhận.
- Link hết hạn, đã dùng, bị thu hồi → thông báo rõ, đề nghị xin link mới.
- Người nhận đã có membership đó → báo "Bạn đã là thành viên".
- Cửa hàng `CLOSED` → không mời được vào cửa hàng đó.

## Tiêu chí chấp nhận
- **Given** Chủ hộ tạo lời mời Nhân viên cho cửa hàng A **When** một người dùng khác chấp nhận **Then** người đó có Membership `STAFF` tại A và thấy A trong danh sách cửa hàng của mình.
- **Given** lời mời đã được chấp nhận **When** người thứ hai mở cùng link **Then** bị từ chối.
- **Given** hộ có 20 thành viên **When** tạo lời mời mới **Then** trả lỗi `BR-27`.
- **Given** người dùng là Nhân viên **When** gọi API tạo lời mời **Then** trả 403.
- Chấp nhận lời mời ghi nhật ký `membership.created`.

## API gợi ý
`POST /api/v1/households/{householdId}/invitations` · `GET /api/v1/invitations/{code}` · `POST /api/v1/invitations/{code}/accept` · `POST /api/v1/households/{householdId}/memberships` · `GET /api/v1/households/{householdId}/members`
