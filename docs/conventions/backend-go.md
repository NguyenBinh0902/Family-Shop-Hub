# Quy ước backend Go

Mục tiêu: code đọc như sách giáo khoa Go. Rõ ràng hơn ngắn gọn, tường minh hơn "thông minh". Khi phân vân, theo [Effective Go](https://go.dev/doc/effective_go) và [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).

## 1. Bố cục package

```
cmd/api/main.go                 wiring: đọc config, tạo pool DB, ráp adapter vào app, chạy server
internal/
  platform/                     hạ tầng dùng chung, không chứa nghiệp vụ
    config/  db/  httpserver/  auth/  tenant/  problem/
  kernel/                       kiểu nhỏ dùng chung giữa các context: Money, ID, Clock
  <context>/                    ví dụ sales/
    domain/                     entity, value object, bất biến, lỗi nghiệp vụ
    app/                        use case + port (interface) mà use case cần
    adapter/
      postgres/                 hiện thực port bằng sqlc
      http/                     handler Gin, chuyển DTO <-> lệnh của app
db/migrations/                  goose, đặt tên NNNN_mo_ta.sql
db/queries/<context>/*.sql      query cho sqlc
api/openapi.yaml
```

Tên context: `identity`, `setup`, `catalog`, `inventory`, `sales`, `customer`, `cashbook`, `reporting`.

## 2. Quy tắc phụ thuộc

```
adapter  →  app  →  domain  →  kernel
```

- `domain` không import Gin, pgx, sqlc, hay `net/http`. Chỉ thư viện chuẩn và `kernel`.
- `app` định nghĩa interface (port) nó cần, ví dụ `type OrderRepository interface`. Interface khai báo ở nơi **dùng**, không ở nơi hiện thực.
- Context này gọi context khác qua interface khai báo trong `app` của mình, hiện thực bằng adapter gọi `app` của context kia. Không import thẳng `adapter/postgres` hay `domain` của context khác.
- `main.go` là nơi duy nhất biết mọi thứ ráp với nhau thế nào. Không dùng DI framework, không biến toàn cục.

## 3. Domain

- Entity có constructor trả lỗi khi vi phạm bất biến: `func NewSalesOrder(...) (*SalesOrder, error)`.
- Trường của entity để unexported khi cần giữ bất biến; thay đổi trạng thái qua method (`order.Complete(now)`, `order.Cancel(reason, by)`).
- Value object là kiểu riêng, so sánh bằng giá trị: `type Money int64`, `type Phone string` có hàm chuẩn hóa.
- Trạng thái là kiểu riêng với hằng số, giá trị 0 là trạng thái không hợp lệ:

```go
type OrderStatus int

const (
	OrderStatusUnknown OrderStatus = iota
	OrderStatusCompleted
	OrderStatusCancelled
)
```

## 4. Lỗi

- Lỗi nghiệp vụ khai báo trong `domain` của từng context, là giá trị kiểu `kernel.Error` (dùng chung để `platform/problem` nhận ra), kèm mã BR và `Kind` quyết định mã HTTP (`Invalid` 422, `Conflict` 409, `Forbidden` 403, `NotFound` 404):

```go
var ErrOrderImmutable = &kernel.Error{Code: "BR-15", Kind: kernel.ErrorKindConflict, Message: "Đơn đã hoàn tất, không thể sửa"}
```

- Bọc lỗi bằng `fmt.Errorf("load order %s: %w", id, err)`; kiểm tra bằng `errors.Is` / `errors.As`. Không so sánh chuỗi lỗi.
- Chuỗi lỗi kỹ thuật viết thường, không dấu chấm cuối. `Message` của lỗi nghiệp vụ là tiếng Việt vì hiển thị cho người dùng.
- Một middleware duy nhất (`platform/problem`) chuyển lỗi thành JSON theo RFC 9457 Problem Details: lỗi nghiệp vụ → 409 hoặc 422 kèm `code`; không tìm thấy → 404; lỗi khác → 500, chỉ log chi tiết, không lộ ra client.
- Không `panic` trong luồng xử lý request.

## 5. Tenant và phân quyền

- Middleware đọc `storeId` từ URL, kiểm tra membership, đặt `tenant.Scope{HouseholdID, StoreID, UserID, Role}` vào request.
- Use case nhận `Scope` là **tham số tường minh**, không lấy ngầm từ `context.Context`. Repository nhận và lọc theo `HouseholdID`/`StoreID` trong mọi query.
- Mỗi transaction chạy `SET LOCAL app.household_id = ...` để Row-Level Security làm lớp chặn thứ hai.

## 6. Transaction

- Một use case ghi dữ liệu = một transaction, do tầng `app` quyết định qua port:

```go
type TxRunner interface {
	InTx(ctx context.Context, householdID uuid.UUID, fn func(ctx context.Context, repos Repos) error) error
}
```

- `Repos` là struct các port mà use case cần, kể cả port sang context khác; `main.go` dựng tất cả từ cùng một `pgx.Tx`. `*db.TxRunner[Repos]` hiện thực interface này và đặt `app.household_id` cho RLS. Chi tiết và lý do ở ADR-0004.

- Repository không tự mở transaction.

## 7. Dữ liệu và thời gian

- Tiền: `kernel.Money` (`int64`, đơn vị đồng). Cột DB `bigint`.
- ID: UUID v7 (sắp xếp được theo thời gian). Mã đơn do client sinh cũng là UUID v7.
- Thời điểm: `time.Time` UTC, cột `timestamptz`. Ngày nghiệp vụ tính bằng `kernel.Clock` với location `Asia/Ho_Chi_Minh`. Không gọi `time.Now()` trong domain; nhận thời điểm qua tham số hoặc `Clock` để test được.
- Thuộc tính động: cột `jsonb`, kiểm tra theo định nghĩa trường của loại hình trước khi lưu.

## 8. HTTP

- Handler mỏng: parse request → gọi use case → map kết quả. Không chứa nghiệp vụ.
- Kiểu request/response sinh từ `api/openapi.yaml`. Sửa API thì sửa spec trước, rồi `make gen`.
- Route nghiệp vụ dạng `/api/v1/stores/{storeId}/...`.

## 9. Logging và context

- `log/slog`, truyền logger qua constructor. Không `fmt.Println`, không `log.Printf`.
- `context.Context` là tham số đầu tiên của mọi hàm có I/O. Không lưu context trong struct.
- Log có `request_id`, `household_id`, `store_id`. Không log mật khẩu, token, số điện thoại đầy đủ.

## 10. Test

- Domain: unit test dạng table-driven, không cần DB.
- Adapter Postgres và use case có transaction: integration test với testcontainers-go trên Postgres thật, migration chạy thật.
- Tên test nêu quy tắc được kiểm: `TestSalesOrder_Cancel_BR15_CreatesReversal`.
- Mỗi BR có ít nhất một test kiểm nó.

## 11. Đặt tên và phong cách

- Package tên ngắn, số ít, chữ thường: `sales`, không `salesService`. Không có package `utils`, `helpers`, `common`.
- Constructor `NewX`; getter không có tiền tố `Get`; interface một method đặt đuôi `-er` khi hợp lý.
- Viết tắt giữ nguyên hoa/thường: `storeID`, `HTTPServer`, `IMEI`.
- Comment giải thích **vì sao**, không nhắc lại code làm gì. Mọi identifier exported có doc comment.
- `gofmt`, `go vet`, `golangci-lint` phải sạch trước khi commit.

## 12. Dependency được phép

`gin`, `pgx/v5`, `sqlc` (công cụ sinh code), `oapi-codegen`, `goose`, `google/uuid`, `golang.org/x/crypto` (argon2id), `testcontainers-go`. Thêm thư viện ngoài danh sách phải nêu lý do và hỏi chủ repo.
