import { useQuery } from '@tanstack/react-query'

import { api, unwrap } from '@/lib/api/client'

/** Asks the server whether it and its database are reachable. */
export function useHealth() {
  return useQuery({
    queryKey: ['healthz'],
    queryFn: async ({ signal }) => unwrap(await api.GET('/healthz', { signal })),
    retry: false,
  })
}
