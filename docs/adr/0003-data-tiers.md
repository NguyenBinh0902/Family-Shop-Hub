# ADR-0003: Phân cấp dữ liệu hộ, chủ thể kinh doanh, cửa hàng

- **Trạng thái:** Đã chấp nhận
- **Ngày:** 2026-11-25

## Bối cảnh

Một gia đình có thể có nhiều hộ kinh doanh đứng tên khác nhau, mỗi hộ có một hoặc nhiều cửa hàng. Chủ hộ cần báo cáo tổng; nghĩa vụ thuế tính theo từng hộ kinh doanh đăng ký.

## Quyết định

- **Hộ (Household)** là tenant. Mọi bảng nghiệp vụ có `household_id`; RLS của PostgreSQL chặn theo `app.household_id` làm lớp thứ hai sau repository.
- **Chủ thể kinh doanh (BusinessEntity)** nằm dưới hộ, giữ MST và tài khoản; mỗi cửa hàng thuộc một chủ thể.
- **Cửa hàng (Store)** là nơi phát sinh và lưu mọi giao dịch (đơn, kho, thu chi, công nợ).
- Dữ liệu dùng chung cấp hộ: thành viên và quyền, chủ thể, khách hàng, nhà cung cấp.
- Báo cáo cấp chủ thể và cấp hộ là **phép gộp** từ dữ liệu cửa hàng, không lưu bản thứ hai.
- Người dùng (User) nằm ngoài tenant, thuộc nhiều hộ qua Membership.

## Hệ quả

- Một nguồn sự thật cho doanh thu, phù hợp yêu cầu sổ sách thuế.
- Khách hàng không bị lưu trùng giữa các cửa hàng; lịch sử và công nợ vẫn tách theo cửa hàng.
- Báo cáo tổng chạy query gộp; nếu chậm thì thêm bảng tổng hợp tính lại được, không thay đổi nguồn sự thật.
