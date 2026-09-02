import { ref } from 'vue'
import { pembayaranSPPService } from '@/services/pembayaranSPPService'

export function usePembayaranSPP() {
  const riwayat = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchRiwayatBySiswa = async (siswaId) => {
    loading.value = true
    try {
      riwayat.value = await pembayaranSPPService.getBySiswaId(siswaId)
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const simpanPembayaran = async (data) => {
    loading.value = true
    try {
      const baru = await pembayaranSPPService.create(data)
      // Jika sedang menampilkan riwayat siswa yang sama, tambahkan ke list
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

  return { riwayat, loading, error, fetchRiwayatBySiswa, simpanPembayaran }
}
