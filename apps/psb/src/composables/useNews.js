import { ref, onMounted } from 'vue'
// Ganti path sesuai struktur project Anda
import { NEWS } from '../mocks/data/schoolData'

export function useNews() {
  const newsList = ref([])
  const loading = ref(true)
  const error = ref(null)

  onMounted(async () => {
    try {
      // Simulasi async (ganti dengan fetch/axios saat siap)
      await new Promise((resolve) => setTimeout(resolve, 600))
      newsList.value = NEWS
    } catch (err) {
      error.value = `Gagal memuat berita. ${err}`
    } finally {
      loading.value = false
    }
  })

  return { newsList, loading, error }
}
