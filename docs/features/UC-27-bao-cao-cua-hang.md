# UC-27 Báo cáo doanh thu, lợi nhuận cửa hàng

- **Milestone:** M8 (bản rút gọn "doanh thu hôm nay" có từ M4) · **Context:** reporting · **Tác nhân:** Quản lý, Chủ hộ
- **Quy tắc:** BR-10, BR-24, BR-25, BR-28, BR-29 · **Quyền:** `report.store`, `cost.view`

## Luồng chính
1. Chọn kỳ: Hôm nay, Hôm qua, 7 ngày, Tháng này, Tháng trước, Tùy chọn.
2. Hiển thị các chỉ số theo định nghĩa ở `domain/reporting.md`: doanh thu, số đơn, số đơn hủy, giá vốn, lợi nhuận gộp, chi phí vận hành, thu nhập khác, lợi nhuận ròng.
3. Biểu đồ doanh thu theo ngày trong kỳ; bảng doanh thu theo hình thức thanh toán; bảng chi phí theo danh mục.
4. Cảnh báo nếu có dòng thiếu giá vốn trong kỳ, kèm danh sách sản phẩm.

## Tiêu chí chấp nhận
- **Given** ngày 01/12: đơn 100.000 (giá vốn 40.000) và đơn 50.000 (giá vốn 20.000), chi mặt bằng 30.000 **Then** doanh thu 150.000, giá vốn 60.000, lợi nhuận gộp 90.000, lợi nhuận ròng 60.000.
- **Given** đơn 01/12 bị hủy ngày 02/12 **Then** kỳ 01/12 không đổi, kỳ 02/12 doanh thu −50.000 (BR-29), kỳ 01–02/12 cộng gộp đúng.
- **Given** trả nhà cung cấp 1.000.000, thu nợ 200.000, nhận cọc 2.000.000 trong kỳ **Then** doanh thu và chi phí vận hành không đổi; phần "Tiền vào/ra" phản ánh các khoản này.
- **Given** đơn lúc 00:30 ngày 02/12 giờ Việt Nam (17:30 UTC ngày 01/12) **Then** tính vào ngày 02/12.
- Số liệu báo cáo khớp với tổng tính tay từ dữ liệu gốc trong test (dữ liệu mẫu cố định).

## API gợi ý
`GET /api/v1/stores/{storeId}/reports/summary?from=&to=` · `.../reports/daily?from=&to=`
