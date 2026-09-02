<template>
  <div>
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-4">
        <q-card class="bg-primary text-white">
          <q-card-section>
            <div class="text-h6">Kelas yang Diampu</div>
            <div class="text-h3">{{ kelasDiampu.length }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-secondary text-white">
          <q-card-section>
            <div class="text-h6">Total Siswa</div>
            <div class="text-h3">{{ totalSiswa }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-accent text-white">
          <q-card-section>
            <div class="text-h6">Deadline Input Nilai</div>
            <div class="text-h3">{{ deadlineCount }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-card>
          <q-card-section>
            <div class="text-h6">Jadwal Hari Ini</div>
            <q-list dense>
              <q-item v-for="jadwal in jadwalHariIni" :key="jadwal.id">
                <q-item-section>
                  <q-item-label>{{ jadwal.mapel_nama }} - Kelas {{ jadwal.kelas_nama }}</q-item-label>
                  <q-item-label caption>{{ jadwal.jam_mulai }} - {{ jadwal.jam_selesai }}</q-item-label>
                </q-item-section>
              </q-item>
              <q-item v-if="jadwalHariIni.length === 0">
                <q-item-section>Tidak ada jadwal hari ini</q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-6">
        <q-card>
          <q-card-section>
            <div class="text-h6">Tugas Perlu Dinilai</div>
            <q-list dense>
              <q-item v-for="tugas in tugasBelumDinilai" :key="tugas.id">
                <q-item-section>
                  <q-item-label>{{ tugas.nama_tugas }} - Kelas {{ tugas.kelas_nama }}</q-item-label>
                  <q-item-label caption>Deadline: {{ tugas.deadline }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-btn flat round dense icon="rate_review" color="primary" />
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from 'stores/auth'
import { useStrukturStore } from 'stores/kurikulum/struktur'
import { useSiswaStore } from 'stores/kesiswaan/siswa'
import { useNilaiStore } from 'stores/kurikulum/nilai'

const authStore = useAuthStore()
const strukturStore = useStrukturStore()
const siswaStore = useSiswaStore()
const nilaiStore = useNilaiStore()

const user = authStore.user
const kelasDiampu = ref([]) // dari strukturStore berdasarkan guru_id
const jadwalHariIni = ref([])
const tugasBelumDinilai = ref([])
const deadlineCount = ref(0)
const totalSiswa = ref(0)

function loadData() {
  if (!user) return
  // ambil kelas yang diampu guru berdasarkan struktur kurikulum
  const allStruktur = strukturStore.list
  kelasDiampu.value = allStruktur.filter(s => s.guru_id === user.id).map(s => s.kelas_nama).filter((v, i, a) => a.indexOf(v) === i)
  // total siswa dari kelas-kelas tersebut
  const siswaAll = siswaStore.list
  totalSiswa.value = siswaAll.filter(s => kelasDiampu.value.includes(s.kelas_nama)).length
  // jadwal hari ini (contoh: ambil dari struktur dengan hari ini)
  const today = new Date().toLocaleDateString('id-ID', { weekday: 'long' }).toLowerCase()
  jadwalHariIni.value = allStruktur.filter(s => s.guru_id === user.id && s.hari === today) // perlu field hari di struktur
  // deadline input nilai: ambil dari ujian yang belum dinilai
  // const nilaiSudah = nilaiStore.list.map(n => `${n.mapel_id}_${n.siswa_id}`)
  const tugas = [] // simulasi
  deadlineCount.value = tugas.length
}
onMounted(() => {
  strukturStore.loadData()
  siswaStore.loadData()
  nilaiStore.loadData()
  loadData()
})
</script>
