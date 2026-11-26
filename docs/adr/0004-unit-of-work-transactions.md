# ADR-0004: Transaction theo Unit of Work tường minh

- **Trạng thái:** Đã chấp nhận
- **Ngày:** 2026-11-26

## Bối cảnh

Nhiều use case ghi dữ liệu của nhiều context trong cùng một transaction. Ví dụ UC-12 (nhập hàng) ghi phiếu nhập và tồn kho (inventory), phiếu chi (cashbook), cờ `has_transactions` của sản phẩm (catalog) và `costing_locked_at` của cửa hàng (setup). `backend-go.md` quy định tầng `app` quyết định transaction, repository không tự mở transaction, và mỗi transaction phải đặt `app.household_id` để RLS có hiệu lực.

Có hai cách phổ biến để repository dùng chung một transaction:

1. Đặt `pgx.Tx` vào `context.Context`, repository tự lấy ra.
2. Truyền tường minh các repository đã gắn với transaction vào hàm của use case (Unit of Work).

## Quyết định

Dùng **Unit of Work tường minh**.

- `internal/platform/db` có kiểu generic `TxRunner[R any]`, tạo bằng `db.NewTxRunner(pool, newRepos func(pgx.Tx) R)`.
- `InTx(ctx, householdID, fn func(ctx, repos R) error)`: mở transaction, gọi `set_config('app.household_id', ..., true)`, chạy `fn`, commit nếu `fn` trả `nil`, ngược lại rollback. `InReadTx` giống vậy nhưng `READ ONLY`. `InSystemTx` không đặt hộ, chỉ dùng cho bảng ngoài tenant (users, sessions).
- Mỗi context khai trong `app` của mình một struct `Repos` (gồm các port nó cần, kể cả port sang context khác) và một interface `TxRunner` nhận `Repos` đó. `*db.TxRunner[Repos]` tự thỏa interface này (interface ngầm của Go), không context nào import `platform/db`.
- `cmd/api/main.go` viết hàm `newRepos(tx pgx.Tx) Repos` dựng mọi adapter từ cùng một `tx`. Port sang context khác (ví dụ `CashRecorder` của inventory) là adapter gọi `app` của context kia với repository dựng từ chính `tx` đó.

```go
// internal/inventory/app
type Repos struct {
    Receipts ReceiptRepository
    Stock    StockRepository
    Cash     CashRecorder // adapter gọi cashbook, cùng tx
    Stores   StoreMarker  // adapter gọi setup, cùng tx
}

type TxRunner interface {
    InTx(ctx context.Context, householdID uuid.UUID, fn func(context.Context, Repos) error) error
}
```

## Hệ quả

- Chữ ký hàm cho thấy hết những gì use case đụng tới; không có phụ thuộc ẩn trong `context.Context`.
- Test use case dễ: truyền một `TxRunner` giả trả về `Repos` gồm các fake.
- Không thể quên đặt `app.household_id`: cách duy nhất để có repository là đi qua `InTx`/`InReadTx`, và hai hàm này từ chối `uuid.Nil`.
- Đọc dữ liệu cũng phải đi qua `InReadTx`, vì biến hộ chỉ tồn tại trong transaction.
- `main.go` dài thêm theo số context. Chấp nhận được: đó là nơi duy nhất biết các phần ráp với nhau thế nào.

## Phương án đã loại

- **`pgx.Tx` trong `context.Context`:** gọn hơn nhưng phụ thuộc ẩn. Quên mở transaction thì repository âm thầm chạy ngoài transaction và không có biến hộ, lỗi chỉ lộ ra lúc chạy.
