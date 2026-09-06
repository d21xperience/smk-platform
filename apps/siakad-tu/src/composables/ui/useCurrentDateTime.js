// apps/siakad-tu/src/composables/ui/useCurrentDateTime.js

import { ref, onMounted, onUnmounted } from 'vue'

/**
 * Composable untuk mendapatkan tanggal dan waktu saat ini secara real-time.
 *
 * @returns {Object} { currentDate, currentTime }
 */
export function useCurrentDateTime() {
  const currentDate = ref('')
  const currentTime = ref('')

  let intervalId = null

  const updateDateTime = () => {
    const now = new Date()

    // Format Tanggal: Senin, 15 Januari 2024
    currentDate.value = now.toLocaleDateString('id-ID', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    })

    // Format Waktu: 14:30:00
    currentTime.value = now.toLocaleTimeString('id-ID', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  }

  onMounted(() => {
    updateDateTime() // Panggil sekali saat mount agar tidak ada delay 1 detik
    intervalId = setInterval(updateDateTime, 1000) // Update setiap 1 detik
  })

  onUnmounted(() => {
    // Bersihkan interval saat komponen dihancurkan untuk mencegah memory leak
    if (intervalId) {
      clearInterval(intervalId)
    }
  })

  return {
    currentDate,
    currentTime,
  }
}
