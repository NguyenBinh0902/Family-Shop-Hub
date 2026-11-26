import { createBrowserRouter } from 'react-router'

import { HomePage } from '@/routes/HomePage'

import { AppLayout } from './AppLayout'

/**
 * All screens and their URLs. Store-level screens will live under
 * /stores/:storeId/... so the current store is always in the URL.
 */
export const router = createBrowserRouter([
  {
    element: <AppLayout />,
    children: [{ index: true, element: <HomePage /> }],
  },
])
