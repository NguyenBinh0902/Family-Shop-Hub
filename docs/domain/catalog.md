# Context: Catalog (`internal/catalog`)

Sản phẩm, nhóm sản phẩm, tùy chọn món, công thức, thuộc tính động.

## Thực thể

**Category**: `id`, `household_id`, `store_id`, `name`, `sort_order`.

**Product**

| Trường | Kiểu | Ghi chú |
| --- | --- | --- |
| id, household_id, store_id | UUID | |
| category_id | UUID, nullable | |
| name | string | 1–120 ký tự, duy nhất trong cửa hàng (không phân biệt hoa thường) với sản phẩm `ACTIVE` |
| sku | string, nullable | Duy nhất trong cửa hàng nếu có |
| unit | string | Đơn vị tính hiển thị: ly, cái, chiếc, g, ml, kg... |
| sell_price | Money | ≥ 0; nguyên liệu có thể 0 |
| is_sellable | bool | `false` cho nguyên liệu: không hiện ở POS |
| stock_tracking | `NONE` \| `QUANTITY` \| `SERIAL` | BR-07 |
| min_stock | Quantity, nullable | Ngưỡng cảnh báo tồn thấp (chỉ `QUANTITY`) |
| warranty_months | int, nullable | Chỉ `SERIAL`, cửa hàng có `serial.warranty` |
| attributes | jsonb | Thuộc tính động theo loại hình |
| status | `ACTIVE` \| `DISCONTINUED` | |
| has_transactions | bool | BR-06 |

Thường gặp: đồ uống pha chế là `NONE` + có công thức; nguyên liệu là `QUANTITY` + `is_sellable=false`; bánh làm sẵn là `QUANTITY`; điện thoại là `SERIAL`; ốp lưng là `QUANTITY`.

**OptionGroup**: `id`, `household_id`, `store_id`, `name`, `selection` (`SINGLE` \| `MULTI`), `required` (bool), `max_select` (cho MULTI, nullable), `sort_order`.

**Option**: `id`, `option_group_id`, `name`, `price_delta` (Money ≥ 0), `ingredients` (danh sách `{product_id, quantity}`, có thể rỗng), `status`.

**ProductOptionGroup**: liên kết sản phẩm với nhóm tùy chọn, có `sort_order`.

**Recipe**: một công thức cho mỗi sản phẩm; dòng `{ingredient_product_id, quantity}`. Nguyên liệu phải là sản phẩm `QUANTITY` cùng cửa hàng, và bản thân không có công thức (không lồng).

## Bất biến

- BR-06, BR-07.
- Đổi giá bán, công thức, tùy chọn chỉ ảnh hưởng đơn tạo sau đó (đơn cũ có snapshot).
- Đổi giá bán ghi nhật ký (BR-26).

## Cung cấp cho context khác

- Sales: lấy thông tin sản phẩm để tạo dòng đơn (giá, tùy chọn hợp lệ, công thức đã "trải phẳng" thành danh sách nguyên liệu × số lượng).
- Inventory: kiểu theo dõi tồn, đơn vị.

## Use case

UC-09, UC-10, UC-11.
