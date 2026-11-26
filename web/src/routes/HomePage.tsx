import { Button } from '@/components/ui/button'
import { useHealth } from '@/features/system/useHealth'
import { errorMessage } from '@/lib/api/client'

export function HomePage() {
  const health = useHealth()

  return (
    <section className="mx-auto flex max-w-md flex-col gap-4 p-4">
      <h1 className="text-xl font-semibold">Chào mừng đến FamilyShop Hub</h1>

      <div className="rounded-lg border p-4" role="status" aria-live="polite">
        {health.isPending && <p className="text-muted-foreground">Đang kiểm tra kết nối…</p>}

        {health.isSuccess && (
          <p className="flex items-center gap-2">
            <span className="size-2.5 rounded-full bg-green-600" aria-hidden />
            Kết nối máy chủ: OK
          </p>
        )}

        {health.isError && (
          <div className="flex flex-col gap-3">
            <p className="flex items-center gap-2">
              <span className="bg-destructive size-2.5 rounded-full" aria-hidden />
              Không kết nối được máy chủ
            </p>
            <p className="text-muted-foreground text-sm">{errorMessage(health.error)}</p>
            <Button
              className="min-h-11 self-start"
              variant="outline"
              onClick={() => void health.refetch()}
              disabled={health.isFetching}
            >
              Thử lại
            </Button>
          </div>
        )}
      </div>
    </section>
  )
}
