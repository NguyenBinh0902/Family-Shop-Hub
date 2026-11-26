# Quyết định và câu hỏi mở

## Đã chốt

| Câu hỏi | Quyết định |
| --- | --- |
| Cách đăng nhập | Số điện thoại + mật khẩu; OTP khi quên mật khẩu |
| Nhân viên giảm giá | Cấu hình mức trần theo cửa hàng, mặc định 0% |
| Nguyên liệu theo lô, hạn dùng | Không có trong MVP |
| Trả hàng một phần | MVP chỉ hủy cả đơn; bút toán đảo ghi theo dòng (BR-15) |
| Thu mua máy cũ | Có, nhưng sau MVP |
| In hóa đơn | In từ trình duyệt, khổ 58 mm và 80 mm |
| Gói dịch vụ | Chưa thu phí; giới hạn 5 cửa hàng, 20 thành viên mỗi hộ |
| Nhân viên xem tồn cửa hàng khác | Không |
| Số repo | Một repo cho cả backend và frontend |

## CHƯA CHỐT – không implement

- Cách gom doanh thu nhiều địa điểm của cùng một chủ thể khi kê khai, và mẫu tờ khai thuế. Cần hỏi kế toán thuế. Ảnh hưởng: UC-29 (chỉ phần tạo tờ khai), TAX-06 (giá trị ngưỡng).
- Nhà cung cấp SMS OTP hoặc Zalo ZNS. Ẩn sau port `OTPSender`; khi phát triển dùng adapter ghi mã ra log.

## Mặc định BA đề xuất – cần chủ repo xem lại

Các mục dưới đây đã được viết vào tài liệu và **được áp dụng** trừ khi chủ repo thay đổi. Chúng là chi tiết cần có để developer không phải tự đoán.

| Mục | Mặc định | Ở đâu |
| --- | --- | --- |
| Giá vốn khi tồn âm hoặc chưa có giá vốn | BR-28 | rules.md |
| Ngày ghi nhận bút toán đảo | Ngày hủy (BR-29) | rules.md |
| Tiền cọc, thu nợ, trả nhà cung cấp | Không phải doanh thu / chi phí vận hành (BR-30, BR-31, BR-32) | rules.md |
| Làm tròn tiền | Tới đồng, nửa lên; số âm làm tròn ra xa số 0 (chủ repo chốt 25/11/2026) (BR-33) | rules.md |
| Đăng ký không xác minh SĐT bằng OTP | Chấp nhận trong MVP vì chưa có nhà cung cấp SMS | UC-01 |
| Mời thành viên | Bằng đường link có mã, Chủ hộ tự gửi qua Zalo; hạn 7 ngày, dùng một lần | UC-03 |
| Phiên đăng nhập | Hết hạn sau 30 ngày không dùng | UC-02 |
| Phiếu nhập sai | Không sửa, không hủy; điều chỉnh bằng kiểm kê | UC-12 |
| Nhập hàng | Coi như trả đủ tiền ngay; công nợ nhà cung cấp ngoài phạm vi | UC-12 |
| Hủy cọc | Quản lý chọn hoàn tiền hoặc giữ tiền cọc | UC-19 |
| Topping trừ nguyên liệu | Tùy chọn có thể mang danh sách nguyên liệu riêng | UC-10 |
| Số đơn hiển thị | `yyMMdd-NNN`, đếm lại mỗi ngày theo từng cửa hàng | domain/sales.md |
| Giảm giá | Chỉ ở cấp đơn (số tiền hoặc %), không giảm theo dòng | UC-16 |
