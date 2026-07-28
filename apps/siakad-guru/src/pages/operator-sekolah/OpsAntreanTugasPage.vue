<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">Antrean Tugas</div>

    <!-- Filter Chip -->
    <div class="q-gutter-sm q-mb-md">
      <q-chip v-for="kat in kategoriList" :key="kat.value" :color="kat.color" text-color="white"
        :outline="filterKategori !== kat.value" clickable @click="filterKategori = kat.value" icon="label">
        {{ kat.label }}
      </q-chip>
      <q-space />
      <q-input dense outlined v-model="searchText" placeholder="Cari nama/NISN/pengusul..." style="width: 300px">
        <template v-slot:append>
          <q-icon name="search" />
        </template>
      </q-input>
    </div>

    <!-- Daftar Tugas -->
    <q-card v-for="tugas in filteredTugas" :key="tugas.id" class="q-mb-md">
      <q-card-section>
        <div class="row items-center">
          <div class="col">
            <div class="text-subtitle1">{{ tugas.judul }}</div>
            <div class="text-caption">NISN: {{ tugas.nisn }} | Pengusul: {{ tugas.pengusul }}</div>
            <div class="text-caption text-grey">Kategori: {{ tugas.kategori }}</div>
          </div>
          <div class="col-auto">
            <q-btn flat round icon="content_copy" @click="salinTeks(tugas.id, tugas.dataBaru)" color="primary" />
            <q-btn flat round icon="check_circle" @click="konfirmasiSelesai(tugas)" color="positive" />
            <q-btn flat round icon="block" @click="konfirmasiAbaikan(tugas)" color="negative" />
          </div>
        </div>
        <q-separator class="q-my-sm" />
        <!-- Komparasi Data Side-by-Side -->
        <div class="row">
          <div class="col-6">
            <div class="text-weight-bold">Data Lama</div>
            <pre>{{ tugas.dataLama }}</pre>
          </div>
          <div class="col-6">
            <div class="text-weight-bold">Data Baru</div>
            <pre class="text-red">{{ tugas.dataBaru }}</pre>
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Dialog Konfirmasi Selesai -->
    <q-dialog v-model="dialogSelesai" persistent>
      <q-card>
        <q-card-section class="row items-center">
          <q-avatar icon="check_circle" color="positive" text-color="white" />
          <span class="q-ml-sm">Selesaikan tugas ini? Data akan dikirim ke Dapodik.</span>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn flat label="Selesai" color="positive" @click="selesaikanTugas" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Abaikan dengan Alasan Penolakan -->
    <q-dialog v-model="dialogAbaikan" persistent>
      <q-card>
        <q-card-section class="row items-center">
          <q-avatar icon="block" color="negative" text-color="white" />
          <span class="q-ml-sm">Abaikan tugas? Berikan alasan penolakan.</span>
        </q-card-section>
        <q-card-section>
          <q-input v-model="alasanPenolakan" type="textarea" label="Alasan penolakan" required />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn flat label="Abaikan" color="negative" @click="abaikanTugas" :disable="!alasanPenolakan" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// Dummy data
const daftarTugas = ref([
  {
    id: 1,
    judul: 'Pembaruan Biodata Siswa',
    kategori: 'Biodata',
    nisn: '1234567890',
    pengusul: 'Ahmad Fauzi',
    dataLama: 'Nama: Ahmad Fauzi, Tempat Lahir: Jakarta',
    dataBaru: 'Nama: Ahmad Fauzi, Tempat Lahir: Bandung',
    status: 'pending'
  },
  {
    id: 2,
    judul: 'Mutasi Masuk',
    kategori: 'Mutasi',
    nisn: '0987654321',
    pengusul: 'Siti Nurhaliza',
    dataLama: 'Asal Sekolah: SMPN 1',
    dataBaru: 'Asal Sekolah: SMPN 2',
    status: 'pending'
  }
])

// Filter
const filterKategori = ref('')
const searchText = ref('')
const kategoriList = [
  { label: 'Semua', value: '', color: 'grey' },
  { label: 'Biodata', value: 'Biodata', color: 'blue' },
  { label: 'Mutasi', value: 'Mutasi', color: 'orange' },
  { label: 'GTK', value: 'GTK', color: 'purple' },
  { label: 'Rombel', value: 'Rombel', color: 'teal' }
]

const filteredTugas = computed(() => {
  let result = daftarTugas.value
  if (filterKategori.value) {
    result = result.filter(t => t.kategori === filterKategori.value)
  }
  if (searchText.value) {
    const needle = searchText.value.toLowerCase()
    result = result.filter(t =>
      t.nisn.includes(needle) ||
      t.pengusul.toLowerCase().includes(needle) ||
      t.judul.toLowerCase().includes(needle)
    )
  }
  return result
})

// Salin teks
function salinTeks(id, teks) {
  navigator.clipboard.writeText(teks)
  $q.notify({ message: '✓ Disalin', color: 'positive', timeout: 1000 })
}

// Dialog Selesai
const dialogSelesai = ref(false)
const selectedTugas = ref(null)

function konfirmasiSelesai(tugas) {
  selectedTugas.value = tugas
  dialogSelesai.value = true
}

function selesaikanTugas() {
  daftarTugas.value = daftarTugas.value.filter(t => t.id !== selectedTugas.value.id)
  $q.notify({ message: 'Tugas diselesaikan', color: 'positive' })
  dialogSelesai.value = false
}

// Dialog Abaikan
const dialogAbaikan = ref(false)
const alasanPenolakan = ref('')

function konfirmasiAbaikan(tugas) {
  selectedTugas.value = tugas
  alasanPenolakan.value = ''
  dialogAbaikan.value = true
}

function abaikanTugas() {
  daftarTugas.value = daftarTugas.value.filter(t => t.id !== selectedTugas.value.id)
  $q.notify({ message: 'Tugas diabaikan', color: 'negative' })
  dialogAbaikan.value = false
  // Catat alasanPenolakan.value ke log (bisa dikirim ke backend)
}
</script>
