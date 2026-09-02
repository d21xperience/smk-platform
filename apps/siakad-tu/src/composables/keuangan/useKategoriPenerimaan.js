import { ref } from 'vue'
import { kategoriPenerimaanService } from '@/services/kategoriPenerimaanService'

export function useKategoriPenerimaan() {
  const kategoriList = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchAll = async () => {
    loading.value = true
    try {
      kategoriList.value = await kategoriPenerimaanService.getAll()
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const tambahKategori = async (data) => {
    loading.value = true
    try {
      const baru = await kategoriPenerimaanService.create(data)
      kategoriList.value.push(baru)
      return baru
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateKategori = async (id, data) => {
    loading.value = true
    try {
      const updated = await kategoriPenerimaanService.update(id, data)
      const index = kategoriList.value.findIndex((k) => k.id === id)
      if (index !== -1) kategoriList.value[index] = updated
      return updated
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const hapusKategori = async (id) => {
    loading.value = true
    try {
      await kategoriPenerimaanService.delete(id)
      kategoriList.value = kategoriList.value.filter((k) => k.id !== id)
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  return { kategoriList, loading, error, fetchAll, tambahKategori, updateKategori, hapusKategori }
}
