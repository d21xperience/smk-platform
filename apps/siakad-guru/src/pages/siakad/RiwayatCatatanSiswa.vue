<template>
  <q-page class="record-page q-pa-md">
    <!-- HEADER -->
    <div class="row items-center justify-between q-mb-md wrap q-gutter-sm">
      <div>
        <div class="text-h6 text-weight-bold text-grey-9">Riwayat Catatan Siswa</div>
        <div class="text-caption text-grey-6">Prestasi & pelanggaran selama masa pendidikan</div>
      </div>
      <q-btn no-caps unelevated color="primary" icon="add" label="Tambah Catatan" @click="openAddDialog" />
    </div>

    <!-- PILIH SISWA -->
    <q-card flat bordered class="q-mb-md">
      <q-card-section>
        <q-select outlined dense v-model="selectedStudentId" :options="studentOptions" emit-value map-options
          label="Pilih Siswa" use-input input-debounce="200" @filter="filterStudentOptions"
          @update:model-value="handleStudentChange">
          <template v-slot:prepend>
            <q-icon name="person_search" />
          </template>
        </q-select>
      </q-card-section>
    </q-card>

    <!-- KARTU IDENTITAS SISWA -->
    <q-card flat bordered class="student-banner q-mb-md overflow-hidden">
      <q-card-section class="row items-center q-col-gutter-md">
        <div class="col-auto">
          <q-avatar size="72px" class="bg-grey-3">
            <img v-if="student.foto" :src="student.foto" :alt="student.nama" />
            <q-icon v-else name="person" size="40px" color="grey-6" />
          </q-avatar>
        </div>
        <div class="col">
          <div class="text-subtitle1 text-weight-bold">{{ student.nama }}</div>
          <div class="text-caption text-grey-7">
            NIS {{ student.nis }} &middot; NISN {{ student.nisn }}
          </div>
          <div class="text-caption text-grey-7">
            {{ student.kelasSaatIni }} &middot; Wali Kelas: {{ student.waliKelas }}
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- RINGKASAN -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-6 col-md-3">
        <q-card flat bordered class="stat-card">
          <q-card-section class="row items-center no-wrap">
            <div class="stat-icon bg-positive-1 q-mr-md">
              <q-icon name="emoji_events" color="positive" size="22px" />
            </div>
            <div>
              <div class="text-h6 text-weight-bold">{{ totalPrestasi }}</div>
              <div class="text-caption text-grey-7">Total Prestasi</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-6 col-md-3">
        <q-card flat bordered class="stat-card">
          <q-card-section class="row items-center no-wrap">
            <div class="stat-icon bg-red-1 q-mr-md">
              <q-icon name="report_problem" color="red" size="22px" />
            </div>
            <div>
              <div class="text-h6 text-weight-bold">{{ totalPelanggaran }}</div>
              <div class="text-caption text-grey-7">Total Pelanggaran</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-6 col-md-3">
        <q-card flat bordered class="stat-card">
          <q-card-section class="row items-center no-wrap">
            <div class="stat-icon bg-orange-1 q-mr-md">
              <q-icon name="speed" color="orange" size="22px" />
            </div>
            <div>
              <div class="text-h6 text-weight-bold">{{ totalPoinPelanggaran }}</div>
              <div class="text-caption text-grey-7">Poin Pelanggaran Berjalan</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-6 col-md-3">
        <q-card flat bordered class="stat-card">
          <q-card-section class="row items-center no-wrap">
            <div class="stat-icon q-mr-md" :class="`bg-${statusKedisiplinan.color}-1`">
              <q-icon name="verified_user" :color="statusKedisiplinan.color" size="22px" />
            </div>
            <div>
              <div class="text-body1 text-weight-bold" :class="`text-${statusKedisiplinan.color}`">
                {{ statusKedisiplinan.label }}
              </div>
              <div class="text-caption text-grey-7">Status Kedisiplinan</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- FILTER -->
    <q-card flat bordered class="q-mb-md">
      <q-card-section class="row q-col-gutter-sm items-center">
        <div class="col-12 col-sm-5">
          <q-input dense outlined v-model="searchQuery" placeholder="Cari catatan..." clearable>
            <template v-slot:prepend>
              <q-icon name="search" />
            </template>
          </q-input>
        </div>
        <div class="col-12 col-sm-4">
          <q-select dense outlined v-model="filterTingkat" :options="tingkatFilterOptions" emit-value map-options
            label="Tingkat Kelas" />
        </div>
        <div class="col-12 col-sm-3">
          <q-btn-toggle v-model="filterJenis" no-caps unelevated dense spread toggle-color="primary" color="white"
            text-color="grey-8" :options="[
              { label: 'Semua', value: 'semua' },
              { label: 'Prestasi', value: 'prestasi' },
              { label: 'Pelanggaran', value: 'pelanggaran' }
            ]" class="jenis-toggle" />
        </div>
      </q-card-section>
    </q-card>

    <!-- RIWAYAT PER TINGKAT / TAHUN AJARAN -->
    <div v-if="groupedRecords.length">
      <q-card flat bordered class="q-mb-md" v-for="group in groupedRecords" :key="group.tingkat">
        <q-expansion-item :label="`Kelas ${group.tingkat} &middot; Tahun Ajaran ${group.tahunAjaran}`"
          :caption="`${group.records.length} catatan`" header-class="text-weight-bold" icon="school"
          :default-opened="group.isOngoing">
          <template v-slot:header>
            <q-item-section avatar>
              <q-icon name="school" color="primary" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-bold">
                Kelas {{ group.tingkat }} &middot; Tahun Ajaran {{ group.tahunAjaran }}
              </q-item-label>
              <q-item-label caption>{{ group.records.length }} catatan</q-item-label>
            </q-item-section>
            <q-item-section side v-if="group.isOngoing">
              <q-badge color="primary" outline>Sedang Berlangsung</q-badge>
            </q-item-section>
          </template>

          <q-separator />
          <q-card-section class="q-gutter-sm">
            <div v-for="record in group.records" :key="record.id" class="record-card"
              :class="record.jenis === 'prestasi' ? 'record-card--prestasi' : 'record-card--pelanggaran'">
              <div class="row items-start no-wrap q-gutter-sm">
                <q-icon :name="record.jenis === 'prestasi' ? 'emoji_events' : 'report_problem'"
                  :color="record.jenis === 'prestasi' ? 'positive' : 'red'" size="24px" class="q-mt-xs" />
                <div class="col">
                  <div class="row items-center justify-between wrap q-gutter-xs">
                    <div class="text-weight-medium">{{ record.judul }}</div>
                    <q-badge :color="record.jenis === 'prestasi' ? 'positive' : categoryBadgeColor(record.kategori)"
                      outline>
                      {{ record.kategori }}
                    </q-badge>
                  </div>
                  <div class="text-caption text-grey-7 q-mt-xs">
                    {{ formatDate(record.tanggal) }}
                    <span v-if="record.jenis === 'pelanggaran'"> &middot; Poin: {{ record.poin }}</span>
                  </div>
                  <div class="text-body2 q-mt-xs">{{ record.keterangan }}</div>
                  <div class="text-caption text-grey-6 q-mt-xs">
                    Dicatat oleh: {{ record.pencatat }}
                  </div>
                  <div class="q-mt-xs" v-if="record.lampiran">
                    <q-btn no-caps dense flat color="primary" icon="attach_file" :label="record.lampiran"
                      @click="handleDownloadLampiran(record)" />
                  </div>
                </div>
              </div>
            </div>
          </q-card-section>
        </q-expansion-item>
      </q-card>
    </div>

    <q-card flat bordered v-else>
      <q-card-section class="text-center text-grey-6 q-pa-xl">
        <q-icon name="folder_off" size="40px" class="q-mb-sm" />
        <div>Tidak ada catatan yang cocok dengan filter</div>
      </q-card-section>
    </q-card>

    <!-- DIALOG TAMBAH CATATAN -->
    <q-dialog v-model="addDialogOpen">
      <q-card style="min-width: 340px; max-width: 480px" class="full-width">
        <q-card-section class="row items-center justify-between">
          <div class="text-subtitle1 text-weight-bold">Tambah Catatan Siswa</div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>
        <q-separator />
        <q-form @submit="submitNewRecord">
          <q-card-section class="q-gutter-md">
            <q-btn-toggle v-model="newRecord.jenis" no-caps unelevated spread toggle-color="primary" color="white"
              text-color="grey-8" :options="[
                { label: 'Prestasi', value: 'prestasi', icon: 'emoji_events' },
                { label: 'Pelanggaran', value: 'pelanggaran', icon: 'report_problem' }
              ]" />

            <q-input outlined dense type="date" v-model="newRecord.tanggal" label="Tanggal"
              :rules="[val => !!val || 'Tanggal wajib diisi']" />

            <q-select outlined dense v-model="newRecord.kategori" :options="kategoriOptionsForJenis" label="Kategori"
              :rules="[val => !!val || 'Kategori wajib dipilih']" />

            <q-input outlined dense v-model="newRecord.judul" label="Judul Catatan"
              :rules="[val => val && val.length > 0 || 'Judul wajib diisi']" />

            <q-input outlined dense type="textarea" autogrow v-model="newRecord.keterangan" label="Keterangan" />

            <q-input v-if="newRecord.jenis === 'pelanggaran'" outlined dense type="number"
              v-model.number="newRecord.poin" label="Poin Pelanggaran"
              :rules="[val => val >= 0 || 'Poin tidak boleh negatif']" />
          </q-card-section>

          <q-card-actions align="right" class="q-pa-md">
            <q-btn no-caps flat label="Batal" v-close-popup />
            <q-btn no-caps unelevated color="primary" label="Simpan" type="submit" />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

