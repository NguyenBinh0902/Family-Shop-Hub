# ADR-0002: Loại hình kinh doanh là cấu hình, capability là code

- **Trạng thái:** Đã chấp nhận
- **Ngày:** 2026-11-25

## Bối cảnh

Mỗi cửa hàng có một loại hình quyết định chức năng. Yêu cầu: thêm ngành mới không viết lại hệ thống; người dùng không tự bật/tắt chức năng; cửa hàng kết hợp nhiều nhóm phải làm được về sau.

## Quyết định

- Chức năng chia thành **capability** (ví dụ `inv.serial`, `fnb.recipe`), mỗi capability là code thật trong các context.
- **Loại hình** là file JSON khai báo tổ hợp capability, trường động, dữ liệu mẫu; nhúng vào binary bằng `go:embed`, có version, không sửa file của version đã dùng.
- Code luôn hỏi "cửa hàng có capability X không" qua `setup`, **không bao giờ** rẽ nhánh theo mã loại hình.
- Thuộc tính động lưu ở cột `jsonb`, kiểm tra theo định nghĩa của loại hình. Trường quan trọng (IMEI, giá, số lượng) luôn là cột thật.

## Hệ quả

- Thêm ngành dùng capability có sẵn = thêm một file JSON và review qua git.
- Loại hình kết hợp (B + C) chỉ là một file khai báo nhiều capability hơn.
- Không có giao diện quản trị loại hình; chấp nhận vì chỉ chủ repo thêm loại hình.

## Phương án đã loại

- Lưu loại hình trong DB: cần màn hình quản trị và migration dữ liệu, không lợi gì ở quy mô này.
- Cho người dùng tự bật/tắt capability: khó dùng với người không rành công nghệ (đã thảo luận và loại).
- Bảng EAV cho thuộc tính động: truy vấn phức tạp, không hơn `jsonb`.
