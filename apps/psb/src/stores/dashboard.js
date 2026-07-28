import { defineStore } from 'pinia'
import { useSiswaStore } from './kesiswaan/siswa'
import { useGuruStore } from './kepegawaian/guru'
import { useKelasStore } from './kesiswaan/kelas'
import { useAbsensiStore } from './kesiswaan/absensi'
import { useNilaiStore } from './kurikulum/nilai'
import { useMapelStore } from './akademik/mapel'
// import { useStrukturStore } from './kurikulum/struktur'

export const useDashboardStore = defineStore('dashboard', {
  actions: {
    getTotalSiswa() {
      const siswaStore = useSiswaStore()
      return siswaStore.list.length
    },
    getTotalGuru() {
      const guruStore = useGuruStore()
      return guruStore.list.length
    },
    getTotalKelas() {
      const kelasStore = useKelasStore()
      return kelasStore.list.length
    },
    getRataRataNilaiKeseluruhan() {
      const nilaiStore = useNilaiStore()
      const nilaiList = nilaiStore.list.filter((n) => n.nilai_akhir)
      if (nilaiList.length === 0) return 0
      const total = nilaiList.reduce((sum, n) => sum + n.nilai_akhir, 0)
      return (total / nilaiList.length).toFixed(1)
    },
    getPersentaseKehadiranHariIni() {
      const absensiStore = useAbsensiStore()
      const today = new Date().toISOString().split('T')[0]
      const absenHariIni = absensiStore.list.filter((a) => a.tanggal === today)
      const hadir = absenHariIni.filter((a) => a.status === 'hadir').length
      const total = absenHariIni.length
      return total === 0 ? 0 : ((hadir / total) * 100).toFixed(1)
    },
    getGrafikNilaiPerMapel() {
      const nilaiStore = useNilaiStore()
      const mapelStore = useMapelStore()
      const result = {}
      nilaiStore.list.forEach((n) => {
        const mapel = mapelStore.list.find((m) => m.id === n.mapel_id)
        if (mapel) {
          if (!result[mapel.nama]) result[mapel.nama] = { total: 0, count: 0 }
          result[mapel.nama].total += n.nilai_akhir
          result[mapel.nama].count++
        }
      })
      return Object.entries(result).map(([mapel, { total, count }]) => ({
        mapel,
        rata: (total / count).toFixed(1),
      }))
    },
    // getUjianTerdekat() {
    //   // ambil dari store ujian jika ada
    //   const ujianStore = useUjianStore?.() // asumsi ada
    //   if (ujianStore && ujianStore.list) {
    //     const now = new Date()
    //     return ujianStore.list.filter(u => new Date(u.tanggal_mulai) > now).sort((a,b) => new Date(a.tanggal_mulai) - new Date(b.tanggal_mulai)).slice(0,5)
    //   }
    //   return []
    // }
  },
})
