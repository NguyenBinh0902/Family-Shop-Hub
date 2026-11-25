# Loại hình kinh doanh và capability

Người dùng chỉ chọn loại hình; hệ thống suy ra capability. **Loại hình là dữ liệu cấu hình, capability là code** (ADR-0002). Thêm ngành mới chỉ cần thêm file cấu hình, trừ khi ngành đó cần capability chưa có.

## Capability MVP

| Mã | Tên | Mô tả | Phụ thuộc | Milestone |
| --- | --- | --- | --- | --- |
| `core.catalog` | Danh mục | Sản phẩm, nhóm, giá bán | (luôn bật) | M3 |
| `core.pos` | Bán tại quầy | Tạo đơn, thanh toán, in hóa đơn | `core.catalog` | M4 |
| `core.cashbook` | Thu chi | Phiếu thu chi | (luôn bật) | M5 |
| `core.customer` | Khách hàng, công nợ | Gắn khách vào đơn, bán nợ, thu nợ | `core.pos` | M5 |
| `inv.quantity` | Tồn theo số lượng | Nhập, kiểm kê, giá vốn bình quân | `core.catalog` | M3 |
| `inv.serial` | Tồn theo serial | Mỗi món có IMEI/số khung, giá vốn riêng | `core.catalog` | M6 |
| `serial.warranty` | Bảo hành | Thời hạn bảo hành theo món, tra cứu theo IMEI | `inv.serial` | M6 |
| `serial.deposit` | Đặt cọc | Nhận cọc giữ hàng, trừ cọc khi bán | `inv.serial`, `core.customer` | M6 |
| `fnb.options` | Tùy chọn món | Size, topping, mức đường; tùy chọn cộng giá | `core.pos` | M7 |
| `fnb.recipe` | Công thức | Bán món tự trừ nguyên liệu | `inv.quantity`, `fnb.options` | M7 |
| `fnb.table` | Đơn theo bàn | Mở đơn theo bàn, gộp, tách | `core.pos` | M7 |

Khi nạp cấu hình, hệ thống kiểm tra phụ thuộc; thiếu phụ thuộc thì không khởi động (lỗi cấu hình, không phải lỗi người dùng).

Một cửa hàng có thể có cả `inv.quantity` và `inv.serial` (điện thoại kèm phụ kiện), nhưng mỗi sản phẩm chỉ một kiểu theo dõi tồn (BR-07).

## Loại hình MVP

`core.*` = `core.catalog`, `core.pos`, `core.cashbook`, `core.customer`.

| Mã | Tên hiển thị | Nhóm | Capability | Bán âm mặc định |
| --- | --- | --- | --- | --- |
| `fnb.milk_tea` | Trà sữa, cà phê | F&B | core.*, inv.quantity, fnb.options, fnb.recipe | Cho phép |
| `fnb.restaurant` | Quán ăn | F&B | core.*, inv.quantity, fnb.options, fnb.recipe, fnb.table | Cho phép |
| `fnb.bakery` | Tiệm bánh | F&B | core.*, inv.quantity, fnb.options | Cho phép |
| `serial.phone` | Cửa hàng điện thoại | Serial | core.*, inv.serial, inv.quantity, serial.warranty, serial.deposit | Chặn |
| `serial.motorbike` | Cửa hàng xe máy | Serial | core.*, inv.serial, serial.warranty, serial.deposit | Chặn |

Trong lúc phát triển MVP, mọi loại hình ở bảng trên đều tạo được ngay từ M2. Capability chưa implement thì chưa có màn hình và API tương ứng; nó tự xuất hiện khi milestone của nó xong. Code luôn kiểm tra capability qua `setup` (không kiểm tra theo mã loại hình), nên không phải sửa gì khi capability mới hoàn thành.

Cờ `available: false` dùng để ẩn loại hình khỏi màn hình tạo cửa hàng (ví dụ loại hình sau MVP đã khai báo nhưng chưa phát hành).

## Lưu trữ cấu hình

- Mỗi version của loại hình là một file JSON: `internal/setup/storetypes/<code>/v<version>.json`, nhúng vào binary bằng `go:embed`.
- Không bao giờ sửa file của version đã phát hành; thay đổi phá vỡ thì thêm file version mới.
- Cửa hàng lưu `store_type_code` và `store_type_version`.

Ví dụ `internal/setup/storetypes/fnb.milk_tea/v1.json`:

```json
{
  "code": "fnb.milk_tea",
  "version": 1,
  "name": "Trà sữa, cà phê",
  "description": "Bán đồ uống có size, topping; tự trừ nguyên liệu theo công thức.",
  "group": "FNB",
  "available": true,
  "capabilities": ["core.catalog", "core.pos", "core.cashbook", "core.customer",
                   "inv.quantity", "fnb.options", "fnb.recipe"],
  "defaults": { "allowNegativeStock": true },
  "productAttributes": [
    { "key": "origin", "label": "Xuất xứ", "type": "text", "required": false }
  ],
  "seed": {
    "categories": ["Trà sữa", "Cà phê", "Topping", "Nguyên liệu"],
    "optionGroups": [
      { "name": "Size", "selection": "SINGLE", "required": true,
        "options": [ { "name": "M", "priceDelta": 0 }, { "name": "L", "priceDelta": 6000 } ] },
      { "name": "Mức đường", "selection": "SINGLE", "required": false,
        "options": [ { "name": "Ít đường", "priceDelta": 0 }, { "name": "Bình thường", "priceDelta": 0 } ] }
    ]
  }
}
```

Dữ liệu `seed` chỉ được nạp khi tạo cửa hàng; người dùng sửa thoải mái sau đó.

## Version

- Thay đổi bổ sung (thêm trường động không bắt buộc, thêm dữ liệu mẫu) có thể sửa trong cùng version **trước khi** có cửa hàng dùng version đó; sau đó phải tạo version mới.
- Thay đổi phá vỡ (bỏ capability, đổi kiểu trường) luôn tạo version mới. Cửa hàng cũ giữ version cũ cho đến khi Chủ hộ chủ động nâng cấp (sau MVP).
- Không đổi loại hình của cửa hàng đã có giao dịch (BR-04).

## Thuộc tính động

- Loại hình khai báo trường động cho sản phẩm trong `productAttributes`. Mỗi trường: `key` (snake_case tiếng Anh), `label` (tiếng Việt), `type` (`text`, `number`, `select`, `multiselect`, `date`), `required`, `options` (cho select).
- Giá trị lưu ở cột `attributes jsonb` của sản phẩm; server kiểm tra theo định nghĩa trước khi lưu; form phía client tự sinh từ định nghĩa (endpoint trả định nghĩa của cửa hàng).
- Key không bao giờ đổi tên. Bỏ trường thì đánh dấu `deprecated: true`; dữ liệu cũ vẫn đọc được.
- Trường cần tìm nhanh, duy nhất hoặc dùng cho báo cáo (IMEI, số khung, giá vốn, số lượng, giá bán) **không** là thuộc tính động.
- Báo cáo tổng hợp chỉ dựa trên trường lõi.
