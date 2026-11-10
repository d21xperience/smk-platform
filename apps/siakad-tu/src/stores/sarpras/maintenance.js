import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios'

export const useMaintenanceStore = defineStore('sarpras-maintenance', () => {
  const damageReports = ref([])
  const schedules = ref([])
  const activeAssets = ref([])
  const isLoading = ref(false)

  const fetchActiveAssets = async () => {
    try {
      const res = await api.get('/sarpras/kir/items')
      activeAssets.value = (res.data.data || []).filter((a) => a.status_aset === 'AKTIF')
    } catch (e) {
      console.error('Gagal mengambil data aset:', e)
    }
  }

  // === DAMAGE REPORTS (TIKET KERUSAKAN) ===
  const fetchDamageReports = async () => {
    isLoading.value = true
    try {
      const res = await api.get('/sarpras/damage-reports')
      damageReports.value = res.data.data || []
    } catch (e) {
      console.error('Gagal mengambil data tiket:', e)
    } finally {
      isLoading.value = false
    }
  }

  const createDamageReport = async (payload) => {
    const res = await api.post('/sarpras/damage-reports', payload)
    await fetchDamageReports()
    return res.data
  }

  const updateDamageReport = async (id, payload) => {
    const res = await api.put(`/sarpras/damage-reports/${id}`, payload)
    await fetchDamageReports()
    return res.data
  }

  // === MAINTENANCE SCHEDULES (JADWAL PERAWATAN) ===
  const fetchSchedules = async () => {
    isLoading.value = true
    try {
      const res = await api.get('/sarpras/maintenance-schedules')
      schedules.value = res.data.data || []
    } catch (e) {
      console.error('Gagal mengambil data jadwal:', e)
    } finally {
      isLoading.value = false
    }
  }

  const createSchedule = async (payload) => {
    const res = await api.post('/sarpras/maintenance-schedules', payload)
    await fetchSchedules()
    return res.data
  }

  const completeSchedule = async (id, payload) => {
    const res = await api.post(`/sarpras/maintenance-schedules/${id}/complete`, payload)
    await fetchSchedules()
    return res.data
  }

  return {
    damageReports,
    schedules,
    activeAssets,
    isLoading,
    fetchActiveAssets,
    fetchDamageReports,
    createDamageReport,
    updateDamageReport,
    fetchSchedules,
    createSchedule,
    completeSchedule,
  }
})