const MONTHS_ID = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
]
function formatDate(iso) {
  const [y, m, d] = iso.split('-').map(Number)
  return `${d} ${MONTHS_ID[m - 1]} ${y}`
}

// TODO: ganti dengan hasil GET /api/kesiswaan/siswa?query= (Go backend) untuk pencarian siswa
const allStudents = reactive([
  {
    id: 'std-1',
    nama: 'Rafi Ardiansyah',
    nis: '2021100234',
    nisn: '0051234567',
    kelasSaatIni: 'XII RPL 1',
    waliKelas: 'Siti Nurhaliza, S.Pd.',
    foto: ''
  },
  { id: 'std-2', nama: 'Nabila Putri', nis: '2021100155', nisn: '0051234321', kelasSaatIni: 'XII RPL 2', waliKelas: 'Dedi Supriadi, S.Kom.', foto: '' },
  { id: 'std-3', nama: 'Fajar Nugraha', nis: '2022100088', nisn: '0059876543', kelasSaatIni: 'XI RPL 1', waliKelas: 'Rina Marlina, S.Pd.', foto: '' }
])

const selectedStudentId = ref('std-1')
const studentOptions = ref(
  allStudents.map(s => ({ label: `${s.nama} - ${s.kelasSaatIni}`, value: s.id }))
)

