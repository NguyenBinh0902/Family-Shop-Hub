# UC-04 Thu hồi quyền thành viên

- **Milestone:** M1 · **Context:** identity · **Tác nhân:** Chủ hộ
- **Quy tắc:** BR-26, BR-36 · **Quyền:** `member.manage`

## Luồng chính
1. Chủ hộ mở danh sách thành viên, chọn một vai trò (một membership) hoặc "Xóa khỏi hộ" (mọi membership).
2. Xác nhận. Membership chuyển `REVOKED`, ghi người thu hồi và thời điểm.
3. Từ request tiếp theo, người đó không còn quyền ở phạm vi bị thu hồi.

Chủ hộ cũng thu hồi được lời mời `PENDING`.

## Ngoại lệ
- Thu hồi Chủ hộ cuối cùng → lỗi `BR-36`.
- Thu hồi chính mình khi là Chủ hộ duy nhất → lỗi `BR-36`.

## Tiêu chí chấp nhận
- **Given** Nhân viên đang mở POS cửa hàng A **When** Chủ hộ thu hồi quyền của họ ở A **Then** request kế tiếp của nhân viên tới A trả 403.
- **Given** hộ có đúng một Chủ hộ **When** thu hồi Chủ hộ đó **Then** trả lỗi `BR-36`, không đổi gì.
- **Given** thành viên bị xóa khỏi hộ **When** xem lại đơn họ đã bán **Then** đơn vẫn hiển thị tên người bán.
- Mỗi lần thu hồi ghi nhật ký với giá trị trước và sau.

## API gợi ý
`DELETE /api/v1/households/{householdId}/memberships/{membershipId}` (đổi trạng thái, không xóa bản ghi) · `DELETE /api/v1/households/{householdId}/invitations/{id}`
