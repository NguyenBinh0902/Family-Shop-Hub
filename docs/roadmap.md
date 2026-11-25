# Roadmap MVP

MVP chia thành 9 milestone theo **lát dọc**: mỗi milestone làm trọn từ migration, domain, API đến màn hình, và khi xong thì chạy, dùng thử được. Làm theo thứ tự; không bắt đầu milestone sau khi milestone trước chưa đạt "Xong khi".

Ưu tiên đi thẳng tới M4 với một cửa hàng F&B đơn giản (bán sản phẩm làm sẵn, chưa có tùy chọn món) để sớm có sản phẩm cho gia đình dùng thử.

## Tổng quan

```
M0 Khung ─► M1 Tài khoản ─► M2 Cửa hàng ─► M3 Sản phẩm & kho ─► M4 Bán tại quầy
                                                                   │
                         ┌──────────────────┬──────────────────────┤
                         ▼                  ▼                      ▼
                   M5 Khách, nợ,      M6 Hàng serial        M7 F&B mở rộng
                      thu chi          (cần M5)                    │
                         └──────────────────┴──────────┬───────────┘
                                                       ▼
                                                  M8 Báo cáo
```

M6 cần M5 vì đặt cọc bắt buộc có khách hàng. M7 chỉ cần M4.

| Milestone | Use case | Xong khi |
| --- | --- | --- |
| M0 Khung chạy được | (hạ tầng) | `make dev` chạy toàn bộ; health check xanh; CI xanh |
| M1 Tài khoản và hộ | UC-01, 02, 03, 04 | Đăng ký, đăng nhập, mời người; dữ liệu tách theo hộ có RLS |
| M2 Cửa hàng | UC-05, 06, 07, 08 | Tạo nhiều cửa hàng theo loại hình, chuyển qua lại |
| M3 Sản phẩm và kho số lượng | UC-09, 12, 14, 15 | Nhập hàng, xem tồn, giá vốn bình quân đúng |
| M4 Bán tại quầy | UC-16 (cơ bản), 20, 21 | Bán, hủy, in hóa đơn thật trên điện thoại; thấy doanh thu hôm nay |
| M5 Khách, công nợ, thu chi | UC-22, 23, 24, 26 (+ UC-16 gắn khách, ghi nợ) | Đủ sổ sách cho một cửa hàng |
| M6 Hàng serial | UC-13, 18, 19, 25 (+ UC-09 bảo hành) | Tiệm điện thoại dùng được |
| M7 F&B mở rộng | UC-10, 11, 17 | Quán trà sữa, quán ăn dùng đầy đủ |
| M8 Báo cáo | UC-27, 28, 29, 30 | Chủ hộ thấy bức tranh tổng, xuất được sổ doanh thu |

## M0 Khung chạy được

Mục tiêu: dựng bộ khung chuẩn mà mọi milestone sau dựa vào. Đây là milestone học được nhiều nhất về cách tổ chức một dự án Go.

Phải có:

- `go.mod`, `cmd/api/main.go` với wiring thủ công (không DI framework), graceful shutdown.
- `internal/platform/config` (đọc biến môi trường, kiểm tra khi khởi động), `db` (pgxpool, `TxRunner`, đặt `app.household_id` cho RLS), `httpserver` (Gin, middleware request ID, log, recover), `problem` (lỗi → RFC 9457), `tenant` và `auth` (khung rỗng, hoàn thiện ở M1).
- `internal/kernel`: `Money`, `Quantity` (phần nghìn) và phép nhân `Quantity × Money`, làm tròn nửa lên (BR-33, BR-34), `Phone` chuẩn hóa, `Clock` với `Asia/Ho_Chi_Minh`, sinh UUID v7, cộng tháng kẹp cuối tháng. Có unit test đầy đủ.
- `db/migrations/0001_init.sql`: extension cần thiết, role ứng dụng không phải superuser (để RLS có hiệu lực).
- `sqlc.yaml`, cấu hình oapi-codegen, `api/openapi.yaml` có `GET /api/v1/healthz` (kiểm tra kết nối DB).
- `web/`: Vite + React + TS strict, Tailwind, shadcn/ui, React Router, TanStack Query, client sinh từ OpenAPI, layout khung, trang gọi healthz, manifest PWA.
- `docker-compose.yml`: postgres, api, caddy (phục vụ `web/dist` và proxy `/api`).
- `Makefile`: `dev`, `test`, `lint`, `gen`, `migrate-up`, `migrate-new`.
- Helper test tích hợp dùng testcontainers (khởi động Postgres, chạy migration).
- CI: `go vet`, `golangci-lint`, `go test ./...`, lint, typecheck, build frontend.
- Cập nhật mục "Lệnh" trong `CLAUDE.md`.

Xong khi: clone mới, `make dev` chạy được; trang web hiện "Kết nối máy chủ: OK"; `make test` và `make lint` sạch.

## M1 – M8

Mỗi milestone: đọc các file use case trong bảng trên. Phần ghi chú "(Mx)" bên trong một use case cho biết đoạn đó thuộc milestone nào; milestone hiện tại chỉ làm phần của mình, nhưng thiết kế không được chặn phần sau (ví dụ: M4 chưa có `DEBT` nhưng enum `PaymentMethod` đã có chỗ cho nó).

Ghi chú riêng:

- **M1:** RLS bật từ đây cho mọi bảng có `household_id`. Có test chứng minh một hộ không đọc được dữ liệu hộ khác kể cả khi repository quên điều kiện lọc.
- **M2:** tạo đủ 5 file cấu hình loại hình MVP (xem `store-types.md`).
- **M3:** `StockService` cho Sales thiết kế luôn ở đây (dù M4 mới dùng).
- **M4:** UC-16 không gồm tùy chọn món, gắn khách, ghi nợ. Trang chủ cửa hàng hiện doanh thu và số đơn hôm nay (bản sớm của UC-27). Có test chạy song song cho BR-17.
- **M5:** bổ sung phần (M5) của UC-16 và UC-20.
- **M6:** bổ sung bảo hành vào UC-09, phần (M6) của UC-20.
- **M8:** khi gặp mục CHƯA CHỐT ở UC-29, dừng ở phạm vi đã mô tả.

## Quy trình mỗi milestone với Claude Code

1. **Lập plan (plan mode, chưa code).** Prompt mẫu:

   > Lập plan cho M3 theo `docs/roadmap.md`. Đọc UC-09, UC-12, UC-14, UC-15, `docs/rules.md`, `docs/domain/catalog.md`, `docs/domain/inventory.md`. Plan gồm: migration, query sqlc, thay đổi OpenAPI, package và kiểu Go chính, test cho từng BR, màn hình frontend. Liệt kê chỗ tài liệu chưa rõ thay vì tự đoán. Chưa viết code.

2. **Chủ repo review plan.** Đây là lúc hỏi "vì sao chia package thế này", "vì sao dùng idiom này".
3. **Chia việc nhỏ.** Mỗi bước một commit có thể chạy được: migration + query → domain + test → app + adapter → API → frontend.
4. **Code kèm test.** Mỗi tiêu chí chấp nhận trong file UC có ít nhất một test tương ứng.
5. **Review diff, commit.** Commit message có mã UC/BR.
6. **Đóng milestone:** đối chiếu "Xong khi"; quyết định kiến trúc mới phát sinh thì ghi ADR; tài liệu sai hoặc thiếu thì sửa `docs/` trong cùng PR.
