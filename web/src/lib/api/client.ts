import createClient from 'openapi-fetch'

import type { components, paths } from './schema'

/**
 * Typed API client generated from api/openapi.yaml. Request and response
 * types come only from the spec: change the spec, run `make gen`, and
 * TypeScript shows every place that must follow.
 */
export const api = createClient<paths>({ baseUrl: '/api/v1' })

export type Problem = components['schemas']['Problem']

const fallbackMessage = 'Đã có lỗi xảy ra, vui lòng thử lại.'

/**
 * Error thrown by API helpers. It keeps the server's Problem Details so the UI
 * can show the Vietnamese `detail` and per-field `errors`, never technical text.
 */
export class ApiError extends Error {
  readonly problem: Problem | undefined

  constructor(problem: Problem | undefined) {
    super(problem?.detail ?? problem?.title ?? fallbackMessage)
    this.name = 'ApiError'
    this.problem = problem
  }
}

/** Returns the message to show a person for any error. */
export function errorMessage(err: unknown): string {
  if (err instanceof ApiError) return err.message
  return fallbackMessage
}

function isProblem(value: unknown): value is Problem {
  return typeof value === 'object' && value !== null && 'status' in value && 'title' in value
}

/**
 * Unwraps an openapi-fetch result: returns the data or throws ApiError. Use it
 * inside TanStack Query functions so failures land in the query's `error`.
 */
export function unwrap<T>(result: { data?: T; error?: unknown }): T {
  if (result.error !== undefined) {
    throw new ApiError(isProblem(result.error) ? result.error : undefined)
  }
  if (result.data === undefined) {
    throw new ApiError(undefined)
  }
  return result.data
}
