# UC-22 Tìm hoặc tạo khách theo số điện thoại

- **Milestone:** M5 · **Context:** customer · **Capability:** `core.customer` · **Tác nhân:** Nhân viên, Quản lý, Chủ hộ
- **Quy tắc:** BR-20 · **Quyền:** `customer.view`

## Luồng chính
1. Ở POS hoặc màn hình khách hàng, gõ số điện thoại hoặc tên. Hệ thống gợi ý khách trong **cả hộ** (tìm theo đầu số, tên không dấu).
2. Chọn khách có sẵn, hoặc "Thêm khách" với tên + số điện thoại (số điện thoại có thể bỏ trống).
3. Khách được gắn vào đơn đang bán.

## Ngoại lệ
- Tạo khách với số đã có trong hộ → trả khách hiện có (HTTP 200, có cờ `existing=true`), không tạo trùng.
- Số sai định dạng → lỗi trường.

## Tiêu chí chấp nhận
- **Given** khách "Chị Lan 0912345678" được tạo ở tiệm điện thoại **When** nhân viên quán trà sữa cùng hộ gõ "0912" **Then** thấy "Chị Lan".
- **Given** tạo khách số "0912 345 678" khi đã có 0912345678 **Then** không có bản ghi mới.
- **Given** hai hộ khác nhau **Then** khách của hộ này không bao giờ hiện ở hộ kia.
- Kết quả tìm kiếm cho nhân viên chỉ gồm tên và số điện thoại, không kèm lịch sử ở cửa hàng khác.

## API gợi ý
`GET /api/v1/stores/{storeId}/customers?q=` · `POST /api/v1/stores/{storeId}/customers` · `PATCH .../customers/{id}`

Đường dẫn qua `storeId` để xác định phạm vi và quyền, dù khách thuộc cấp hộ.
