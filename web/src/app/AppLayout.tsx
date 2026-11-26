import { Outlet } from 'react-router'

/**
 * Root layout: a thin header with the app name. M2 puts the store switcher
 * in the header (store name and color always visible, frontend.md §3).
 */
export function AppLayout() {
  return (
    <div className="bg-background text-foreground flex min-h-dvh flex-col">
      <header className="flex h-14 items-center gap-3 border-b px-4">
        <span className="font-semibold">FamilyShop Hub</span>
        {/* Store switcher goes here (M2). */}
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  )
}
