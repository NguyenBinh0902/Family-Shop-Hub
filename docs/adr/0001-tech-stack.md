# ADR-0001: Kiến trúc và tech stack

- **Trạng thái:** Đã chấp nhận
- **Ngày:** 2026-11-20

## Bối cảnh

- Dự án cá nhân, một người phát triển cùng các agent Claude Code.
- Chủ repo thạo Go, Gin, PostgreSQL, Kubernetes và kiến trúc hexagonal/DDD; không chuyên frontend.
- Mục tiêu học hỏi: backend phải chuẩn, gọn, idiomatic.
- Nghiệp vụ cần transaction chặt (bán hàng trừ kho và ghi sổ cùng lúc), dữ liệu đa tenant, thuộc tính động theo loại hình, và báo cáo tổng hợp.
- Người dùng cuối dùng chủ yếu trên điện thoại và máy tính bảng tại quầy.

## Quyết định

### Kiến trúc

**Modular monolith, API-first.** Một binary Go duy nhất. Mỗi bounded context (Identity & Access, Business Setup, Catalog, Inventory, Sales, Customer & Debt, Cashbook, Reporting) là một package với tầng domain, app, adapter theo hexagonal. `api/openapi.yaml` là hợp đồng giữa backend và frontend, viết trước khi code.

### Stack

| Tầng | Lựa chọn |
| --- | --- |
| Ngôn ngữ backend | Go, bản stable mới nhất lúc scaffold (tối thiểu 1.26) |
| HTTP | Gin, handler sinh interface từ OpenAPI bằng oapi-codegen |
| Database | PostgreSQL: transaction, JSONB cho thuộc tính động, Row-Level Security |
| Truy cập DB | pgx v5 + sqlc |
| Migration | goose (file SQL) |
| Xác thực | Session token ngẫu nhiên lưu DB (thu hồi được), mật khẩu băm argon2id |
| Frontend | React + Vite + TypeScript (SPA) |
| UI | Tailwind CSS + shadcn/ui |
| Dữ liệu frontend | TanStack Query + client sinh từ OpenAPI |
| Ứng dụng cài đặt | PWA (vite-plugin-pwa) |
| In hóa đơn | CSS print khổ 58 mm và 80 mm từ trình duyệt |
| Test | `go test`, testcontainers-go (Postgres thật), Playwright cho luồng chính |
| Log | `log/slog` dạng JSON |
| Triển khai MVP | Docker Compose trên một VPS: API, Postgres, Caddy (HTTPS); backup `pg_dump` hằng ngày ra object storage |

## Hệ quả

- Chỉ một tiến trình, một database: debug, deploy, test đơn giản. Ranh giới context được giữ bằng package và interface, không bằng mạng.
- sqlc buộc viết SQL thật; báo cáo phức tạp viết thẳng bằng SQL. Đổi lại phải chạy `make gen` sau mỗi thay đổi query.
- OpenAPI sinh code ở cả hai đầu nên backend và frontend luôn khớp kiểu; mọi thay đổi API bắt đầu từ file spec.
- Session lưu DB cho phép thu hồi tức thì (NFR-06) với giá một lần đọc DB mỗi request; chấp nhận được ở quy mô MVP.
- Image Docker chuẩn từ đầu, chuyển sang Kubernetes sau mà không phải sửa ứng dụng.

## Phương án đã loại

| Phương án | Lý do loại |
| --- | --- |
| Microservices | Một người phát triển, chưa có nhu cầu scale riêng từng phần; chi phí vận hành cao. |
| Aerospike | Không hợp dữ liệu quan hệ, transaction nhiều bảng và báo cáo SQL. |
| ORM (GORM...) | Che SQL, khó tối ưu báo cáo, khó đọc khi review. |
| JWT không trạng thái | Khó thu hồi phiên ngay khi rút quyền thành viên. |
| Next.js | MVP không cần SEO hay render phía server; SPA đơn giản hơn khi backend là Go. |
| Ứng dụng di động native | PWA đủ cho MVP; app riêng nằm ngoài phạm vi. |
| Kubernetes ngay từ MVP | Quá mức cần cho một VPS; để dành khi có nhiều hộ. |

## Còn mở

- Nhà cung cấp SMS OTP (hoặc Zalo ZNS) cho chức năng quên mật khẩu: chưa chọn. Ẩn sau port `OTPSender`; MVP dùng adapter ghi log khi phát triển.
- Nơi đặt máy chủ và dữ liệu: cần kiểm tra quy định bảo vệ dữ liệu cá nhân hiện hành trước khi đưa người dùng thật vào.
