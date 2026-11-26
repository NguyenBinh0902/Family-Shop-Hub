# Frontend (web/)

Chỉ dẫn này bổ sung cho `CLAUDE.md` ở gốc repo khi làm việc trong `web/`.

Chủ repo làm backend, không chuyên frontend. Agent tự quyết chi tiết triển khai, nhưng:

- **Phải hỏi trước** các quyết định trong mục 6 của `docs/conventions/frontend.md` (bố cục màn hình mới, điều hướng, phong cách chung, thư viện mới, chỗ làm khác use case). Đưa 1–2 phương án, mỗi phương án một câu lý do.
- Khi tóm tắt thay đổi, giải thích bằng ngôn ngữ dễ hiểu với người làm backend.

## Đọc trước khi làm

- `docs/conventions/frontend.md`: stack, bố cục thư mục, quy tắc giao diện, dữ liệu.
- File use case liên quan trong `docs/features/` (mục "Ghi chú UI").
- `docs/roles-permissions.md`: ẩn thao tác người dùng không có quyền; server vẫn là nơi kiểm tra cuối.

## Nguyên tắc cốt lõi

- Mobile-first từ 360 px; vùng bấm ≥ 44 px; POS thao tác một tay.
- Toàn bộ chữ trên giao diện tiếng Việt, từ ngữ đời thường. Code, comment tiếng Anh.
- Ngữ cảnh cửa hàng nằm trong URL `/stores/:storeId/...`; không giữ cửa hàng hiện tại trong state toàn cục.
- Kiểu dữ liệu API chỉ lấy từ client sinh từ `api/openapi.yaml`. Đổi API thì sửa spec ở gốc, chạy `make gen`, không sửa file sinh tự động.
- Tạo đơn: sinh UUID v7 trước khi gửi, gửi lại cùng mã khi lỗi mạng (BR-17), khóa nút khi đang gửi.
- Tiền định dạng `vi-VN`/`VND`; ngày `dd/MM/yyyy`, múi giờ `Asia/Ho_Chi_Minh`.
- Ẩn chức năng theo capability của cửa hàng (lấy từ `GET /api/v1/stores/{storeId}`), không theo mã loại hình.

## Lệnh

Chạy qua `Makefile` ở gốc repo: `make dev-web` (Vite hot reload, cần `make dev-db` và `make dev-api`), `make lint-web`, `make test-web`, `make gen` (sinh lại `src/lib/api/schema.d.ts`). Trong `web/`: `npm run lint`, `npm run typecheck`, `npm run build`, `npm run format`.

TypeScript pin ở 5.9 vì `openapi-typescript` và `typescript-eslint` chưa hỗ trợ TypeScript 7.