// Filter opsi dropdown siswa saat mengetik (client-side sederhana)
// TODO: ganti dengan pemanggilan API pencarian siswa saat data siswa sudah banyak
function filterStudentOptions(query, update) {
  update(() => {
    if (!query) {
      studentOptions.value = allStudents.map(s => ({ label: `${s.nama} - ${s.kelasSaatIni}`, value: s.id }))
      return
    }
    const q = query.toLowerCase()
    studentOptions.value = allStudents
      .filter(s => s.nama.toLowerCase().includes(q) || s.nis.includes(q))
      .map(s => ({ label: `${s.nama} - ${s.kelasSaatIni}`, value: s.id }))
  })
}

const student = computed(() => {
  return allStudents.find(s => s.id === selectedStudentId.value) || allStudents[0]
})

function handleStudentChange() {
  // TODO: panggil GET /api/kesiswaan/siswa/:id/catatan untuk mengambil riwayat catatan siswa terpilih
  if (selectedStudentId.value !== 'std-1') {
    $q.notify({
      type: 'info',
      message: 'Data riwayat untuk siswa ini belum tersedia pada contoh ini.',
      position: 'top',
      timeout: 2000
    })
  }
}

// TODO: ganti seluruh data berikut dengan hasil GET /api/kesiswaan/siswa/:id/catatan (Go backend)
// Dikelompokkan berdasarkan tingkat kelas & tahun ajaran saat catatan itu dibuat
const records = reactive([
  {
    id: 'r1',
    tingkat: 'X',
    tahunAjaran: '2023/2024',
    tanggal: '2023-10-05',
    jenis: 'prestasi',
    kategori: 'Akademik',
    judul: 'Juara 2 Lomba LKS Web Design Tingkat Kabupaten',
    keterangan: 'Mewakili sekolah pada ajang Lomba Kompetensi Siswa (LKS) bidang Web Design.',
    pencatat: 'Kaprog RPL',
    lampiran: 'Sertifikat-LKS-2023.pdf'
  },
  {
    id: 'r2',
    tingkat: 'X',
    tahunAjaran: '2023/2024',
    tanggal: '2024-02-12',
    jenis: 'pelanggaran',
    kategori: 'Ringan',
    poin: 5,
    judul: 'Terlambat Masuk Sekolah 3 Kali dalam Sebulan',
    keterangan: 'Ditegur lisan dan diminta membuat surat pernyataan.',
    pencatat: 'Guru BK',
    lampiran: ''
  },
  {
    id: 'r3',
    tingkat: 'X',
    tahunAjaran: '2023/2024',
    tanggal: '2024-04-01',
    jenis: 'pelanggaran',
    kategori: 'Ringan',
    poin: 3,
    judul: 'Tidak Mengikuti Upacara Bendera',
    keterangan: 'Tanpa keterangan yang jelas.',
    pencatat: 'Guru Piket',
    lampiran: ''
  },
  {
    id: 'r4',
    tingkat: 'XI',
    tahunAjaran: '2024/2025',
    tanggal: '2024-08-20',
    jenis: 'prestasi',
    kategori: 'Non-Akademik',
    judul: 'Anggota Terbaik Ekstrakurikuler Pramuka',
    keterangan: 'Aktif dan berprestasi dalam kegiatan kepramukaan tingkat gugus depan.',
    pencatat: 'Pembina Pramuka',
    lampiran: ''
  },
  {
    id: 'r5',
    tingkat: 'XI',
    tahunAjaran: '2024/2025',
    tanggal: '2024-11-15',
    jenis: 'pelanggaran',
    kategori: 'Berat',
    poin: 25,
    judul: 'Berkelahi dengan Teman Sekelas',
    keterangan: 'Dipanggil orang tua/wali untuk pembinaan bersama pihak sekolah.',
    pencatat: 'Guru BK',
    lampiran: 'BA-Pelanggaran-2024-11.pdf'
  },
  {
    id: 'r6',
    tingkat: 'XI',
    tahunAjaran: '2024/2025',
    tanggal: '2025-03-10',
    jenis: 'prestasi',
    kategori: 'Akademik',
    judul: 'Juara 1 Lomba Debat Bahasa Inggris Tingkat Provinsi',
    keterangan: 'Mewakili sekolah dan meraih peringkat pertama tingkat provinsi.',
    pencatat: 'Guru Bahasa Inggris',
    lampiran: 'Sertifikat-Debat-2025.pdf'
  },
  {
    id: 'r7',
    tingkat: 'XII',
    tahunAjaran: '2025/2026',
    tanggal: '2025-08-05',
    jenis: 'pelanggaran',
    kategori: 'Ringan',
    poin: 5,
    judul: 'Tidak Mengerjakan Tugas Berturut-turut',
    keterangan: 'Sudah diingatkan oleh guru mata pelajaran terkait.',
    pencatat: 'Wali Kelas',
    lampiran: ''
  },
  {
    id: 'r8',
    tingkat: 'XII',
    tahunAjaran: '2025/2026',
    tanggal: '2025-09-01',
    jenis: 'prestasi',
    kategori: 'Akademik',
    judul: 'Lolos Seleksi PKL di Perusahaan Mitra Unggulan',
    keterangan: 'Diterima Praktik Kerja Lapangan di salah satu mitra industri utama sekolah.',
    pencatat: 'Kaprog RPL',
    lampiran: ''
  }
])

