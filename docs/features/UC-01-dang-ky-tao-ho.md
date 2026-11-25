# UC-01 Đăng ký tài khoản và tạo hộ

- **Milestone:** M1 · **Context:** identity · **Tác nhân:** Khách truy cập
- **Quy tắc:** BR-27, BR-36, NFR-06

## Mục tiêu
Người mới tạo tài khoản và hộ kinh doanh của mình trong một bước, sau đó đăng nhập luôn.

## Luồng chính
1. Người dùng nhập: họ tên, số điện thoại, mật khẩu (≥ 8 ký tự), tên hộ.
2. Hệ thống chuẩn hóa số điện thoại, kiểm tra chưa có tài khoản.
3. Trong một transaction: tạo User, Household (gói `FREE`, giới hạn 5 cửa hàng, 20 thành viên), Membership `OWNER`.
4. Tạo Session, đặt cookie, chuyển tới màn hình tạo cửa hàng đầu tiên (UC-06).

## Ngoại lệ
- Số điện thoại sai định dạng → lỗi trường `phone`.
- Số điện thoại đã có tài khoản → lỗi "Số điện thoại đã được đăng ký", gợi ý đăng nhập.
- Mật khẩu < 8 ký tự → lỗi trường `password`.

## Tiêu chí chấp nhận
- **Given** số 0901234567 chưa đăng ký **When** đăng ký với số "+84 901 234 567" **Then** User được tạo với phone `0901234567`, có một Household và một Membership `OWNER` (`store_id` rỗng), và phản hồi có cookie phiên `HttpOnly`.
- **Given** số 0901234567 đã đăng ký **When** đăng ký lại số đó **Then** trả 409, không tạo Household mới.
- **Given** đăng ký thành công **When** gọi `GET /api/v1/me` **Then** trả thông tin người dùng và danh sách hộ gồm hộ vừa tạo với vai trò `OWNER`.
- Mật khẩu không bao giờ xuất hiện trong log hay phản hồi.

## API gợi ý
`POST /api/v1/auth/register` · `GET /api/v1/me`

## Ghi chú
- MVP không xác minh số điện thoại bằng OTP khi đăng ký (xem `open-questions.md`).
- Đăng ký qua lời mời (UC-03) dùng cùng form nhưng không tạo hộ.
