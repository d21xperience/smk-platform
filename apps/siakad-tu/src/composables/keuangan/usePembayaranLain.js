import { ref } from 'vue'
import { pembayaranLainService } from '@/services/pembayaranLainService'

export function usePembayaranLain() {
  const riwayat = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchRiwayatBySiswa = async (siswaId) => {
    loading.value = true
    try {
      riwayat.value = await pembayaranLainService.getBySiswaId(siswaId)
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const simpanPembayaran = async (data) => {
    loading.value = true
    try {
      const baru = await pembayaranLainService.create(data)
      if (riwayat.value.length && riwayat.value[0]?.siswaId === data.siswaId) {
        riwayat.value.push(baru)
      }
      return baru
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const hapusPembayaran = async (id) => {
    loading.value = true
    try {
      await pembayaranLainService.delete(id)
      riwayat.value = riwayat.value.filter((p) => p.id !== id)
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  return { riwayat, loading, error, fetchRiwayatBySiswa, simpanPembayaran, hapusPembayaran }
}
