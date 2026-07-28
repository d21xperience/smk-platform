// src/composables/siakad/useSchedule.js
import { ref } from 'vue'
import { useApi } from '../useApi'

export function useSchedule() {
  const { api, request, loading } = useApi()
  const schedule = ref([])

  /**
   * Mengambil jadwal mengajar guru untuk periode aktif
   * GET /api/v1/siakad/schedule/teacher
   * (Parameter 'periode' otomatis di-inject oleh axios interceptor)
   */
  const fetchMySchedule = async () => {
    const data = await request(() => api.get('/siakad/schedule/teacher'))
    if (data) {
      // Kita group data berdasarkan Hari agar mudah dirender di UI
      schedule.value = groupByDay(data)
    }
  }

  // Helper function untuk mengelompokkan data berdasarkan hari
  const groupByDay = (data) => {
    const daysOrder = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']
    const grouped = {}

    // Inisialisasi semua hari (agar tetap muncul meski kosong)
    daysOrder.forEach((day) => {
      grouped[day] = []
    })

    data.forEach((item) => {
      if (grouped[item.day]) {
        grouped[item.day].push(item)
      }
    })

    // Urutkan jadwal dalam satu hari berdasarkan jam mulai
    Object.keys(grouped).forEach((day) => {
      grouped[day].sort((a, b) => a.start_time.localeCompare(b.start_time))
    })

    return grouped
  }

  return {
    loading,
    schedule,
    fetchMySchedule,
  }
}
