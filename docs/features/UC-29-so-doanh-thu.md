# UC-29 Sổ doanh thu theo chủ thể kinh doanh

- **Milestone:** M8 · **Context:** reporting · **Tác nhân:** Chủ hộ
- **Quy tắc:** BR-15, BR-24, BR-29, TAX-01, TAX-02, TAX-06 · **Quyền:** `report.household`

## Luồng chính
1. Chọn chủ thể kinh doanh và kỳ.
2. Hiển thị danh sách chứng từ theo thời gian: mỗi đơn hoàn tất là một dòng dương, mỗi bút toán đảo là một dòng âm. Cột: ngày giờ, cửa hàng, số đơn, loại (bán / hủy), tổng tiền, các hình thức thanh toán (tiền mặt, chuyển khoản kèm tài khoản, ví kèm ví, nợ, cọc).
3. Tổng theo kỳ; tổng theo hình thức thanh toán; tổng theo tài khoản ngân hàng/ví.
4. Doanh thu lũy kế từ đầu năm của chủ thể; cảnh báo khi đạt 80% ngưỡng (TAX-06).
5. Xuất CSV (UTF-8 có BOM để Excel đọc đúng tiếng Việt).

## CHƯA CHỐT – không implement
Tạo tờ khai thuế, cách gom doanh thu nhiều địa điểm khi kê khai. Chỉ làm phần liệt kê và xuất dữ liệu ở trên. Giá trị ngưỡng là cấu hình, chưa phát hành tính năng cảnh báo cho người dùng thật trước khi kế toán xác nhận.

## Tiêu chí chấp nhận
- **Given** đơn 500.000 ngày 01/12 bị hủy ngày 03/12 **Then** sổ tháng 12 có hai dòng: +500.000 ngày 01/12 và −500.000 ngày 03/12.
- **Given** chủ thể X có hai cửa hàng **Then** sổ của X gồm chứng từ của cả hai, cột cửa hàng ghi rõ.
- **Given** lũy kế năm 820.000.000 với ngưỡng 1.000.000.000 **Then** hiện cảnh báo.
- File CSV mở bằng Excel hiển thị đúng dấu tiếng Việt; số tiền không có dấu phân cách hàng nghìn.

## API gợi ý
`GET /api/v1/households/{householdId}/business-entities/{id}/sales-ledger?from=&to=` · `.../sales-ledger.csv?from=&to=`
