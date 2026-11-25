# UC-05 Tạo, sửa chủ thể kinh doanh

- **Milestone:** M2 · **Context:** setup · **Tác nhân:** Chủ hộ
- **Quy tắc:** BR-03, BR-26, TAX-02, TAX-03 · **Quyền:** `entity.manage`

## Luồng chính
1. Chủ hộ nhập tên hộ kinh doanh, người đứng tên, MST (có thể để trống), địa chỉ.
2. Thêm tài khoản ngân hàng (ngân hàng, số tài khoản, chủ tài khoản) và ví điện tử (nhà cung cấp, số tài khoản ví).
3. Lưu. Sửa sau được; mọi thay đổi ghi nhật ký.

## Ngoại lệ
- MST có ký tự không phải số hoặc độ dài ngoài 10–13 → lỗi trường.
- Xóa chủ thể còn cửa hàng → không cho phép.
- Xóa tài khoản/ví đã được dùng trong thanh toán → chuyển sang ngưng dùng thay vì xóa.

## Tiêu chí chấp nhận
- **Given** hộ chưa có chủ thể **When** tạo chủ thể không có MST **Then** lưu thành công, hiển thị nhắc "Bổ sung mã số thuế".
- **Given** tài khoản ngân hàng đã có Payment tham chiếu **When** xóa **Then** tài khoản chuyển ngưng dùng, Payment cũ vẫn hiển thị đúng tên tài khoản.
- **Given** người dùng là Quản lý **When** gọi API sửa chủ thể **Then** trả 403.
- Sửa MST ghi nhật ký giá trị trước và sau.

## API gợi ý
`GET/POST /api/v1/households/{householdId}/business-entities` · `GET/PATCH /api/v1/households/{householdId}/business-entities/{id}`
