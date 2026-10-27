import { useContextStore } from '../../stores/contextStore.js'
/*
 * Axios Interceptor yang OTOMATIS meng-inject Operational Context
 * ke setiap request API.
 *
 * Ini mencegah developer lupa mengirim schoolId/periodId.
 * Backend (edge-bff) akan membaca header ini untuk scope query.
 */
export function setupContextInterceptor(axiosInstance) {
  axiosInstance.interceptors.request.use(
    (config) => {
      const store = useContextStore()
      if (store.hasContext) {
        const headers = store.current.toHeaders()
        config.headers['X-School-Id'] = headers['X-School-Id']
        config.headers['X-Academic-Period'] = headers['X-Academic-Period']
      }
      return config
    },
    (error) => {
      return Promise.reject(error)
    },
  )
  return axiosInstance
}
