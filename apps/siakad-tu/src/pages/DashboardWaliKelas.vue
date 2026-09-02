<template>
  <div>
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-4">
        <q-card class="bg-primary text-white">
          <q-card-section>
            <div class="text-h6">Rata-rata Kehadiran</div>
            <div class="text-h3">{{ rataKehadiran }}%</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-warning text-black">
          <q-card-section>
            <div class="text-h6">Siswa At-Risk</div>
            <div class="text-h3">{{ siswaAtRisk.length }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-secondary text-white">
          <q-card-section>
            <div class="text-h6">Progres Nilai</div>
            <div class="text-h3">{{ progres }}%</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-card>
          <q-card-section>
            <div class="text-h6">Siswa dengan Alpha Terbanyak Bulan Ini</div>
            <q-table :rows="siswaAtRisk" :columns="riskColumns" dense flat />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-6">
        <q-card>
          <q-card-section>
            <div class="text-h6">Progres Pengisian Nilai</div>
            <div class="text-center q-pa-md">
              <q-circular-progress :value="progres" size="150px" :thickness="0.2" color="primary" track-color="grey-3"
                class="q-ma-md" show-value font-size="24px" />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from 'stores/auth'
import { useKelasStore } from '@/stores/kesiswaan/kelasStore'
import { useSiswaStore } from 'stores/kesiswaan/siswa'
import { useAbsensiStore } from '@/stores/kesiswaan/absensiStore'
import { useNilaiStore } from 'stores/kurikulum/nilai'
import { useStrukturStore } from 'stores/kurikulum/struktur'

const authStore = useAuthStore()
const kelasStore = useKelasStore()
const siswaStore = useSiswaStore()
const absensiStore = useAbsensiStore()
const nilaiStore = useNilaiStore()
const strukturStore = useStrukturStore()

const user = authStore.user
const kelasWali = ref(null)
const siswaAtRisk = ref([])
const rataKehadiran = ref(0)
const progres = ref(0)

const riskColumns = [
  { name: 'nama', label: 'Nama', field: 'nama' },
  { name: 'alpha', label: 'Jumlah Alpha', field: 'alpha' }
]

function loadData() {
  // asumsi wali kelas memiliki relasi ke kelas_id
  kelasWali.value = kelasStore.list.find(k => k.wali_kelas_id === user?.id) // perlu field wali_kelas_id
  if (!kelasWali.value) return
  const siswaDiKelas = siswaStore.list.filter(s => s.kelas_id === kelasWali.value.id)
  const bulanIni = `${new Date().getFullYear()}-${String(new Date().getMonth() + 1).padStart(2, '0')}`
  const absenBulanIni = absensiStore.list.filter(a => a.tanggal.startsWith(bulanIni) && siswaDiKelas.some(s => s.id === a.siswa_id))
  // const totalHari = [...new Set(absenBulanIni.map(a => a.tanggal))].length
  const hadirPerSiswa = {}
  siswaDiKelas.forEach(s => { hadirPerSiswa[s.id] = { hadir: 0, alpha: 0, total: 0 } })
  absenBulanIni.forEach(a => {
    if (hadirPerSiswa[a.siswa_id]) {
      hadirPerSiswa[a.siswa_id].total++
      if (a.status === 'hadir') hadirPerSiswa[a.siswa_id].hadir++
      else if (a.status === 'alpha') hadirPerSiswa[a.siswa_id].alpha++
    }
  })
  const totalHadir = Object.values(hadirPerSiswa).reduce((sum, h) => sum + h.hadir, 0)
  const totalKehadiran = Object.values(hadirPerSiswa).reduce((sum, h) => sum + h.total, 0)
  rataKehadiran.value = totalKehadiran === 0 ? 0 : ((totalHadir / totalKehadiran) * 100).toFixed(1)
  siswaAtRisk.value = Object.entries(hadirPerSiswa)
    .filter(([, h]) => h.alpha >= 3)
    .map(([id, h]) => ({ id, nama: siswaDiKelas.find(s => s.id == id)?.nama, alpha: h.alpha }))

  // Progres nilai: jumlah mapel yang sudah input nilai minimal 1 kali per siswa
  const mapelDiKelas = strukturStore.list.filter(s => s.kelas_id === kelasWali.value.id).map(s => s.mapel_id)
  const totalTarget = siswaDiKelas.length * mapelDiKelas.length
  const sudahDiinput = nilaiStore.list.filter(n => siswaDiKelas.some(s => s.id === n.siswa_id) && mapelDiKelas.includes(n.mapel_id)).length
  progres.value = totalTarget === 0 ? 0 : ((sudahDiinput / totalTarget) * 100).toFixed(1)
}
onMounted(() => {
  kelasStore.loadData()
  siswaStore.loadData()
  absensiStore.loadData()
  nilaiStore.loadData()
  strukturStore.loadData()
  loadData()
})
</script>
