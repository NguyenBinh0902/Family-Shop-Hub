# Quy tắc nghiệp vụ

Mỗi quy tắc có mã cố định. Code, test, commit dẫn chiếu mã này. Không đánh số lại; quy tắc bỏ đi thì ghi "Đã bỏ", không xóa dòng.

Vi phạm quy tắc phải trả lỗi nghiệp vụ có mã BR và thông điệp tiếng Việt; không bỏ qua lặng lẽ.

BR-28 đến BR-37 là quy tắc BA bổ sung khi chi tiết hóa cho developer. Chúng là mặc định được áp dụng; chủ repo có thể điều chỉnh (xem `open-questions.md`).

## Tenant

| Mã | Quy tắc |
| --- | --- |
| BR-01 | Mọi dữ liệu nghiệp vụ mang `household_id`; dữ liệu cấp cửa hàng mang thêm `store_id`. Không truy vấn nào được trả dữ liệu của hộ khác. |
| BR-02 | Một người dùng có thể thuộc nhiều hộ; ngữ cảnh hộ và cửa hàng xác định từ URL và được kiểm tra với membership. |
| BR-36 | Mỗi hộ luôn có ít nhất một Chủ hộ còn hiệu lực. Không thu hồi được Chủ hộ cuối cùng. |

## Cửa hàng

| Mã | Quy tắc |
| --- | --- |
| BR-03 | Mỗi cửa hàng thuộc đúng một chủ thể kinh doanh và đúng một loại hình (kèm version). |
| BR-04 | Không đổi loại hình hoặc chủ thể của cửa hàng đã có giao dịch. |
| BR-05 | Đóng cửa hàng không xóa dữ liệu; cửa hàng đã đóng chỉ xem được, không phát sinh giao dịch mới. Không mở lại được. |

## Sản phẩm

| Mã | Quy tắc |
| --- | --- |
| BR-06 | Sản phẩm đã có giao dịch không được xóa cứng, chỉ được ngừng bán. |
| BR-07 | Kiểu theo dõi tồn (`NONE`, `QUANTITY`, `SERIAL`) chọn lúc tạo sản phẩm và không đổi được. Chỉ chọn được kiểu mà loại hình của cửa hàng có capability. |

## Giá vốn

| Mã | Quy tắc |
| --- | --- |
| BR-08 | Phương pháp giá vốn chọn lúc tạo cửa hàng, MVP chỉ có bình quân gia quyền; bị khóa sau giao dịch kho đầu tiên. |
| BR-09 | Giá vốn bình quân tính lại sau mỗi phiếu nhập. Hàng serial dùng giá vốn của chính món đó. |
| BR-10 | Mỗi dòng đơn bán lưu snapshot giá vốn tại thời điểm bán; báo cáo lợi nhuận dùng snapshot, không tính lại. |
| BR-28 | Khi nhập hàng mà tồn hiện tại ≤ 0, giá vốn bình quân mới bằng giá nhập của phiếu đó. Bán sản phẩm chưa từng có giá vốn thì snapshot giá vốn là 0 và dòng đơn được đánh dấu `cost_missing` để báo cáo cảnh báo. |

## Serial

| Mã | Quy tắc |
| --- | --- |
| BR-11 | IMEI/số khung là duy nhất trong hộ đối với các đơn vị đang `IN_STOCK` hoặc `RESERVED`. Trạng thái: `IN_STOCK`, `RESERVED`, `SOLD`, `RETURNED` (`RETURNED` dành cho sau MVP). |
| BR-12 | Bảo hành bắt đầu từ ngày bán; thời hạn lấy theo số tháng bảo hành của sản phẩm tại thời điểm bán. |

## Tồn kho

| Mã | Quy tắc |
| --- | --- |
| BR-13 | Cấu hình cửa hàng quyết định cho bán khi tồn không đủ hay chặn. Mặc định F&B cho phép, Serial chặn. Hàng serial luôn chặn. |
| BR-14 | Công thức trừ nguyên liệu theo đơn vị tính của nguyên liệu, tại thời điểm hoàn tất đơn. Tùy chọn có thể mang thêm nguyên liệu riêng (ví dụ topping). |
| BR-34 | Số lượng lưu dạng số nguyên phần nghìn (`kernel.Quantity`), không dùng số thực. |
| BR-37 | Mọi thay đổi tồn kho tạo một `StockMovement`; biến động kho chỉ thêm, không sửa, không xóa. |