const searchQuery = ref('')
const filterJenis = ref('semua')
const filterTingkat = ref('semua')

const tingkatFilterOptions = [
  { label: 'Semua Tingkat', value: 'semua' },
  { label: 'Kelas X', value: 'X' },
  { label: 'Kelas XI', value: 'XI' },
  { label: 'Kelas XII', value: 'XII' }
]

function categoryBadgeColor(kategori) {
  if (kategori === 'Ringan') return 'orange'
  if (kategori === 'Sedang') return 'deep-orange'
  if (kategori === 'Berat') return 'red'
  return 'grey'
}

// Catatan terfilter berdasarkan pencarian, jenis, dan tingkat kelas
const filteredRecords = computed(() => {
  const query = (searchQuery.value || '').toLowerCase().trim()
  return records
    .filter(r => {
      const matchQuery = !query
        || r.judul.toLowerCase().includes(query)
        || r.keterangan.toLowerCase().includes(query)
      const matchJenis = filterJenis.value === 'semua' || r.jenis === filterJenis.value
      const matchTingkat = filterTingkat.value === 'semua' || r.tingkat === filterTingkat.value
      return matchQuery && matchJenis && matchTingkat
    })
    .slice()
    .sort((a, b) => (a.tanggal < b.tanggal ? 1 : -1))
})

