# Tổng quan sản phẩm

FamilyShop Hub là ứng dụng web SaaS cho phép một hộ kinh doanh quản lý nhiều cửa hàng thuộc nhiều ngành trong một tài khoản. Loại hình kinh doanh chọn khi tạo cửa hàng quyết định chức năng, trường dữ liệu và giao diện của cửa hàng đó.

## Vấn đề

- Mỗi cửa hàng dùng một công cụ khác nhau (sổ tay, Excel, Zalo), dữ liệu rời rạc.
- Chủ hộ không nhìn được tổng doanh thu, lợi nhuận, tồn kho của cả gia đình.
- Phần mềm hiện có chỉ phục vụ một ngành, hoặc quá phức tạp và đắt với hộ nhỏ.
- Từ 01/01/2026 thuế khoán bị bãi bỏ; hộ kinh doanh tự kê khai theo doanh thu thực tế, nên cần sổ sách bán hàng đầy đủ và tin cậy.

## Mục tiêu

1. Một đăng nhập quản lý nhiều cửa hàng, có bảng điều khiển tổng hợp toàn hộ.
2. Chức năng gắn theo loại hình kinh doanh; thêm ngành mới bằng cấu hình, không viết lại hệ thống.
3. Người không rành công nghệ dùng được: tiếng Việt, tối ưu điện thoại và máy tính bảng.
4. Dữ liệu bán hàng đủ tin cậy để làm căn cứ kê khai thuế.

## Người dùng

| Đối tượng | Mô tả | Thiết bị chính |
| --- | --- | --- |
| Chủ hộ | Đứng đầu hộ, sở hữu 1 đến vài cửa hàng, xem báo cáo tổng | Điện thoại |
| Quản lý cửa hàng | Thành viên gia đình hoặc người được giao điều hành một cửa hàng | Điện thoại, máy tính bảng |
| Nhân viên | Đứng quầy, bán hàng, nhập kho | Máy tính bảng, điện thoại tại quầy |

## Phạm vi MVP

- Đăng ký hộ, đăng nhập, phân quyền theo từng cửa hàng.
- Nhiều cửa hàng, chuyển giữa các cửa hàng, góc nhìn Toàn hộ.
- Hai nhóm loại hình: F&B và Hàng có serial (xem `store-types.md`).
- Sản phẩm, tồn kho, nhập hàng, bán tại quầy (POS), thu chi, khách hàng và công nợ.
- Báo cáo doanh thu, lợi nhuận theo cửa hàng, theo chủ thể kinh doanh và toàn hộ.
- Sổ doanh thu bất biến, tách hình thức thanh toán, sẵn điểm tích hợp hóa đơn điện tử.

## Ngoài phạm vi MVP

Không implement các mục sau, kể cả khi thấy "tiện thì làm luôn":

- Thanh toán online, trang bán hàng công khai, ứng dụng di động native.
- Tích hợp hóa đơn điện tử và kê khai thuế thực tế (chỉ chừa port).
- Nhóm loại hình Dịch vụ & sửa chữa (C) và Bán lẻ tổng hợp (D); loại hình kết hợp.
- Chuyển kho giữa cửa hàng, POS offline đầy đủ, đa ngôn ngữ, thu phí gói dịch vụ.
- Quản lý nguyên liệu theo lô và hạn sử dụng.
- Thu mua hàng cũ từ khách (phiếu thu mua), trả hàng một phần, công nợ nhà cung cấp.
- Kiểm kê hàng serial.
