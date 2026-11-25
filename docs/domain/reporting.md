# Context: Reporting (`internal/reporting`)

Chỉ đọc. Gộp số liệu từ Sales, Inventory, Cashbook theo ba phạm vi: cửa hàng, chủ thể kinh doanh, toàn hộ (BR-24). Không ghi dữ liệu nghiệp vụ; được phép đọc thẳng bảng của context khác bằng query SQL riêng (sqlc) vì đây là read model.

## Định nghĩa chỉ số

Khoảng thời gian tính theo ngày nghiệp vụ `Asia/Ho_Chi_Minh` (BR-25), từ 00:00 ngày đầu đến hết ngày cuối.

| Chỉ số | Công thức |
| --- | --- |
| Doanh thu | Σ `total` đơn `COMPLETED` có `completed_at` trong kỳ − Σ `total` Reversal có `created_at` trong kỳ (BR-29) |
| Số đơn | Số đơn hoàn tất trong kỳ (đơn bị hủy vẫn đếm; số đơn hủy hiển thị riêng) |
| Giá vốn hàng bán | Σ `unit_cost × quantity` của dòng đơn hoàn tất trong kỳ − phần tương ứng của Reversal trong kỳ |
| Lợi nhuận gộp | Doanh thu − Giá vốn hàng bán |
| Chi phí vận hành | Σ CashEntry `EXPENSE` thuộc danh mục tính chi phí vận hành, `occurred_on` trong kỳ (kể cả phiếu điều chỉnh) |
| Thu nhập khác | Σ CashEntry `OTHER_INCOME` trong kỳ + Σ `amount` phiếu cọc hủy kiểu `FORFEIT` có ngày hủy trong kỳ |
| Lợi nhuận ròng | Lợi nhuận gộp − Chi phí vận hành + Thu nhập khác |
| Doanh thu theo hình thức | Σ Payment của đơn hoàn tất theo `method` (trừ phần đơn đã hủy trong kỳ) |
| Tiền vào / ra | Payment `CASH`, `BANK_TRANSFER`, `EWALLET` của đơn + CashEntry `INCOME`/`EXPENSE` (dòng tiền thực) |
| Giá trị tồn kho | Σ `on_hand × avg_cost` (chỉ `on_hand > 0`) + Σ `unit_cost` serial `IN_STOCK`/`RESERVED` |
| Cảnh báo thiếu giá vốn | Số dòng `cost_missing` trong kỳ (BR-28) |

Người không có `cost.view` không nhận các chỉ số giá vốn, lợi nhuận, giá trị tồn.

## Phạm vi

- Cửa hàng: `store_id`.
- Chủ thể: mọi cửa hàng có `business_entity_id`, kể cả cửa hàng đã đóng.
- Toàn hộ: mọi cửa hàng của hộ, kể cả đã đóng; có bảng tách theo cửa hàng và theo chủ thể.

## Hiệu năng

MVP query trực tiếp với index `(store_id, completed_at)`, `(store_id, created_at)`, `(store_id, occurred_on)`. Chỉ thêm bảng tổng hợp theo ngày khi đo được chậm, và bảng đó phải tính lại được từ dữ liệu gốc (vẫn tuân BR-24).

## Use case

UC-27, UC-28, UC-29, UC-30.
