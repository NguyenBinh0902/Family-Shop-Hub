# Context: Inventory (`internal/inventory`)

Nhà cung cấp, nhập hàng, tồn kho theo số lượng, đơn vị serial, giá vốn, kiểm kê.

## Thực thể

**Supplier** (cấp hộ): `id`, `household_id`, `name`, `phone` (nullable), `address`, `note`.

**PurchaseReceipt**: `id`, `household_id`, `store_id`, `supplier_id` (nullable), `received_at`, `created_by`, `note`, `total` (Money), `payment_method` (`CASH` \| `BANK_TRANSFER` \| `EWALLET`), `payment_account_id` (nullable). Bất biến sau khi ghi.

**PurchaseReceiptLine**: `product_id`, `quantity`, `unit_cost`, `line_total`; với hàng serial là danh sách `{serial, unit_cost}` (mỗi serial quantity = 1).

**StockLevel**: khóa `(store_id, product_id)`; `on_hand` (Quantity, có thể âm nếu cửa hàng cho bán âm), `avg_cost` (Money, nullable khi chưa từng nhập), `updated_at`.

**StockMovement** (chỉ thêm, BR-37): `id`, `household_id`, `store_id`, `product_id`, `type`, `quantity_delta` (Quantity, có dấu), `unit_cost` (Money), `ref_type`, `ref_id`, `occurred_at`, `created_by`.

| type | Khi nào | Dấu |
| --- | --- | --- |
| `RECEIPT` | Nhập hàng | + |
| `SALE` | Bán sản phẩm `QUANTITY` | − |
| `RECIPE_CONSUMPTION` | Trừ nguyên liệu khi bán món | − |
| `SALE_REVERSAL` | Hủy đơn, hoàn kho | + |
| `ADJUSTMENT` | Kiểm kê | ± |

**SerialUnit**: `id`, `household_id`, `store_id`, `product_id`, `serial`, `unit_cost`, `status`, `receipt_id`, `sold_order_line_id`, `sold_at`, `warranty_until` (date), `reserved_deposit_id`.

**StockCount**: `id`, `store_id`, `counted_at`, `created_by`, `reason`; dòng `{product_id, system_qty, counted_qty, delta}`.

## Trạng thái SerialUnit

```
IN_STOCK → RESERVED (nhận cọc) → IN_STOCK (hủy cọc)
IN_STOCK → SOLD (bán)          RESERVED → SOLD (bán từ cọc)
SOLD → IN_STOCK (hủy đơn)
RETURNED: để dành cho trả hàng sau MVP
```

Duy nhất `(household_id, serial)` với các đơn vị `IN_STOCK`, `RESERVED` (BR-11) — dùng partial unique index. Serial chuẩn hóa: bỏ khoảng trắng, chữ hoa.

## Giá vốn bình quân (BR-09, BR-28, BR-33)

Khi nhập `q` đơn vị giá `c` vào sản phẩm đang tồn `Q` với giá vốn `C`:

- Nếu `Q ≤ 0` hoặc `C` rỗng: `C_mới = c`.
- Ngược lại: `C_mới = round((Q × C + q × c) / (Q + q))`, làm tròn tới đồng, nửa lên.

Ví dụ: tồn 10 kg giá 50.000 đ/kg, nhập 5 kg giá 56.000 đ/kg → (500.000 + 280.000) / 15 = 52.000 đ/kg.

Ví dụ phần lẻ: tồn 3 ly giá 10.000, nhập 4 ly giá 11.000 → 74.000 / 7 = 10.571,43 → 10.571 đ.

Bán và hoàn kho không thay đổi giá vốn bình quân. Hoàn kho khi hủy đơn dùng giá vốn snapshot của dòng đơn làm `unit_cost` của movement.

`Quantity` là phần nghìn: phép nhân `Quantity × Money` phải chia 1000 sau cùng và làm tròn một lần (đặt trong `kernel`, có test).

## Cung cấp cho Sales (trong cùng transaction)

```go
type StockService interface {
    // Trừ kho cho các dòng QUANTITY và nguyên liệu; trả giá vốn đơn vị cho từng dòng.
    Consume(ctx, scope, lines []ConsumeLine) ([]ConsumeResult, error)
    // Chuyển serial sang SOLD, trả unit_cost và warranty_until.
    SellSerial(ctx, scope, serialUnitID, orderLineID, soldAt) (SerialSale, error)
    // Hoàn kho khi hủy đơn.
    Restore(ctx, scope, reversal ReversalLines) error
    ReserveSerial / ReleaseSerial (cho cọc)
}
```

Không đủ tồn và cửa hàng chặn bán âm → lỗi `BR-13`.

## Use case

UC-12, UC-13, UC-14, UC-15.
