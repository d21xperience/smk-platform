import { authRealAdapter } from './authRealAdapter.js'

export const authApi = authRealAdapter;

export function __forceMockAuthAdapter() {
  return authRealAdapter;
}
