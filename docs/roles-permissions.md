# Vai trò và phân quyền

Quyền gán theo từng cửa hàng: một người có thể là Quản lý ở cửa hàng này và Nhân viên ở cửa hàng khác. Chủ hộ là vai trò cấp hộ, tự động có toàn quyền trên mọi cửa hàng của hộ. Một hộ có thể có nhiều Chủ hộ.

## Vai trò

| Vai trò | Mã | Cấp | Phạm vi |
| --- | --- | --- | --- |
| Chủ hộ | `OWNER` | Hộ (`store_id` rỗng) | Mọi cửa hàng, mọi chủ thể, góc nhìn Toàn hộ, cài đặt hộ, mời thành viên |
| Quản lý | `MANAGER` | Cửa hàng | Mọi nghiệp vụ trong cửa hàng được giao, xem báo cáo cửa hàng đó |
| Nhân viên | `STAFF` | Cửa hàng | Bán hàng, nhập kho, thu nợ; không xem lợi nhuận, không hủy đơn |

## Ma trận quyền

Mã quyền ở cột đầu được dùng trong code (`permission.Can(scope, perm)`).

| Mã quyền | Mô tả | OWNER | MANAGER | STAFF |
| --- | --- | --- | --- | --- |
| `store.manage` | Tạo, sửa, tạm ngưng, đóng cửa hàng | Có | Không | Không |
| `entity.manage` | Quản lý chủ thể kinh doanh, thông tin thuế | Có | Không | Không |
| `member.manage` | Mời thành viên, gán và thu hồi vai trò | Có | Không | Không |
| `report.household` | Báo cáo Toàn hộ, sổ doanh thu theo chủ thể | Có | Không | Không |
| `report.store` | Báo cáo doanh thu cửa hàng | Có | Có | Không |
| `cost.view` | Xem giá vốn, lợi nhuận, giá nhập | Có | Có | Không |
| `catalog.manage` | Sản phẩm, giá bán, tùy chọn, công thức | Có | Có | Không |
| `sale.create` | Tạo đơn, thu tiền, in hóa đơn | Có | Có | Có |
| `sale.discount` | Giảm giá trên đơn | Có | Có | Đến mức trần của cửa hàng (mặc định 0%) |
| `sale.cancel` | Hủy đơn đã hoàn tất | Có | Có | Không |
| `deposit.create` | Nhận cọc | Có | Có | Có |
| `deposit.cancel` | Hủy cọc (hoàn hoặc giữ tiền) | Có | Có | Không |
| `stock.receive` | Nhập kho | Có | Có | Có |
| `stock.adjust` | Kiểm kê, điều chỉnh tồn | Có | Có | Không |
| `stock.view` | Xem tồn kho của cửa hàng mình | Có | Có | Có |
| `cash.entry` | Ghi phiếu thu chi | Có | Có | Không |
| `debt.collect` | Thu nợ | Có | Có | Có |
| `customer.view` | Xem thông tin liên hệ khách trong hộ | Có | Có | Có |
| `customer.history.store` | Xem lịch sử, công nợ của khách tại cửa hàng mình | Có | Có | Có |
| `customer.history.household` | Xem lịch sử, công nợ của khách ở mọi cửa hàng | Có | Không | Không |

Nhân viên nhập kho thấy được giá nhập mình đang nhập, nhưng không xem được giá vốn bình quân hay lợi nhuận.

## Quy tắc

- Mỗi request nghiệp vụ mang `storeId` trong URL. Server kiểm tra membership còn hiệu lực và quyền; không tin dữ liệu ngữ cảnh từ client.
- Quản lý và Nhân viên chỉ xem được dữ liệu (kể cả tồn kho) của cửa hàng mình được giao. Chỉ Chủ hộ xem được mọi cửa hàng.
- Quyền phụ thuộc capability: cửa hàng không có `serial.deposit` thì không tồn tại thao tác nhận cọc, dù vai trò có quyền.
- Cửa hàng `SUSPENDED` hoặc `CLOSED`: chỉ cho xem, mọi quyền ghi đều bị từ chối.
- Thu hồi membership có hiệu lực từ request tiếp theo.
- Hành động nhạy cảm ghi nhật ký (BR-26).
