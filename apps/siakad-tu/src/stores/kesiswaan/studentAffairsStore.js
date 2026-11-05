import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios' // Gunakan axios untuk memanggil API backend

export const useStudentAffairsStore = defineStore('studentAffairs', () => {
  const students = ref([])
  const mutations = ref([])
  const dashboardData = ref(null)
  const loading = ref(false)
  const error = ref(null)

  const filters = ref({
    search: '',
    status: '',
    jurusan: 'Semua',
    tingkat: 'Semua',
  })

  // Actions: Dashboard (FIX: Memanggil API Statistik Backend)
  async function loadDashboard() {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/statistics')
      const stats = response.data

      // Mapping format backend ke format yang diharapkan UI (DashboardPage.vue)
      dashboardData.value = {
        metrics: {
          totalSiswa: stats.total_siswa,
          totalGuru: stats.total_guru,
          totalRombel: stats.total_rombel,
          totalJurusan: stats.total_jurusan,
        },
        // Format untuk ApexCharts
        chartSeries: [
          {
            name: 'Jumlah Siswa',
            data: stats.distribusi_jurusan.map((j) => j.total),
          },
        ],
        chartCategories: stats.distribusi_jurusan.map((j) => j.jurusan),
        reminders: [], // Kosongkan sementara (bisa diisi nanti)
      }
    } catch (err) {
      error.value = err.message
      console.error('Gagal memuat dashboard:', err)
    } finally {
      loading.value = false
    }
  }

  // === Actions: Students (Tetap seperti sebelumnya) ===
  async function loadStudents() {
    loading.value = true
    error.value = null
    try {
      // Jika API student belum diintegrasikan, tetap gunakan mock/kosongkan
      students.value = []
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  function setFilter(key, value) {
    filters.value[key] = value
    loadStudents()
  }

  function clearFilters() {
    filters.value = { search: '', status: '', jurusan: 'Semua', tingkat: 'Semua' }
    loadStudents()
  }

  async function loadMutations() {
    mutations.value = []
  }

  async function approveMutationAction() {
    // Placeholder
  }

  async function rejectMutationAction() {
    // Placeholder
  }

  return {
    students,
    mutations,
    dashboardData,
    loading,
    error,
    filters,
    loadDashboard,
    loadStudents,
    loadMutations,
    setFilter,
    clearFilters,
    approveMutationAction,
    rejectMutationAction,
  }
})
