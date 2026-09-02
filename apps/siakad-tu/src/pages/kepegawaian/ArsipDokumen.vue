<template>
  <q-page class="q-pa-md">
    <div class="row justify-between q-mb-md">
      <div class="text-h5">Arsip Dokumen GTK</div>
      <q-btn color="primary" icon="add" label="Upload Dokumen" @click="openFormDialog" />
    </div>

    <!-- Alert dokumen mendekati kadaluarsa -->
    <q-banner v-if="expiringDocs.length" rounded class="bg-warning text-black q-mb-md">
      <template v-slot:avatar>
        <q-icon name="warning" />
      </template>
      Terdapat {{ expiringDocs.length }} dokumen yang akan kadaluarsa dalam 30 hari ke depan.
      <q-btn flat label="Lihat" color="black" @click="showExpiringList" />
    </q-banner>

    <q-table :rows="arsipStore.list" :columns="columns" row-key="id" flat bordered dense>
      <template v-slot:body-cell-tanggal_kadaluarsa="props">
        <q-td :props="props">
          <span :class="{ 'text-negative': isExpiringSoon(props.row.tanggal_kadaluarsa) }">
            {{ props.row.tanggal_kadaluarsa || '-' }}
          </span>
        </q-td>
      </template>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <q-btn flat round dense icon="edit" color="info" @click="openFormDialog(props.row)" class="q-mr-sm">
            <q-tooltip>Edit</q-tooltip>
          </q-btn>
          <q-btn flat round dense icon="download" color="primary" @click="downloadFile(props.row)" class="q-mr-sm">
            <q-tooltip>Download</q-tooltip>
          </q-btn>
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)">
            <q-tooltip>Hapus</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- Dialog form tambah/edit -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 550px">
        <q-card-section>
          <div class="text-h6">{{ editMode ? 'Edit Dokumen' : 'Upload Dokumen' }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitForm" class="q-gutter-md">
            <q-select v-model="form.guru_id" :options="guruOptions" label="Pilih Guru" option-value="id"
              option-label="nama" outlined :rules="[v => !!v]" />
            <q-input v-model="form.jenis_dokumen" label="Jenis Dokumen" outlined :rules="[v => !!v]" />
            <q-input v-model="form.nomor_dokumen" label="Nomor Dokumen" outlined />
            <q-input v-model="form.tanggal_terbit" label="Tanggal Terbit" type="date" outlined />
            <q-input v-model="form.tanggal_kadaluarsa" label="Tanggal Kadaluarsa (kosongkan jika tidak ada)" type="date"
              outlined />
            <q-input v-model="form.catatan" label="Catatan" type="textarea" rows="2" outlined />
            <q-file v-model="form.file" label="File (PDF/Image)" accept=".pdf,.jpg,.jpeg,.png" outlined
              @update:model-value="onFileSelected">
              <template v-slot:prepend>
                <q-icon name="attach_file" />
              </template>
            </q-file>
            <div v-if="form.file_name" class="text-caption">File saat ini: {{ form.file_name }}</div>
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup @click="resetForm" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog daftar dokumen kadaluarsa -->
    <q-dialog v-model="expiryDialog">
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">Dokumen Mendekati Kadaluarsa</div>
        </q-card-section>
        <q-card-section>
          <q-list dense>
            <q-item v-for="doc in expiringDocs" :key="doc.id">
              <q-item-section>
                <q-item-label>{{ doc.jenis_dokumen }} - {{ doc.guru_nama }}</q-item-label>
                <q-item-label caption>Kadaluarsa: {{ formatDate(doc.tanggal_kadaluarsa) }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Tutup" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useArsipStore } from 'stores/kepegawaian/arsip'
import { useGuruStore } from 'stores/kepegawaian/guru'

const $q = useQuasar()
const arsipStore = useArsipStore()
const guruStore = useGuruStore()

const formDialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const expiryDialog = ref(false)

const form = ref({
  guru_id: null,
  jenis_dokumen: '',
  nomor_dokumen: '',
  tanggal_terbit: '',
  tanggal_kadaluarsa: null,
  catatan: '',
  file: null,
  file_base64: '',
  file_name: ''
})

const columns = [
  { name: 'guru_nama', label: 'Guru', field: 'guru_nama', align: 'left' },
  { name: 'jenis_dokumen', label: 'Jenis', field: 'jenis_dokumen', align: 'left' },
  { name: 'nomor_dokumen', label: 'Nomor', field: 'nomor_dokumen', align: 'left' },
  { name: 'tanggal_terbit', label: 'Terbit', field: 'tanggal_terbit', align: 'left' },
  { name: 'tanggal_kadaluarsa', label: 'Kadaluarsa', field: 'tanggal_kadaluarsa', align: 'left' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

const guruOptions = computed(() => guruStore.list.map(g => ({ id: g.id, nama: g.nama })))
const expiringDocs = computed(() => arsipStore.getExpiringDocs(30))

function isExpiringSoon(dateStr) {
  if (!dateStr) return false
  const expiry = new Date(dateStr)
  const today = new Date()
  const diffDays = Math.ceil((expiry - today) / (1000 * 60 * 60 * 24))
  return diffDays <= 30 && diffDays >= 0
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID')
}

function onFileSelected(file) {
  if (!file) return
  const reader = new FileReader()
  reader.onload = (e) => {
    form.value.file_base64 = e.target.result
    form.value.file_name = file.name
  }
  reader.readAsDataURL(file)
}

function openFormDialog(row = null) {
  if (row) {
    editMode.value = true
    editId.value = row.id
    form.value = {
      guru_id: row.guru_id,
      jenis_dokumen: row.jenis_dokumen,
      nomor_dokumen: row.nomor_dokumen,
      tanggal_terbit: row.tanggal_terbit,
      tanggal_kadaluarsa: row.tanggal_kadaluarsa,
      catatan: row.catatan,
      file: null,
      file_base64: row.file_base64,
      file_name: row.file_name
    }
  } else {
    editMode.value = false
    editId.value = null
    resetForm()
  }
  formDialog.value = true
}

function resetForm() {
  form.value = {
    guru_id: null,
    jenis_dokumen: '',
    nomor_dokumen: '',
    tanggal_terbit: '',
    tanggal_kadaluarsa: null,
    catatan: '',
    file: null,
    file_base64: '',
    file_name: ''
  }
}

function submitForm() {
  if (!form.value.guru_id || !form.value.jenis_dokumen) {
    $q.notify({ type: 'negative', message: 'Lengkapi field wajib' })
    return
  }
  const guru = guruStore.list.find(g => g.id === form.value.guru_id)
  if (!guru) return

  const dokumenData = {
    guru_id: form.value.guru_id,
    guru_nama: guru.nama,
    jenis_dokumen: form.value.jenis_dokumen,
    nomor_dokumen: form.value.nomor_dokumen,
    tanggal_terbit: form.value.tanggal_terbit,
    tanggal_kadaluarsa: form.value.tanggal_kadaluarsa || null,
    catatan: form.value.catatan,
    file_name: form.value.file_name,
    file_base64: form.value.file_base64
  }

  if (editMode.value) {
    arsipStore.update(editId.value, dokumenData)
    $q.notify({ type: 'positive', message: 'Dokumen diupdate' })
  } else {
    arsipStore.tambah(dokumenData)
    $q.notify({ type: 'positive', message: 'Dokumen diupload' })
  }
  formDialog.value = false
  resetForm()
}

function downloadFile(dokumen) {
  if (!dokumen.file_base64) {
    $q.notify({ type: 'warning', message: 'File tidak tersedia' })
    return
  }
  const link = document.createElement('a')
  link.href = dokumen.file_base64
  link.download = dokumen.file_name
  link.click()
}

function confirmDelete(id) {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus dokumen ini?',
    cancel: true
  }).onOk(() => {
    arsipStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Dokumen dihapus' })
  })
}

function showExpiringList() {
  expiryDialog.value = true
}

onMounted(() => {
  guruStore.loadData()
  arsipStore.loadData()
})
</script>
