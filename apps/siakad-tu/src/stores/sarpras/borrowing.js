import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios'

export const useBorrowingStore = defineStore('sarpras-borrowing', () => {
  const borrowings = ref([])
  const activeAssets = ref([])
  const isLoading = ref(false)

  const fetchBorrowings = async () => {
    isLoading.value = true
    try {
      const res = await api.get('/sarpras/borrowings')
      borrowings.value = res.data.data || []
    } catch (e) {
      console.error('Gagal mengambil data peminjaman:', e)
    } finally {
      isLoading.value = false
    }
  }

  const fetchActiveAssets = async () => {
    try {
      // Ambil dari endpoint KIR (karena sudah menyediakan semua aset)
      const res = await api.get('/sarpras/kir/items')
      // Filter hanya yang berstatus AKTIF
      activeAssets.value = (res.data.data || []).filter((a) => a.status_aset === 'AKTIF')
    } catch (e) {
      console.error('Gagal mengambil data aset:', e)
    }
  }

  const createBorrowing = async (payload) => {
    const res = await api.post('/sarpras/borrowings', payload)
    await fetchBorrowings()
    return res.data
  }

  const approveBorrowing = async (id, payload) => {
    const res = await api.post(`/sarpras/borrowings/${id}/approve`, payload)
    await fetchBorrowings()
    return res.data
  }

  const pickupBorrowing = async (id, payload) => {
    const res = await api.post(`/sarpras/borrowings/${id}/pickup`, payload)
    await fetchBorrowings()
    return res.data
  }

  const returnBorrowing = async (id, payload) => {
    const res = await api.post(`/sarpras/borrowings/${id}/return`, payload)
    await fetchBorrowings()
    return res.data
  }

  return {
    borrowings,
    activeAssets,
    isLoading,
    fetchBorrowings,
    fetchActiveAssets,
    createBorrowing,
    approveBorrowing,
    pickupBorrowing,
    returnBorrowing,
  }
})
