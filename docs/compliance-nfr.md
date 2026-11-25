# Tuân thủ thuế và yêu cầu phi chức năng

## Bối cảnh thuế

Từ 01/01/2026 hộ kinh doanh không còn nộp thuế khoán mà tự kê khai theo doanh thu thực tế; hộ có doanh thu từ 1 tỷ đồng/năm trở lên phải dùng hóa đơn điện tử khởi tạo từ máy tính tiền ([MISA eShop](https://www.misaeshop.vn/?p=30699), [VOH](https://voh.com.vn/kinh-te/tu-112026-chinh-thuc-bo-thue-khoan-chuyen-sang-ke-khai-thue-642709.html)). Hộ phải thông báo cho cơ quan thuế các tài khoản ngân hàng, ví điện tử dùng cho kinh doanh ([MISA meInvoice](https://www.meinvoice.vn/?p=47137)).

MVP chưa tích hợp kê khai hay hóa đơn điện tử, nhưng dữ liệu phải đủ làm căn cứ khi tích hợp.

## Yêu cầu dữ liệu phục vụ thuế

| Mã | Yêu cầu |
| --- | --- |
| TAX-01 | Sổ doanh thu bất biến (BR-15, BR-23, BR-29); xuất được theo chủ thể kinh doanh và khoảng thời gian (UC-29). |
| TAX-02 | Mỗi khoản thanh toán ghi rõ hình thức: tiền mặt, chuyển khoản (tài khoản nào), ví (ví nào), nợ, cọc. |
| TAX-03 | Chủ thể kinh doanh lưu MST, danh sách tài khoản ngân hàng và ví đã đăng ký. |
| TAX-04 | Có port `InvoiceIssuer` trong `internal/sales/app`; MVP dùng adapter không làm gì (no-op), gọi sau khi đơn hoàn tất. |
| TAX-05 | Đơn bán có trường `einvoice_ref` (rỗng trong MVP). |
| TAX-06 | Báo cáo doanh thu lũy kế năm theo chủ thể, cảnh báo khi đạt 80% ngưỡng. Ngưỡng là cấu hình, mặc định 1.000.000.000 đồng; phải được kế toán xác nhận trước khi phát hành. |

**CHƯA CHỐT – không implement:** cách gom doanh thu nhiều địa điểm của cùng một chủ thể khi kê khai, và mẫu tờ khai. UC-29 chỉ liệt kê và xuất dữ liệu, không tạo tờ khai.

## Yêu cầu phi chức năng

| Mã | Nhóm | Yêu cầu |
| --- | --- | --- |
| NFR-01 | Giao diện | Mobile-first; dùng tốt từ độ rộng 360 px và trên máy tính bảng. |
| NFR-02 | Giao diện | Toàn bộ tiếng Việt, từ ngữ của người kinh doanh nhỏ. |
| NFR-03 | Hiệu năng | Tạo đơn tại quầy phản hồi dưới 1 giây (p95 phía server dưới 300 ms) với mạng 4G thông thường. |
| NFR-04 | Tin cậy | Gửi lại đơn khi mạng chập chờn không tạo đơn trùng (BR-17). POS offline đầy đủ để sau MVP. |
| NFR-05 | Bảo mật | Cách ly dữ liệu giữa các hộ ở tầng repository, cộng Row-Level Security của PostgreSQL làm lớp chặn thứ hai. |
| NFR-06 | Bảo mật | Mật khẩu băm argon2id; phiên đăng nhập lưu DB, thu hồi được; cookie `HttpOnly`, `Secure`, `SameSite=Strict`. |
| NFR-07 | Dữ liệu | Sao lưu hằng ngày; dữ liệu giao dịch không bị xóa vật lý. |
| NFR-08 | Thời gian | Lưu UTC (`timestamptz`), hiển thị và chốt ngày theo `Asia/Ho_Chi_Minh`. |
| NFR-09 | Tiền tệ | Tiền là số nguyên đồng (`int64`), không dùng số thực. |
| NFR-10 | Mở rộng | Thêm loại hình mới chỉ bằng cấu hình nếu dùng capability có sẵn. |
| NFR-11 | Kiểm toán | Nhật ký hành động nhạy cảm (BR-26) giữ tối thiểu 5 năm. |
| NFR-12 | Bảo mật | Giới hạn đăng nhập sai: 5 lần trong 15 phút cho mỗi số điện thoại thì khóa tạm 15 phút. |
| NFR-13 | Riêng tư | Log không chứa mật khẩu, token, mã OTP; số điện thoại trong log che 4 số giữa. |
