<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Logbook Digital PKL</div>

    <!-- Filter siswa -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-4">
        <q-select v-model="selectedSiswa" :options="siswaOptions" label="Pilih Siswa" option-value="id"
          option-label="nama" outlined dense clearable />
      </div>
      <div class="col-12 col-md-2">
        <q-btn color="primary" label="Tambah Entri" @click="openForm" :disable="!selectedSiswa" />
      </div>
    </div>

    <q-table :rows="filteredLogbook" :columns="columns" row-key="id" flat bordered dense>
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-badge
            :color="props.row.status === 'disetujui' ? 'positive' : (props.row.status === 'ditolak' ? 'negative' : 'warning')">
            {{ props.row.status }}
          </q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <div v-if="authStore.user?.role === 'guru' || authStore.user?.role === 'pembimbing'">
            <q-btn v-if="props.row.status === 'pending'" flat round dense icon="check_circle" color="positive"
              @click="approveLog(props.row.id)" class="q-mr-sm">
              <q-tooltip>Setujui</q-tooltip>
            </q-btn>
            <q-btn v-if="props.row.status === 'pending'" flat round dense icon="cancel" color="negative"
              @click="rejectLog(props.row.id)">
              <q-tooltip>Tolak</q-tooltip>
            </q-btn>
          </div>
          <q-btn flat round dense icon="edit" color="info" @click="openForm(props.row)" class="q-mr-sm" />
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)" />
        </q-td>
      </template>
    </q-table>

    <!-- Form dialog -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">{{ editMode ? 'Edit Logbook' : 'Tambah Logbook' }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitLog" class="q-gutter-md">
            <q-input v-model="form.tanggal" label="Tanggal" type="date" outlined :rules="[v => !!v]" />
            <q-input v-model="form.aktivitas" label="Aktivitas / Kegiatan" type="textarea" rows="3" outlined
              :rules="[v => !!v]" />
            <q-input v-model="form.hasil" label="Hasil / Output" type="textarea" rows="2" outlined />
            <q-input v-model="form.keterangan" label="Keterangan (opsional)" outlined />
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup @click="resetForm" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog penolakan -->
    <q-dialog v-model="rejectDialog" persistent>
      <q-card>
        <q-card-section>
          <div class="text-h6">Alasan Penolakan</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="rejectReason" label="Alasan" type="textarea" outlined autofocus />
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
import { useLogbookStore } from 'stores/pkl/logbook'
import { useSiswaStore } from 'stores/kesiswaan/siswa'
import { useAuthStore } from 'stores/auth-store'

const $q = useQuasar()
const logbookStore = useLogbookStore()
const siswaStore = useSiswaStore()
const authStore = useAuthStore()

const selectedSiswa = ref(null)
const formDialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const rejectDialog = ref(false)
let currentRejectId = null
const rejectReason = ref('')

const form = ref({
  siswa_id: null,
  tanggal: new Date().toISOString().split('T')[0],
  aktivitas: '',
  hasil: '',
  keterangan: ''
})

const columns = [
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal', align: 'left' },
  { name: 'aktivitas', label: 'Aktivitas', field: 'aktivitas', align: 'left' },
  { name: 'hasil', label: 'Hasil', field: 'hasil', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

const siswaOptions = computed(() => siswaStore.list.map(s => ({ id: s.id, nama: s.nama })))

const filteredLogbook = computed(() => {
  if (!selectedSiswa.value) return []
  return logbookStore.list.filter(l => l.siswa_id === selectedSiswa.value.id)
})

function openForm(row = null) {
  if (!selectedSiswa.value && !row) {
    $q.notify({ type: 'warning', message: 'Pilih siswa terlebih dahulu' })
    return
  }
  if (row) {
    editMode.value = true
    editId.value = row.id
    form.value = {
      siswa_id: row.siswa_id,
      tanggal: row.tanggal,
      aktivitas: row.aktivitas,
      hasil: row.hasil,
      keterangan: row.keterangan
    }
  } else {
    editMode.value = false
    resetForm()
    form.value.siswa_id = selectedSiswa.value.id
  }
  formDialog.value = true
}

function resetForm() {
  form.value = {
    siswa_id: null,
    tanggal: new Date().toISOString().split('T')[0],
    aktivitas: '',
    hasil: '',
    keterangan: ''
  }
}

function submitLog() {
  if (editMode.value) {
    logbookStore.update(editId.value, form.value)
    $q.notify({ type: 'positive', message: 'Logbook diupdate' })
  } else {
    logbookStore.tambah({ ...form.value, status: 'pending' })
    $q.notify({ type: 'positive', message: 'Logbook ditambahkan' })
  }
  formDialog.value = false
  resetForm()
}

function approveLog(id) {
  $q.dialog({ title: 'Konfirmasi', message: 'Setujui logbook ini?', cancel: true }).onOk(() => {
    logbookStore.approve(id)
    $q.notify({ type: 'positive', message: 'Logbook disetujui' })
  })
}

function rejectLog(id) {
  currentRejectId = id
  rejectReason.value = ''
  rejectDialog.value = true
}

function executeReject() {
  if (!rejectReason.value.trim()) {
    $q.notify({ type: 'negative', message: 'Alasan penolakan harus diisi' })
    return
  }
  logbookStore.reject(currentRejectId, rejectReason.value)
  $q.notify({ type: 'negative', message: 'Logbook ditolak' })
  rejectDialog.value = false
  currentRejectId = null
  rejectReason.value = ''
}

function confirmDelete(id) {
  $q.dialog({ title: 'Konfirmasi', message: 'Hapus logbook ini?', cancel: true }).onOk(() => {
    logbookStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Logbook dihapus' })
  })
}

onMounted(() => {
  siswaStore.loadData()
  logbookStore.loadData()
})
</script>
