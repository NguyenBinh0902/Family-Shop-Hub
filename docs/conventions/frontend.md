# Quy ước frontend

Chủ repo không chuyên frontend. Agent tự quyết các chi tiết triển khai trong khuôn khổ dưới đây, nhưng **phải hỏi trước** các quyết định ở mục 6. Khi tóm tắt thay đổi, giải thích bằng ngôn ngữ dễ hiểu với người làm backend.

## 1. Stack

React + Vite + TypeScript (`strict: true`) · React Router · TanStack Query · Tailwind CSS · shadcn/ui · react-hook-form + zod · openapi-typescript + openapi-fetch (client sinh từ `api/openapi.yaml`) · vite-plugin-pwa · Playwright.

## 2. Bố cục

```
web/src/
  app/              khởi tạo: router, QueryClient, layout gốc
  routes/           mỗi màn hình một file, map với URL
  features/<context>/   component, hook, schema của từng nghiệp vụ (sales, catalog...)
  components/ui/    component shadcn/ui
  lib/api/          client sinh tự động + hàm gọi API dùng chung
  lib/format/       định dạng tiền, ngày, số điện thoại
```

## 3. Quy tắc giao diện

- Mobile-first: thiết kế cho màn hình 360 px trước, sau đó mới máy tính bảng và desktop.
- Vùng bấm tối thiểu 44 × 44 px. Màn hình POS ưu tiên nút lớn, ít chữ, thao tác một tay.
- Ngữ cảnh cửa hàng nằm trong URL: `/stores/:storeId/...`. Tên và màu cửa hàng hiển thị cố định ở đầu màn hình.
- Toàn bộ chữ tiếng Việt, từ ngữ đời thường ("Bán hàng", "Nhập hàng", "Khách nợ"), không dùng thuật ngữ kế toán.
- Tiền hiển thị bằng `Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' })`; ngày theo `dd/MM/yyyy`, múi giờ `Asia/Ho_Chi_Minh`.
- Hóa đơn in qua CSS `@media print` với hai khổ 58 mm và 80 mm.

## 4. Dữ liệu

- Mọi gọi API qua client sinh từ OpenAPI. Không tự viết kiểu dữ liệu trùng với API.
- Dữ liệu server quản lý bằng TanStack Query; không chép vào state toàn cục.
- Tạo đơn: sinh mã đơn (UUID v7) ở client trước khi gửi, gửi lại cùng mã khi mạng lỗi (BR-17). Khóa nút thanh toán khi đang gửi.
- Lỗi từ server ở dạng Problem Details: hiển thị trường `detail`/thông điệp tiếng Việt, không hiển thị lỗi kỹ thuật.
- Form kiểm tra bằng zod; server vẫn là nơi kiểm tra cuối cùng.

## 5. Chất lượng

- ESLint + Prettier sạch trước khi commit; TypeScript không dùng `any`.
- Luồng chính (đăng nhập, tạo cửa hàng, bán tại quầy, hủy đơn) có test Playwright.
- Kiểm tra mỗi màn hình mới ở độ rộng 360 px.

## 6. Cần hỏi chủ repo trước khi làm

Đưa ra 1–2 phương án, mỗi phương án một câu lý do, rồi chờ trả lời:

- Thêm thư viện npm ngoài danh sách ở mục 1.
- Bố cục và luồng của một màn hình mới (đặc biệt POS, tạo cửa hàng, báo cáo).
- Thay đổi điều hướng: menu, bộ chuyển cửa hàng, cách vào góc nhìn Toàn hộ.
- Màu sắc, font, phong cách chung của ứng dụng.
- Bất cứ chỗ nào làm khác luồng đã mô tả trong use case.

Không cần hỏi: chi tiết căn chỉnh, khoảng cách, tách component, đặt tên file, cách viết hook.
