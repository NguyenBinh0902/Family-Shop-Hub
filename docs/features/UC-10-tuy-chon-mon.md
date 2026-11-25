# UC-10 Cấu hình tùy chọn món

- **Milestone:** M7 · **Context:** catalog · **Capability:** `fnb.options` · **Tác nhân:** Quản lý, Chủ hộ
- **Quy tắc:** BR-14, BR-16 · **Quyền:** `catalog.manage`

## Luồng chính
1. Tạo nhóm tùy chọn: tên (Size, Topping...), chọn một hay nhiều, bắt buộc hay không, số lựa chọn tối đa (nếu chọn nhiều).
2. Thêm lựa chọn: tên, giá cộng thêm (≥ 0), nguyên liệu kèm theo (tùy chọn, ví dụ "Trân châu" trừ 30 g trân châu).
3. Gắn nhóm tùy chọn vào sản phẩm, sắp thứ tự hiển thị.

## Ngoại lệ
- Nguyên liệu kèm theo không phải sản phẩm `QUANTITY` của cùng cửa hàng → lỗi.
- Xóa nhóm/lựa chọn đã được dùng trong đơn → chuyển ngưng dùng (đơn cũ có snapshot).

## Tiêu chí chấp nhận
- **Given** "Trà sữa" 25.000 có Size (bắt buộc, chọn một: M +0, L +6.000) và Topping (chọn nhiều, tối đa 3: Trân châu +5.000) **When** POS chọn L + Trân châu **Then** đơn giá dòng là 36.000 và snapshot ghi cả hai lựa chọn.
- **Given** Size bắt buộc **When** thêm món không chọn size **Then** server từ chối dòng đó.
- **Given** chọn 4 topping khi tối đa 3 **Then** server từ chối.
- **Given** "Trân châu" mang 30 g nguyên liệu trân châu **When** hoàn tất đơn có 2 ly kèm trân châu **Then** kho trân châu giảm 60 g.

## API gợi ý
`GET/POST /api/v1/stores/{storeId}/option-groups` · `PATCH .../option-groups/{id}` · `PUT .../products/{id}/option-groups`