## Đơn bán

| Mã | Quy tắc |
| --- | --- |
| BR-15 | Đơn hoàn tất là bất biến: không sửa, không xóa. Hủy bằng bút toán đảo có lý do và người thực hiện. MVP chỉ hủy cả đơn, nhưng bút toán đảo ghi theo từng dòng để sau này hỗ trợ trả hàng một phần. |
| BR-16 | Dòng đơn lưu snapshot tên sản phẩm, tùy chọn, giá bán, giảm giá. |
| BR-17 | Mã đơn (UUID v7) do client sinh; gửi lại cùng mã thì server trả về đơn đã có, không tạo trùng. |
| BR-18 | Tổng các khoản thanh toán của đơn phải bằng tổng tiền đơn (phần bán nợ và phần trừ cọc tính là một khoản). |
| BR-19 | Bán nợ bắt buộc gắn khách hàng. |
| BR-29 | Bút toán đảo ghi nhận vào ngày hủy; không sửa số liệu của ngày bán. |
| BR-30 | Tiền cọc không phải doanh thu; chỉ thành doanh thu khi được trừ vào một đơn bán hoàn tất. |
| BR-33 | Phép tính tiền có phần lẻ (giảm giá theo %, giá vốn bình quân, giá trị nguyên liệu) làm tròn tới đồng, nửa đồng làm tròn lên. |
| BR-35 | Đơn mở (`OPEN`, chỉ có ở cửa hàng có `fnb.table`) được sửa tự do; chỉ khi hoàn tất mới trừ kho, ghi doanh thu, và trở thành bất biến. |

## Khách hàng và công nợ

| Mã | Quy tắc |
| --- | --- |
| BR-20 | Khách hàng thuộc cấp hộ; số điện thoại (đã chuẩn hóa) là duy nhất trong hộ nếu có. Cho phép khách không có số điện thoại. |
| BR-21 | Lịch sử giao dịch và công nợ thuộc cấp cửa hàng. Màn hình cấp hộ chỉ hiển thị giá trị gộp. |
| BR-22 | Thu nợ ghi vào đúng cửa hàng đã cho nợ, phân bổ vào đơn nợ cũ nhất trước; không thu vượt số nợ. |
| BR-31 | Thu nợ không phải doanh thu (doanh thu đã ghi khi bán nợ); chỉ là tiền vào. |

## Thu chi

| Mã | Quy tắc |
| --- | --- |
| BR-23 | Phiếu thu chi đã ghi là bất biến; sửa bằng phiếu điều chỉnh liên kết với phiếu gốc. |
| BR-32 | Tiền trả nhà cung cấp khi nhập hàng là tiền ra, không tính vào chi phí vận hành (đã phản ánh qua giá vốn hàng bán). |

## Báo cáo

| Mã | Quy tắc |
| --- | --- |
| BR-24 | Báo cáo cấp hộ và cấp chủ thể là phép gộp từ dữ liệu cửa hàng, không lưu bản dữ liệu thứ hai. |
| BR-25 | Doanh thu ghi nhận theo thời điểm hoàn tất đơn, gom theo ngày nghiệp vụ (`Asia/Ho_Chi_Minh`). |

## Kiểm toán và gói dịch vụ

| Mã | Quy tắc |
| --- | --- |
| BR-26 | Hủy đơn, điều chỉnh tồn, đổi giá bán, hủy cọc, thay đổi quyền, sửa chủ thể kinh doanh đều ghi nhật ký: ai, lúc nào, cửa hàng, giá trị trước và sau, lý do (nếu có). |
| BR-27 | Số cửa hàng và số thành viên của hộ bị giới hạn theo gói; kiểm tra tập trung khi tạo mới. MVP chưa thu phí, giới hạn mặc định 5 cửa hàng (không tính cửa hàng đã đóng) và 20 thành viên (đếm người dùng khác nhau còn hiệu lực). |
