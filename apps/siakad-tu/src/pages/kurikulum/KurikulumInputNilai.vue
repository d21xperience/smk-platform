<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Input Nilai (Harian, UTS, UAS, Praktik, Portofolio)</div>

    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-3">
        <q-select v-model="selectedKelas" :options="kelasOptions" label="Kelas" outlined dense />
      </div>
      <div class="col-12 col-md-3">
        <q-select v-model="selectedMapel" :options="mapelOptions" label="Mata Pelajaran" outlined dense />
      </div>
      <div class="col-12 col-md-2">
        <q-select v-model="selectedSemester" :options="['Ganjil', 'Genap']" label="Semester" outlined dense />
      </div>
      <div class="col-12 col-md-2">
        <q-input v-model="tahunAjaran" label="Tahun Ajaran" outlined dense />
      </div>
      <div class="col-12 col-md-2">
        <q-btn color="primary" label="Tampilkan" @click="loadData" />
      </div>
    </div>

    <div v-if="siswaNilaiList.length">
      <q-table :rows="siswaNilaiList" :columns="columns" row-key="id" flat bordered dense>
        <template v-slot:body-cell-nilai_harian_avg="props">
          <q-td :props="props">
            <q-input dense type="number" step="0.5" v-model="props.row.nilai_harian_avg"
              @blur="updateNilai(props.row, 'nilai_harian_avg', $event)" style="width: 80px" />
          </q-td>
        </template>
        <template v-slot:body-cell-nilai_uts="props">
          <q-td :props="props">
            <q-input dense type="number" step="0.5" v-model="props.row.nilai_uts"
              @blur="updateNilai(props.row, 'nilai_uts', $event)" style="width: 80px" />
          </q-td>
        </template>
        <template v-slot:body-cell-nilai_uas="props">
          <q-td :props="props">
            <q-input dense type="number" step="0.5" v-model="props.row.nilai_uas"
              @blur="updateNilai(props.row, 'nilai_uas', $event)" style="width: 80px" />
          </q-td>
        </template>
        <template v-slot:body-cell-nilai_praktik="props">
          <q-td :props="props">
            <q-input dense type="number" step="0.5" v-model="props.row.nilai_praktik"
              @blur="updateNilai(props.row, 'nilai_praktik', $event)" style="width: 80px" />
          </q-td>
        </template>
        <template v-slot:body-cell-nilai_portofolio="props">
          <q-td :props="props">
            <q-input dense type="number" step="0.5" v-model="props.row.nilai_portofolio"
              @blur="updateNilai(props.row, 'nilai_portofolio', $event)" style="width: 80px" />
          </q-td>
        </template>
        <template v-slot:body-cell-nilai_akhir="props">
          <q-td :props="props">
            <strong>{{ props.row.nilai_akhir }}</strong>
          </q-td>
        </template>
      </q-table>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useNilaiStore } from '@/stores/kurikulum/nilai'
import { useKelasStore } from '@/stores/kesiswaan/kelasStore'
import { useMapelStore } from '@/stores/akademik/mapel'
import { useSiswaStore } from '@/stores/kesiswaan/siswa'

const $q = useQuasar()
const nilaiStore = useNilaiStore()
const kelasStore = useKelasStore()
const mapelStore = useMapelStore()
const siswaStore = useSiswaStore()

const selectedKelas = ref(null)
const selectedMapel = ref(null)
const selectedSemester = ref('Ganjil')
const tahunAjaran = ref('2025/2026')
const siswaNilaiList = ref([])

const kelasOptions = computed(() => kelasStore.list)
const mapelOptions = computed(() => mapelStore.list)

const columns = [
  { name: 'nama', label: 'Siswa', field: 'nama' },
  { name: 'nis', label: 'NIS', field: 'nis' },
  { name: 'nilai_harian_avg', label: 'Rata Harian', field: 'nilai_harian_avg' },
  { name: 'nilai_uts', label: 'UTS', field: 'nilai_uts' },
  { name: 'nilai_uas', label: 'UAS', field: 'nilai_uas' },
  { name: 'nilai_praktik', label: 'Praktik', field: 'nilai_praktik' },
  { name: 'nilai_portofolio', label: 'Portofolio', field: 'nilai_portofolio' },
  { name: 'nilai_akhir', label: 'Nilai Akhir', field: 'nilai_akhir' }
]

function loadData() {
  if (!selectedKelas.value || !selectedMapel.value) {
    $q.notify({ type: 'warning', message: 'Pilih kelas dan mata pelajaran' })
    return
  }
  const siswaDiKelas = siswaStore.list.filter(s => s.kelas_id === selectedKelas.value.id)
  const list = siswaDiKelas.map(s => {
    let nilai = nilaiStore.getNilaiBySiswaMapel(s.id, selectedMapel.value.id, selectedSemester.value, tahunAjaran.value)
    if (!nilai) {
      nilai = nilaiStore.tambahNilai(s.id, selectedMapel.value.id, selectedKelas.value.id, selectedSemester.value, tahunAjaran.value)
    }
    return { ...s, ...nilai }
  })
  siswaNilaiList.value = list
}

function updateNilai(row, field, event) {
  const value = parseFloat(event.target.value)
  if (isNaN(value)) return
  const nilaiEntry = nilaiStore.getNilaiBySiswaMapel(row.id, selectedMapel.value.id, selectedSemester.value, tahunAjaran.value)
  if (nilaiEntry) {
    nilaiStore.setNilai(nilaiEntry.id, field, value)
    // refresh nilai_akhir di tampilan
    const updated = nilaiStore.getNilaiBySiswaMapel(row.id, selectedMapel.value.id, selectedSemester.value, tahunAjaran.value)
    row.nilai_akhir = updated.nilai_akhir
  }
}

onMounted(() => {
  nilaiStore.loadData()
  kelasStore.loadData()
  mapelStore.loadData()
  siswaStore.loadData()
})
</script>
