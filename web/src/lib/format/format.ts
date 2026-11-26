/** Business time zone: every date shown to people uses Vietnam time. */
export const TIME_ZONE = 'Asia/Ho_Chi_Minh'

const moneyFormat = new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' })

/** Formats an amount in đồng: 25000 → "25.000 ₫". Amounts are integers. */
export function formatMoney(amount: number): string {
  return moneyFormat.format(amount)
}

const dateFormat = new Intl.DateTimeFormat('vi-VN', {
  timeZone: TIME_ZONE,
  day: '2-digit',
  month: '2-digit',
  year: 'numeric',
})

/** Formats an instant as dd/MM/yyyy in Vietnam time. */
export function formatDate(value: Date | string): string {
  return dateFormat.format(typeof value === 'string' ? new Date(value) : value)
}

const dateTimeFormat = new Intl.DateTimeFormat('vi-VN', {
  timeZone: TIME_ZONE,
  day: '2-digit',
  month: '2-digit',
  year: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
  hour12: false,
})

/** Formats an instant as "HH:mm dd/MM/yyyy" in Vietnam time. */
export function formatDateTime(value: Date | string): string {
  return dateTimeFormat.format(typeof value === 'string' ? new Date(value) : value)
}
