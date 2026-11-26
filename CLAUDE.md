# FamilyShop Hub

Ứng dụng web SaaS cho hộ kinh doanh gia đình quản lý nhiều cửa hàng khác ngành (F&B, hàng có serial...) trong một tài khoản. Loại hình kinh doanh của cửa hàng quyết định chức năng, trường dữ liệu và giao diện.

## Nguồn sự thật

| Cần biết | Đọc |
| --- | --- |
| Làm gì tiếp theo, thứ tự milestone | `docs/roadmap.md` |
| Đọc tài liệu nào khi làm một feature | `docs/README.md` |
| Use case và tiêu chí chấp nhận | `docs/features/UC-xx-*.md` |
| Quy tắc nghiệp vụ BR-xx | `docs/rules.md` |
| Thực thể, trạng thái theo context | `docs/domain/` |
| Thuật ngữ và tên trong code | `docs/glossary.md` |
| Stack và quyết định kiến trúc | `docs/adr/` |
| Quy ước code backend Go | `docs/conventions/backend-go.md` |
| Frontend | `web/CLAUDE.md` |

Khi code dẫn chiếu một quy tắc nghiệp vụ, ghi mã của nó (`BR-15`, `UC-16`) trong comment, tên test và commit message.

## Stack

Go (modular monolith, hexagonal) · Gin · PostgreSQL · pgx + sqlc · goose · OpenAPI + oapi-codegen · React + Vite + TypeScript · Tailwind + shadcn/ui · TanStack Query · PWA · Docker Compose.

## Cấu trúc repo

```
cmd/api/            entrypoint của backend
internal/<context>/ mỗi bounded context: domain/, app/, adapter/
internal/platform/  hạ tầng dùng chung: config, db, http server, auth, tenant
api/openapi.yaml    hợp đồng API, nguồn sinh code cho cả backend và frontend
db/migrations/      migration SQL (goose)
db/queries/         query SQL cho sqlc, chia theo context
web/                frontend
docs/               tài liệu nghiệp vụ, ADR, quy ước
```

## Lệnh

Chạy `make help` để xem đủ danh sách. Cần Go 1.27, Node 22, Docker (test tích hợp dùng testcontainers).

| Lệnh | Việc |
| --- | --- |
| `make dev` | Chạy cả hệ thống bằng Docker Compose tại http://localhost:8080 (tự tạo `.env` từ `.env.example`) |
| `make dev-db` rồi `make dev-api` + `make dev-web` | Phát triển nhanh: Postgres trong Docker, API chạy trên máy (cổng 8081), Vite hot reload proxy `/api` |
| `make test` | `go test -race ./...` (gồm test tích hợp với Postgres thật) và typecheck frontend |
| `go test -short ./...` | Chỉ unit test, bỏ qua test cần Docker |
| `make lint` | `go vet`, `golangci-lint` (có depguard chặn `domain` import hạ tầng), ESLint, Prettier |
| `make gen` | Sinh lại code sqlc, oapi-codegen và client TypeScript sau khi sửa `db/queries/` hoặc `api/openapi.yaml` |
| `make gen-check` | Báo lỗi nếu code sinh ra chưa khớp (CI chạy lệnh này) |
| `make migrate-new name=add_products` | Tạo `db/migrations/NNNN_add_products.sql` |
| `make migrate-up` / `migrate-down` / `migrate-status` | Áp dụng, lùi, xem migration trên DB dev |

Công cụ sinh code được pin trong `tools/go.mod`, gọi qua `go tool -modfile=tools/go.mod`.

## Luật không được vi phạm

- Mọi bảng nghiệp vụ có `household_id`; bảng cấp cửa hàng có thêm `store_id`. Không query nào được trả dữ liệu của hộ khác (BR-01).
- Tiền là số nguyên VND (`int64`), không bao giờ dùng số thực.
- Đơn bán, phiếu thu chi đã hoàn tất là bất biến. Hủy hoặc sửa bằng bút toán đảo, không UPDATE/DELETE (BR-15, BR-23).
- Dòng đơn lưu snapshot tên, giá bán, giá vốn tại thời điểm bán (BR-10, BR-16).
- Tạo đơn idempotent theo mã đơn do client sinh (BR-17).
- Thời điểm lưu UTC; ngày nghiệp vụ tính theo `Asia/Ho_Chi_Minh`.
- Báo cáo cấp hộ là phép gộp, không lưu bản dữ liệu thứ hai (BR-24).
- Mục đánh dấu CHƯA CHỐT trong `docs/` thì không implement; hỏi chủ repo.

## Cách làm việc

- **Backend:** đây là dự án để chủ repo học Go. Viết code idiomatic, rõ ràng, theo `docs/conventions/backend-go.md`. Khi tóm tắt thay đổi, giải thích ngắn gọn các lựa chọn thiết kế không hiển nhiên và idiom Go được dùng.
- **Frontend:** xem `web/CLAUDE.md`. Chủ repo không chuyên frontend; có những quyết định phải hỏi trước.
- **Plan trước, code sau:** mỗi milestone hoặc feature lớn bắt đầu bằng plan được chủ repo duyệt (quy trình ở cuối `docs/roadmap.md`). Tài liệu chưa rõ thì liệt kê câu hỏi, không tự đoán.
- **Ngôn ngữ:** UI tiếng Việt. Identifier, comment, commit message bằng tiếng Anh. Tài liệu trong `docs/` bằng tiếng Việt.
- Thêm dependency mới (Go hoặc npm) phải nêu lý do và hỏi trước.
