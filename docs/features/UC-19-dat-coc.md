# UC-19 Nhận cọc giữ hàng, chuyển cọc thành đơn

- **Milestone:** M6 · **Context:** sales (+ inventory, customer, cashbook) · **Capability:** `serial.deposit` · **Tác nhân:** Nhân viên (nhận cọc), Quản lý, Chủ hộ
- **Quy tắc:** BR-11, BR-18, BR-26, BR-30 · **Quyền:** `deposit.create`, `deposit.cancel`

## Luồng chính: nhận cọc
1. Chọn serial `IN_STOCK`, chọn hoặc tạo khách (bắt buộc), nhập số tiền cọc (> 0, ≤ giá bán hiện tại), hình thức thanh toán.
2. Trong một transaction: tạo Deposit `ACTIVE`, serial → `RESERVED`, CashEntry `INCOME` loại `DEPOSIT_RECEIVED` (ghi hình thức, tài khoản). Không ghi doanh thu (BR-30).
3. In phiếu nhận cọc (dùng chung mẫu in UC-21).

## Luồng chính: chuyển cọc thành đơn
1. Mở phiếu cọc → "Bán". POS mở sẵn đơn có dòng serial đó và khách.
2. Thanh toán: hệ thống tự thêm khoản `DEPOSIT` bằng số tiền cọc; phần còn lại trả bằng hình thức khác.
3. Hoàn tất như UC-16/UC-18; phiếu cọc → `APPLIED`, `applied_order_id` = đơn.

## Luồng chính: hủy cọc (Quản lý, Chủ hộ)
1. Chọn hoàn tiền hoặc giữ tiền cọc, nhập lý do.
2. Serial → `IN_STOCK`, phiếu cọc → `CANCELLED`.
3. Hoàn tiền → CashEntry `EXPENSE` loại `DEPOSIT_REFUND`. Giữ tiền → không tạo phiếu thu chi; báo cáo tính vào "Thu nhập khác". Ghi nhật ký.

## Ngoại lệ
- Serial không `IN_STOCK` → lỗi `BR-11`.
- Khoản `DEPOSIT` của đơn không khớp phiếu cọc đang `ACTIVE` của đúng serial → lỗi `BR-18`.

## Tiêu chí chấp nhận
- **Given** IMEI A `IN_STOCK` **When** nhận cọc 2.000.000 tiền mặt **Then** A `RESERVED`, có CashEntry `DEPOSIT_RECEIVED` 2.000.000, doanh thu ngày không đổi.
- **Given** phiếu cọc 2.000.000, giá bán 21.000.000 **When** chuyển thành đơn, trả thêm 19.000.000 chuyển khoản **Then** đơn có Payment `DEPOSIT` 2.000.000 + `BANK_TRANSFER` 19.000.000, doanh thu ngày tăng 21.000.000, phiếu cọc `APPLIED`.
- **Given** A `RESERVED` **When** bán A trực tiếp ở POS không qua phiếu cọc **Then** lỗi `BR-11`.
- **Given** hủy cọc kiểu giữ tiền **Then** A `IN_STOCK`, không có CashEntry mới, báo cáo "Thu nhập khác" tăng 2.000.000 vào ngày hủy.
- **Given** đơn bán từ cọc bị hủy (UC-20) **Then** phiếu cọc quay lại `ACTIVE`, A quay lại `RESERVED`.

## API gợi ý
`POST /api/v1/stores/{storeId}/deposits` · `GET .../deposits?status=ACTIVE` · `POST .../deposits/{id}/cancel`
