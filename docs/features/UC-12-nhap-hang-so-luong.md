# UC-12 Nhập hàng theo số lượng

- **Milestone:** M3 · **Context:** inventory (+ cashbook) · **Capability:** `inv.quantity` · **Tác nhân:** Nhân viên, Quản lý, Chủ hộ
- **Quy tắc:** BR-08, BR-09, BR-28, BR-32, BR-33, BR-34, BR-37 · **Quyền:** `stock.receive`

## Luồng chính
1. Chọn nhà cung cấp (tìm trong hộ, hoặc tạo nhanh) hoặc bỏ trống.
2. Thêm dòng: sản phẩm `QUANTITY`, số lượng, giá nhập đơn vị. Hệ thống tính thành tiền và tổng.
3. Chọn hình thức trả tiền nhà cung cấp (tiền mặt, chuyển khoản, ví).
4. Xác nhận. Trong một transaction:
   - tạo PurchaseReceipt và dòng;
   - mỗi dòng: tạo `StockMovement` `RECEIPT`, cập nhật `on_hand` và `avg_cost` theo công thức ở `domain/inventory.md`;
   - tạo CashEntry `EXPENSE` loại `PURCHASE_PAYMENT` bằng tổng phiếu (BR-32);
   - đặt `costing_locked_at` cho cửa hàng nếu chưa có (BR-08), đánh dấu `has_transactions`.

## Ngoại lệ
- Sản phẩm không phải `QUANTITY`, hoặc thuộc cửa hàng khác → lỗi.
- Số lượng ≤ 0 hoặc giá nhập < 0 → lỗi trường.
- Cửa hàng không `ACTIVE` → từ chối.
- Phiếu đã ghi không sửa, không hủy; sai thì điều chỉnh bằng kiểm kê (UC-14).

## Tiêu chí chấp nhận
- **Given** sữa tồn 10.000 ml giá vốn 30 đ/ml **When** nhập 5.000 ml giá 36 đ/ml **Then** tồn 15.000 ml, giá vốn 32 đ/ml.
- **Given** tồn hiện tại −200 ml (đã bán âm) **When** nhập 1.000 ml giá 35 đ/ml **Then** tồn 800 ml, giá vốn 35 đ/ml (BR-28).
- **Given** tồn 3 giá 10.000 **When** nhập 4 giá 11.000 **Then** giá vốn 10.571 (BR-33).
- **Given** phiếu tổng 1.250.000 đ trả tiền mặt **Then** có CashEntry `PURCHASE_PAYMENT` 1.250.000 đ, và báo cáo chi phí vận hành không tăng.
- **Given** phiếu nhập đầu tiên của cửa hàng **Then** `costing_locked_at` được đặt.
- Nhân viên nhập được và thấy giá nhập trên phiếu mình nhập, nhưng màn hình tồn kho không hiện giá vốn bình quân.

## API gợi ý
`POST /api/v1/stores/{storeId}/purchase-receipts` · `GET .../purchase-receipts` · `GET .../purchase-receipts/{id}` · `GET/POST /api/v1/households/{householdId}/suppliers`
