<template>
  <q-page class="q-pa-md">
    <div class="row justify-between q-mb-md">
      <div class="text-h5">Pengajuan Cuti / SPT</div>
      <q-btn color="primary" icon="add" label="Ajukan Cuti/SPT" @click="openFormDialog" />
    </div>

    <!-- Filter status -->
    <q-select v-model="filterStatus" :options="['semua', 'pending', 'disetujui', 'ditolak']" label="Filter Status" dense
      outlined class="q-mb-md" style="max-width: 200px" />

    <q-table :rows="filteredCuti" :columns="columns" row-key="id" flat bordered dense>
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-badge
            :color="props.row.status === 'pending' ? 'warning' : (props.row.status === 'disetujui' ? 'positive' : 'negative')">
            {{ props.row.status }}
          </q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <!-- Tombol approval hanya jika role kepsek -->
          <div v-if="authStore.user?.role === 'kepsek' && props.row.status === 'pending'">
            <q-btn flat round dense icon="check_circle" color="positive" @click="approve(props.row.id)" class="q-mr-sm">
              <q-tooltip>Setujui</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="cancel" color="negative" @click="reject(props.row.id)">
              <q-tooltip>Tolak</q-tooltip>
            </q-btn>
          </div>
          <!-- Hapus hanya untuk role admin/kepsek -->
          <div v-if="authStore.user?.role === 'kepsek' || authStore.user?.role === 'admin'">
            <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)" />
          </div>
        </q-td>
      </template>
    </q-table>

    <!-- Dialog form pengajuan -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">Form Pengajuan Cuti / SPT</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitCuti" class="q-gutter-md">
            <q-select v-model="form.guru_id" :options="guruOptions" label="Pilih Guru" option-value="id"
              option-label="nama" outlined :rules="[v => !!v]" />
            <q-input v-model="form.tanggal_mulai" label="Tanggal Mulai" type="date" outlined :rules="[v => !!v]" />
            <q-input v-model="form.tanggal_selesai" label="Tanggal Selesai" type="date" outlined :rules="[v => !!v]" />
            <q-select v-model="form.jenis"
              :options="['Cuti Tahunan', 'Cuti Sakit', 'SPT (Surat Perintah Tugas)', 'Izin']" label="Jenis" outlined />
            <q-input v-model="form.alasan" label="Alasan" type="textarea" rows="2" outlined />
            <div class="q-gutter-sm">
              <q-btn label="Ajukan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup @click="resetForm" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog catatan penolakan -->
    <q-dialog v-model="rejectDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Alasan Penolakan</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="rejectNote" label="Catatan" type="textarea" outlined autofocus />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Tolak" color="primary" @click="executeReject" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useCutiStore } from '@/stores/kepegawaian/cuti'
import { useGuruStore } from '@/stores/kepegawaian/guru'
import { useAuthStore } from '@/stores/authStore'

const $q = useQuasar()
const cutiStore = useCutiStore()
const guruStore = useGuruStore()
const authStore = useAuthStore()

const filterStatus = ref('semua')
const formDialog = ref(false)
const rejectDialog = ref(false)
let currentRejectId = null
const rejectNote = ref('')
const form = ref({
  guru_id: null,
  tanggal_mulai: '',
  tanggal_selesai: '',
  jenis: 'Cuti Tahunan',
  alasan: ''
})

const columns = [
  { name: 'guru_nama', label: 'Guru', field: 'guru_nama', align: 'left' },
  { name: 'tanggal_mulai', label: 'Mulai', field: 'tanggal_mulai', align: 'left' },
  { name: 'tanggal_selesai', label: 'Selesai', field: 'tanggal_selesai', align: 'left' },
  { name: 'jenis', label: 'Jenis', field: 'jenis', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

const guruOptions = computed(() => guruStore.list.map(g => ({ id: g.id, nama: g.nama })))

const filteredCuti = computed(() => {
  if (filterStatus.value === 'semua') return cutiStore.list
  return cutiStore.list.filter(c => c.status === filterStatus.value)
})

function openFormDialog() {
  resetForm()
  formDialog.value = true
}

function resetForm() {
  form.value = {
    guru_id: null,
    tanggal_mulai: '',
    tanggal_selesai: '',
    jenis: 'Cuti Tahunan',
    alasan: ''
  }
}

function submitCuti() {
  if (!form.value.guru_id || !form.value.tanggal_mulai || !form.value.tanggal_selesai) {
    $q.notify({ type: 'negative', message: 'Lengkapi semua field' })
    return
  }
  const guru = guruStore.list.find(g => g.id === form.value.guru_id)
  if (!guru) return
  const newCuti = {
    guru_id: form.value.guru_id,
    guru_nama: guru.nama,
    tanggal_mulai: form.value.tanggal_mulai,
    tanggal_selesai: form.value.tanggal_selesai,
    jenis: form.value.jenis,
    alasan: form.value.alasan
  }
  cutiStore.tambah(newCuti)
  $q.notify({ type: 'positive', message: 'Pengajuan dikirim' })
  formDialog.value = false
}

function approve(id) {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Setujui pengajuan ini?',
    cancel: true
  }).onOk(() => {
    cutiStore.updateStatus(id, 'disetujui', 'Disetujui oleh Kepala Sekolah')
    $q.notify({ type: 'positive', message: 'Pengajuan disetujui' })
  })
}

function reject(id) {
  currentRejectId = id
  rejectNote.value = ''
  rejectDialog.value = true
}

function executeReject() {
  if (!rejectNote.value.trim()) {
    $q.notify({ type: 'negative', message: 'Alasan penolakan harus diisi' })
    return
  }
  cutiStore.updateStatus(currentRejectId, 'ditolak', rejectNote.value)
  $q.notify({ type: 'negative', message: 'Pengajuan ditolak' })
  rejectDialog.value = false
  currentRejectId = null
  rejectNote.value = ''
}

function confirmDelete(id) {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus pengajuan ini?',
    cancel: true
  }).onOk(() => {
    cutiStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Data dihapus' })
  })
}

onMounted(() => {
  guruStore.loadData()
  cutiStore.loadData()
})
</script>
