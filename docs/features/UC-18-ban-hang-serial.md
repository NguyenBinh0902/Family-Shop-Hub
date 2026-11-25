# UC-18 Bán hàng serial

- **Milestone:** M6 · **Context:** sales (+ inventory) · **Capability:** `inv.serial` (bảo hành cần `serial.warranty`) · **Tác nhân:** Nhân viên, Quản lý, Chủ hộ
- **Quy tắc:** BR-10, BR-11, BR-12, BR-13 · **Quyền:** `sale.create`

## Luồng chính
1. Ở POS, chọn sản phẩm `SERIAL` → hệ thống hiển thị danh sách serial `IN_STOCK` của sản phẩm đó; hoặc quét/gõ serial trực tiếp ở ô tìm kiếm.
2. Mỗi serial là một dòng, số lượng 1. Đơn có thể kèm dòng số lượng (phụ kiện).
3. Thanh toán như UC-16.
4. Khi hoàn tất: serial → `SOLD`; dòng đơn lưu `serial_unit_id`, `serial`, `unit_cost` = giá vốn của chiếc đó; nếu có `serial.warranty` và sản phẩm có `warranty_months` thì `warranty_until` = ngày bán + số tháng.

## Ngoại lệ
- Serial không `IN_STOCK` (đã bán, đang được cọc bởi phiếu cọc khác) → lỗi `BR-11`. Serial `RESERVED` chỉ bán được qua chuyển cọc thành đơn (UC-19).
- Cùng serial xuất hiện hai lần trong đơn → lỗi.

## Tiêu chí chấp nhận
- **Given** IMEI A giá vốn 18.000.000, sản phẩm bảo hành 12 tháng **When** bán ngày 15/01/2027 giá 21.000.000 **Then** dòng đơn có `unit_cost` 18.000.000, serial A `SOLD`, `warranty_until` 15/01/2028.
- **Given** hai nhân viên cùng bán IMEI A một lúc **Then** chỉ một đơn thành công, đơn kia lỗi `BR-11` (khóa hàng `SELECT ... FOR UPDATE` hoặc cập nhật có điều kiện trạng thái).
- **Given** bán ngày 31/01/2027, bảo hành 1 tháng **Then** `warranty_until` là 28/02/2027. Lưu ý `time.AddDate(0, 1, 0)` của Go cho ra 03/03/2027; phải kẹp về ngày cuối của tháng đích. Viết hàm riêng trong `kernel` kèm test.

## API gợi ý
Dùng chung `POST /api/v1/stores/{storeId}/orders` (dòng có `serialUnitId`) · `GET .../products/{id}/serial-units?status=IN_STOCK`
