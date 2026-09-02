<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Log Aktivitas</div>

    <!-- Filter Bar -->
    <q-card flat bordered class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-3">
            <q-input v-model="filters.search" label="Cari (pengguna, aksi, detail)" outlined dense debounce="400"
              clearable @update:model-value="loadLogs">
              <template v-slot:append>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>
          <div class="col-12 col-md-2">
            <q-select v-model="filters.role" :options="roleOptions" label="Role" outlined dense clearable
              @update:model-value="loadLogs" />
          </div>
          <div class="col-12 col-md-2">
            <q-select v-model="filters.aksi" :options="aksiOptions" label="Aksi" outlined dense clearable
              @update:model-value="loadLogs" />
          </div>
          <div class="col-12 col-md-2">
            <q-input v-model="filters.tanggalMulai" label="Tanggal Mulai" type="date" outlined dense
              @update:model-value="loadLogs" />
          </div>
          <div class="col-12 col-md-2">
            <q-input v-model="filters.tanggalAkhir" label="Tanggal Akhir" type="date" outlined dense
              @update:model-value="loadLogs" />
          </div>
          <div class="col-12 col-md-1">
            <q-btn color="negative" label="Reset" flat @click="resetFilters" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Tabel Log -->
    <q-table :rows="logs" :columns="columns" row-key="id" flat bordered dense :loading="loading"
      :pagination="pagination" @request="onRequest" rows-per-page-options="[10, 25, 50]" @row-click="(evt, row) => showDetail(row)">
      <template v-slot:body-cell-role="props">
        <q-td :props="props">
          <q-badge :color="getRoleColor(props.row.role)" :label="props.row.role" />
        </q-td>
      </template>
      <template v-slot:body-cell-waktu="props">
        <q-td :props="props">{{ formatDateTime(props.row.waktu) }}</q-td>
      </template>
      <template v-slot:body-cell-detail="props">
        <q-td :props="props">
          <div class="ellipsis" style="max-width: 300px;">{{ props.row.detail }}</div>
        </q-td>
      </template>

    </q-table>

    <!-- Dialog detail log -->
    <q-dialog v-model="detailDialog">
      <q-card style="min-width: 500px;">
        <q-card-section>
          <div class="text-h6">Detail Aktivitas</div>
        </q-card-section>
        <q-card-section v-if="selectedLog">
          <div><strong>Pengguna:</strong> {{ selectedLog.pengguna }} ({{ selectedLog.role }})</div>
          <div><strong>Aksi:</strong> {{ selectedLog.aksi }}</div>
          <div><strong>Waktu:</strong> {{ formatDateTime(selectedLog.waktu) }}</div>
          <div><strong>Detail:</strong></div>
          <pre class="detail-pre">{{ selectedLog.detail }}</pre>
          <div v-if="selectedLog.ip_address"><strong>IP Address:</strong> {{ selectedLog.ip_address }}</div>
          <div v-if="selectedLog.user_agent"><strong>User Agent:</strong> {{ selectedLog.user_agent }}</div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Tutup" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'

// Data dummy logs
const allLogs = ref([
  { id: 1, pengguna: 'Admin Utama', role: 'admin', aksi: 'LOGIN', detail: 'Login ke sistem dari IP 192.168.1.1', waktu: '2025-06-15T08:30:00', ip_address: '192.168.1.1', user_agent: 'Chrome/120' },
  { id: 2, pengguna: 'Budi Santoso', role: 'guru', aksi: 'BUAT_SOAL', detail: 'Membuat soal baru (ID: 15) untuk kategori Matematika', waktu: '2025-06-15T09:15:00', ip_address: '192.168.1.45', user_agent: 'Firefox/115' },
  { id: 3, pengguna: 'Ahmad Faizal', role: 'siswa', aksi: 'MULAI_UJIAN', detail: 'Memulai ujian "UTS Matematika" (ID: 3)', waktu: '2025-06-15T10:00:00', ip_address: '192.168.1.78' },
  { id: 4, pengguna: 'Siti Nurhaliza', role: 'siswa', aksi: 'SELESAI_UJIAN', detail: 'Menyelesaikan ujian "UAS Fisika" dengan nilai 88', waktu: '2025-06-15T11:30:00' },
  { id: 5, pengguna: 'Dewi Lestari', role: 'guru', aksi: 'EDIT_UJIAN', detail: 'Mengubah durasi ujian "UTS Matematika" dari 90 menjadi 120 menit', waktu: '2025-06-14T14:20:00' },
  { id: 6, pengguna: 'Admin Utama', role: 'admin', aksi: 'TAMBAH_PENGGUNA', detail: 'Menambahkan user baru: siswa "Rizki Pratama"', waktu: '2025-06-14T10:05:00' },
  { id: 7, pengguna: 'Budi Santoso', role: 'guru', aksi: 'HAPUS_SOAL', detail: 'Menghapus soal ID: 8 (kategori Fisika)', waktu: '2025-06-13T13:45:00' },
  { id: 8, pengguna: 'Ahmad Faizal', role: 'siswa', aksi: 'LOGIN', detail: 'Login dari perangkat mobile', waktu: '2025-06-13T07:55:00' }
])