// Dikelompokkan per tingkat & tahun ajaran, diurutkan dari yang terbaru
const groupedRecords = computed(() => {
  const groups = []
  filteredRecords.value.forEach(record => {
    let group = groups.find(g => g.tingkat === record.tingkat && g.tahunAjaran === record.tahunAjaran)
    if (!group) {
      group = {
        tingkat: record.tingkat,
        tahunAjaran: record.tahunAjaran,
        // TODO: tandai tahun ajaran yang sedang berjalan berdasarkan data resmi dari backend, bukan hardcode
        isOngoing: record.tahunAjaran === '2025/2026',
        records: []
      }
      groups.push(group)
    }
    group.records.push(record)
  })
  return groups.sort((a, b) => (a.tahunAjaran < b.tahunAjaran ? 1 : -1))
})

const totalPrestasi = computed(() => records.filter(r => r.jenis === 'prestasi').length)
const totalPelanggaran = computed(() => records.filter(r => r.jenis === 'pelanggaran').length)
const totalPoinPelanggaran = computed(() =>
  records.filter(r => r.jenis === 'pelanggaran').reduce((sum, r) => sum + (r.poin || 0), 0)
)

// TODO: ganti ambang batas poin ini sesuai aturan tata tertib resmi sekolah dari backend
const statusKedisiplinan = computed(() => {
  const poin = totalPoinPelanggaran.value
  if (poin <= 10) return { label: 'Baik', color: 'positive' }
  if (poin <= 25) return { label: 'Perlu Perhatian', color: 'orange' }
  return { label: 'Perlu Pembinaan Intensif', color: 'red' }
})

function handleDownloadLampiran(record) {
  // TODO: ganti dengan permintaan unduh sesungguhnya ke Go backend,
  // misalnya GET /api/kesiswaan/catatan/:id/lampiran
  $q.notify({
    type: 'positive',
    message: `Mengunduh "${record.lampiran}"...`,
    position: 'top',
    timeout: 1500
  })
}

// ---------- dialog tambah catatan ----------
const addDialogOpen = ref(false)
const newRecord = reactive({
  jenis: 'prestasi',
  tanggal: '',
  kategori: null,
  judul: '',
  keterangan: '',
  poin: 0
})

const kategoriOptionsForJenis = computed(() => {
  return newRecord.jenis === 'prestasi'
    ? ['Akademik', 'Non-Akademik']
    : ['Ringan', 'Sedang', 'Berat']
})

function openAddDialog() {
  newRecord.jenis = 'prestasi'
  newRecord.tanggal = ''
  newRecord.kategori = null
  newRecord.judul = ''
  newRecord.keterangan = ''
  newRecord.poin = 0
  addDialogOpen.value = true
}

function submitNewRecord() {
  // TODO: ganti dengan POST /api/kesiswaan/siswa/:id/catatan (Go backend)
  // const now = new Date(newRecord.tanggal)
  const tingkatSaatIni = 'XII' // TODO: tentukan tingkat aktif siswa dari data backend, bukan hardcode
  const tahunAjaranSaatIni = '2025/2026'

  records.push({
    id: `r-${Date.now()}`,
    tingkat: tingkatSaatIni,
    tahunAjaran: tahunAjaranSaatIni,
    tanggal: newRecord.tanggal || formatDate,
    jenis: newRecord.jenis,
    kategori: newRecord.kategori,
    poin: newRecord.jenis === 'pelanggaran' ? newRecord.poin : undefined,
    judul: newRecord.judul,
    keterangan: newRecord.keterangan,
    pencatat: 'Anda', // TODO: isi otomatis dari authStore.user.name
    lampiran: ''
  })

  $q.notify({
    type: 'positive',
    message: 'Catatan berhasil ditambahkan.',
    position: 'top',
    timeout: 2000
  })

  addDialogOpen.value = false
}
</script>

<style scoped>
.record-page {
  max-width: 1200px;
  margin: 0 auto;
}

.student-banner {
  border-radius: 8px;
}

.stat-card {
  border-radius: 8px;
  height: 100%;
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.jenis-toggle {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
}

.record-card {
  border-radius: 8px;
  border: 1px solid #ebebeb;
  border-left-width: 4px;
  padding: 12px;
}

.record-card--prestasi {
  border-left-color: var(--q-positive, #21ba45);
}

.record-card--pelanggaran {
  border-left-color: var(--q-red, #c10015);
}
</style>
