<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Monitoring Ketidakhadiran (Alert & Surat Peringatan)</div>

    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-3">
        <q-select v-model="selectedKelas" :options="kelasOptions" label="Kelas" outlined dense />
      </div>
      <div class="col-12 col-md-2">
        <q-input v-model="bulan" label="Bulan" type="month" outlined dense />
      </div>
      <div class="col-12 col-md-2">
        <q-btn color="primary" label="Tampilkan" @click="loadRekap" />
      </div>
    </div>

    <div v-if="rekapSiswa.length">
      <q-table :rows="rekapSiswa" :columns="rekapColumns" row-key="id" flat bordered dense>
        <template v-slot:body-cell-alpha="props">
          <q-td :props="props" :class="{ 'bg-negative text-white': props.row.alpha >= 3 }">
            {{ props.row.alpha }}
          </q-td>
        </template>
        <template v-slot:body-cell-aksi="props">
          <q-td :props="props">
            <q-btn v-if="props.row.alpha >= 3" color="warning" label="Surat Peringatan"
              @click="generateSurat(props.row)" />
          </q-td>
        </template>
      </q-table>
    </div>

    <q-dialog v-model="suratDialog">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Surat Peringatan</div>
        </q-card-section>
        <q-card-section id="surat-content">
          <div class="text-center">
            <h5>SURAT PERINGATAN</h5>
            <p>Nomor: {{ nomorSurat }}</p>
          </div>
          <p>Kepada Yth. Orang Tua/Wali dari siswa:</p>
          <p><strong>{{ selectedSiswa?.nama }}</strong> (NIS: {{ selectedSiswa?.nis }})</p>
          <p>Berdasarkan rekap absensi bulan {{ bulan }}, siswa tersebut tercatat tidak hadir (alpha) sebanyak
            <strong>{{
              selectedSiswa?.alpha }}</strong> hari.
          </p>
          <p>Mohon perhatian dan bimbingan. Terima kasih.</p>
          <p class="q-mt-lg">Bandung, {{ tanggalCetak }}</p>
          <p>Wali Kelas</p>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Cetak" color="primary" @click="printSurat" />
          <q-btn flat label="Tutup" color="negative" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
// import { useQuasar } from 'quasar'
import { useSiswaStore } from 'stores/kesiswaan/siswa'
import { useKelasStore } from '@/stores/kesiswaan/kelasStore'
import { useAbsensiStore } from '@/stores/kesiswaan/absensiStore'

// const $q = useQuasar()
const siswaStore = useSiswaStore()
const kelasStore = useKelasStore()
const absensiStore = useAbsensiStore()

const selectedKelas = ref(null)
const bulan = ref(new Date().toISOString().slice(0, 7)) // YYYY-MM
const rekapSiswa = ref([])
const suratDialog = ref(false)
const selectedSiswa = ref(null)
const nomorSurat = ref('001/SP/' + new Date().getFullYear())
const tanggalCetak = new Date().toLocaleDateString('id-ID')

const kelasOptions = computed(() => kelasStore.list)
const rekapColumns = [
  { name: 'nama', label: 'Nama', field: 'nama' },
  { name: 'nis', label: 'NIS', field: 'nis' },
  { name: 'hadir', label: 'Hadir', field: 'hadir' },
  { name: 'izin', label: 'Izin', field: 'izin' },
  { name: 'sakit', label: 'Sakit', field: 'sakit' },
  { name: 'alpha', label: 'Alpha', field: 'alpha' },
  { name: 'aksi', label: 'Aksi', field: 'aksi' }
]

function loadRekap() {
  if (!selectedKelas.value) return
  const siswaDiKelas = siswaStore.list.filter(s => s.kelas_id === selectedKelas.value.id)
  const [tahun, month] = bulan.value.split('-')
  const data = siswaDiKelas.map(s => {
    const absenSiswa = absensiStore.list.filter(a => a.siswa_id === s.id && a.tanggal.startsWith(`${tahun}-${month}`))
    const hadir = absenSiswa.filter(a => a.status === 'hadir').length
    const izin = absenSiswa.filter(a => a.status === 'izin').length
    const sakit = absenSiswa.filter(a => a.status === 'sakit').length
    const alpha = absenSiswa.filter(a => a.status === 'alpha').length
    return { ...s, hadir, izin, sakit, alpha }
  })
  rekapSiswa.value = data
}

function generateSurat(siswa) {
  selectedSiswa.value = siswa
  suratDialog.value = true
}

function printSurat() {
  const content = document.getElementById('surat-content').innerHTML
  const win = window.open('', '_blank')
  win.document.write(`<html><head><title>Surat Peringatan</title></head><body>${content}</body></html>`)
  win.document.close()
  win.print()
}

onMounted(() => {
  siswaStore.loadData()
  kelasStore.loadData()
  absensiStore.loadData()
})
</script>