const roleOptions = ['admin', 'guru', 'siswa']
const aksiOptions = ['LOGIN', 'LOGOUT', 'BUAT_SOAL', 'EDIT_SOAL', 'HAPUS_SOAL', 'BUAT_UJIAN', 'EDIT_UJIAN', 'HAPUS_UJIAN', 'MULAI_UJIAN', 'SELESAI_UJIAN', 'TAMBAH_PENGGUNA', 'EDIT_PENGGUNA', 'HAPUS_PENGGUNA', 'BACKUP', 'RESTORE', 'EXPORT']

const filters = reactive({
  search: '',
  role: null,
  aksi: null,
  tanggalMulai: null,
  tanggalAkhir: null
})

const loading = ref(false)
const logs = ref([])
const pagination = ref({
  sortBy: 'waktu',
  descending: true,
  page: 1,
  rowsPerPage: 10,
  rowsNumber: 0
})

const detailDialog = ref(false)
const selectedLog = ref(null)

const columns = [
  { name: 'waktu', label: 'Waktu', field: 'waktu', align: 'left', sortable: true },
  { name: 'pengguna', label: 'Pengguna', field: 'pengguna', align: 'left' },
  { name: 'role', label: 'Role', field: 'role', align: 'center' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'left' },
  { name: 'detail', label: 'Detail', field: 'detail', align: 'left' }
]

function getRoleColor(role) {
  if (role === 'admin') return 'red'
  if (role === 'guru') return 'blue'
  return 'green'
}

function formatDateTime(isoString) {
  if (!isoString) return '-'
  const date = new Date(isoString)
  return date.toLocaleString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function filterLogs() {
  let filtered = [...allLogs.value]
  if (filters.search) {
    const searchLower = filters.search.toLowerCase()
    filtered = filtered.filter(log =>
      log.pengguna.toLowerCase().includes(searchLower) ||
      log.aksi.toLowerCase().includes(searchLower) ||
      log.detail.toLowerCase().includes(searchLower)
    )
  }
  if (filters.role) {
    filtered = filtered.filter(log => log.role === filters.role)
  }
  if (filters.aksi) {
    filtered = filtered.filter(log => log.aksi === filters.aksi)
  }
  if (filters.tanggalMulai) {
    const start = new Date(filters.tanggalMulai)
    start.setHours(0, 0, 0, 0)
    filtered = filtered.filter(log => new Date(log.waktu) >= start)
  }
  if (filters.tanggalAkhir) {
    const end = new Date(filters.tanggalAkhir)
    end.setHours(23, 59, 59, 999)
    filtered = filtered.filter(log => new Date(log.waktu) <= end)
  }
  // sorting by waktu descending (default)
  filtered.sort((a, b) => new Date(b.waktu) - new Date(a.waktu))
  pagination.value.rowsNumber = filtered.length
  const start = (pagination.value.page - 1) * pagination.value.rowsPerPage
  const end = start + pagination.value.rowsPerPage
  logs.value = filtered.slice(start, end)
}

function loadLogs() {
  loading.value = true
  setTimeout(() => {
    filterLogs()
    loading.value = false
  }, 300)
}

function onRequest(props) {
  const { page, rowsPerPage, sortBy, descending } = props.pagination
  pagination.value.page = page
  pagination.value.rowsPerPage = rowsPerPage
  pagination.value.sortBy = sortBy
  pagination.value.descending = descending
  loadLogs()
}

function resetFilters() {
  filters.search = ''
  filters.role = null
  filters.aksi = null
  filters.tanggalMulai = null
  filters.tanggalAkhir = null
  pagination.value.page = 1
  loadLogs()
}

function showDetail(log) {
  selectedLog.value = log
  detailDialog.value = true
}

// Override column click to show detail
// Bisa tambahkan tombol detail di kolom aksi, tapi untuk sederhana, kita buat baris bisa diklik
// Kita tambahkan @row-click pada q-table
// Saya akan tambahkan di template nanti

onMounted(() => {
  loadLogs()
})
</script>

<style scoped>
.detail-pre {
  background: #f5f5f5;
  padding: 8px;
  border-radius: 4px;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 200px;
  overflow: auto;
}
</style>
