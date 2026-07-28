import { useQuasar } from "quasar";

/**
 * useNotification - Shared composable untuk notifikasi
 *
 * @example
 * const { success, error, warning, info } = useNotification()
 * success('Data berhasil disimpan!')
 * error('Gagal menyimpan data')
 */
export function useNotification() {
  const $q = useQuasar();

  const success = (message, options = {}) => {
    $q.notify({
      type: "positive",
      message,
      position: options.position || "top",
      timeout: options.timeout || 3000,
      ...options,
    });
  };

  const error = (message, options = {}) => {
    $q.notify({
      type: "negative",
      message,
      position: options.position || "top",
      timeout: options.timeout || 5000,
      ...options,
    });
  };

  const warning = (message, options = {}) => {
    $q.notify({
      type: "warning",
      message,
      position: options.position || "top",
      timeout: options.timeout || 4000,
      ...options,
    });
  };

  const info = (message, options = {}) => {
    $q.notify({
      type: "info",
      message,
      position: options.position || "top",
      timeout: options.timeout || 3000,
      ...options,
    });
  };

  return {
    success,
    error,
    warning,
    info,
  };
}
