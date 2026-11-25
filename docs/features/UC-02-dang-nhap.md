# UC-02 Đăng nhập, đăng xuất, quên mật khẩu

- **Milestone:** M1 · **Context:** identity · **Tác nhân:** Mọi người dùng
- **Quy tắc:** NFR-06, NFR-12, NFR-13

## Luồng chính: đăng nhập
1. Nhập số điện thoại, mật khẩu.
2. Hệ thống kiểm tra, tạo Session (hết hạn 30 ngày kể từ lần dùng cuối), đặt cookie.
3. Điều hướng: thuộc nhiều hộ → chọn hộ; là Chủ hộ → trang Toàn hộ; chỉ thuộc một cửa hàng → thẳng vào cửa hàng đó; nhiều cửa hàng → chọn cửa hàng.

## Luồng phụ: đăng xuất
Thu hồi Session hiện tại, xóa cookie.

## Luồng phụ: quên mật khẩu
1. Nhập số điện thoại → hệ thống gửi OTP 6 chữ số qua port `OTPSender` (hạn 5 phút).
2. Nhập OTP và mật khẩu mới → đổi mật khẩu, thu hồi **mọi** phiên của người đó, đăng nhập phiên mới.

## Ngoại lệ
- Sai số điện thoại hoặc mật khẩu → một thông báo chung "Số điện thoại hoặc mật khẩu không đúng" (không tiết lộ số có tồn tại hay không).
- Sai 5 lần trong 15 phút → khóa đăng nhập số đó 15 phút (NFR-12).
- Yêu cầu OTP cho số không tồn tại → vẫn trả thành công (không tiết lộ), không gửi gì.
- Nhập sai OTP 5 lần → mã bị vô hiệu, phải yêu cầu mã mới. Tối đa 3 lần yêu cầu mã trong 1 giờ.

## Tiêu chí chấp nhận
- **Given** tài khoản hợp lệ **When** đăng nhập đúng **Then** có cookie phiên, `GET /api/v1/me` trả 200.
- **Given** đã đăng xuất **When** dùng lại token cũ **Then** trả 401.
- **Given** 5 lần sai liên tiếp **When** đăng nhập lần 6 dù đúng mật khẩu **Then** trả 429 cho tới hết 15 phút.
- **Given** người dùng có 2 phiên đang mở **When** đổi mật khẩu qua OTP **Then** cả 2 phiên cũ trả 401.
- **Given** Session không dùng 30 ngày **When** gọi API **Then** trả 401.
- Mã OTP không xuất hiện trong log ở môi trường production; adapter phát triển được phép ghi mã ra log.

## API gợi ý
`POST /api/v1/auth/login` · `POST /api/v1/auth/logout` · `POST /api/v1/auth/password-reset/request` · `POST /api/v1/auth/password-reset/confirm`
