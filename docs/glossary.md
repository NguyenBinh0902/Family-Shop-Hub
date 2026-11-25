# Thuật ngữ

UI và tài liệu dùng tên tiếng Việt; code dùng tên ở cột thứ hai. Không đặt tên khác cho cùng một khái niệm.

| Thuật ngữ | Tên trong code | Định nghĩa |
| --- | --- | --- |
| Hộ | `Household` | Tenant của hệ thống. Một gia đình hoặc một chủ kinh doanh; sở hữu chủ thể kinh doanh và cửa hàng. |
| Chủ thể kinh doanh | `BusinessEntity` | Hộ kinh doanh đăng ký với nhà nước, có mã số thuế. Một hộ có thể có nhiều chủ thể. |
| Cửa hàng | `Store` | Điểm kinh doanh cụ thể. Thuộc đúng một chủ thể, có đúng một loại hình. |
| Loại hình kinh doanh | `StoreType` | Gói cấu hình (trà sữa, điện thoại...) gồm danh sách capability, trường động, dữ liệu mẫu. Có version. |
| Capability | `Capability` | Một khối chức năng nghiệp vụ (tồn kho, serial, topping...). Người dùng chỉ thấy qua loại hình. |
| Người dùng | `User` | Một tài khoản đăng nhập, định danh bằng số điện thoại. Có thể thuộc nhiều hộ. |
| Thành viên | `Membership` | Quan hệ người dùng với hộ (Chủ hộ) hoặc với cửa hàng (Quản lý, Nhân viên). |
| Vai trò | `Role` | `OWNER` (Chủ hộ), `MANAGER` (Quản lý), `STAFF` (Nhân viên). |
| Lời mời | `Invitation` | Mã mời một người vào hộ với vai trò cho trước, gửi qua đường link. |
| Phạm vi | `tenant.Scope` | Hộ, cửa hàng, người dùng, vai trò của request hiện tại. |
| Sản phẩm | `Product` | Hàng hóa, món hoặc nguyên liệu, thuộc một cửa hàng. |
| Kiểu theo dõi tồn | `StockTracking` | `NONE` (không theo dõi), `QUANTITY` (theo số lượng), `SERIAL` (theo từng món). |
| Nguyên liệu | (Product, `IsSellable=false`) | Sản phẩm theo dõi tồn theo số lượng, không bán trực tiếp, dùng trong công thức. |
| Nhóm tùy chọn | `OptionGroup` | Size, Topping, Mức đường... gồm các `Option` có giá cộng thêm. |
| Công thức | `Recipe` | Danh sách nguyên liệu và định lượng trừ kho khi bán một món. |
| Đơn vị serial | `SerialUnit` | Một món hàng cụ thể có IMEI/số khung, giá vốn riêng. |
| Tồn kho | `StockLevel` | Số lượng tồn và giá vốn bình quân của một sản phẩm theo số lượng. |
| Biến động kho | `StockMovement` | Bản ghi chỉ-thêm của mọi thay đổi tồn kho. |
| Nhà cung cấp | `Supplier` | Nơi nhập hàng, dùng chung trong hộ. |
| Phiếu nhập | `PurchaseReceipt` | Ghi nhận hàng nhập kho. Đã ghi thì bất biến. |
| Phiếu kiểm kê | `StockCount` | Ghi số đếm thực tế và chênh lệch so với tồn trên hệ thống. |
| Đơn bán | `SalesOrder` | Giao dịch bán. Trạng thái `OPEN` (đơn mở theo bàn), `COMPLETED`, `CANCELLED`. |
| Dòng đơn | `OrderLine` | Một dòng trong đơn, lưu snapshot tên, tùy chọn, giá bán, giá vốn. |
| Thanh toán | `Payment` | Một khoản trả cho đơn: tiền mặt, chuyển khoản, ví, nợ, cọc. |
| Bút toán đảo | `Reversal` | Bản ghi đối ứng hủy một đơn đã hoàn tất, thay cho việc xóa. |
| Phiếu cọc | `Deposit` | Tiền khách đặt trước để giữ một đơn vị serial. |
| Bàn | `DiningTable` | Bàn trong quán ăn, gắn với tối đa một đơn mở. |
| Khách hàng | `Customer` | Người mua, dùng chung trong hộ, định danh bằng số điện thoại. |
| Công nợ | `Receivable` | Số tiền khách còn nợ một cửa hàng, gắn với một đơn bán nợ. |
| Phiếu thu chi | `CashEntry` | Khoản thu hoặc chi ngoài bán hàng (lương, mặt bằng, thu nợ, trả nhà cung cấp...). |
| Giá vốn bình quân | weighted average cost | (tồn cũ × giá vốn cũ + số nhập × giá nhập) / tổng tồn sau nhập. |
| Tiền | `kernel.Money` | `int64`, đơn vị đồng. |
| Số lượng | `kernel.Quantity` | `int64`, đơn vị phần nghìn (1 ly = 1000, 0,5 kg = 500). |
| Ngày nghiệp vụ | business date | Ngày tính theo múi giờ `Asia/Ho_Chi_Minh`. |
| Toàn hộ | household view | Góc nhìn báo cáo gộp mọi cửa hàng, chỉ dành cho Chủ hộ. |
| Nhật ký | `AuditLog` | Bản ghi chỉ-thêm các hành động nhạy cảm. |
