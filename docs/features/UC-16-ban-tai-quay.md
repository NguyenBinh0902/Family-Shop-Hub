# UC-16 Bán tại quầy

- **Milestone:** M4 (gắn khách và bán nợ bổ sung ở M5, tùy chọn món ở M7) · **Context:** sales (+ inventory, customer, setup) · **Capability:** `core.pos` · **Tác nhân:** Nhân viên, Quản lý, Chủ hộ
- **Quy tắc:** BR-10, BR-13, BR-14, BR-16, BR-17, BR-18, BR-19, BR-25, BR-28, BR-33, TAX-02, TAX-04, NFR-03, NFR-04 · **Quyền:** `sale.create`, `sale.discount`

## Mục tiêu
Bán một đơn nhanh nhất có thể; một đơn thường dưới 15 giây.

## Luồng chính
1. Màn hình POS hiển thị sản phẩm `is_sellable`, `ACTIVE`, nhóm theo danh mục, có ô tìm kiếm.
2. Nhân viên chạm sản phẩm để thêm (chạm lại tăng số lượng). Nếu có nhóm tùy chọn → hộp chọn tùy chọn (M7).
3. (M5) Gắn khách theo số điện thoại (UC-22).
4. (Tùy chọn) Giảm giá cấp đơn: số tiền hoặc %.
5. Bấm "Thanh toán", chọn một hoặc nhiều hình thức: tiền mặt (nhập tiền khách đưa → hiện tiền thối), chuyển khoản (chọn tài khoản của chủ thể), ví, (M5) ghi nợ.
6. Client gửi đơn với `id` UUID v7 sinh từ trước (BR-17).
7. Server hoàn tất đơn theo các bước ở `domain/sales.md`.
8. Hiển thị "Đã bán" kèm nút In (UC-21) và "Đơn mới".

## Ngoại lệ
- Tồn không đủ: cửa hàng cho bán âm → vẫn bán, tồn âm; cửa hàng chặn → lỗi `BR-13` chỉ rõ sản phẩm thiếu.
- Tổng thanh toán ≠ tổng tiền → lỗi `BR-18`.
- Ghi nợ không có khách → lỗi `BR-19`.
- Nhân viên giảm giá vượt mức trần → lỗi 403 kèm mức trần.
- Mạng lỗi khi gửi → client gửi lại cùng `id`; server trả đơn đã có nếu đã tạo.
- Cửa hàng không `ACTIVE` → từ chối.

## Tiêu chí chấp nhận
- **Given** "Bánh mì" `QUANTITY` 15.000 tồn 10 **When** bán 2, trả tiền mặt 50.000 **Then** đơn `COMPLETED` tổng 30.000, Payment `CASH` 30.000 với `tendered` 50.000, tồn còn 8, có movement `SALE` −2.
- **Given** cùng đơn gửi hai lần với cùng `id` **Then** chỉ có một đơn, tồn chỉ giảm một lần, lần hai trả 200 với đơn cũ.
- **Given** cửa hàng chặn bán âm, tồn 1 **When** bán 2 **Then** lỗi `BR-13`, không ghi gì.
- **Given** đơn 100.000, giảm 15% **Then** `discount_amount` 15.000, `total` 85.000.
- **Given** đơn 33.333, giảm 10% **Then** `discount_amount` 3.333 (BR-33).
- **Given** nhân viên, mức trần 0% **When** giảm 5.000 **Then** 403.
- **Given** trả 50.000 chuyển khoản + 35.000 tiền mặt cho đơn 85.000 **Then** hai Payment đúng hình thức, chuyển khoản gắn `account_id`.
- **Given** sản phẩm chưa từng nhập (chưa có giá vốn) **When** bán **Then** dòng có `unit_cost` 0 và `cost_missing=true` (BR-28).
- **Given** đơn hoàn tất lúc 23:30 ngày 01/12 giờ Việt Nam (16:30 UTC) **Then** báo cáo tính vào ngày 01/12 và số đơn có dạng `261201-NNN`.
- Có test chạy song song 20 request cùng `id` chỉ tạo một đơn.

## API gợi ý
`POST /api/v1/stores/{storeId}/orders` (tạo và hoàn tất) · `GET .../orders/{id}` · `GET .../pos/catalog` (dữ liệu tối giản cho POS: sản phẩm, giá, tùy chọn)

## Ghi chú UI
Lưới sản phẩm nút lớn; giỏ hàng luôn thấy trên máy tính bảng, kéo lên từ dưới trên điện thoại. Hỏi chủ repo bố cục trước khi làm (frontend.md mục 6).
